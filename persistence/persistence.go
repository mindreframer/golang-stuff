package persistence

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var logfileReg = regexp.MustCompile("^(\\d+)\\.(snap|log)$")

const (
	stopped = iota
	recording
	playing
)

const (
	snapSuffix       = "snap"
	logSuffix        = "log"
	unfinishedSuffix = "unfinished"
)

// Op is a simple get/put/clear or configuration operation to log or replay.
type Op struct {
	Key           []byte
	SubKey        []byte
	Value         []byte
	Timestamp     int64
	Put           bool
	Clear         bool
	Configuration map[string]string
}

type logfile struct {
	timestamp time.Time
	filename  string
	suffix    string
	file      *os.File
	encoder   *gob.Encoder
	decoder   *gob.Decoder
}

func createLogfile(dir, suffix string) (rval *logfile) {
	rval = &logfile{}
	rval.timestamp = time.Now()
	rval.suffix = suffix
	rval.filename = filepath.Join(dir, fmt.Sprintf("%v.%v", rval.timestamp.UnixNano(), suffix))
	return
}

func parseLogfile(file string) (rval *logfile, err error) {
	match := logfileReg.FindStringSubmatch(filepath.Base(file))
	if match == nil {
		err = fmt.Errorf("%v does not match %v", file, logfileReg)
		return
	}
	nanos, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return
	}
	rval = &logfile{}
	rval.suffix = match[2]
	rval.timestamp = time.Unix(0, nanos)
	rval.filename = file
	return
}

func (self *logfile) play(operate Operate) {
	if self == nil {
		return
	}
	self.read()
	defer self.close()
	var err error
	for {
		var op Op
		err = self.decoder.Decode(&op)
		if err != nil {
			break
		}
		operate(op)
	}
	if err != io.EOF {
		panic(err)
	}
}

func (self *logfile) read() *logfile {
	var err error
	self.file, err = os.Open(self.filename)
	if err != nil {
		panic(err)
	}
	self.decoder = gob.NewDecoder(self.file)
	return self
}

func (self *logfile) write() *logfile {
	var err error
	self.file, err = os.Create(self.filename)
	if err != nil {
		panic(err)
	}
	self.encoder = gob.NewEncoder(self.file)
	return self
}

func (self *logfile) close() {
	self.file.Close()
}

type logfiles []*logfile

func (self logfiles) Len() int {
	return len(self)
}
func (self logfiles) Less(i, j int) bool {
	return self[i].timestamp.Before(self[j].timestamp)
}
func (self logfiles) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// Operate is a function that operates on an Op, for replay purposes.
// It is supposed to insert Ops with the Put flag, Clear data if the Clear flag is set, handle configuration changes or delete data.
type Operate func(o Op)

// Logger is something that can log or replay Ops.
type Logger struct {
	ops      chan Op
	stops    chan chan bool
	dir      string
	state    int32
	snapping int32
	maxSize  int64
	suffix   string
	cond     *sync.Cond
	lock     *sync.Mutex
}

// NewLogger will return a Logger that will dump data into dir, or replay data from dir.
func NewLogger(dir string) *Logger {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	lock := new(sync.Mutex)
	return &Logger{
		ops:    make(chan Op),
		stops:  make(chan chan bool),
		dir:    dir,
		suffix: logSuffix,
		lock:   lock,
		cond:   sync.NewCond(lock),
	}
}

func (self *Logger) hasState(s int32) bool {
	return atomic.LoadInt32(&self.state) == s
}
func (self *Logger) changeState(old, neu int32) bool {
	return atomic.CompareAndSwapInt32(&self.state, old, neu)
}

func (self *Logger) setSuffix(s string) *Logger {
	self.suffix = s
	return self
}

// Limit will limit the size of the last logfile to maxSize bytes.
// When the last logfile is bigger than maxSize, it will merge the last snapshot and any logfile created after it into a new snapshot, 
// and start a new logfile to continue. All this will happen transparently in a separate goroutine.
func (self *Logger) Limit(maxSize int64) *Logger {
	self.maxSize = maxSize
	return self
}

func (self *Logger) logfiles() (result logfiles) {
	dir, err := os.Open(self.dir)
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	files, err := dir.Readdirnames(0)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		var logf *logfile
		logf, err = parseLogfile(filepath.Join(self.dir, file))
		if err == nil {
			result = append(result, logf)
		}
	}
	return
}

func (self *Logger) latest() (latestSnapshot *logfile, logs logfiles) {
	for _, logf := range self.logfiles() {
		if logf.suffix == snapSuffix {
			if latestSnapshot == nil || latestSnapshot.timestamp.After(logf.timestamp) {
				latestSnapshot = logf
			}
		}
	}
	for _, logf := range self.logfiles() {
		if logf.suffix == logSuffix {
			if latestSnapshot == nil || latestSnapshot.timestamp.Before(logf.timestamp) {
				logs = append(logs, logf)
			}
		}
	}
	sort.Sort(logs)
	return
}

// Recording returns true if this Logger is currently recording (as opposed to replaying or idling).
func (self *Logger) Recording() bool {
	return self.hasState(recording)
}

// Play will replay the latest snapshot and all logfiles created after it using the provided operate.
func (self *Logger) Play(operate Operate) {
	if self.changeState(stopped, playing) {
		defer self.changeState(playing, stopped)
		snapshot, logs := self.latest()
		snapshot.play(operate)
		for _, logf := range logs {
			logf.play(operate)
		}
	}
}

// Stop will stop this Logger. It will not return until all running recordings or snaphots are finished.
func (self *Logger) Stop() *Logger {
	if self.hasState(recording) {
		stop := make(chan bool)
		self.stops <- stop
		<-stop
		for atomic.LoadInt32(&self.snapping) == 1 {
			self.lock.Lock()
			self.cond.Wait()
			self.lock.Unlock()
		}
	} else {
		panic(fmt.Errorf("%v is not in state recording", self))
	}
	return self
}

func (self *Logger) clearOlderThan(t time.Time) {
	for _, logf := range self.logfiles() {
		if logf.timestamp.Before(t) {
			if err := os.Remove(logf.filename); err != nil {
				log.Printf("failed removing %v: %v", logf.filename, err)
			}
		}
	}
}

// Clear will stop this Logger and remove all snapshots or logfiles older than now.
func (self *Logger) Clear() {
	self.Stop()
	self.clearOlderThan(time.Now())
	<-self.Record()
}

func (self *Logger) snapshot(snap *logfile, files logfiles) {
	byteCompressor := make(map[string]Op)
	treeCompressor := make(map[string]map[string]Op)
	var latestConf *Op
	confCompressor := make(map[string]Op)
	var subMap map[string]Op
	var ok bool
	operate := func(op Op) {
		if op.Configuration != nil {
			if op.Key == nil {
				latestConf = &op
			} else {
				confCompressor[string(op.Key)] = op
			}
		} else if op.Put {
			if op.SubKey == nil {
				byteCompressor[string(op.Key)] = op
			} else {
				subMap, ok = treeCompressor[string(op.Key)]
				if !ok {
					subMap = make(map[string]Op)
					treeCompressor[string(op.Key)] = subMap
				}
				subMap[string(op.SubKey)] = op
			}
		} else {
			if op.SubKey == nil {
				if op.Clear {
					if op.Key == nil {
						byteCompressor = make(map[string]Op)
					} else {
						delete(treeCompressor, string(op.Key))
					}
				} else {
					delete(byteCompressor, string(op.Key))
				}
			} else {
				subMap, ok = treeCompressor[string(op.Key)]
				if ok {
					delete(subMap, string(op.SubKey))
					if len(subMap) == 0 {
						delete(treeCompressor, string(op.Key))
					}
				}
			}
		}
	}
	snap.play(operate)
	for _, logf := range files {
		logf.play(operate)
	}
	if latestConf != nil {
		self.Dump(*latestConf)
	}
	for _, op := range confCompressor {
		self.Dump(op)
	}
	for _, op := range byteCompressor {
		self.Dump(op)
	}
	for _, subMap := range treeCompressor {
		for _, op := range subMap {
			self.Dump(op)
		}
	}
}

func (self *Logger) snapshotAndDelete(oldrec *logfile, p chan *logfile, snapping *int32) {
	defer atomic.StoreInt32(snapping, 0)
	defer self.cond.Broadcast()
	latestSnapshot, logfiles := self.latest()
	snapshotter := NewLogger(self.dir).setSuffix(unfinishedSuffix)
	snapshotfile := <-snapshotter.Record()
	p <- snapshotfile
	snapshotter.snapshot(latestSnapshot, logfiles)
	snapshotter.Stop()
	if err := os.Rename(snapshotfile.filename, filepath.Join(self.dir, fmt.Sprintf("%v.%v", snapshotfile.timestamp.UnixNano(), snapSuffix))); err != nil {
		panic(err)
	}
	self.clearOlderThan(snapshotfile.timestamp)
}

func (self *Logger) swap(fi *os.FileInfo, err *error, rec *logfile) *logfile {
	if atomic.LoadInt32(&self.snapping) == 0 {
		if *fi, *err = os.Stat(rec.filename); *err != nil {
			panic(*err)
		}
		if (*fi).Size() > self.maxSize {
			rec.close()
			started := make(chan *logfile)
			atomic.StoreInt32(&self.snapping, 1)
			go self.snapshotAndDelete(rec, started, &self.snapping)
			<-started
			rec = createLogfile(self.dir, self.suffix)
			rec.write()
		}
	}
	return rec
}

// Record will make this Logger start recording.
func (self *Logger) Record() (rval chan *logfile) {
	if !self.changeState(stopped, recording) {
		panic(fmt.Errorf("%v unable to change state from stopped to recording", self))
	}
	rval = make(chan *logfile, 1)
	go self.record(rval)
	return
}

func (self *Logger) record(p chan *logfile) {
	var err error
	var op Op
	var fi os.FileInfo
	var stop chan bool

	rec := createLogfile(self.dir, self.suffix)
	rec.write()
	p <- rec
	defer rec.close()

	for {
		if self.maxSize != 0 {
			rec = self.swap(&fi, &err, rec)
		}

		select {
		case op = <-self.ops:
			if err = rec.encoder.Encode(op); err != nil {
				panic(err)
			}
		case stop = <-self.stops:
			if !self.changeState(recording, stopped) {
				panic(fmt.Errorf("%v unable to change state from recording to stopped", self))
			}
			stop <- true
			return
		}
		select {
		case stop = <-self.stops:
			if !self.changeState(recording, stopped) {
				panic(fmt.Errorf("%v unable to change state from recording to stopped", self))
			}
			stop <- true
			return
		default:
		}
	}
}

// Dump will accept an operation if this Logger is recording, and dump it into a logfile.
func (self *Logger) Dump(o Op) {
	if !self.hasState(recording) {
		panic(fmt.Errorf("%v is not recording", self))
	}
	self.ops <- o
}

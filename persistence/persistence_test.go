package persistence

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
)

type testmap struct {
	l *sync.RWMutex
	m map[string]string
	p *Logger
}

func newTestmap() (rval testmap) {
	rval.m = make(map[string]string)
	rval.l = new(sync.RWMutex)
	rval.p = NewLogger("test3")
	return
}

func (self testmap) playback() {
	self.p.Play(self.operator())
}

func (self testmap) record() {
	self.p.Limit(1024)
	self.p.Record()
}
func (self testmap) put(k, v string) {
	self.l.Lock()
	defer self.l.Unlock()
	self.p.Dump(Op{
		Key:   []byte(k),
		Value: []byte(v),
		Put:   true,
	})
	self.m[k] = v
}
func (self testmap) del(k string) {
	self.l.Lock()
	defer self.l.Unlock()
	self.p.Dump(Op{
		Key: []byte(k),
	})
	delete(self.m, k)
}
func (self testmap) operator() Operate {
	return func(o Op) {
		self.l.Lock()
		defer self.l.Unlock()
		if o.Put {
			self.m[string(o.Key)] = string(o.Value)
		} else {
			delete(self.m, string(o.Key))
		}
	}
}

func operator(ary *[]Op) Operate {
	return func(o Op) {
		*ary = append(*ary, o)
	}
}

func TestRecordPlay(t *testing.T) {
	os.RemoveAll("test1")
	p := NewLogger("test1")
	p.Record()
	op := Op{
		Key:       []byte("a"),
		Value:     []byte("1"),
		Timestamp: 1,
	}
	p.Dump(op)
	p.Stop()
	var ary []Op
	p.Play(operator(&ary))
	if !reflect.DeepEqual(ary, []Op{op}) {
		t.Errorf("%+v should be %+v", ary, []Op{op})
	}
}

func TestSwap(t *testing.T) {
	os.RemoveAll("test3")
	tm := newTestmap()
	tm.record()
	for i := 0; i < 1000; i++ {
		tm.put(fmt.Sprint(i), fmt.Sprint(i))
	}
	for i := 0; i < 1000; i += 3 {
		tm.del(fmt.Sprint(i))
	}
	tm.p.Stop()

	dir, err := os.Open("test3")
	if err != nil {
		t.Fatal(err)
	}
	files, err := dir.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("%v should only be two files", files)
	}

	tm2 := newTestmap()
	tm2.playback()
	if !reflect.DeepEqual(tm.m, tm2.m) {
		t.Errorf("%v should be equal to %v", tm2.m, tm.m)
	}
}

func BenchmarkRecord(b *testing.B) {
	b.StopTimer()
	os.RemoveAll("test2")
	p := NewLogger("test2")
	p.Record()
	op := Op{
		Key:       []byte("a"),
		Value:     []byte("1"),
		Timestamp: 1,
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p.Dump(op)
	}
}

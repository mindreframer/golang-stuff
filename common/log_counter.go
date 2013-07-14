package common

import (
	"encoding/json"
	steno "github.com/cloudfoundry/gosteno"
	"sync"
)

type LogCounter struct {
	sync.Mutex
	counts map[string]int
}

func NewLogCounter() *LogCounter {
	lc := &LogCounter{
		counts: make(map[string]int),
	}
	return lc
}

func (lc *LogCounter) AddRecord(record *steno.Record) {
	lc.Lock()
	defer lc.Unlock()

	lc.counts[record.Level.Name] += 1
}

func (lc *LogCounter) GetCount(key string) int {
	lc.Lock()
	defer lc.Unlock()

	return lc.counts[key]
}

func (lc *LogCounter) Flush()                     {}
func (lc *LogCounter) SetCodec(codec steno.Codec) {}

func (lc *LogCounter) GetCodec() steno.Codec {
	return nil
}

func (lc *LogCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(lc.counts)
}

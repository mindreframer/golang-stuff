package common

import (
	"encoding/json"
	"sync"
)

type GenericVarz struct {
	// Static common metrics
	NumCores int `json:"num_cores"`

	// Dynamic common metrics
	MemStat int64   `json:"mem"`
	Cpu     float64 `json:"cpu"`

	Uptime    Duration    `json:"uptime"`
	LogCounts *LogCounter `json:"log_counts"`
}

type Varz struct {
	sync.Mutex
	GenericVarz
	UniqueVarz interface{} // Every component's unique metrics
}

func transform(x interface{}, y *map[string]interface{}) error {
	var b []byte
	var err error

	b, err = json.Marshal(x)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, y)
	if err != nil {
		return err
	}

	return nil
}

func (v *Varz) MarshalJSON() ([]byte, error) {
	r := make(map[string]interface{})

	var err error

	err = transform(v.UniqueVarz, &r)
	if err != nil {
		return nil, err
	}

	err = transform(v.GenericVarz, &r)
	if err != nil {
		return nil, err
	}

	err = transform(Component, &r)
	if err != nil {
		return nil, err
	}

	return json.Marshal(r)
}

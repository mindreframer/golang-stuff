package common

import (
	"encoding/json"
	steno "github.com/cloudfoundry/gosteno"
	. "launchpad.net/gocheck"
)

type VarzSuite struct {
}

var _ = Suite(&VarzSuite{})

func (s *VarzSuite) SetUpTest(c *C) {
	Component = VcapComponent{
		Credentials: []string{"foo", "bar"},
		Config:      map[string]interface{}{"ip": "localhost", "port": 8080},
	}
}

func (s *VarzSuite) TearDownTest(c *C) {
	Component = VcapComponent{}
}

func (s *VarzSuite) TestEmptyVarz(c *C) {
	varz := &Varz{}
	varz.LogCounts = NewLogCounter()

	bytes, err := json.Marshal(varz)
	c.Assert(err, IsNil)

	data := make(map[string]interface{})
	err = json.Unmarshal(bytes, &data)
	c.Assert(err, IsNil)

	members := []string{
		"type",
		"index",
		"host",
		"credentials",
		"config",
		"start",
		"uuid",
		"uptime",
		"num_cores",
		"mem",
		"cpu",
		"log_counts",
	}

	for _, key := range members {
		if _, ok := data[key]; !ok {
			c.Fatalf(`member "%s" not found`, key)
		}
	}
}

func (s *VarzSuite) TestLogCounts(c *C) {
	varz := &Varz{}
	varz.LogCounts = NewLogCounter()

	varz.LogCounts.AddRecord(&steno.Record{Level: steno.LOG_INFO})

	bytes, _ := json.Marshal(varz)
	data := make(map[string]interface{})
	json.Unmarshal(bytes, &data)

	counts := data["log_counts"].(map[string]interface{})
	count := counts["info"]

	c.Assert(count, Equals, 1.0)
}

func (s *VarzSuite) TestTransformStruct(c *C) {
	component := struct {
		Type  string `json:"type"`
		Index int    `json:"index"`
	}{
		Type:  "Router",
		Index: 1,
	}

	m := make(map[string]interface{})
	transform(component, &m)
	c.Assert(m["type"], Equals, "Router")
	c.Assert(m["index"], Equals, float64(1))
}

func (s *VarzSuite) TestTransformMap(c *C) {
	data := map[string]interface{}{"type": "Dea", "index": 1}

	m := make(map[string]interface{})
	transform(data, &m)
	c.Assert(m["type"], Equals, "Dea")
	c.Assert(m["index"], Equals, float64(1))
}

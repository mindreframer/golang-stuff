package common

import (
	"encoding/json"
	"fmt"
	. "launchpad.net/gocheck"
	"time"
)

type DurationSuite struct {
}

var _ = Suite(&DurationSuite{})

func (s *DurationSuite) TestJsonInterface(c *C) {
	d := Duration(123456)
	var i interface{} = &d

	_, ok := i.(json.Marshaler)
	c.Assert(ok, Equals, true)

	_, ok = i.(json.Unmarshaler)
	c.Assert(ok, Equals, true)
}

func (s *DurationSuite) TestMarshalJSON(c *C) {
	d := Duration(time.Hour*36 + time.Second*10)
	b, err := json.Marshal(d)
	c.Assert(err, IsNil)
	c.Assert(string(b), Equals, `"1d:12h:0m:10s"`)
}

func (s *DurationSuite) TestUnmarshalJSON(c *C) {
	d := Duration(time.Hour*36 + time.Second*20)
	b, err := json.Marshal(d)
	c.Assert(err, IsNil)

	var dd Duration
	dd.UnmarshalJSON(b)
	c.Assert(dd, Equals, d)
}

func (s *DurationSuite) TestTimeMarshalJSON(c *C) {
	n := time.Now()
	f := "2006-01-02 15:04:05 -0700"

	t := Time(n)
	b, e := json.Marshal(t)
	c.Assert(e, IsNil)
	c.Assert(string(b), Equals, fmt.Sprintf(`"%s"`, n.Format(f)))
}

func (s *DurationSuite) TestTimeUnmarshalJSON(c *C) {
	t := Time(time.Unix(time.Now().Unix(), 0)) // The precision of Time is 'second'
	b, err := json.Marshal(t)
	c.Assert(err, IsNil)

	var tt Time
	err = tt.UnmarshalJSON(b)
	c.Assert(err, IsNil)
	c.Assert(tt, Equals, t)
}

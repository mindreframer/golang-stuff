package util

import (
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
	"strconv"
)

type PidfileSuite struct {
	path    string
	pidfile string
}

var _ = Suite(&PidfileSuite{})

func (s *PidfileSuite) SetUpTest(c *C) {
	x, err := ioutil.TempDir("", "PidFileSuite")
	c.Assert(err, IsNil)

	s.path = x
	s.pidfile = path.Join(s.path, "pidfile")
}

func (s *PidfileSuite) TearDownTest(c *C) {
	err := os.RemoveAll(s.path)
	c.Assert(err, IsNil)
}

func (s *PidfileSuite) assertPidfileNonzero(c *C) {
	x, err := ioutil.ReadFile(s.pidfile)
	c.Assert(err, IsNil)

	y, err := strconv.Atoi(string(x))
	c.Assert(err, IsNil)
	c.Assert(y, Not(Equals), 0)
}

func (s *PidfileSuite) TestWritePidfile(c *C) {
	err := WritePidFile(s.pidfile)
	c.Assert(err, IsNil)

	s.assertPidfileNonzero(c)
}

func (s *PidfileSuite) TestWritePidfileOverwrites(c *C) {
	err := ioutil.WriteFile(s.pidfile, []byte("0"), 0644)
	c.Assert(err, IsNil)

	err = WritePidFile(s.pidfile)
	c.Assert(err, IsNil)

	s.assertPidfileNonzero(c)
}

func (s *PidfileSuite) TestWritePidfileReturnsError(c *C) {
	err := os.RemoveAll(s.path)
	c.Assert(err, IsNil)

	err = WritePidFile(s.pidfile)
	c.Assert(err, Not(IsNil))
}

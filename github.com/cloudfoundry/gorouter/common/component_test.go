package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net"
	"net/http"
	"runtime"
)

type ComponentSuite struct {
	Component *VcapComponent
}

var _ = Suite(&ComponentSuite{})

type MarshalableValue struct {
	Value map[string]string
}

func (m *MarshalableValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

func (s *ComponentSuite) SetUpTest(c *C) {
	port, err := GrabEphemeralPort()
	c.Assert(err, IsNil)

	s.Component = &VcapComponent{
		Host:        fmt.Sprintf("127.0.0.1:%d", port),
		Credentials: []string{"username", "password"},
	}
}

func (s *ComponentSuite) TearDownTest(c *C) {}

func (s *ComponentSuite) TestInfoRouteAccessUnauthorized(c *C) {
	path := "/test"

	s.Component.InfoRoutes = map[string]json.Marshaler{
		path: &MarshalableValue{Value: map[string]string{"key": "value"}},
	}
	s.serveComponent(c)

	req := s.buildGetRequest(c, path)
	code, _, _ := s.doGetRequest(c, req)
	c.Check(code, Equals, 401)

	req = s.buildGetRequest(c, path)
	req.SetBasicAuth("username", "incorrect-password")
	code, _, _ = s.doGetRequest(c, req)
	c.Check(code, Equals, 401)

	req = s.buildGetRequest(c, path)
	req.SetBasicAuth("incorrect-username", "password")
	code, _, _ = s.doGetRequest(c, req)
	c.Check(code, Equals, 401)
}

func (s *ComponentSuite) TestInfoRouteAccessAuthorized(c *C) {
	path := "/test"

	s.Component.InfoRoutes = map[string]json.Marshaler{
		path: &MarshalableValue{Value: map[string]string{"key": "value"}},
	}
	s.serveComponent(c)

	req := s.buildGetRequest(c, path)
	req.SetBasicAuth("username", "password")

	code, header, body := s.doGetRequest(c, req)
	c.Check(code, Equals, 200)
	c.Check(header.Get("Content-Type"), Equals, "application/json")
	c.Check(body, Equals, `{"key":"value"}`+"\n")
}

func (s *ComponentSuite) TestInfoRouteAccessNonExistent(c *C) {
	s.serveComponent(c)

	req := s.buildGetRequest(c, "/non-existent-path")
	req.SetBasicAuth("username", "password")

	code, _, _ := s.doGetRequest(c, req)
	c.Check(code, Equals, 404)
}

func (s *ComponentSuite) serveComponent(c *C) {
	go s.Component.ListenAndServe()

	for i := 0; i < 200; i++ {
		// Yield to component's server listen goroutine
		runtime.Gosched()

		conn, err := net.Dial("tcp", s.Component.Host)
		if err == nil {
			conn.Close()
			return
		}
	}
	panic("Could not connect to vcap.Component")
}

func (s *ComponentSuite) buildGetRequest(c *C, path string) *http.Request {
	req, err := http.NewRequest("GET", "http://"+s.Component.Host+path, nil)
	c.Assert(err, IsNil)
	return req
}

func (s *ComponentSuite) doGetRequest(c *C, req *http.Request) (int, http.Header, string) {
	var client http.Client
	var resp *http.Response
	var err error

	resp, err = client.Do(req)
	c.Assert(err, IsNil)
	c.Assert(resp, Not(IsNil))

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	c.Assert(err, IsNil)

	return resp.StatusCode, resp.Header, string(body)
}

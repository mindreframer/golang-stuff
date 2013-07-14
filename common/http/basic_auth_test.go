package http

import (
	. "launchpad.net/gocheck"
	"net"
	"net/http"
)

type BasicAuthSuite struct {
	Listener net.Listener
}

var _ = Suite(&BasicAuthSuite{})

func (s *BasicAuthSuite) TearDownTest(c *C) {
	if s.Listener != nil {
		s.Listener.Close()
	}
}

func (s *BasicAuthSuite) Bootstrap(x Authenticator) *http.Request {
	var err error

	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	y := &BasicAuth{http.HandlerFunc(h), x}

	z := &http.Server{Handler: y}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	go z.Serve(l)

	// Keep listener around such that test teardown can close it
	s.Listener = l

	r, err := http.NewRequest("GET", "http://"+l.Addr().String(), nil)
	if err != nil {
		panic(err)
	}

	return r
}

func (s *BasicAuthSuite) TestNoCredentials(c *C) {
	req := s.Bootstrap(nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	c.Check(resp.StatusCode, Equals, http.StatusUnauthorized)
}

func (s *BasicAuthSuite) TestInvalidHeader(c *C) {
	req := s.Bootstrap(nil)

	req.Header.Set("Authorization", "invalid")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	c.Check(resp.StatusCode, Equals, http.StatusUnauthorized)
}

func (s *BasicAuthSuite) TestBadCredentials(c *C) {
	f := func(u, p string) bool {
		c.Check(u, Equals, "user")
		c.Check(p, Equals, "bad")
		return false
	}

	req := s.Bootstrap(f)

	req.SetBasicAuth("user", "bad")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	c.Check(resp.StatusCode, Equals, http.StatusUnauthorized)
}

func (s *BasicAuthSuite) TestGoodCredentials(c *C) {
	f := func(u, p string) bool {
		c.Check(u, Equals, "user")
		c.Check(p, Equals, "good")
		return true
	}

	req := s.Bootstrap(f)

	req.SetBasicAuth("user", "good")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	c.Check(resp.StatusCode, Equals, http.StatusOK)
}

package router

import (
	"bufio"
	"code.google.com/p/gomock/gomock"
	"encoding/json"
	"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net"
	"net/http"
	"github.com/cloudfoundry/gorouter/test"
	"strconv"
	"strings"
	"time"
)

type connHandler func(*conn)

type nullVarz struct{}

func (_ nullVarz) MarshalJSON() ([]byte, error) { return json.Marshal(nil) }

func (_ nullVarz) CaptureBadRequest(req *http.Request)                                    {}
func (_ nullVarz) CaptureBackendRequest(b *Backend, req *http.Request)                    {}
func (_ nullVarz) CaptureBackendResponse(b *Backend, res *http.Response, d time.Duration) {}

type conn struct {
	net.Conn

	c *C

	br *bufio.Reader
	bw *bufio.Writer
}

func newConn(x net.Conn, c *C) *conn {
	return &conn{
		Conn: x,
		c:    c,
		br:   bufio.NewReader(x),
		bw:   bufio.NewWriter(x),
	}
}

func (x *conn) ReadRequest() (*http.Request, string) {
	req, err := http.ReadRequest(x.br)
	x.c.Assert(err, IsNil)

	b, err := ioutil.ReadAll(req.Body)
	x.c.Assert(err, IsNil)

	return req, string(b)
}

func (x *conn) NewRequest(method, urlStr string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, urlStr, body)
	x.c.Assert(err, IsNil)
	return req
}

func (x *conn) WriteRequest(req *http.Request) {
	err := req.Write(x.bw)
	x.c.Assert(err, IsNil)
	x.bw.Flush()
}

func (x *conn) ReadResponse() (*http.Response, string) {
	resp, err := http.ReadResponse(x.br, &http.Request{})
	x.c.Assert(err, IsNil)

	b, err := ioutil.ReadAll(resp.Body)
	x.c.Assert(err, IsNil)

	return resp, string(b)
}

func newResponse(status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}
}

func (x *conn) WriteResponse(resp *http.Response) {
	err := resp.Write(x.bw)
	x.c.Assert(err, IsNil)
	x.bw.Flush()
}

func (x *conn) CheckLine(expected string) {
	l, err := x.br.ReadString('\n')
	x.c.Check(err, IsNil)
	x.c.Check(strings.TrimRight(l, "\r\n"), Equals, expected)
}

func (x *conn) CheckLines(expected []string) {
	for _, e := range expected {
		x.CheckLine(e)
	}

	x.CheckLine("")
}

func (x *conn) WriteLine(line string) {
	x.bw.WriteString(line)
	x.bw.WriteString("\r\n")
	x.bw.Flush()
}

func (x *conn) WriteLines(lines []string) {
	for _, e := range lines {
		x.WriteLine(e)
	}

	x.WriteLine("")
}

type ProxySuite struct {
	*C

	r *Registry
	p *Proxy

	// This channel is closed when the test is done
	done chan bool
}

var _ = Suite(&ProxySuite{})

func (s *ProxySuite) SetUpTest(c *C) {
	x := DefaultConfig()
	x.TraceKey = "my_trace_key"

	mocksController := gomock.NewController(c)
	s.r = NewRegistry(x, test.NewMockCFMessageBus(mocksController))
	s.p = NewProxy(x, s.r, nullVarz{})

	s.done = make(chan bool)
}

func (s *ProxySuite) TearDownTest(c *C) {
	close(s.done)
}

func (s *ProxySuite) registerAddr(u string, a net.Addr) {
	h, p, err := net.SplitHostPort(a.String())
	if err != nil {
		panic(err)
	}

	x, err := strconv.Atoi(p)
	if err != nil {
		panic(err)
	}

	m := registryMessage{
		Host: h,
		Port: uint16(x),
		Uris: []Uri{Uri(u)},
	}

	s.r.Register(&m)
}

func (s *ProxySuite) RegisterHandler(u string, h connHandler) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	// Close listener when test is done
	go func() {
		<-s.done
		ln.Close()
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}

			go h(newConn(conn, s.C))
		}
	}()

	s.registerAddr(u, ln.Addr())

	return ln
}

func (s *ProxySuite) StartProxy() net.Addr {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	// Close listener when test is done
	go func() {
		<-s.done
		ln.Close()
	}()

	go func() {
		http.Serve(ln, s.p)
	}()

	return ln.Addr()
}

func (s *ProxySuite) DialProxy() *conn {
	y := s.StartProxy()

	x, err := net.Dial("tcp", y.String())
	if err != nil {
		panic(err)
	}

	return newConn(x, s.C)
}

func (s *ProxySuite) TestRespondsToHttp10(c *C) {
	s.C = c

	s.RegisterHandler("test", func(x *conn) {
		x.CheckLine("GET / HTTP/1.1")

		x.WriteLines([]string{
			"HTTP/1.1 200 OK",
			"Content-Length: 0",
		})
	})

	x := s.DialProxy()

	x.WriteLines([]string{
		"GET / HTTP/1.0",
		"Host: test",
	})

	x.CheckLine("HTTP/1.0 200 OK")
}

func (s *ProxySuite) TestRespondsToHttp11(c *C) {
	s.C = c

	s.RegisterHandler("test", func(x *conn) {
		x.CheckLine("GET / HTTP/1.1")

		x.WriteLines([]string{
			"HTTP/1.1 200 OK",
			"Content-Length: 0",
		})
	})

	x := s.DialProxy()

	x.WriteLines([]string{
		"GET / HTTP/1.1",
		"Host: test",
	})

	x.CheckLine("HTTP/1.1 200 OK")
}

func (s *ProxySuite) TestDoesNotRespondToUnsupportedHttp(c *C) {
	s.C = c

	x := s.DialProxy()

	x.WriteLines([]string{
		"GET / HTTP/0.9",
		"Host: test",
	})

	x.CheckLine("HTTP/1.0 400 Bad Request")
}

func (s *ProxySuite) TestRespondsToLoadBalancerCheck(c *C) {
	s.C = c

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "HTTP-Monitor/1.1")
	x.WriteRequest(req)

	_, body := x.ReadResponse()
	s.Check(body, Equals, "ok\n")
}

func (s *ProxySuite) TestRespondsToUnknownHostWith404(c *C) {
	s.C = c

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Header.Set("Host", "unknown")
	x.WriteRequest(req)

	resp, body := x.ReadResponse()
	s.Check(resp.StatusCode, Equals, http.StatusNotFound)
	s.Check(body, Equals, "404 Not Found\n")
}

func (s *ProxySuite) TestRespondsToMisbehavingHostWith502(c *C) {
	s.C = c

	s.RegisterHandler("enfant-terrible", func(x *conn) {
		x.Close()
	})

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Host = "enfant-terrible"
	x.WriteRequest(req)

	resp, body := x.ReadResponse()
	s.Check(resp.StatusCode, Equals, http.StatusBadGateway)
	s.Check(body, Equals, "502 Bad Gateway\n")
}

func (s *ProxySuite) TestTraceHeadersAddedOnCorrectTraceKey(c *C) {
	s.C = c

	ln := s.RegisterHandler("trace-test", func(x *conn) {
		resp := newResponse(http.StatusOK)
		x.WriteResponse(resp)
		x.Close()
	})

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Host = "trace-test"
	req.Header.Set("X-Vcap-Trace", "my_trace_key")
	x.WriteRequest(req)

	resp, _ := x.ReadResponse()
	c.Check(resp.Header.Get("X-Vcap-Backend"), Equals, ln.Addr().String())
	c.Check(resp.Header.Get("X-Vcap-Router"), Equals, s.p.Config.Ip)
}

func (s *ProxySuite) TestTraceHeadersNotAddedOnIncorrectTraceKey(c *C) {
	s.C = c

	s.RegisterHandler("trace-test", func(x *conn) {
		resp := newResponse(http.StatusOK)
		x.WriteResponse(resp)
		x.Close()
	})

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Host = "trace-test"
	req.Header.Set("X-Vcap-Trace", "a_bad_trace_key")
	x.WriteRequest(req)

	resp, _ := x.ReadResponse()
	c.Check(resp.Header.Get("X-Vcap-Backend"), Equals, "")
	c.Check(resp.Header.Get("X-Vcap-Router"), Equals, "")
}

func (s *ProxySuite) TestXFFIsAdded(c *C) {
	s.C = c

	done := make(chan bool)

	s.RegisterHandler("app", func(x *conn) {
		req, _ := x.ReadRequest()
		c.Check(req.Header.Get("X-Forwarded-For"), Equals, "127.0.0.1")
		done <- true
	})

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Host = "app"
	x.WriteRequest(req)

	<-done
}

func (s *ProxySuite) TestXFFIsAppended(c *C) {
	s.C = c

	done := make(chan bool)

	s.RegisterHandler("app", func(x *conn) {
		req, _ := x.ReadRequest()
		c.Check(req.Header.Get("X-Forwarded-For"), Equals, "1.2.3.4, 127.0.0.1")
		done <- true
	})

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Host = "app"
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	x.WriteRequest(req)

	<-done
}

func (s *ProxySuite) TestWebSocketUpgrade(c *C) {
	s.C = c

	s.RegisterHandler("ws", func(x *conn) {
		req, _ := x.ReadRequest()
		c.Check(req.Header.Get("Upgrade"), Equals, "websocket")
		c.Check(req.Header.Get("Connection"), Equals, "Upgrade")

		resp := newResponse(http.StatusSwitchingProtocols)
		resp.Header.Set("Upgrade", "websocket")
		resp.Header.Set("Connection", "Upgrade")

		x.WriteResponse(resp)

		x.CheckLine("hello from client")
		x.WriteLine("hello from server")
	})

	x := s.DialProxy()

	req := x.NewRequest("GET", "/chat", nil)
	req.Host = "ws"
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")

	x.WriteRequest(req)

	resp, _ := x.ReadResponse()
	c.Check(resp.StatusCode, Equals, http.StatusSwitchingProtocols)
	c.Check(resp.Header.Get("Upgrade"), Equals, "websocket")
	c.Check(resp.Header.Get("Connection"), Equals, "Upgrade")

	x.WriteLine("hello from client")
	x.CheckLine("hello from server")
}

func (s *ProxySuite) TestTransferEncodingChunked(c *C) {
	s.C = c

	s.RegisterHandler("chunk", func(x *conn) {
		r, w := io.Pipe()

		// Write 3 times on a 100ms interval
		go func() {
			t := time.NewTicker(100 * time.Millisecond)
			defer t.Stop()
			defer w.Close()
			for i := 0; i < 3; i++ {
				select {
				case <-t.C:
					w.Write([]byte("hello"))
				}
			}
		}()

		resp := newResponse(http.StatusOK)
		resp.TransferEncoding = []string{"chunked"}
		resp.Body = r
		resp.Write(x)
	})

	x := s.DialProxy()

	req := x.NewRequest("GET", "/", nil)
	req.Host = "chunk"
	req.Write(x)

	resp, err := http.ReadResponse(x.br, &http.Request{})
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)
	c.Assert(resp.TransferEncoding, DeepEquals, []string{"chunked"})

	// Expect 3 individual reads to complete
	for i := 0; i < 3; i++ {
		var b [16]byte

		n, err := resp.Body.Read(b[0:])
		c.Assert(err, IsNil)
		c.Check(n, Equals, 5)
		c.Check(string(b[0:n]), Equals, "hello")
	}
}

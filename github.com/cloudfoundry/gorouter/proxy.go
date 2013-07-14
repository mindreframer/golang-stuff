package router

import (
	"bufio"
	"fmt"
	steno "github.com/cloudfoundry/gosteno"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	VcapBackendHeader = "X-Vcap-Backend"
	VcapRouterHeader  = "X-Vcap-Router"
	VcapTraceHeader   = "X-Vcap-Trace"

	VcapCookieId    = "__VCAP_ID__"
	StickyCookieKey = "JSESSIONID"
)

type Proxy struct {
	sync.RWMutex
	*steno.Logger
	*Config
	*Registry
	Varz
	*AccessLogger
}

type responseWriter struct {
	http.ResponseWriter
	*steno.Logger
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := rw.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (rw *responseWriter) WriteStatus(code int) {
	body := fmt.Sprintf("%d %s", code, http.StatusText(code))
	rw.Warn(body)
	http.Error(rw, body, code)
}

func (rw *responseWriter) CopyFrom(src io.Reader) (int64, error) {
	if src == nil {
		return 0, nil
	}

	var dst io.Writer = rw

	// Use MaxLatencyFlusher if needed
	if v, ok := rw.ResponseWriter.(writeFlusher); ok {
		u := NewMaxLatencyWriter(v, 50*time.Millisecond)
		defer u.Stop()
		dst = u
	}

	return io.Copy(dst, src)
}

func NewProxy(c *Config, r *Registry, v Varz) *Proxy {
	p := &Proxy{
		Config:   c,
		Logger:   steno.NewLogger("router.proxy"),
		Registry: r,
		Varz:     v,
	}

	if c.AccessLog != "" {
		f, err := os.OpenFile(c.AccessLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		p.AccessLogger = NewAccessLogger(f)
		go p.AccessLogger.Run()
	}

	return p
}

func hostWithoutPort(req *http.Request) string {
	host := req.Host

	// Remove :<port>
	pos := strings.Index(host, ":")
	if pos >= 0 {
		host = host[0:pos]
	}

	return host
}

func (p *Proxy) Lookup(req *http.Request) (*Backend, bool) {
	h := hostWithoutPort(req)

	// Try choosing a backend using sticky session
	if _, err := req.Cookie(StickyCookieKey); err == nil {
		if sticky, err := req.Cookie(VcapCookieId); err == nil {
			b, ok := p.Registry.LookupByPrivateInstanceId(h, sticky.Value)
			if ok {
				return b, ok
			}
		}
	}

	// Choose backend using host alone
	return p.Registry.Lookup(h)
}

func (p *Proxy) ServeHTTP(hrw http.ResponseWriter, req *http.Request) {
	rw := responseWriter{
		ResponseWriter: hrw,
		Logger:         p.Logger.Copy(),
	}

	rw.Set("RemoteAddr", req.RemoteAddr)
	rw.Set("Host", req.Host)
	rw.Set("X-Forwarded-For", req.Header["X-Forwarded-For"])
	rw.Set("X-Forwarded-Proto", req.Header["X-Forwarded-Proto"])

	a := AccessLogRecord{
		Request:   req,
		StartedAt: time.Now(),
	}

	if req.ProtoMajor != 1 && (req.ProtoMinor != 0 || req.ProtoMinor != 1) {
		c, brw, err := rw.Hijack()
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(brw, "HTTP/1.0 400 Bad Request\r\n\r\n")
		brw.Flush()
		c.Close()
		return
	}

	start := time.Now()

	// Return 200 OK for heartbeats from LB
	if req.UserAgent() == "HTTP-Monitor/1.1" {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintln(rw, "ok")
		return
	}

	x, ok := p.Lookup(req)
	if !ok {
		p.Varz.CaptureBadRequest(req)
		rw.WriteStatus(http.StatusNotFound)
		return
	}

	rw.Set("Backend", x.ToLogData())

	a.Backend = x

	p.Registry.CaptureBackendRequest(x, start)
	p.Varz.CaptureBackendRequest(x, req)

	req.URL.Scheme = "http"
	req.URL.Host = x.CanonicalAddr()

	// Add X-Forwarded-For
	if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// We assume there is a trusted upstream (L7 LB) that properly
		// strips client's XFF header

		// This is sloppy but fine since we don't share this request or
		// headers. Otherwise we should copy the underlying header and
		// append
		xff := append(req.Header["X-Forwarded-For"], host)
		req.Header.Set("X-Forwarded-For", strings.Join(xff, ", "))
	}

	// Check if the connection is going to be upgraded to a WebSocket connection
	if p.CheckWebSocket(rw, req) {
		p.ServeWebSocket(rw, req)
		return
	}

	// Use a new connection for every request
	// Keep-alive can be bolted on later, if we want to
	req.Close = true
	req.Header.Del("Connection")

	res, err := http.DefaultTransport.RoundTrip(req)

	latency := time.Since(start)

	a.FirstByteAt = time.Now()
	a.Response = res

	if err != nil {
		p.Varz.CaptureBackendResponse(x, res, latency)
		rw.Warnf("Error reading from upstream: %s", err)
		rw.WriteStatus(http.StatusBadGateway)
		return
	}

	p.Varz.CaptureBackendResponse(x, res, latency)

	for k, vv := range res.Header {
		for _, v := range vv {
			rw.Header().Add(k, v)
		}
	}

	if p.Config.TraceKey != "" && req.Header.Get(VcapTraceHeader) == p.Config.TraceKey {
		rw.Header().Set(VcapRouterHeader, p.Config.Ip)
		rw.Header().Set(VcapBackendHeader, x.CanonicalAddr())
	}

	needSticky := false
	for _, v := range res.Cookies() {
		if v.Name == StickyCookieKey {
			needSticky = true
			break
		}
	}

	if needSticky && x.PrivateInstanceId != "" {
		cookie := &http.Cookie{
			Name:  VcapCookieId,
			Value: x.PrivateInstanceId,
			Path:  "/",
		}
		http.SetCookie(rw, cookie)
	}

	rw.WriteHeader(res.StatusCode)
	n, _ := rw.CopyFrom(res.Body)

	a.FinishedAt = time.Now()
	a.BodyBytesSent = n

	if p.AccessLogger != nil {
		p.AccessLogger.Log(a)
	}
}

func (p *Proxy) CheckWebSocket(rw http.ResponseWriter, req *http.Request) bool {
	return req.Header.Get("Connection") == "Upgrade" && req.Header.Get("Upgrade") == "websocket"
}

func (p *Proxy) ServeWebSocket(rw responseWriter, req *http.Request) {
	var err error

	rw.Set("Upgrade", "websocket")

	dc, _, err := rw.Hijack()
	if err != nil {
		rw.Warnf("hj.Hijack: %s", err)
		rw.WriteStatus(http.StatusBadRequest)
		return
	}

	defer dc.Close()

	// Dial backend
	uc, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		rw.Warnf("net.Dial: %s", err)
		rw.WriteStatus(http.StatusBadRequest)
		return
	}

	defer uc.Close()

	// Write request
	err = req.Write(uc)
	if err != nil {
		rw.Warnf("Writing request: %s", err)
		rw.WriteStatus(http.StatusBadRequest)
		return
	}

	errch := make(chan error, 2)

	copy := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		if err != nil {
			errch <- err
		}
	}

	go copy(uc, dc)
	go copy(dc, uc)

	// Don't care about error, both connections will be closed if necessary
	<-errch
}

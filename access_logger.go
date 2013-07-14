package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AccessLogRecord struct {
	Request       *http.Request
	Response      *http.Response
	Backend       *Backend
	StartedAt     time.Time
	FirstByteAt   time.Time
	FinishedAt    time.Time
	BodyBytesSent int64
}

func (r *AccessLogRecord) FormatStartedAt() string {
	return r.StartedAt.Format("02/01/2006:15:04:05 -0700")
}

func (r *AccessLogRecord) FormatRequestHeader(k string) (v string) {
	v = r.Request.Header.Get(k)
	if v == "" {
		v = "-"
	}
	return
}

func (r *AccessLogRecord) ResponseTime() float64 {
	return float64(r.FinishedAt.UnixNano()-r.StartedAt.UnixNano()) / float64(time.Second)
}

func (r *AccessLogRecord) WriteTo(w io.Writer) (int64, error) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, `%s - `, r.Request.Host)
	fmt.Fprintf(b, `[%s] `, r.FormatStartedAt())
	fmt.Fprintf(b, `"%s %s %s" `, r.Request.Method, r.Request.URL.RequestURI(), r.Request.Proto)
	fmt.Fprintf(b, `%d `, r.Response.StatusCode)
	fmt.Fprintf(b, `%d `, r.BodyBytesSent)
	fmt.Fprintf(b, `"%s" `, r.FormatRequestHeader("Referer"))
	fmt.Fprintf(b, `"%s" `, r.FormatRequestHeader("User-Agent"))
	fmt.Fprintf(b, `%s `, r.Request.RemoteAddr)
	fmt.Fprintf(b, `response_time:%.9f `, r.ResponseTime())
	fmt.Fprintf(b, `app_id:%s`, r.Backend.ApplicationId)
	fmt.Fprint(b, "\n")
	return b.WriteTo(w)
}

type AccessLogger struct {
	c chan AccessLogRecord
	w io.Writer
}

func NewAccessLogger(f *os.File) *AccessLogger {
	return &AccessLogger{
		w: f,
		c: make(chan AccessLogRecord, 128),
	}
}

func (x *AccessLogger) Run() {
	for r := range x.c {
		r.WriteTo(x.w)
	}
}

func (x *AccessLogger) Stop() {
	close(x.c)
}

func (x *AccessLogger) Log(r AccessLogRecord) {
	x.c <- r
}

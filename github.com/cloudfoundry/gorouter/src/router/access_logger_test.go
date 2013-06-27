package router

import (
	"bytes"
	. "launchpad.net/gocheck"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type AccessLoggerSuite struct{}

var _ = Suite(&AccessLoggerSuite{})

func (s *AccessLoggerSuite) CreateAccessLogRecord() *AccessLogRecord {
	u, err := url.Parse("http://foo.bar:1234/quz?wat")
	if err != nil {
		panic(err)
	}

	req := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		Header:     make(http.Header),
		Host:       "foo.bar",
		RemoteAddr: "1.2.3.4:5678",
	}

	req.Header.Set("Referer", "referer")
	req.Header.Set("User-Agent", "user-agent")

	res := &http.Response{
		StatusCode: http.StatusOK,
	}

	b := &Backend{
		ApplicationId: "my_awesome_id",
		Host:          "127.0.0.1",
		Port:          4567,
	}

	r := AccessLogRecord{
		Request:       req,
		Response:      res,
		Backend:       b,
		StartedAt:     time.Unix(10, 100000000),
		FirstByteAt:   time.Unix(10, 200000000),
		FinishedAt:    time.Unix(10, 300000000),
		BodyBytesSent: 42,
	}

	return &r
}

func (s *AccessLoggerSuite) TestAccessLogRecordEncode(c *C) {
	r := s.CreateAccessLogRecord()

	p := `` +
		regexp.QuoteMeta(`foo.bar `) +
		regexp.QuoteMeta(`- `) +
		`\[\d{2}/\d{2}/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4}\] ` +
		regexp.QuoteMeta(`"GET /quz?wat HTTP/1.1" `) +
		regexp.QuoteMeta(`200 `) +
		regexp.QuoteMeta(`42 `) +
		regexp.QuoteMeta(`"referer" `) +
		regexp.QuoteMeta(`"user-agent" `) +
		regexp.QuoteMeta(`1.2.3.4:5678 `) +
		regexp.QuoteMeta(`response_time:0.200000000 `) +
		regexp.QuoteMeta(`app_id:my_awesome_id`)

	b := &bytes.Buffer{}
	_, err := r.WriteTo(b)
	c.Assert(err, IsNil)

	c.Check(b.String(), Matches, "^"+p+"\n")
}

type nullWriter struct{}

func (n nullWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (s *AccessLoggerSuite) BenchmarkAccessLogRecordWriteTo(c *C) {
	r := s.CreateAccessLogRecord()
	w := nullWriter{}

	for i := 0; i < c.N; i++ {
		r.WriteTo(w)
	}
}

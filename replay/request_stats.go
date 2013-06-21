package replay

import (
	"time"
)

// Stats stores in context of current timestamp
type RequestStat struct {
	timestamp int64

	Codes map[int]int // { 200: 10, 404:2, 500:1 }

	Count  int // All requests including errors
	Errors int // Rquests with errors (timeout or host not reachable). Not include 50x errors.

	host *ForwardHost
}

// Ensure that current stats is actual (for current timestamp)
func (s *RequestStat) Touch() {
	if s.timestamp != time.Now().Unix() {
		s.reset()
	}
}

// Called on request start
func (s *RequestStat) IncReq() {
	s.Count++
}

// Called after response
func (s *RequestStat) IncResp(resp *HttpResponse) {
	s.Touch()

	if resp.err != nil {
		s.Errors++
		return
	}

	s.Codes[resp.resp.StatusCode]++
}

// Updated stats timestamp to current time and reset to zero all stats values
// TODO: Further on reset it should write stats to file
func (s *RequestStat) reset() {
	if s.timestamp != 0 {        
	   Debug("Host:", s.host.Url, "Requests:", s.Count, "Errors:", s.Errors, "Status codes:", s.Codes)
	}

	s.timestamp = time.Now().Unix()

	s.Codes = make(map[int]int)
	s.Count = 0
	s.Errors = 0
}

// RequestStat constructor
func NewRequestStats(host *ForwardHost) (stat *RequestStat) {
	stat = &RequestStat{host: host}
	stat.reset()

	return
}

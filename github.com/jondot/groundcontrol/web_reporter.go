package main

/*
   Web Reporter

   Will hold a series of Health reports, so that it will be easy to
   plot things like time-series.

   Every `ReportHealth` invocation will tuck an additional Health value
   onto the series.

   The series width is fixed and old measurements are dropped (this is a
   moving window).

   A couple of configuration parameters are relevant here:

   historyInterval - at what distance to take measurements in seconds
   (every minute, hour, etc)

   historyBacklog - how many measurement points to hold. This affects
   memory usage
*/

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type WebReporter struct {
	health          Health
	Mount           string
	history         *ring.Ring
	lastReport      time.Time
	historyInterval int
}

func NewWebReporter(historyInter int, historyBacklog int) (h *WebReporter) {
	return &WebReporter{Mount: "/health", history: ring.New(historyBacklog), historyInterval: historyInter}
}

func (self *WebReporter) ReportHealth(h *Health) {
	self.health = *h
	self.addReading(*h)
}

func (self *WebReporter) addReading(h Health) {
	t := time.Now()

	if t.Sub(self.lastReport).Seconds() > float64(self.historyInterval) {
		self.lastReport = t

		p := self.history.Prev()
		p.Value = map[string]Health{fmt.Sprintf("%d", t.Unix()): h}
		self.history = p
	}
}

func (self *WebReporter) Handler(w http.ResponseWriter, r *http.Request) {
	hdr := w.Header()
	hdr.Add("Access-Control-Allow-Origin", "*")
	series := []map[string]Health{}
	self.history.Do(func(x interface{}) {
		if x != nil {
			series = append(series, x.(map[string]Health))
		}
	})
	log.Println("series", series)
	enc := json.NewEncoder(w)
	err := enc.Encode(&series)
	if err != nil {
		log.Println("encoding error", err)
	}
}

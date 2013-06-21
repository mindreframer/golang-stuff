package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LibratoBulk struct {
	Gauges []map[string]interface{} `json:"gauges"`
}

type LibratoReporter struct {
	Credentials ReporterCredentials
}

// lifted from https://github.com/rcrowley/go-librato/blob/master/simple.go
// thanks!
type tbody map[string]tibody
type tibody []tmetric
type tmetric map[string]interface{}

var libratoReporterUA = func() string {
	return fmt.Sprintf("groundcontrol/%s", VERSION)
}()

func NewLibratoReporter(creds ReporterCredentials) (h *LibratoReporter) {
	return &LibratoReporter{Credentials: creds}
}

func (self *LibratoReporter) ReportHealth(h *Health) {

	bulk := LibratoBulk{}

	hmap := h.Map()
	for k, v := range hmap {
		bulk.Gauges = append(bulk.Gauges, map[string]interface{}{"name": k, "value": v, "source": "pi"})
	}

	b, _ := json.Marshal(bulk)

	req, err := http.NewRequest(
		"POST",
		"https://metrics-api.librato.com/v1/metrics",
		bytes.NewBuffer(b),
	)

	if nil != err {
		log.Println("Error creating request", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", libratoReporterUA)
	req.SetBasicAuth(self.Credentials.User, self.Credentials.Key)
	resp, err := http.DefaultClient.Do(req)
	
	if nil != err {
		log.Println("Error receiving response", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("Error: Librato API Error: ", resp)
	}
}

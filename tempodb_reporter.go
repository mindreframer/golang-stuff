package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Bulk struct {
	Time string                   `json:"t"`
	Data []map[string]interface{} `json:"data"`
}

type TempoDBReporter struct {
	Credentials ReporterCredentials
}

func NewTempoDBReporter(creds ReporterCredentials) (h *TempoDBReporter) {
	return &TempoDBReporter{Credentials: creds}
}

func (self *TempoDBReporter) ReportHealth(h *Health) {

	base := "https://api.tempo-db.com/v1/data"

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	http.DefaultTransport = tr

	purl, _ := url.Parse(base)
	purl.User = url.UserPassword(self.Credentials.User,
		self.Credentials.Key)

	t := time.Now()
	now := t.Format("2006-01-02T15:04:05.999Z0700")
	blk := &Bulk{Time: now}

	blk.Data = append(blk.Data,
		map[string]interface{}{"key": "avg1", "v": h.LoadAvg1},
		map[string]interface{}{"key": "avg5", "v": h.LoadAvg5},
		map[string]interface{}{"key": "avg15", "v": h.LoadAvg15},
		map[string]interface{}{"key": "memfree", "v": h.MemActualFree},
		map[string]interface{}{"key": "memused", "v": h.MemActualUsed},
		map[string]interface{}{"key": "CPU temp", "v": h.CPUTemp},
	)

	for _, disk := range h.Disks {
		blk.Data = append(blk.Data, map[string]interface{}{"key": fmt.Sprintf("disk_used: %s", disk.DeviceName), "v": disk.Used})
		blk.Data = append(blk.Data, map[string]interface{}{"key": fmt.Sprintf("disk_used_pcent: %s", disk.DeviceName), "v": disk.UsedPcent})
	}

	b, _ := json.Marshal(blk)
	r := bytes.NewReader(b)
	resp, err := http.Post(purl.String(), "application/json", r)
	if nil != err || resp.StatusCode != 200 {
		log.Println("Error: TempoDB API Error: ", err, resp)
	}
}

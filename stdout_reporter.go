package main

import (
	"fmt"
	"log"
	"sort"
)

type StdoutReporter struct {
}

func NewStdoutReporter() (h *StdoutReporter) {
	return &StdoutReporter{}
}

func (self *StdoutReporter) ReportHealth(h *Health) {
	var keys []string
	m := h.Map()

	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		log.Println(fmt.Sprintf("%s: %v", k, m[k]))
	}
}

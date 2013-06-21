package main

/*
  Health

  Hold the general structure to be used for gathering health metrics.

  Health should know how to transform itself into a map, which is
  useful for those who don't want to deal with strongly-typed value
  with their reporters
*/

import (
	"fmt"
	"github.com/jondot/gosigar"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type DiskInfo struct {
	DeviceName string  `json:"deviceName"`
	Used       uint64  `json:"used"`
	UsedPcent  float64 `json:"usedPercent"`
}

type Health struct {
	LoadAvg1      float64    `json:"loadAvg1"`
	LoadAvg5      float64    `json:"loadAvg5"`
	LoadAvg15     float64    `json:"loadAvg15"`
	MemActualFree uint64     `json:"freeMemory"`
	MemActualUsed uint64     `json:"usedMemory"`
	Disks         []DiskInfo `json:"disks"`
	CPUTemp       float64    `json:"cpuTemp"`
}

func GetHealth(tempPath string) (h *Health, err error) {
	avg := sigar.LoadAverage{}
	err = avg.Get()
	if err != nil {
		return
	}

	mem := sigar.Mem{}
	err = mem.Get()
	if err != nil {
		return
	}

	health := &Health{
		LoadAvg1:      avg.One,
		LoadAvg5:      avg.Five,
		LoadAvg15:     avg.Fifteen,
		MemActualFree: mem.ActualFree,
		MemActualUsed: mem.ActualUsed,
	}

	fslist := sigar.FileSystemList{}
	err = fslist.Get()
	if err != nil {
		return
	}

	for _, fs := range fslist.List {
		usage := sigar.FileSystemUsage{}
		usage.Get(fs.DirName)
		health.Disks = append(health.Disks, DiskInfo{DeviceName: fs.DevName, Used: usage.Used, UsedPcent: usage.UsePercent()})
	}

	health.CPUTemp = ABS_ZERO
	if tempPath != "" {
		health.CPUTemp = getCpuTemp(tempPath)
	}

	h = health
	return
}

const ABS_ZERO = -273.15 // doh. http://en.wikipedia.org/wiki/Absolute_zero

func getCpuTemp(path string) float64 {
	text, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read sensor data at %s: %s\n", path, err)
		return ABS_ZERO
	}

	v, err := strconv.Atoi(strings.TrimRight(string(text), "\n"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid temp: %s\n", err)
		return ABS_ZERO
	}
	return float64(v) / 1000.0
}

func (self *Health) Map() map[string]interface{} {
	m := map[string]interface{}{
		"load.avg1":       self.LoadAvg1,
		"load.avg5":       self.LoadAvg5,
		"load.avg15":      self.LoadAvg15,
		"mem.actualfree":  self.MemActualFree,
		"mem.actualused":  self.MemActualUsed,
		"sensors.cputemp": self.CPUTemp,
	}

	for _, disk := range self.Disks {
		m[deviceToKey(fmt.Sprintf("disks.%s.used", disk.DeviceName))] = disk.Used
		m[deviceToKey(fmt.Sprintf("disks.%s.used_pcent", disk.DeviceName))] = disk.UsedPcent
	}

	return m
}

func deviceToKey(dev string) string {
	re := regexp.MustCompile("[-_/:\\s]")
	re2 := regexp.MustCompile("\\.+")
	s := re.ReplaceAllString(dev, ".")
	s = re2.ReplaceAllString(s, ".")
	return s
}

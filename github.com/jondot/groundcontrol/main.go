package main

/*
  Ground Control

  The main outline should read the configuration file, set up
  various reporters and a timer with a goroutine that handles grabbing
  the health data and reporting under each kind of reporter.

  It should also set up the various HTTP handlers on the proper
  mountpoints.
*/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	configPath := flag.String("config", "/etc/groundcontrol.json", "Ground control config file.")
	version := flag.Bool("version", false, "Output version and exit")

	flag.Parse()

	if *version {
		fmt.Println("Ground Control", VERSION)
		os.Exit(0)
	}

	config := loadConfiguration(*configPath)

	erroroutAndUsage := false

	if config.Interval < 10 {
		println("Error: Interval cannot be smaller than 10")
		erroroutAndUsage = true
	}

	if erroroutAndUsage {
		println("") // intentionally blank
		flag.Usage()
		os.Exit(-1)
	}

	control := NewControl(config.Controls)

	//
	// set up reporters
	//
	webreporter := NewWebReporter(config.HistoryInterval, config.HistoryBacklog)
	reporters := []Reporter{webreporter}

	if config.TempoDB.User == "" || config.TempoDB.Key == "" {
		log.Println("Reporters: No TempoDB credentials, skipping.")
	} else {
		reporters = append(reporters, NewTempoDBReporter(config.TempoDB))
		log.Println("Reporters: TempoDB OK.")
	}

	if config.Librato.User == "" || config.Librato.Key == "" {
		log.Println("Reporters: No Librato credentials, skipping.")
	} else {
		reporters = append(reporters, NewLibratoReporter(config.Librato))
		log.Println("Reporters: Librato OK.")
	}

	if config.Graphite.LineRec == "" {
		log.Println("Reporters: No Graphite config, skipping.")
	} else {
		reporters = append(reporters, NewGraphiteReporter(config.Graphite))
		log.Println("Reporters: Graphite OK.")
	}

	if config.Stdout {
		log.Println("Reporters: Showing you output (you said -stdout=true).")
		reporters = append(reporters, NewStdoutReporter())
	}

	log.Println("Lauching Health")

	report(config, &reporters)

	// set up a periodic report
	ticker := time.NewTicker(time.Second * time.Duration(config.Interval))
	go func() {
		for _ = range ticker.C {
			report(config, &reporters)
		}
	}()

	log.Println("Launching Control")

	//statics (UI)
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	//control endpoint
	http.HandleFunc(control.Mount, control.Handler)
	http.HandleFunc(webreporter.Mount, webreporter.Handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil)
}

func report(config *GroundControlConfig, reporters *[]Reporter) {
	h, err := GetHealth(config.Temperature)
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	for _, r := range *reporters {
		r.ReportHealth(h)
	}
}

func loadConfiguration(configPath string) (c *GroundControlConfig) {
	if configPath == "" {
		println("Error: You should specify a config file.")
		flag.Usage()
		os.Exit(-1)
	}

	text, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln("Cannot read config file: ", configPath)
	}

	config := &GroundControlConfig{}
	err = json.Unmarshal(text, &config)

	if err != nil {
		log.Fatalln("Cannot parse config file: ", configPath)
	}
	return config
}

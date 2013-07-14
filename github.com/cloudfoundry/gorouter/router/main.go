package main

import (
	"flag"
	"github.com/cloudfoundry/gorouter"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration File")

	flag.Parse()
}

func main() {
	c := router.DefaultConfig()
	if configFile != "" {
		c = router.InitConfigFromFile(configFile)
	}

	router.SetupLoggerFromConfig(c)

	router.NewRouter(c).Run()
}

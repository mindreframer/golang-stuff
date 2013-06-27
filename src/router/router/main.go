package main

import (
	"flag"
	"router"
	"router/config"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration File")

	flag.Parse()
}

func main() {
	c := config.DefaultConfig()
	if configFile != "" {
		c = config.InitConfigFromFile(configFile)
	}

	router.SetupLoggerFromConfig(c)

	router.NewRouter(c).Run()
}

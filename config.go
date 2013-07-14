package router

import (
	"io/ioutil"
	"launchpad.net/goyaml"
	vcap "github.com/cloudfoundry/gorouter/common"
	"time"
)

type StatusConfig struct {
	Port uint16 "port"
	User string "user"
	Pass string "pass"
}

var defaultStatusConfig = StatusConfig{
	Port: 8082,
	User: "",
	Pass: "",
}

type NatsConfig struct {
	Host string "host"
	Port uint16 "port"
	User string "user"
	Pass string "pass"
}

var defaultNatsConfig = NatsConfig{
	Host: "localhost",
	Port: 4222,
	User: "",
	Pass: "",
}

type LoggingConfig struct {
	File   string "file"
	Syslog string "syslog"
	Level  string "level"
}

var defaultLoggingConfig = LoggingConfig{
	Level: "debug",
}

type Config struct {
	Status  StatusConfig  "status"
	Nats    NatsConfig    "nats"
	Logging LoggingConfig "logging"

	Port       uint16 "port"
	Index      uint   "index"
	Pidfile    string "pidfile"
	GoMaxProcs int    "go_max_procs,omitempty"
	TraceKey   string "trace_key"
	AccessLog  string "access_log"

	PublishStartMessageIntervalInSeconds int "publish_start_message_interval"
	PruneStaleDropletsIntervalInSeconds  int "prune_stale_droplets_interval"
	DropletStaleThresholdInSeconds       int "droplet_stale_threshold"
	PublishActiveAppsIntervalInSeconds   int "publish_active_apps_interval"
	StartResponseDelayIntervalInSeconds  int "start_response_delay_interval"

	// These fields are populated by the `Process` function.
	PruneStaleDropletsInterval time.Duration
	DropletStaleThreshold      time.Duration
	PublishActiveAppsInterval  time.Duration
	StartResponseDelayInterval time.Duration

	Ip string
}

var defaultConfig = Config{
	Status:  defaultStatusConfig,
	Nats:    defaultNatsConfig,
	Logging: defaultLoggingConfig,

	Port:       8081,
	Index:      0,
	Pidfile:    "",
	GoMaxProcs: 8,

	PublishStartMessageIntervalInSeconds: 30,
	PruneStaleDropletsIntervalInSeconds:  30,
	DropletStaleThresholdInSeconds:       120,
	PublishActiveAppsIntervalInSeconds:   0,
	StartResponseDelayIntervalInSeconds:  5,
}

func DefaultConfig() *Config {
	c := defaultConfig

	c.Process()

	return &c
}

func (c *Config) Process() {
	var err error

	c.PruneStaleDropletsInterval = time.Duration(c.PruneStaleDropletsIntervalInSeconds) * time.Second
	c.DropletStaleThreshold = time.Duration(c.DropletStaleThresholdInSeconds) * time.Second
	c.PublishActiveAppsInterval = time.Duration(c.PublishActiveAppsIntervalInSeconds) * time.Second
	c.StartResponseDelayInterval = time.Duration(c.StartResponseDelayIntervalInSeconds) * time.Second

	c.Ip, err = vcap.LocalIP()
	if err != nil {
		panic(err)
	}
}

func InitConfigFromFile(path string) *Config {
	var c *Config = DefaultConfig()
	var e error

	b, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e.Error())
	}

	e = goyaml.Unmarshal(b, c)
	if e != nil {
		panic(e.Error())
	}

	c.Process()

	return c
}

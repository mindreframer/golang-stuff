package router

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	mbus "github.com/cloudfoundry/go_cfmessagebus"
	steno "github.com/cloudfoundry/gosteno"
	"net"
	vcap "github.com/cloudfoundry/gorouter/common"
	"github.com/cloudfoundry/gorouter/proxy"
	"github.com/cloudfoundry/gorouter/util"
	"runtime"
	"time"
)

type Router struct {
	config     *Config
	proxy      *Proxy
	mbusClient mbus.CFMessageBus
	registry   *Registry
	varz       Varz
	component  *vcap.VcapComponent
}

func NewRouter(c *Config) *Router {
	router := &Router{
		config: c,
	}

	// setup number of procs
	if router.config.GoMaxProcs != 0 {
		runtime.GOMAXPROCS(router.config.GoMaxProcs)
	}

	router.establishMBus()

	router.registry = NewRegistry(router.config, router.mbusClient)
	router.registry.StartPruningCycle()

	router.varz = NewVarz(router.registry)
	router.proxy = NewProxy(router.config, router.registry, router.varz)

	var host string
	if router.config.Status.Port != 0 {
		host = fmt.Sprintf("%s:%d", router.config.Ip, router.config.Status.Port)
	}

	varz := &vcap.Varz{
		UniqueVarz: router.varz,
	}
	varz.LogCounts = logCounter

	healthz := &vcap.Healthz{
		LockableObject: router.registry,
	}

	router.component = &vcap.VcapComponent{
		Type:        "Router",
		Index:       router.config.Index,
		Host:        host,
		Credentials: []string{router.config.Status.User, router.config.Status.Pass},
		Config:      router.config,
		Logger:      log,
		Varz:        varz,
		Healthz:     healthz,
		InfoRoutes: map[string]json.Marshaler{
			"/routes": router.registry,
		},
	}

	vcap.StartComponent(router.component)

	return router
}

func (r *Router) RegisterComponent() {
	vcap.Register(r.component, r.mbusClient)
}

func (r *Router) subscribeRegistry(subject string, successCallback func(*registryMessage)) {
	callback := func(payload []byte) {
		var msg registryMessage

		err := json.Unmarshal(payload, &msg)
		if err != nil {
			logMessage := fmt.Sprintf("%s: Error unmarshalling JSON (%d; %s): %s", subject, len(payload), payload, err)
			log.Log(steno.LOG_WARN, logMessage, map[string]interface{}{"payload": string(payload)})
		}

		logMessage := fmt.Sprintf("%s: Received message", subject)
		log.Log(steno.LOG_DEBUG, logMessage, map[string]interface{}{"message": msg})

		successCallback(&msg)
	}
	err := r.mbusClient.Subscribe(subject, callback)
	if err != nil {
		log.Errorf("Error subscribing to %s: %s", subject, err.Error())
	}
}

func (router *Router) SubscribeRegister() {
	router.subscribeRegistry("router.register", func(registryMessage *registryMessage) {
		log.Infof("Got router.register: %v", registryMessage)
		router.registry.Register(registryMessage)
	})
}

func (r *Router) SubscribeUnregister() {
	r.subscribeRegistry("router.unregister", func(rm *registryMessage) {
		log.Infof("Got router.unregister: %v", rm)
		r.registry.Unregister(rm)
	})
}

func (r *Router) HandleGreetings() {
	r.mbusClient.RespondToChannel("router.greet", func(_ []byte) []byte {
		response, _ := r.greetMessage()
		return response
	})
}

func (r *Router) flushApps(t time.Time) {
	x := r.registry.ActiveSince(t)

	y, err := json.Marshal(x)
	if err != nil {
		log.Warnf("flushApps: Error marshalling JSON: %s", err)
		return
	}

	b := bytes.Buffer{}
	w := zlib.NewWriter(&b)
	w.Write(y)
	w.Close()

	z := b.Bytes()

	log.Debugf("Active apps: %d, message size: %d", len(x), len(z))

	r.mbusClient.Publish("router.active_apps", z)
}

func (r *Router) ScheduleFlushApps() {
	if r.config.PublishActiveAppsInterval == 0 {
		return
	}

	go func() {
		t := time.NewTicker(r.config.PublishActiveAppsInterval)
		x := time.Now()

		for {
			select {
			case <-t.C:
				y := time.Now()
				r.flushApps(x)
				x = y
			}
		}
	}()
}

func (r *Router) SendStartMessage() {
	b, err := r.greetMessage()
	if err != nil {
		panic(err)
	}

	// Send start message once at start
	r.mbusClient.Publish("router.start", b)
}

func (r *Router) greetMessage() ([]byte, error) {
	host, err := vcap.LocalIP()
	if err != nil {
		return nil, err
	}

	d := vcap.RouterStart{
		vcap.GenerateUUID(),
		[]string{host},
		r.config.StartResponseDelayIntervalInSeconds,
	}

	return json.Marshal(d)
}

func (router *Router) Run() {
	var err error

	for {
		err = router.mbusClient.Connect()
		if err == nil {
			break
		}
		log.Errorf("Could not connect to NATS: ", err.Error())
		time.Sleep(500 * time.Millisecond)
	}

	router.mbusClient.OnConnect(func() {
		router.SendStartMessage()
	})

	router.RegisterComponent()

	// Subscribe register/unregister router
	router.SubscribeRegister()
	router.HandleGreetings()
	router.SubscribeUnregister()

	// Kickstart sending start messages
	router.SendStartMessage()

	// Schedule flushing active app's app_id
	router.ScheduleFlushApps()

	// Wait for one start message send interval, such that the router's registry
	// can be populated before serving requests.
	if router.config.StartResponseDelayInterval != 0 {
		log.Infof("Waiting %s before listening...", router.config.StartResponseDelayInterval)
		time.Sleep(router.config.StartResponseDelayInterval)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", router.config.Port))
	if err != nil {
		log.Fatalf("net.Listen: %s", err)
	}

	util.WritePidFile(router.config.Pidfile)

	log.Infof("Listening on %s", listen.Addr())

	server := proxy.Server{Handler: router.proxy}

	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("proxy.Serve: %s", err)
	}
}

func (r *Router) establishMBus() {
	mbusClient, err := mbus.NewCFMessageBus("NATS")
	r.mbusClient = mbusClient
	if err != nil {
		panic("Could not connect to NATS")
	}

	host := r.config.Nats.Host
	user := r.config.Nats.User
	pass := r.config.Nats.Pass
	port := r.config.Nats.Port

	r.mbusClient.Configure(host, int(port), user, pass)
	r.mbusClient.SetLogger(log)
}

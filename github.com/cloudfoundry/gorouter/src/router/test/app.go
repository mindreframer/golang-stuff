package test

import (
	"encoding/json"
	"fmt"
	mbus "github.com/cloudfoundry/go_cfmessagebus"
	. "launchpad.net/gocheck"
	"net/http"
	"router/common"
)

type TestApp struct {
	port       uint16   // app listening port
	rPort      uint16   // router listening port
	urls       []string // host registered host name
	mbusClient mbus.CFMessageBus
	tags       map[string]string
	mux        *http.ServeMux
}

func NewTestApp(urls []string, rPort uint16, mbusClient mbus.CFMessageBus, tags map[string]string) *TestApp {
	app := new(TestApp)

	port, _ := common.GrabEphemeralPort()

	app.port = port
	app.rPort = rPort
	app.urls = urls
	app.mbusClient = mbusClient
	app.tags = tags

	app.mux = http.NewServeMux()

	return app
}

func (a *TestApp) AddHandler(path string, handler func(http.ResponseWriter, *http.Request)) {
	a.mux.HandleFunc(path, handler)
}

func (a *TestApp) Urls() []string {
	return a.urls
}

func (a *TestApp) Endpoint() string {
	return fmt.Sprintf("http://%s:%d/", a.urls[0], a.rPort)
}

func (a *TestApp) Listen() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.port),
		Handler: a.mux,
	}

	a.Register()

	go server.ListenAndServe()
}

func (a *TestApp) Register() {
	rm := registerMessage{
		Host: "localhost",
		Port: a.port,
		Uris: a.urls,
		Tags: a.tags,
		Dea:  "dea",
		App:  "0",

		PrivateInstanceId: common.GenerateUUID(),
	}

	b, _ := json.Marshal(rm)
	a.mbusClient.Publish("router.register", b)
}

func (a *TestApp) Unregister() {
	rm := registerMessage{
		Host: "localhost",
		Port: a.port,
		Uris: a.urls,
		Tags: nil,
		Dea:  "dea",
		App:  "0",
	}

	b, _ := json.Marshal(rm)
	a.mbusClient.Publish("router.unregister", b)
}

func (a *TestApp) VerifyAppStatus(status int, c *C) {
	for _, url := range a.urls {
		uri := fmt.Sprintf("http://%s:%d", url, a.rPort)
		resp, err := http.Get(uri)
		c.Assert(err, IsNil)
		c.Check(resp.StatusCode, Equals, status)
	}
}

// Types imported from router
const (
	VcapBackendHeader = "X-Vcap-Backend"
	VcapRouterHeader  = "X-Vcap-Router"
	VcapTraceHeader   = "X-Vcap-Trace"

	VcapCookieId    = "__VCAP_ID__"
	StickyCookieKey = "JSESSIONID"
)

type registerMessage struct {
	Host string            `json:"host"`
	Port uint16            `json:"port"`
	Uris []string          `json:"uris"`
	Tags map[string]string `json:"tags"`
	Dea  string            `json:"dea"`
	App  string            `json:"app"`

	PrivateInstanceId string `json:"private_instance_id"`
}

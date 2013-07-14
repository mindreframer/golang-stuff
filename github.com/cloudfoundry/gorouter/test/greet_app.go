package test

import (
	mbus "github.com/cloudfoundry/go_cfmessagebus"
	"io"
	"net/http"
)

func NewGreetApp(urls []string, rPort uint16, mbusClient mbus.CFMessageBus, tags map[string]string) *TestApp {
	app := NewTestApp(urls, rPort, mbusClient, tags)
	app.AddHandler("/", greetHandler)

	return app
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world")
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/bmizerany/pat"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/api"
	"github.com/globocom/gandalf/db"
	"log"
	"net/http"
)

func main() {
	dry := flag.Bool("dry", false, "dry-run: does not start the server (for testing purpose)")
	configFile := flag.String("config", "/etc/gandalf.conf", "Gandalf configuration file")
	flag.Parse()

	err := config.ReadAndWatchConfigFile(*configFile)
	if err != nil {
		msg := `Could not find gandalf config file. Searched on %s.
For an example conf check gandalf/etc/gandalf.conf file.\n %s`
		log.Panicf(msg, *configFile, err)
	}
	db.Connect()
	router := pat.New()
	router.Post("/user/:name/key", http.HandlerFunc(api.AddKey))
	router.Del("/user/:name/key/:keyname", http.HandlerFunc(api.RemoveKey))
	router.Get("/user/:name/keys", http.HandlerFunc(api.ListKeys))
	router.Post("/user", http.HandlerFunc(api.NewUser))
	router.Del("/user/:name", http.HandlerFunc(api.RemoveUser))
	router.Post("/repository", http.HandlerFunc(api.NewRepository))
	router.Post("/repository/grant", http.HandlerFunc(api.GrantAccess))
	router.Del("/repository/revoke", http.HandlerFunc(api.RevokeAccess))
	router.Del("/repository/:name", http.HandlerFunc(api.RemoveRepository))
	router.Get("/repository/:name", http.HandlerFunc(api.GetRepository))

	port, err := config.GetString("webserver:port")
	if err != nil {
		panic(err)
	}
	if !*dry {
		log.Fatal(http.ListenAndServe(port, router))
	}
}

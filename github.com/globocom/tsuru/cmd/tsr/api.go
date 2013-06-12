// Copyright 2013 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/globocom/tsuru/api"
	"github.com/globocom/tsuru/cmd"
	"launchpad.net/gnuflag"
)

type apiCmd struct {
	fs     *gnuflag.FlagSet
	config string
	dry    bool
}

func (c *apiCmd) Run(context *cmd.Context, client *cmd.Client) error {
	flags := map[string]interface{}{}
	flags["dry"] = c.dry
	flags["config"] = c.config
	api.RunServer(flags)
	return nil
}

func (apiCmd) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "api",
		Usage:   "api",
		Desc:    "Starts the tsuru api webserver.",
		MinArgs: 0,
	}
}

func (c *apiCmd) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("api", gnuflag.ExitOnError)
		c.fs.BoolVar(&c.dry, "dry", false, "dry-run: does not start the server (for testing purpose)")
		c.fs.BoolVar(&c.dry, "d", false, "dry-run: does not start the server (for testing purpose)")
		c.fs.StringVar(&c.config, "config", "/etc/tsuru/tsuru.conf", "tsr api server config file.")
		c.fs.StringVar(&c.config, "c", "/etc/tsuru/tsuru.conf", "tsr api server config file.")
	}
	return c.fs
}

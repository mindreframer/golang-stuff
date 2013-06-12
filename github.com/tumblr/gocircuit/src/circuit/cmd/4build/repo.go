// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func parseRepo(s string) (schema, key, value, url string) {
	switch {
	case strings.HasPrefix(s, "{hg}"):
		schema, s = "hg", s[len("{hg}"):]
	case strings.HasPrefix(s, "{git}"):
		schema, s = "git", s[len("{git}"):]
	case strings.HasPrefix(s, "{rsync}"):
		schema, s = "rsync", s[len("{rsync}"):]
	default:
		Fatalf("Repo '%s' has unrecognizable schema\n", s)
	}
	i := strings.Index(s, "}")
	if len(s) > 1 && s[0] == '{' && i > 0 {
		var arg string
		arg, s = s[1:i], s[i+1:]
		key, value = parseArg(arg)
	}
	url = s
	return
}

// {git}{rev:51e592253000600d586408f3e36a3f4692011086}
func parseArg(arg string) (key, value string) {
	part := strings.SplitN(arg, ":", 2)
	if len(part) > 0 {
		key = part[0]
	}
	if len(part) > 1 {
		value = part[1]
	}
	return
}

//
//	{hg}{changeset:4ad21a3b23a4}
//	{hg}{id:4ad21a3b23a4}
//	{hg}{rev:3452}
//	{hg}{tag:weekly}
//	{hg}{tip}
//	{hg}{branch:master}
//
func cloneMercurialRepo(dir, key, value, url string) {
	var opt string
	switch key {
	case "changeset":
		opt = "-u " + value
	case "rev":
		opt = "-u " + value
	case "id":
		opt = "-u " + value
	case "tag":
		opt = "-u " + value
	case "branch":
		opt = "-u " + value
	case "tip":
		opt = "-u tip"
	case "":
	default:
		Fatalf("unknown hg option key '%s'", key)
	}
	if err := Shell(x.env, "", fmt.Sprintf("hg clone %s %s %s", opt, url, dir)); err != nil {
		Fatalf("Problem cloning repository '%s' (%s)", url, err)
	}
}

func syncMercurialRepo(dir string) {
	if err := Shell(x.env, dir, "hg pull"); err != nil {
		Fatalf("Problem pulling repo in %s (%s)", dir, err)
	}
	if err := Shell(x.env, dir, "hg update"); err != nil {
		Fatalf("Problem pulling repo in %s (%s)", dir, err)
	}
}

//
//	{git}{changeset:4ad21a3b23a4}
//	{git}{id:4ad21a3b23a4}
//	{git}{rev:3452}
//	{git}{tag:weekly}
//	{git}{tip}
//	{git}{branch:master}
//
func cloneGitRepo(dir, key, value, url string) {
	if err := Shell(x.env, "", fmt.Sprintf("git clone %s %s", url, dir)); err != nil {
		Fatalf("Problem cloning repo '%s' (%s)", url, err)
	}
	var opt string
	switch key {
	case "changeset":
		opt = "git checkout " + value
	case "rev":
		opt = "git checkout " + value
	case "id":
		opt = "git checkout " + value
	case "tag":
		opt = "git checkout " + value
	case "branch":
		opt = "git checkout " + value
	case "tip", "":
	default:
		Fatalf("unknown git option key '%s'", key)
	}
	if err := Shell(x.env, dir, opt); err != nil {
		Fatalf("Problem executing '%s' in repo '%s' (%s)", opt, url, err)
	}
}

func syncGitRepo(dir string) {
	if err := Shell(x.env, dir, "git pull origin master"); err != nil {
		Fatalf("Problem pulling repo in %s (%s)", dir, err)
	}
}

func syncRsyncRepo(dir, url string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		Fatalf("Problem creating repo directory (%s)\n", err)
	}
	if err := Shell(x.env, "", "rsync -acrv --delete --exclude .git --exclude .hg --exclude *.a "+url+"/* "+dir+"/"); err != nil {
		Fatalf("Problem rsyncing dir '%s' to within '%s' (%s)", url, dir, err)
	}
}

func SyncRepo(namespace, repo, relGoPath string, fetchFresh, updateGoPath bool) (clonePath string) {

	schema, key, value, url := parseRepo(repo)

	// If fetching fresh, remove pre-existing clones
	if fetchFresh {
		if err := os.RemoveAll(path.Join(x.jail, namespace)); err != nil {
			Fatalf("Problem removing old repo clone (%s)\n", err)
		}
	}

	clonePath = path.Join(x.jail, namespace)

	// Check whether repo clone directory exists
	ok, err := Exists(clonePath)
	if err != nil {
		Fatalf("Problem stat'ing %s (%s)", clonePath, err)
	}
	switch schema {
	case "hg":
		if !ok {
			cloneMercurialRepo(clonePath, key, value, url)
		} else {
			syncMercurialRepo(clonePath)
		}
	case "git":
		if !ok {
			cloneGitRepo(clonePath, key, value, url)
		} else {
			syncGitRepo(clonePath)
		}
	case "rsync":
		syncRsyncRepo(clonePath, url)
	default:
		Fatalf("Unrecognized repo schema: %s\n", schema)
	}

	// Create build environment for building in this repo
	p := path.Join(clonePath, relGoPath)
	x.goPath[namespace] = p
	if updateGoPath {
		oldGoPath := x.env.Get("GOPATH")
		x.env.Set("GOPATH", p+":"+oldGoPath)
	}
	return
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/fs"
	"os/exec"
	"path"
)

var bare string

func bareLocation() string {
	if bare != "" {
		return bare
	}
	var err error
	bare, err = config.GetString("git:bare:location")
	if err != nil {
		panic("You should configure a git:bare:location for gandalf.")
	}
	return bare
}

func newBare(name string) error {
	args := []string{"init", path.Join(bareLocation(), formatName(name)), "--bare"}
	if bareTempl, err := config.GetString("git:bare:template"); err == nil {
		args = append(args, "--template="+bareTempl)
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Could not create git bare repository: %s. %s", err, string(out))
	}
	return nil
}

func removeBare(name string) error {
	err := fs.Filesystem().RemoveAll(path.Join(bareLocation(), formatName(name)))
	if err != nil {
		return fmt.Errorf("Could not remove git bare repository: %s", err)
	}
	return nil
}

func formatName(name string) string {
	return fmt.Sprintf("%s.git", name)
}

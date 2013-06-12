// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/repository"
	"github.com/globocom/gandalf/user"
	"io"
	"labix.org/v2/mgo/bson"
	"log/syslog"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var log *syslog.Writer

func hasWritePermission(u *user.User, r *repository.Repository) (allowed bool) {
	for _, userName := range r.Users {
		if u.Name == userName {
			return true
		}
	}
	return false
}

func hasReadPermission(u *user.User, r *repository.Repository) (allowed bool) {
	if r.IsPublic {
		return true
	}
	for _, userName := range r.Users {
		if u.Name == userName {
			return true
		}
	}
	return false
}

// Returns the command being executed by ssh.
// When a user runs `$ git push` from his/her machine, the server
// receives a ssh command, identified by this user (by the ssh key).
// The command and it's parameters are available through the SSH_ORIGINAL_COMMAND
// environment variable. In the git push example, it would have the following value:
// SSH_ORIGINAL_COMMAND=git-receive-pack 'foo.git'
// This function is responsible for retrieving the `git-receive-pack` part of SSH_ORIGINAL_COMMAND
func action() string {
	return strings.Split(os.Getenv("SSH_ORIGINAL_COMMAND"), " ")[0]
}

// Get the repository name requested in SSH_ORIGINAL_COMMAND and retrieves
// the related document on the database and returns it.
// this function does two distinct things (maybe it should'n), it
// parses the SSH_ORIGINAL_COMMAND and returns a "validation" error if it doesn't
// matches the expected format and gets the repository from the database based on the info
// obtained by the SSH_ORIGINAL_COMMAND parse.
func requestedRepository() (repository.Repository, error) {
	repoName, err := requestedRepositoryName()
	if err != nil {
		return repository.Repository{}, err
	}
	var repo repository.Repository
	if err := db.Session.Repository().Find(bson.M{"_id": repoName}).One(&repo); err != nil {
		return repository.Repository{}, errors.New("Repository not found")
	}
	return repo, nil
}

func requestedRepositoryName() (string, error) {
	r, err := regexp.Compile(`[\w-]+ '([\w-]+)\.git'`)
	if err != nil {
		panic(err)
	}
	m := r.FindStringSubmatch(os.Getenv("SSH_ORIGINAL_COMMAND"))
	if len(m) < 2 {
		return "", errors.New("Cannot deduce repository name from command. You are probably trying to do something nasty")
	}
	return m[1], nil
}

// Checks whether a command is a valid git command
// The following format is allowed:
//  git-([\w-]+) '([\w-]+)\.git'
func validateCmd() error {
	r, err := regexp.Compile(`git-([\w-]+) '([\w-]+)\.git'`)
	if err != nil {
		panic(err)
	}
	if m := r.FindStringSubmatch(os.Getenv("SSH_ORIGINAL_COMMAND")); len(m) < 3 {
		return errors.New("You've tried to execute some weird command, I'm deliberately denying you to do that, get over it.")
	}
	return nil
}

// Executes the SSH_ORIGINAL_COMMAND based on the condition
// defined by the `f` parameter.
// Also receives a custom error message to print to the end user and a
// stdout object, where the SSH_ORIGINAL_COMMAND output is going to be written
func executeAction(f func(*user.User, *repository.Repository) bool, errMsg string, stdout io.Writer) {
	var u user.User
	if err := db.Session.User().Find(bson.M{"_id": os.Args[1]}).One(&u); err != nil {
		log.Err("Error obtaining user. Gandalf database is probably in an inconsistent state.")
		fmt.Fprintln(os.Stderr, "Error obtaining user. Gandalf database is probably in an inconsistent state.")
		return
	}
	repo, err := requestedRepository()
	if err != nil {
		log.Err(err.Error())
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if f(&u, &repo) {
		// split into a function (maybe executeCmd)
		c, err := formatCommand()
		if err != nil {
			log.Err(err.Error())
			fmt.Fprintln(os.Stderr, err.Error())
		}
		log.Info("Executing " + strings.Join(c, " "))
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = stdout
		stderr := &bytes.Buffer{}
		cmd.Stderr = stderr
		err = cmd.Run()
		if err != nil {
			log.Err("Got error while executing original command: " + err.Error())
			log.Err(stderr.String())
			fmt.Fprintln(os.Stderr, "Got error while executing original command: "+err.Error())
			fmt.Fprintln(os.Stderr, stderr.String())
		}
		return
	}
	log.Err("Permission denied.")
	log.Err(errMsg)
	fmt.Fprintln(os.Stderr, "Permission denied.")
	fmt.Fprintln(os.Stderr, errMsg)
}

func formatCommand() ([]string, error) {
	p, err := config.GetString("git:bare:location")
	if err != nil {
		log.Err(err.Error())
		return []string{}, err
	}
	repoName, err := requestedRepositoryName()
	if err != nil {
		log.Err(err.Error())
		return []string{}, err
	}
	repoName += ".git"
	cmdList := strings.Split(os.Getenv("SSH_ORIGINAL_COMMAND"), " ")
	for i, c := range cmdList {
		if c == "'"+repoName+"'" {
			cmdList[i] = path.Join(p, repoName)
			break
		}
	}
	return cmdList, nil
}

func main() {
	var err error
	log, err = syslog.New(syslog.LOG_INFO, "gandalf-listener")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		panic(err.Error())
	}
	err = config.ReadConfigFile("/etc/gandalf.conf")
	if err != nil {
		log.Err(err.Error())
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	db.Connect()
	err = validateCmd()
	if err != nil {
		log.Err(err.Error())
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	a := action()
	if a == "git-receive-pack" {
		executeAction(hasWritePermission, "You don't have access to write in this repository.", os.Stdout)
		return
	}
	if a == "git-upload-pack" {
		executeAction(hasReadPermission, "You don't have access to read this repository.", os.Stdout)
		return
	}
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"regexp"
)

type Repository struct {
	Name     string `bson:"_id"`
	Users    []string
	IsPublic bool
}

// MarshalJSON marshals the Repository in json format.
func (r *Repository) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"name":    r.Name,
		"public":  r.IsPublic,
		"ssh_url": r.SshUrl(),
		"git_url": r.GitUrl(),
	}
	return json.Marshal(&data)
}

// Creates a representation of a git repository
// This function creates a git repository using the "bare-dir" config
// and saves repository's meta data in the database
func New(name string, users []string, isPublic bool) (*Repository, error) {
	r := &Repository{Name: name, Users: users, IsPublic: isPublic}
	if v, err := r.isValid(); !v {
		return r, err
	}
	if err := newBare(name); err != nil {
		return r, err
	}
	err := db.Session.Repository().Insert(&r)
	if mgo.IsDup(err) {
		return r, fmt.Errorf("A repository with this name already exists.")
	}
	return r, err
}

// Get find a repository by name
func Get(name string) (Repository, error) {
	var r Repository
	err := db.Session.Repository().FindId(name).One(&r)
	return r, err
}

// Deletes the repository from the database and
// removes it's bare git repository
func Remove(r *Repository) error {
	// maybe it should receive only a name, to standardize the api (user.Remove already does that)
	if err := removeBare(r.Name); err != nil {
		return err
	}
	if err := db.Session.Repository().RemoveId(r.Name); err != nil {
		return fmt.Errorf("Could not remove repository: %s", err)
	}
	return nil
}

// SshUrl formats the git ssh url and return it
// If no remote is configured in gandalf.conf SshUrl will panic
func (r *Repository) SshUrl() string {
	host, err := config.GetString("host")
	if err != nil {
		panic(err.Error())
	}
	uid, err := config.GetString("uid")
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s@%s:%s", uid, host, formatName(r.Name))
}

// GitUrl formats the git url and return it
// If no host is configured in gandalf.conf GitUrl will panic
func (r *Repository) GitUrl() string {
	host, err := config.GetString("host")
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("git://%s/%s", host, formatName(r.Name))
}

// Validates a repository
// A valid repository must have:
//  - a name without any special chars only alphanumeric and underlines are allowed.
//  - at least one user in users array
func (r *Repository) isValid() (bool, error) {
	m, e := regexp.Match(`^[\w-]+$`, []byte(r.Name))
	if e != nil {
		panic(e)
	}
	if !m {
		return false, errors.New("Validation Error: repository name is not valid")
	}
	if len(r.Users) == 0 {
		return false, errors.New("Validation Error: repository should have at least one user")
	}
	return true, nil
}

// Gives write permission for users (uNames) in all specified repositories (rNames)
// If any of the repositories/users do not exists, just skip it.
func GrantAccess(rNames, uNames []string) error {
	_, err := db.Session.Repository().UpdateAll(bson.M{"_id": bson.M{"$in": rNames}}, bson.M{"$addToSet": bson.M{"users": bson.M{"$each": uNames}}})
	return err
}

func RevokeAccess(rNames, uNames []string) error {
	_, err := db.Session.Repository().UpdateAll(bson.M{"_id": bson.M{"$in": rNames}}, bson.M{"$pullAll": bson.M{"users": uNames}})
	return err
}

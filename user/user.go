// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"errors"
	"fmt"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/repository"
	"labix.org/v2/mgo/bson"
	"regexp"
)

var ErrUserNotFound = errors.New("User not found")

type User struct {
	Name string `bson:"_id"`
}

// Creates a new user and write his/her keys into authorized_keys file.
//
// The authorized_keys file belongs to the user running the process.
func New(name string, keys map[string]string) (*User, error) {
	u := &User{Name: name}
	if v, err := u.isValid(); !v {
		return u, err
	}
	if err := db.Session.User().Insert(&u); err != nil {
		return u, err
	}
	return u, addKeys(keys, u.Name)
}

func (u *User) isValid() (isValid bool, err error) {
	m, err := regexp.Match(`\s|[^aA-zZ0-9\.@]|(^$)`, []byte(u.Name))
	if err != nil {
		panic(err)
	}
	if m {
		return false, errors.New("Validation Error: user name is not valid")
	}
	return true, nil
}

// Removes a user.
// Also removes it's associated keys from authorized_keys and repositories
// It handles user with repositories specially when:
// - a user has at least one repository:
//     - if he/she is the only one with access to the repository, the removal will stop and return an error
//     - if there are more than one user with access to the repository, gandalf will first revoke user's access and then remove the user permanently
// - a user has no repositories: gandalf will simply remove the user
func Remove(name string) error {
	var u *User
	if err := db.Session.User().Find(bson.M{"_id": name}).One(&u); err != nil {
		return fmt.Errorf("Could not remove user: %s", err)
	}
	if err := u.handleAssociatedRepositories(); err != nil {
		return err
	}
	if err := db.Session.User().RemoveId(u.Name); err != nil {
		return fmt.Errorf("Could not remove user: %s", err.Error())
	}
	return removeUserKeys(u.Name)
}

func (u *User) handleAssociatedRepositories() error {
	var repos []repository.Repository
	if err := db.Session.Repository().Find(bson.M{"users": u.Name}).All(&repos); err != nil {
		return err
	}
	for _, r := range repos {
		if len(r.Users) == 1 {
			return errors.New("Could not remove user: user is the only one with access to at least one of it's repositories")
		}
	}
	for _, r := range repos {
		for i, v := range r.Users {
			if v == u.Name {
				r.Users[i], r.Users = r.Users[len(r.Users)-1], r.Users[:len(r.Users)-1]
				if err := db.Session.Repository().Update(bson.M{"_id": r.Name}, r); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// Adds a key into a user.
//
// Stores the key in the user's document and write it in authorized_keys.
//
// Returns an error in case the user does not exists.
func AddKey(uName string, k map[string]string) error {
	var u User
	if err := db.Session.User().FindId(uName).One(&u); err != nil {
		return ErrUserNotFound
	}
	return addKeys(k, u.Name)
}

// RemoveKey removes the key from the database and from authorized_keys file.
//
// If the user or the key is not found, returns an error.
func RemoveKey(uName, kName string) error {
	var u User
	if err := db.Session.User().FindId(uName).One(&u); err != nil {
		return ErrUserNotFound
	}
	return removeKey(kName, uName)
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globocom/gandalf/repository"
	"github.com/globocom/gandalf/user"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

func accessParameters(body io.ReadCloser) (repositories, users []string, err error) {
	var params map[string][]string
	if err := parseBody(body, &params); err != nil {
		return []string{}, []string{}, err
	}
	users, ok := params["users"]
	if !ok {
		return []string{}, []string{}, errors.New("It is need a user list")
	}
	repositories, ok = params["repositories"]
	if !ok {
		return []string{}, []string{}, errors.New("It is need a repository list")
	}
	return repositories, users, nil
}

func GrantAccess(w http.ResponseWriter, r *http.Request) {
	// TODO: update README
	repositories, users, err := accessParameters(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := repository.GrantAccess(repositories, users); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Successfully granted access to users \"%s\" into repository \"%s\"", users, repositories)
}

func RevokeAccess(w http.ResponseWriter, r *http.Request) {
	repositories, users, err := accessParameters(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := repository.RevokeAccess(repositories, users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Successfully revoked access to users \"%s\" into repositories \"%s\"", users, repositories)
}

func AddKey(w http.ResponseWriter, r *http.Request) {
	keys := map[string]string{}
	if err := parseBody(r.Body, &keys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(keys) == 0 {
		http.Error(w, "A key is needed", http.StatusBadRequest)
		return
	}
	uName := r.URL.Query().Get(":name")
	if err := user.AddKey(uName, keys); err != nil {
		switch err {
		case user.ErrInvalidKey:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case user.ErrDuplicateKey:
			http.Error(w, "Key already exists.", http.StatusConflict)
		case user.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprint(w, "Key(s) successfully created")
}

func RemoveKey(w http.ResponseWriter, r *http.Request) {
	uName := r.URL.Query().Get(":name")
	kName := r.URL.Query().Get(":keyname")
	if err := user.RemoveKey(uName, kName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Key \"%s\" successfully removed", kName)
}

func ListKeys(w http.ResponseWriter, r *http.Request) {
	uName := r.URL.Query().Get(":name")
	keys, err := user.ListKeys(uName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	out, err := json.Marshal(&keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

type jsonUser struct {
	Name string
	Keys map[string]string
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	var usr jsonUser
	if err := parseBody(r.Body, &usr); err != nil {
		http.Error(w, "Got error while parsing body: "+err.Error(), http.StatusBadRequest)
		return
	}
	u, err := user.New(usr.Name, usr.Keys)
	if err != nil {
		http.Error(w, "Got error while creating user: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "User \"%s\" successfully created\n", u.Name)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(":name")
	if err := user.Remove(name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "User \"%s\" successfully removed\n", name)
}

func NewRepository(w http.ResponseWriter, r *http.Request) {
	var repo repository.Repository
	if err := parseBody(r.Body, &repo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rep, err := repository.New(repo.Name, repo.Users, repo.IsPublic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Repository \"%s\" successfully created\n", rep.Name)
}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	repo, err := repository.Get(r.URL.Query().Get(":name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := json.Marshal(&repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

func RemoveRepository(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(":name")
	if err := repository.Remove(name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Repository \"%s\" successfully removed\n", name)
}

func RenameRepository(w http.ResponseWriter, r *http.Request) {
	var p struct{ Name string }
	defer r.Body.Close()
	err := parseBody(r.Body, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	name := r.URL.Query().Get(":name")
	err = repository.Rename(name, p.Name)
	if err != nil && err.Error() == "not found" {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseBody(body io.ReadCloser, result interface{}) error {
	if reflect.ValueOf(result).Kind() == reflect.Struct {
		return errors.New("parseBody function cannot deal with struct. Use pointer")
	}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return errors.New(fmt.Sprintf("Could not parse json: %s", err.Error()))
	}
	return nil
}

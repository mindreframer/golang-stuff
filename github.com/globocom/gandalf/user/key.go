// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"bufio"
	"code.google.com/p/go.crypto/ssh"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/fs"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"
)

var (
	ErrDuplicateKey = errors.New("Duplicate key")
	ErrInvalidKey   = errors.New("Invalid key")
	ErrKeyNotFound  = errors.New("Key not found")
)

type Key struct {
	Name     string
	Body     string
	Comment  string
	UserName string
}

func newKey(name, user, raw string) (*Key, error) {
	key, comment, _, _, ok := ssh.ParseAuthorizedKey([]byte(raw))
	if !ok {
		return nil, ErrInvalidKey
	}
	body := ssh.MarshalAuthorizedKey(key)
	k := Key{
		Name:     name,
		Body:     string(body),
		Comment:  comment,
		UserName: user,
	}
	return &k, nil
}

func (k *Key) String() string {
	parts := make([]string, 1, 2)
	parts[0] = strings.TrimSpace(k.Body)
	if k.Comment != "" {
		parts = append(parts, k.Comment)
	}
	return strings.Join(parts, " ")
}

func (k *Key) format() string {
	binPath, err := config.GetString("bin-path")
	if err != nil {
		panic(err)
	}
	keyFmt := `no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty,command="%s %s" %s` + "\n"
	return fmt.Sprintf(keyFmt, binPath, k.UserName, k)
}

func (k *Key) dump(w io.Writer) error {
	formatted := k.format()
	n, err := fmt.Fprint(w, formatted)
	if err != nil {
		return err
	}
	if n != len(formatted) {
		return io.ErrShortWrite
	}
	return nil
}

// authKey returns the file to write user's keys.
func authKey() string {
	var home string
	if current, err := user.Current(); err == nil {
		home = current.HomeDir
	} else {
		home = os.ExpandEnv("$HOME")
	}
	return path.Join(home, ".ssh", "authorized_keys")
}

// writeKeys serializes the given key in the authorized_keys file (of the
// current user).
func writeKey(k *Key) error {
	file, err := fs.Filesystem().OpenFile(authKey(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	return k.dump(file)
}

// Writes `key` in authorized_keys file (from current user)
// It does not writes in the database, there is no need for that since the key
// object is embedded on the user's document
func addKey(name, body, username string) error {
	key, err := newKey(name, username, body)
	if err != nil {
		return err
	}
	err = db.Session.Key().Insert(key)
	if err != nil {
		if e, ok := err.(*mgo.LastError); ok && e.Code == 11000 {
			return ErrDuplicateKey
		}
		return err
	}
	return writeKey(key)
}

func addKeys(keys map[string]string, username string) error {
	for name, k := range keys {
		err := addKey(name, k, username)
		if err != nil {
			return err
		}
	}
	return nil
}

func remove(k *Key) error {
	formatted := k.format()
	file, err := fs.Filesystem().OpenFile(authKey(), os.O_RDWR|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	lines := make([]string, 0, 10)
	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('\n')
	for line != "" {
		if line != formatted {
			lines = append(lines, line)
		}
		line, _ = reader.ReadString('\n')
	}
	file.Truncate(0)
	file.Seek(0, 0)
	content := strings.Join(lines, "")
	n, err := file.WriteString(content)
	if err != nil {
		return err
	}
	if n != len(content) {
		return io.ErrShortWrite
	}
	return nil
}

func removeUserKeys(username string) error {
	var keys []Key
	q := bson.M{"username": username}
	err := db.Session.Key().Find(q).All(&keys)
	if err != nil {
		return err
	}
	db.Session.Key().RemoveAll(q)
	for _, k := range keys {
		remove(&k)
	}
	return nil
}

// removes a key from the database and the authorized_keys file.
func removeKey(name, username string) error {
	var k Key
	err := db.Session.Key().Find(bson.M{"name": name, "username": username}).One(&k)
	if err != nil {
		return ErrKeyNotFound
	}
	db.Session.Key().Remove(k)
	return remove(&k)
}

type KeyList []Key

func (keys KeyList) MarshalJSON() ([]byte, error) {
	m := make(map[string]string, len(keys))
	for _, key := range keys {
		m[key.Name] = key.String()
	}
	return json.Marshal(m)
}

// ListKeys lists all user's keys.
//
// If the user is not found, returns an error
func ListKeys(uName string) (KeyList, error) {
	if n, err := db.Session.User().FindId(uName).Count(); err != nil || n != 1 {
		return nil, ErrUserNotFound
	}
	var keys []Key
	err := db.Session.Key().Find(bson.M{"username": uName}).All(&keys)
	return KeyList(keys), err
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/db"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
	"os"
	"path"
)

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) {
	return len(p) / 2, nil
}

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) {
	return 0, errors.New("Failed")
}

const rawKey = "ssh-dss AAAAB3NzaC1kc3MAAACBAIHfSDLpSCfIIVEJ/Is3RFMQhsCi7WZtFQeeyfi+DzVP0NGX4j/rMoQEHgXgNlOKVCJvPk5e00tukSv6iVzJPFcozArvVaoCc5jCoDi5Ef8k3Jil4Q7qNjcoRDDyqjqLcaviJEz5GrtmqAyXEIzJ447BxeEdw3Z7UrIWYcw2YyArAAAAFQD7wiOGZIoxu4XIOoeEe5aToTxN1QAAAIAZNAbJyOnNceGcgRRgBUPfY5ChX+9A29n2MGnyJ/Cxrhuh8d7B0J8UkvEBlfgQICq1UDZbC9q5NQprwD47cGwTjUZ0Z6hGpRmEEZdzsoj9T6vkLiteKH3qLo7IPVx4mV6TTF6PWQbQMUsuxjuDErwS9nhtTM4nkxYSmUbnWb6wfwAAAIB2qm/1J6Jl8bByBaMQ/ptbm4wQCvJ9Ll9u6qtKy18D4ldoXM0E9a1q49swml5CPFGyU+cgPRhEjN5oUr5psdtaY8CHa2WKuyIVH3B8UhNzqkjpdTFSpHs6tGluNVC+SQg1MVwfG2wsZUdkUGyn+6j8ZZarUfpAmbb5qJJpgMFEKQ== f@xikinbook.local"
const body = "ssh-dss AAAAB3NzaC1kc3MAAACBAIHfSDLpSCfIIVEJ/Is3RFMQhsCi7WZtFQeeyfi+DzVP0NGX4j/rMoQEHgXgNlOKVCJvPk5e00tukSv6iVzJPFcozArvVaoCc5jCoDi5Ef8k3Jil4Q7qNjcoRDDyqjqLcaviJEz5GrtmqAyXEIzJ447BxeEdw3Z7UrIWYcw2YyArAAAAFQD7wiOGZIoxu4XIOoeEe5aToTxN1QAAAIAZNAbJyOnNceGcgRRgBUPfY5ChX+9A29n2MGnyJ/Cxrhuh8d7B0J8UkvEBlfgQICq1UDZbC9q5NQprwD47cGwTjUZ0Z6hGpRmEEZdzsoj9T6vkLiteKH3qLo7IPVx4mV6TTF6PWQbQMUsuxjuDErwS9nhtTM4nkxYSmUbnWb6wfwAAAIB2qm/1J6Jl8bByBaMQ/ptbm4wQCvJ9Ll9u6qtKy18D4ldoXM0E9a1q49swml5CPFGyU+cgPRhEjN5oUr5psdtaY8CHa2WKuyIVH3B8UhNzqkjpdTFSpHs6tGluNVC+SQg1MVwfG2wsZUdkUGyn+6j8ZZarUfpAmbb5qJJpgMFEKQ==\n"
const comment = "f@xikinbook.local"
const otherKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCaNZSIEyP6FSdCX0WHDcUFTvebNbvqKiiLEiC7NTGvKrT15r2MtCDi4EPi4Ul+UyxWqb2D7FBnK1UmIcEFHd/ZCnBod2/FSplGOIbIb2UVVbqPX5Alv7IBCMyZJD14ex5cFh16zoqOsPOkOD803LMIlNvXPDDwKjY4TVOQV1JtA2tbZXvYUchqhTcKPxt5BDBZbeQkMMgUgHIEz6IueglFB3+dIZfrzlmM8CVSElKZOpucnJ5JOpGh3paSO/px2ZEcvY8WvjFdipvAWsis75GG/04F641I6XmYlo9fib/YytBXS23szqmvOqEqAopFnnGkDEo+LWI0+FXgPE8lc5BD"

func (s *S) TestNewKey(c *gocheck.C) {
	k, err := newKey("key1", "me@tsuru.io", rawKey)
	c.Assert(err, gocheck.IsNil)
	c.Assert(k.Name, gocheck.Equals, "key1")
	c.Assert(k.Body, gocheck.Equals, body)
	c.Assert(k.Comment, gocheck.Equals, comment)
	c.Assert(k.UserName, gocheck.Equals, "me@tsuru.io")
}

func (s *S) TestNewKeyInvalidKey(c *gocheck.C) {
	raw := "ssh-dss ASCCDD== invalid@tsuru.io"
	k, err := newKey("key1", "me@tsuru.io", raw)
	c.Assert(k, gocheck.IsNil)
	c.Assert(err, gocheck.Equals, ErrInvalidKey)
}

func (s *S) TestKeyString(c *gocheck.C) {
	k := Key{Body: "ssh-dss not-secret", Comment: "me@host"}
	c.Assert(k.String(), gocheck.Equals, k.Body+" "+k.Comment)
}

func (s *S) TestKeyStringNewLine(c *gocheck.C) {
	k := Key{Body: "ssh-dss not-secret\n", Comment: "me@host"}
	c.Assert(k.String(), gocheck.Equals, "ssh-dss not-secret me@host")
}

func (s *S) TestKeyStringNoComment(c *gocheck.C) {
	k := Key{Body: "ssh-dss not-secret"}
	c.Assert(k.String(), gocheck.Equals, k.Body)
}

func (s *S) TestFormatKeyShouldAddSshLoginRestrictionsAtBegining(c *gocheck.C) {
	key := Key{
		Name:     "my-key",
		Body:     "somekey\n",
		Comment:  "me@host",
		UserName: "brain",
	}
	got := key.format()
	expected := fmt.Sprintf("no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty,command=.* %s\n", &key)
	c.Assert(got, gocheck.Matches, expected)
}

func (s *S) TestFormatKeyShouldAddCommandAfterSshRestrictions(c *gocheck.C) {
	key := Key{
		Name:     "my-key",
		Body:     "somekey\n",
		Comment:  "me@host",
		UserName: "brain",
	}
	got := key.format()
	p, err := config.GetString("bin-path")
	c.Assert(err, gocheck.IsNil)
	expected := fmt.Sprintf(`no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty,command="%s brain" %s`+"\n", p, &key)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestFormatKeyShouldGetCommandPathFromGandalfConf(c *gocheck.C) {
	oldConf, err := config.GetString("bin-path")
	c.Assert(err, gocheck.IsNil)
	config.Set("bin-path", "/foo/bar/hi.go")
	defer config.Set("bin-path", oldConf)
	key := Key{
		Name:     "my-key",
		Body:     "somekey\n",
		Comment:  "me@host",
		UserName: "dash",
	}
	got := key.format()
	expected := fmt.Sprintf(`no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty,command="/foo/bar/hi.go dash" %s`+"\n", &key)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestFormatKeyShouldAppendUserNameAsCommandParameter(c *gocheck.C) {
	p, err := config.GetString("bin-path")
	c.Assert(err, gocheck.IsNil)
	key := Key{
		Name:     "my-key",
		Body:     "somekey\n",
		Comment:  "me@host",
		UserName: "someuser",
	}
	got := key.format()
	expected := fmt.Sprintf(`no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty,command="%s someuser" %s`+"\n", p, &key)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestDump(c *gocheck.C) {
	var buf bytes.Buffer
	key := Key{
		Name:     "my-key",
		Body:     "somekey\n",
		Comment:  "me@host",
		UserName: "someuser",
	}
	err := key.dump(&buf)
	c.Assert(err, gocheck.IsNil)
	c.Assert(buf.String(), gocheck.Equals, key.format())
}

func (s *S) TestDumpShortWrite(c *gocheck.C) {
	key := Key{
		Name:     "my-key",
		Body:     "somekey\n",
		Comment:  "me@host",
		UserName: "someuser",
	}
	err := key.dump(shortWriter{})
	c.Assert(err, gocheck.Equals, io.ErrShortWrite)
}

func (s *S) TestDumpWriteFailure(c *gocheck.C) {
	key := Key{
		Name:     "my-key",
		Body:     "somekey\n",
		Comment:  "me@host",
		UserName: "someuser",
	}
	err := key.dump(failWriter{})
	c.Assert(err, gocheck.NotNil)
}

func (s *S) TestAuthKeysShouldBeAbsolutePathToUsersAuthorizedKeysByDefault(c *gocheck.C) {
	home := os.Getenv("HOME")
	expected := path.Join(home, ".ssh", "authorized_keys")
	c.Assert(authKey(), gocheck.Equals, expected)
}

func (s *S) TestWriteKey(c *gocheck.C) {
	key, err := newKey("my-key", "me@tsuru.io", rawKey)
	c.Assert(err, gocheck.IsNil)
	writeKey(key)
	f, err := s.rfs.Open(authKey())
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	got := string(b)
	c.Assert(got, gocheck.Equals, key.format())
}

func (s *S) TestWriteTwoKeys(c *gocheck.C) {
	key1 := Key{
		Name:     "my-key",
		Body:     "ssh-dss mykeys-not-secret",
		Comment:  "me@machine",
		UserName: "gopher",
	}
	key2 := Key{
		Name:     "your-key",
		Body:     "ssh-dss yourkeys-not-secret",
		Comment:  "me@machine",
		UserName: "glenda",
	}
	writeKey(&key1)
	writeKey(&key2)
	expected := key1.format() + key2.format()
	f, err := s.rfs.Open(authKey())
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	got := string(b)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestAddKeyStoresKeyInTheDatabase(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	var k Key
	err = db.Session.Key().Find(bson.M{"name": "key1"}).One(&k)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Key().Remove(bson.M{"name": "key1"})
	c.Assert(k.Name, gocheck.Equals, "key1")
	c.Assert(k.UserName, gocheck.Equals, "gopher")
	c.Assert(k.Comment, gocheck.Equals, comment)
	c.Assert(k.Body, gocheck.Equals, body)
}

func (s *S) TestAddKeyShouldSaveTheKeyInTheAuthorizedKeys(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Key().Remove(bson.M{"name": "key1"})
	var k Key
	err = db.Session.Key().Find(bson.M{"name": "key1"}).One(&k)
	c.Assert(err, gocheck.IsNil)
	f, err := s.rfs.Open(authKey())
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	c.Assert(string(b), gocheck.Equals, k.format())
}

func (s *S) TestAddKeyDuplicate(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Key().Remove(bson.M{"name": "key1"})
	err = addKey("key2", rawKey, "gopher")
	c.Assert(err, gocheck.Equals, ErrDuplicateKey)
}

func (s *S) TestAddKeyInvalidKey(c *gocheck.C) {
	err := addKey("key1", "something-invalid", "gopher")
	c.Assert(err, gocheck.Equals, ErrInvalidKey)
}

func (s *S) TestRemoveKeyDeletesFromDB(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	err = removeKey("key1", "gopher")
	c.Assert(err, gocheck.IsNil)
	count, err := db.Session.Key().Find(bson.M{"name": "key1"}).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(count, gocheck.Equals, 0)
}

func (s *S) TestRemoveKeyDeletesOnlyTheRightKey(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	defer removeKey("key1", "gopher")
	err = addKey("key1", otherKey, "glenda")
	c.Assert(err, gocheck.IsNil)
	err = removeKey("key1", "glenda")
	c.Assert(err, gocheck.IsNil)
	count, err := db.Session.Key().Find(bson.M{"name": "key1", "username": "gopher"}).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(count, gocheck.Equals, 1)
}

func (s *S) TestRemoveUnknownKey(c *gocheck.C) {
	err := removeKey("wut", "glenda")
	c.Assert(err, gocheck.Equals, ErrKeyNotFound)
}

func (s *S) TestRemoveKeyRemovesFromAuthorizedKeys(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	err = removeKey("key1", "gopher")
	f, err := s.rfs.Open(authKey())
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	got := string(b)
	c.Assert(got, gocheck.Equals, "")
}

func (s *S) TestRemoveKeyKeepOtherKeys(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	defer removeKey("key1", "gopher")
	err = addKey("key2", otherKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	err = removeKey("key2", "gopher")
	c.Assert(err, gocheck.IsNil)
	var key Key
	err = db.Session.Key().Find(bson.M{"name": "key1"}).One(&key)
	c.Assert(err, gocheck.IsNil)
	f, err := s.rfs.Open(authKey())
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	got := string(b)
	c.Assert(got, gocheck.Equals, key.format())
}

func (s *S) TestRemoveUserKeys(c *gocheck.C) {
	err := addKey("key1", rawKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	defer removeKey("key1", "gopher")
	err = addKey("key1", otherKey, "glenda")
	c.Assert(err, gocheck.IsNil)
	err = removeUserKeys("glenda")
	c.Assert(err, gocheck.IsNil)
	var key Key
	err = db.Session.Key().Find(bson.M{"name": "key1"}).One(&key)
	c.Assert(err, gocheck.IsNil)
	f, err := s.rfs.Open(authKey())
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	got := string(b)
	c.Assert(got, gocheck.Equals, key.format())
}

func (s *S) TestRemoveUserMultipleKeys(c *gocheck.C) {
	err := addKey("key1", rawKey, "glenda")
	c.Assert(err, gocheck.IsNil)
	err = addKey("key1", otherKey, "glenda")
	c.Assert(err, gocheck.IsNil)
	err = removeUserKeys("glenda")
	c.Assert(err, gocheck.IsNil)
	count, err := db.Session.Key().Find(nil).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(count, gocheck.Equals, 0)
	f, err := s.rfs.Open(authKey())
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	got := string(b)
	c.Assert(got, gocheck.Equals, "")
}

func (s *S) TestKeyListJSON(c *gocheck.C) {
	keys := []Key{
		{Name: "key1", Body: "ssh-dss not-secret", Comment: "me@host1"},
		{Name: "key2", Body: "ssh-dss not-secret1", Comment: "me@host2"},
		{Name: "another-key", Body: "ssh-rsa not-secret", Comment: "me@work"},
	}
	expected := map[string]string{
		keys[0].Name: keys[0].String(),
		keys[1].Name: keys[1].String(),
		keys[2].Name: keys[2].String(),
	}
	var got map[string]string
	b, err := KeyList(keys).MarshalJSON()
	c.Assert(err, gocheck.IsNil)
	err = json.Unmarshal(b, &got)
	c.Assert(err, gocheck.IsNil)
	c.Assert(got, gocheck.DeepEquals, expected)
}

func (s *S) TestListKeys(c *gocheck.C) {
	user := map[string]string{"_id": "glenda"}
	err := db.Session.User().Insert(user)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(user)
	err = addKey("key1", rawKey, "glenda")
	c.Assert(err, gocheck.IsNil)
	err = addKey("key2", otherKey, "glenda")
	c.Assert(err, gocheck.IsNil)
	defer removeUserKeys("glenda")
	var expected []Key
	err = db.Session.Key().Find(nil).All(&expected)
	c.Assert(err, gocheck.IsNil)
	got, err := ListKeys("glenda")
	c.Assert(err, gocheck.IsNil)
	c.Assert(got, gocheck.DeepEquals, KeyList(expected))
}

func (s *S) TestListKeysUnknownUser(c *gocheck.C) {
	got, err := ListKeys("glenda")
	c.Assert(got, gocheck.IsNil)
	c.Assert(err, gocheck.Equals, ErrUserNotFound)
}

func (s *S) TestListKeysEmpty(c *gocheck.C) {
	user := map[string]string{"_id": "gopher"}
	err := db.Session.User().Insert(user)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(user)
	got, err := ListKeys("gopher")
	c.Assert(err, gocheck.IsNil)
	c.Assert(got, gocheck.HasLen, 0)
}

func (s *S) TestListKeysFromTheUserOnly(c *gocheck.C) {
	user := map[string]string{"_id": "gopher"}
	err := db.Session.User().Insert(user)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(user)
	user2 := map[string]string{"_id": "glenda"}
	err = db.Session.User().Insert(user2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(user2)
	err = addKey("key1", rawKey, "glenda")
	c.Assert(err, gocheck.IsNil)
	err = addKey("key1", otherKey, "gopher")
	c.Assert(err, gocheck.IsNil)
	defer removeUserKeys("glenda")
	defer removeUserKeys("gopher")
	var expected []Key
	err = db.Session.Key().Find(bson.M{"username": "gopher"}).All(&expected)
	c.Assert(err, gocheck.IsNil)
	got, err := ListKeys("gopher")
	c.Assert(err, gocheck.IsNil)
	c.Assert(got, gocheck.DeepEquals, KeyList(expected))
}

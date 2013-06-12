// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/globocom/config"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/fs"
	"github.com/globocom/gandalf/repository"
	fstesting "github.com/globocom/tsuru/fs/testing"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
	"os"
	"path"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct {
	rfs *fstesting.RecordingFs
}

var _ = gocheck.Suite(&S{})

func (s *S) authKeysContent(c *gocheck.C) string {
	authFile := path.Join(os.Getenv("HOME"), ".ssh", "authorized_keys")
	f, err := fs.Filesystem().OpenFile(authFile, os.O_RDWR, 0755)
	c.Assert(err, gocheck.IsNil)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	c.Assert(err, gocheck.IsNil)
	return string(b)
}

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/gandalf.conf")
	c.Check(err, gocheck.IsNil)
	config.Set("database:name", "gandalf_user_tests")
	db.Connect()
}

func (s *S) SetUpTest(c *gocheck.C) {
	s.rfs = &fstesting.RecordingFs{}
	fs.Fsystem = s.rfs
}

func (s *S) TearDownTest(c *gocheck.C) {
	s.rfs.Remove(authKey())
}

func (s *S) TearDownSuite(c *gocheck.C) {
	fs.Fsystem = nil
	db.Session.DB.DropDatabase()
}

func (s *S) TestNewUserReturnsAStructFilled(c *gocheck.C) {
	u, err := New("someuser", map[string]string{"somekey": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(bson.M{"_id": u.Name})
	defer db.Session.Key().Remove(bson.M{"name": "somekey"})
	c.Assert(u.Name, gocheck.Equals, "someuser")
	var key Key
	err = db.Session.Key().Find(bson.M{"name": "somekey"}).One(&key)
	c.Assert(err, gocheck.IsNil)
	c.Assert(key.Name, gocheck.Equals, "somekey")
	c.Assert(key.Body, gocheck.Equals, body)
	c.Assert(key.Comment, gocheck.Equals, comment)
	c.Assert(key.UserName, gocheck.Equals, u.Name)
}

func (s *S) TestNewUserShouldStoreUserInDatabase(c *gocheck.C) {
	u, err := New("someuser", map[string]string{"somekey": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(bson.M{"_id": u.Name})
	defer db.Session.Key().Remove(bson.M{"name": "somekey"})
	err = db.Session.User().FindId(u.Name).One(&u)
	c.Assert(err, gocheck.IsNil)
	c.Assert(u.Name, gocheck.Equals, "someuser")
	n, err := db.Session.Key().Find(bson.M{"name": "somekey"}).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(n, gocheck.Equals, 1)
}

func (s *S) TestNewChecksIfUserIsValidBeforeStoring(c *gocheck.C) {
	_, err := New("", map[string]string{})
	c.Assert(err, gocheck.NotNil)
	got := err.Error()
	expected := "Validation Error: user name is not valid"
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestNewWritesKeyInAuthorizedKeys(c *gocheck.C) {
	u, err := New("piccolo", map[string]string{"somekey": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(bson.M{"_id": u.Name})
	defer db.Session.Key().Remove(bson.M{"name": "somekey"})
	var key Key
	err = db.Session.Key().Find(bson.M{"name": "somekey"}).One(&key)
	c.Assert(err, gocheck.IsNil)
	keys := s.authKeysContent(c)
	c.Assert(keys, gocheck.Equals, key.format())
}

func (s *S) TestIsValidReturnsErrorWhenUserDoesNotHaveAName(c *gocheck.C) {
	u := User{}
	v, err := u.isValid()
	c.Assert(v, gocheck.Equals, false)
	c.Assert(err, gocheck.NotNil)
	expected := "Validation Error: user name is not valid"
	got := err.Error()
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestIsValidShouldAcceptEmailsAsUserName(c *gocheck.C) {
	u := User{Name: "r2d2@gmail.com"}
	v, err := u.isValid()
	c.Assert(err, gocheck.IsNil)
	c.Assert(v, gocheck.Equals, true)
}

func (s *S) TestRemove(c *gocheck.C) {
	u, err := New("someuser", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	err = Remove(u.Name)
	c.Assert(err, gocheck.IsNil)
	lenght, err := db.Session.User().FindId(u.Name).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)
}

func (s *S) TestRemoveRemovesKeyFromAuthorizedKeysFile(c *gocheck.C) {
	u, err := New("gandalf", map[string]string{"somekey": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Key().Remove(bson.M{"name": "somekey"})
	err = Remove(u.Name)
	c.Assert(err, gocheck.IsNil)
	got := s.authKeysContent(c)
	c.Assert(got, gocheck.Equals, "")
}

func (s *S) TestRemoveInexistentUserReturnsDescriptiveMessage(c *gocheck.C) {
	err := Remove("otheruser")
	c.Assert(err, gocheck.ErrorMatches, "Could not remove user: not found")
}

func (s *S) TestRemoveDoesNotRemovesUserWhenUserIsTheOnlyOneAssciatedWithOneRepository(c *gocheck.C) {
	u, err := New("silver", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	r := s.createRepo("run", []string{u.Name}, c)
	defer db.Session.Repository().Remove(bson.M{"_id": r.Name})
	defer db.Session.User().Remove(bson.M{"_id": u.Name})
	err = Remove(u.Name)
	c.Assert(err, gocheck.ErrorMatches, "^Could not remove user: user is the only one with access to at least one of it's repositories$")
}

func (s *S) TestRemoveRevokesAccessToReposWithMoreThanOneUserAssociated(c *gocheck.C) {
	u, r, r2 := s.userPlusRepos(c)
	defer db.Session.Repository().Remove(bson.M{"_id": r.Name})
	defer db.Session.Repository().Remove(bson.M{"_id": r2.Name})
	defer db.Session.User().Remove(bson.M{"_id": u.Name})
	err := Remove(u.Name)
	c.Assert(err, gocheck.IsNil)
	s.retrieveRepos(r, r2, c)
	c.Assert(r.Users, gocheck.DeepEquals, []string{"slot"})
	c.Assert(r2.Users, gocheck.DeepEquals, []string{"cnot"})
}

func (s *S) retrieveRepos(r, r2 *repository.Repository, c *gocheck.C) {
	err := db.Session.Repository().FindId(r.Name).One(&r)
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r2.Name).One(&r2)
	c.Assert(err, gocheck.IsNil)
}

func (s *S) userPlusRepos(c *gocheck.C) (*User, *repository.Repository, *repository.Repository) {
	u, err := New("silver", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	r := s.createRepo("run", []string{u.Name, "slot"}, c)
	r2 := s.createRepo("stay", []string{u.Name, "cnot"}, c)
	return u, &r, &r2
}

func (s *S) createRepo(name string, users []string, c *gocheck.C) repository.Repository {
	r := repository.Repository{Name: name, Users: users}
	err := db.Session.Repository().Insert(&r)
	c.Assert(err, gocheck.IsNil)
	return r
}

func (s *S) TestHandleAssociatedRepositoriesShouldRevokeAccessToRepoWithMoreThanOneUserAssociated(c *gocheck.C) {
	u, r, r2 := s.userPlusRepos(c)
	defer db.Session.Repository().RemoveId(r.Name)
	defer db.Session.Repository().RemoveId(r2.Name)
	defer db.Session.User().RemoveId(u.Name)
	err := u.handleAssociatedRepositories()
	c.Assert(err, gocheck.IsNil)
	s.retrieveRepos(r, r2, c)
	c.Assert(r.Users, gocheck.DeepEquals, []string{"slot"})
	c.Assert(r2.Users, gocheck.DeepEquals, []string{"cnot"})
}

func (s *S) TestHandleAssociateRepositoriesReturnsErrorWhenUserIsOnlyOneWithAccessToAtLeastOneRepo(c *gocheck.C) {
	u, err := New("umi", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	r := s.createRepo("proj1", []string{"umi"}, c)
	defer db.Session.User().RemoveId(u.Name)
	defer db.Session.Repository().RemoveId(r.Name)
	err = u.handleAssociatedRepositories()
	expected := "^Could not remove user: user is the only one with access to at least one of it's repositories$"
	c.Assert(err, gocheck.ErrorMatches, expected)
}

func (s *S) TestAddKeyShouldSaveTheKeyInTheDatabase(c *gocheck.C) {
	u, err := New("umi", map[string]string{})
	defer db.Session.User().RemoveId(u.Name)
	k := map[string]string{"somekey": rawKey}
	err = AddKey("umi", k)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Key().Remove(bson.M{"name": "somekey"})
	var key Key
	err = db.Session.Key().Find(bson.M{"name": "somekey"}).One(&key)
	c.Assert(err, gocheck.IsNil)
	c.Assert(key.Name, gocheck.Equals, "somekey")
	c.Assert(key.Body, gocheck.Equals, body)
	c.Assert(key.Comment, gocheck.Equals, comment)
	c.Assert(key.UserName, gocheck.Equals, u.Name)
}

func (s *S) TestAddKeyShouldWriteKeyInAuthorizedKeys(c *gocheck.C) {
	u, err := New("umi", map[string]string{})
	defer db.Session.User().RemoveId(u.Name)
	defer db.Session.Key().Remove(bson.M{"name": "somekey"})
	k := map[string]string{"somekey": rawKey}
	err = AddKey("umi", k)
	c.Assert(err, gocheck.IsNil)
	var key Key
	err = db.Session.Key().Find(bson.M{"name": "somekey"}).One(&key)
	content := s.authKeysContent(c)
	c.Assert(content, gocheck.Equals, key.format())
}

func (s *S) TestAddKeyShouldReturnCustomErrorWhenUserDoesNotExists(c *gocheck.C) {
	err := AddKey("umi", map[string]string{"somekey": "ssh-rsa mykey umi@host"})
	c.Assert(err, gocheck.Equals, ErrUserNotFound)
}

func (s *S) TestRemoveKeyShouldRemoveKeyFromTheDatabase(c *gocheck.C) {
	u, err := New("luke", map[string]string{"homekey": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId(u.Name)
	err = RemoveKey("luke", "homekey")
	c.Assert(err, gocheck.IsNil)
	count, err := db.Session.Key().Find(bson.M{"name": "homekey", "username": u.Name}).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(count, gocheck.Equals, 0)
}

func (s *S) TestRemoveKeyShouldRemoveFromAuthorizedKeysFile(c *gocheck.C) {
	u, err := New("luke", map[string]string{"homekey": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId(u.Name)
	defer db.Session.Key().Remove(bson.M{"name": "homekey"})
	err = RemoveKey("luke", "homekey")
	c.Assert(err, gocheck.IsNil)
	content := s.authKeysContent(c)
	c.Assert(content, gocheck.Equals, "")
}

func (s *S) TestRemoveUnknownKeyFromUser(c *gocheck.C) {
	u, err := New("luke", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId(u.Name)
	err = RemoveKey("luke", "homekey")
	c.Assert(err, gocheck.Equals, ErrKeyNotFound)
}

func (s *S) TestRemoveKeyShouldReturnFormatedErrorMsgWhenUserDoesNotExists(c *gocheck.C) {
	err := RemoveKey("luke", "homekey")
	c.Assert(err, gocheck.Equals, ErrUserNotFound)
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"encoding/json"
	"fmt"
	"github.com/globocom/commandmocker"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/fs"
	fstesting "github.com/globocom/tsuru/fs/testing"
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
	"path"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct {
	tmpdir string
}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/gandalf.conf")
	c.Assert(err, gocheck.IsNil)
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "gandalf_repository_tests")
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestNewShouldCreateANewRepository(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	users := []string{"smeagol", "saruman"}
	r, err := New("myRepo", users, true)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().Remove(bson.M{"_id": "myRepo"})
	c.Assert(r.Name, gocheck.Equals, "myRepo")
	c.Assert(r.Users, gocheck.DeepEquals, users)
	c.Assert(r.IsPublic, gocheck.Equals, true)
}

func (s *S) TestNewShouldRecordItOnDatabase(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	r, err := New("someRepo", []string{"smeagol"}, true)
	defer db.Session.Repository().Remove(bson.M{"_id": "someRepo"})
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().Find(bson.M{"_id": "someRepo"}).One(&r)
	c.Assert(err, gocheck.IsNil)
	c.Assert(r.Name, gocheck.Equals, "someRepo")
	c.Assert(r.Users, gocheck.DeepEquals, []string{"smeagol"})
	c.Assert(r.IsPublic, gocheck.Equals, true)
}

func (s *S) TestNewBreaksOnValidationError(c *gocheck.C) {
	_, err := New("", []string{"smeagol"}, false)
	c.Check(err, gocheck.NotNil)
	expected := "Validation Error: repository name is not valid"
	got := err.Error()
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestRepositoryIsNotValidWithoutAName(c *gocheck.C) {
	r := Repository{Users: []string{"gollum"}, IsPublic: true}
	v, err := r.isValid()
	c.Assert(v, gocheck.Equals, false)
	c.Check(err, gocheck.NotNil)
	got := err.Error()
	expected := "Validation Error: repository name is not valid"
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestRepositoryIsNotValidWithInvalidName(c *gocheck.C) {
	r := Repository{Name: "foo bar", Users: []string{"gollum"}, IsPublic: true}
	v, err := r.isValid()
	c.Assert(v, gocheck.Equals, false)
	c.Check(err, gocheck.NotNil)
	got := err.Error()
	expected := "Validation Error: repository name is not valid"
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestRepositoryShoudBeInvalidWIthoutAnyUsers(c *gocheck.C) {
	r := Repository{Name: "foo_bar", Users: []string{}, IsPublic: true}
	v, err := r.isValid()
	c.Assert(v, gocheck.Equals, false)
	c.Assert(err, gocheck.NotNil)
	got := err.Error()
	expected := "Validation Error: repository should have at least one user"
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestRepositoryShouldBeValidWithoutIsPublic(c *gocheck.C) {
	r := Repository{Name: "someName", Users: []string{"smeagol"}}
	v, _ := r.isValid()
	c.Assert(v, gocheck.Equals, true)
}

func (s *S) TestNewShouldCreateNewGitBareRepository(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	_, err = New("myRepo", []string{"pumpkin"}, true)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().Remove(bson.M{"_id": "myRepo"})
	c.Assert(commandmocker.Ran(tmpdir), gocheck.Equals, true)
}

func (s *S) TestNewShouldNotStoreRepoInDbWhenBareCreationFails(c *gocheck.C) {
	dir, err := commandmocker.Error("git", "", 1)
	c.Check(err, gocheck.IsNil)
	defer commandmocker.Remove(dir)
	r, err := New("myRepo", []string{"pumpkin"}, true)
	c.Check(err, gocheck.NotNil)
	err = db.Session.Repository().Find(bson.M{"_id": r.Name}).One(&r)
	c.Assert(err, gocheck.ErrorMatches, "^not found$")
}

func (s *S) TestRemoveShouldRemoveBareRepositoryFromFileSystem(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	rfs := &fstesting.RecordingFs{FileContent: "foo"}
	fs.Fsystem = rfs
	defer func() { fs.Fsystem = nil }()
	r, err := New("myRepo", []string{"pumpkin"}, false)
	c.Assert(err, gocheck.IsNil)
	err = Remove(r.Name)
	c.Assert(err, gocheck.IsNil)
	action := "removeall " + path.Join(bareLocation(), "myRepo.git")
	c.Assert(rfs.HasAction(action), gocheck.Equals, true)
}

func (s *S) TestRemoveShouldRemoveRepositoryFromDatabase(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	rfs := &fstesting.RecordingFs{FileContent: "foo"}
	fs.Fsystem = rfs
	defer func() { fs.Fsystem = nil }()
	r, err := New("myRepo", []string{"pumpkin"}, false)
	c.Assert(err, gocheck.IsNil)
	err = Remove(r.Name)
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().Find(bson.M{"_id": r.Name}).One(&r)
	c.Assert(err, gocheck.ErrorMatches, "^not found$")
}

func (s *S) TestRemoveShouldReturnMeaningfulErrorWhenRepositoryDoesNotExistsInDatabase(c *gocheck.C) {
	rfs := &fstesting.RecordingFs{FileContent: "foo"}
	fs.Fsystem = rfs
	defer func() { fs.Fsystem = nil }()
	r := &Repository{Name: "fooBar"}
	err := Remove(r.Name)
	c.Assert(err, gocheck.ErrorMatches, "^Could not remove repository: not found$")
}

func (s *S) TestRename(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	repository, err := New("freedom", []string{"fss@corp.globo.com", "andrews@corp.globo.com"}, true)
	c.Check(err, gocheck.IsNil)
	commandmocker.Remove(tmpdir)
	rfs := &fstesting.RecordingFs{}
	fs.Fsystem = rfs
	defer func() { fs.Fsystem = nil }()
	err = Rename(repository.Name, "free")
	c.Assert(err, gocheck.IsNil)
	_, err = Get("freedom")
	c.Assert(err, gocheck.NotNil)
	repo, err := Get("free")
	c.Assert(err, gocheck.IsNil)
	repository.Name = "free"
	c.Assert(repo, gocheck.DeepEquals, *repository)
	action := "rename " + barePath("freedom") + " " + barePath("free")
	c.Assert(rfs.HasAction(action), gocheck.Equals, true)
}

func (s *S) TestRenameNotFound(c *gocheck.C) {
	err := Rename("something", "free")
	c.Assert(err, gocheck.NotNil)
}

func (s *S) TestGitURL(c *gocheck.C) {
	host, err := config.GetString("host")
	c.Assert(err, gocheck.IsNil)
	remote := (&Repository{Name: "lol"}).GitURL()
	c.Assert(remote, gocheck.Equals, fmt.Sprintf("git://%s/lol.git", host))
}

func (s *S) TestSshURL(c *gocheck.C) {
	host, err := config.GetString("host")
	c.Assert(err, gocheck.IsNil)
	remote := (&Repository{Name: "lol"}).SshURL()
	c.Assert(remote, gocheck.Equals, fmt.Sprintf("git@%s:lol.git", host))
}

func (s *S) TestSshURLUseUidFromConfigFile(c *gocheck.C) {
	uid, err := config.GetString("uid")
	c.Assert(err, gocheck.IsNil)
	host, err := config.GetString("host")
	c.Assert(err, gocheck.IsNil)
	config.Set("uid", "test")
	defer config.Set("uid", uid)
	remote := (&Repository{Name: "f#"}).SshURL()
	c.Assert(remote, gocheck.Equals, fmt.Sprintf("test@%s:f#.git", host))
}

func (s *S) TestGrantAccessShouldAddUserToListOfRepositories(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	r, err := New("proj1", []string{"someuser"}, true)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().RemoveId(r.Name)
	r2, err := New("proj2", []string{"otheruser"}, true)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().RemoveId(r2.Name)
	u := struct {
		Name string `bson:"_id"`
	}{Name: "lolcat"}
	err = db.Session.User().Insert(&u)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId(u.Name)
	err = GrantAccess([]string{r.Name, r2.Name}, []string{u.Name})
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r.Name).One(&r)
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r2.Name).One(&r2)
	c.Assert(err, gocheck.IsNil)
	c.Assert(r.Users, gocheck.DeepEquals, []string{"someuser", u.Name})
	c.Assert(r2.Users, gocheck.DeepEquals, []string{"otheruser", u.Name})
}

func (s *S) TestGrantAccessShouldAddFirstUserIntoRepositoryDocument(c *gocheck.C) {
	r := Repository{Name: "proj1"}
	err := db.Session.Repository().Insert(&r)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().RemoveId(r.Name)
	r2 := Repository{Name: "proj2"}
	err = db.Session.Repository().Insert(&r2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().RemoveId(r2.Name)
	err = GrantAccess([]string{r.Name, r2.Name}, []string{"Umi"})
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r.Name).One(&r)
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r.Name).One(&r2)
	c.Assert(err, gocheck.IsNil)
	c.Assert(r.Users, gocheck.DeepEquals, []string{"Umi"})
	c.Assert(r2.Users, gocheck.DeepEquals, []string{"Umi"})
}

func (s *S) TestGrantAccessShouldSkipDuplicatedUsers(c *gocheck.C) {
	r := Repository{Name: "proj1", Users: []string{"umi", "luke", "pade"}}
	err := db.Session.Repository().Insert(&r)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().RemoveId(r.Name)
	err = GrantAccess([]string{r.Name}, []string{"pade"})
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r.Name).One(&r)
	c.Assert(err, gocheck.IsNil)
	c.Assert(r.Users, gocheck.DeepEquals, []string{"umi", "luke", "pade"})
}

func (s *S) TestRevokeAccessShouldRemoveUserFromAllRepositories(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	r, err := New("proj1", []string{"someuser", "umi"}, true)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().RemoveId(r.Name)
	r2, err := New("proj2", []string{"otheruser", "umi"}, true)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().RemoveId(r2.Name)
	err = RevokeAccess([]string{r.Name, r2.Name}, []string{"umi"})
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r.Name).One(&r)
	c.Assert(err, gocheck.IsNil)
	err = db.Session.Repository().FindId(r2.Name).One(&r2)
	c.Assert(err, gocheck.IsNil)
	c.Assert(r.Users, gocheck.DeepEquals, []string{"someuser"})
	c.Assert(r2.Users, gocheck.DeepEquals, []string{"otheruser"})
}

func (s *S) TestConflictingRepositoryNameShouldReturnExplicitError(c *gocheck.C) {
	tmpdir, err := commandmocker.Add("git", "$*")
	c.Assert(err, gocheck.IsNil)
	defer commandmocker.Remove(tmpdir)
	_, err = New("someRepo", []string{"gollum"}, true)
	defer db.Session.Repository().Remove(bson.M{"_id": "someRepo"})
	c.Assert(err, gocheck.IsNil)
	_, err = New("someRepo", []string{"gollum"}, true)
	c.Assert(err, gocheck.ErrorMatches, "A repository with this name already exists.")
}

func (s *S) TestGet(c *gocheck.C) {
	repo := Repository{Name: "somerepo", Users: []string{}}
	err := db.Session.Repository().Insert(repo)
	c.Assert(err, gocheck.IsNil)
	r, err := Get("somerepo")
	c.Assert(err, gocheck.IsNil)
	c.Assert(r, gocheck.DeepEquals, repo)
}

func (s *S) TestMarshalJSON(c *gocheck.C) {
	repo := Repository{Name: "somerepo", Users: []string{}}
	expected := map[string]interface{}{
		"name":    repo.Name,
		"public":  repo.IsPublic,
		"ssh_url": repo.SshURL(),
		"git_url": repo.GitURL(),
	}
	data, err := json.Marshal(&repo)
	c.Assert(err, gocheck.IsNil)
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	c.Assert(err, gocheck.IsNil)
	c.Assert(result, gocheck.DeepEquals, expected)
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/fs"
	"github.com/globocom/gandalf/repository"
	"github.com/globocom/gandalf/user"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
)

type bufferCloser struct {
	*bytes.Buffer
}

func (b bufferCloser) Close() error { return nil }

func get(url string, b io.Reader, c *gocheck.C) (*httptest.ResponseRecorder, *http.Request) {
	return request("GET", url, b, c)
}

func post(url string, b io.Reader, c *gocheck.C) (*httptest.ResponseRecorder, *http.Request) {
	return request("POST", url, b, c)
}

func del(url string, b io.Reader, c *gocheck.C) (*httptest.ResponseRecorder, *http.Request) {
	return request("DELETE", url, b, c)
}

func request(method, url string, b io.Reader, c *gocheck.C) (*httptest.ResponseRecorder, *http.Request) {
	request, err := http.NewRequest(method, url, b)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	return recorder, request
}

func readBody(b io.Reader, c *gocheck.C) string {
	body, err := ioutil.ReadAll(b)
	c.Assert(err, gocheck.IsNil)
	return string(body)
}

func (s *S) authKeysContent(c *gocheck.C) string {
	authKeysPath := path.Join(os.Getenv("HOME"), ".ssh", "authorized_keys")
	f, err := fs.Filesystem().OpenFile(authKeysPath, os.O_RDWR|os.O_EXCL, 0755)
	c.Assert(err, gocheck.IsNil)
	content, err := ioutil.ReadAll(f)
	return string(content)
}

func (s *S) TestNewUser(c *gocheck.C) {
	b := strings.NewReader(fmt.Sprintf(`{"name": "brain", "keys": {"keyname": %q}}`, rawKey))
	recorder, request := post("/user", b, c)
	NewUser(recorder, request)
	defer db.Session.User().Remove(bson.M{"_id": "brain"})
	defer db.Session.Key().Remove(bson.M{"username": "brain"})
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	c.Assert(string(body), gocheck.Equals, "User \"brain\" successfully created\n")
	c.Assert(recorder.Code, gocheck.Equals, 200)
}

func (s *S) TestNewUserShouldSaveInDB(c *gocheck.C) {
	b := strings.NewReader(`{"name": "brain", "keys": {"content": "some id_rsa.pub key.. use your imagination!", "name": "somekey"}}`)
	recorder, request := post("/user", b, c)
	NewUser(recorder, request)
	defer db.Session.User().Remove(bson.M{"_id": "brain"})
	defer db.Session.Key().Remove(bson.M{"username": "brain"})
	var u user.User
	err := db.Session.User().Find(bson.M{"_id": "brain"}).One(&u)
	c.Assert(err, gocheck.IsNil)
	c.Assert(u.Name, gocheck.Equals, "brain")
}

func (s *S) TestNewUserShouldRepassParseBodyErrors(c *gocheck.C) {
	b := strings.NewReader("{]9afe}")
	recorder, request := post("/user", b, c)
	NewUser(recorder, request)
	body := readBody(recorder.Body, c)
	expected := "Got error while parsing body: Could not parse json: invalid character ']' looking for beginning of object key string"
	got := strings.Replace(body, "\n", "", -1)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestNewUserShouldRequireUserName(c *gocheck.C) {
	b := strings.NewReader(`{"name": ""}`)
	recorder, request := post("/user", b, c)
	NewUser(recorder, request)
	body := readBody(recorder.Body, c)
	expected := "Got error while creating user: Validation Error: user name is not valid"
	got := strings.Replace(body, "\n", "", -1)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestNewUserWihoutKeys(c *gocheck.C) {
	b := strings.NewReader(`{"name": "brain"}`)
	recorder, request := post("/user", b, c)
	NewUser(recorder, request)
	defer db.Session.User().Remove(bson.M{"_id": "brain"})
	c.Assert(recorder.Code, gocheck.Equals, 200)
}

func (s *S) TestGetRepository(c *gocheck.C) {
	r := repository.Repository{Name: "onerepo"}
	err := db.Session.Repository().Insert(&r)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().Remove(bson.M{"_id": r.Name})
	recorder, request := get("/repository/onerepo?:name=onerepo", nil, c)
	GetRepository(recorder, request)
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	c.Assert(err, gocheck.IsNil)
	expected := map[string]interface{}{
		"name":    r.Name,
		"public":  r.IsPublic,
		"ssh_url": r.SshURL(),
		"git_url": r.GitURL(),
	}
	c.Assert(data, gocheck.DeepEquals, expected)
}

func (s *S) TestGetRepositoryDoesNotExist(c *gocheck.C) {
	recorder, request := get("/repository/doesnotexists?:name=doesnotexists", nil, c)
	GetRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 500)
}

func (s *S) TestNewRepository(c *gocheck.C) {
	defer db.Session.Repository().Remove(bson.M{"_id": "some_repository"})
	b := strings.NewReader(`{"name": "some_repository", "users": ["r2d2"]}`)
	recorder, request := post("/repository", b, c)
	NewRepository(recorder, request)
	got := readBody(recorder.Body, c)
	expected := "Repository \"some_repository\" successfully created\n"
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestNewRepositoryShouldSaveInDB(c *gocheck.C) {
	b := strings.NewReader(`{"name": "myRepository", "users": ["r2d2"]}`)
	recorder, request := post("/repository", b, c)
	NewRepository(recorder, request)
	collection := db.Session.Repository()
	defer collection.Remove(bson.M{"_id": "myRepository"})
	var p repository.Repository
	err := collection.Find(bson.M{"_id": "myRepository"}).One(&p)
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestNewRepositoryShouldSaveUserIdInRepository(c *gocheck.C) {
	b := strings.NewReader(`{"name": "myRepository", "users": ["r2d2", "brain"]}`)
	recorder, request := post("/repository", b, c)
	NewRepository(recorder, request)
	collection := db.Session.Repository()
	defer collection.Remove(bson.M{"_id": "myRepository"})
	var p repository.Repository
	err := collection.Find(bson.M{"_id": "myRepository"}).One(&p)
	c.Assert(err, gocheck.IsNil)
	c.Assert(len(p.Users), gocheck.Not(gocheck.Equals), 0)
}

func (s *S) TestNewRepositoryShouldReturnErrorWhenNoUserIsPassed(c *gocheck.C) {
	b := strings.NewReader(`{"name": "myRepository"}`)
	recorder, request := post("/repository", b, c)
	NewRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 400)
	body := readBody(recorder.Body, c)
	expected := "Validation Error: repository should have at least one user"
	got := strings.Replace(body, "\n", "", -1)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestNewRepositoryShouldReturnErrorWhenNoParametersArePassed(c *gocheck.C) {
	b := strings.NewReader("{}")
	recorder, request := post("/repository", b, c)
	NewRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 400)
	body := readBody(recorder.Body, c)
	expected := "Validation Error: repository name is not valid"
	got := strings.Replace(body, "\n", "", -1)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestParseBodyShouldMapBodyJsonToGivenStruct(c *gocheck.C) {
	var p repository.Repository
	b := bufferCloser{bytes.NewBufferString(`{"name": "Dummy Repository"}`)}
	err := parseBody(b, &p)
	c.Assert(err, gocheck.IsNil)
	expected := "Dummy Repository"
	c.Assert(p.Name, gocheck.Equals, expected)
}

func (s *S) TestParseBodyShouldReturnErrorWhenJsonIsInvalid(c *gocheck.C) {
	var p repository.Repository
	b := bufferCloser{bytes.NewBufferString("{]ja9aW}")}
	err := parseBody(b, &p)
	c.Assert(err, gocheck.NotNil)
	expected := "Could not parse json: invalid character ']' looking for beginning of object key string"
	c.Assert(err.Error(), gocheck.Equals, expected)
}

func (s *S) TestParseBodyShouldReturnErrorWhenBodyIsEmpty(c *gocheck.C) {
	var p repository.Repository
	b := bufferCloser{bytes.NewBufferString("")}
	err := parseBody(b, &p)
	c.Assert(err, gocheck.NotNil)
	c.Assert(err, gocheck.ErrorMatches, `^Could not parse json:.*$`)
}

func (s *S) TestParseBodyShouldReturnErrorWhenResultParamIsNotAPointer(c *gocheck.C) {
	var p repository.Repository
	b := bufferCloser{bytes.NewBufferString(`{"name": "something"}`)}
	err := parseBody(b, p)
	c.Assert(err, gocheck.NotNil)
	expected := "parseBody function cannot deal with struct. Use pointer"
	c.Assert(err.Error(), gocheck.Equals, expected)
}

func (s *S) TestNewRepositoryShouldReturnErrorWhenBodyIsEmpty(c *gocheck.C) {
	b := strings.NewReader("")
	recorder, request := post("/repository", b, c)
	NewRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 400)
}

func (s *S) TestGrantAccessUpdatesReposDocument(c *gocheck.C) {
	u, err := user.New("pippin", map[string]string{})
	defer db.Session.User().Remove(bson.M{"_id": "pippin"})
	c.Assert(err, gocheck.IsNil)
	r := repository.Repository{Name: "onerepo"}
	err = db.Session.Repository().Insert(&r)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().Remove(bson.M{"_id": r.Name})
	r2 := repository.Repository{Name: "otherepo"}
	err = db.Session.Repository().Insert(&r2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().Remove(bson.M{"_id": r2.Name})
	b := bytes.NewBufferString(fmt.Sprintf(`{"repositories": ["%s", "%s"], "users": ["%s"]}`, r.Name, r2.Name, u.Name))
	rec, req := del("/repository/grant", b, c)
	GrantAccess(rec, req)
	var repos []repository.Repository
	err = db.Session.Repository().Find(bson.M{"_id": bson.M{"$in": []string{r.Name, r2.Name}}}).All(&repos)
	c.Assert(err, gocheck.IsNil)
	c.Assert(rec.Code, gocheck.Equals, 200)
	for _, repo := range repos {
		c.Assert(repo.Users, gocheck.DeepEquals, []string{u.Name})
	}
}

func (s *S) TestRevokeAccessUpdatesReposDocument(c *gocheck.C) {
	r := repository.Repository{Name: "onerepo", Users: []string{"Umi", "Luke"}}
	err := db.Session.Repository().Insert(&r)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().Remove(bson.M{"_id": r.Name})
	r2 := repository.Repository{Name: "otherepo", Users: []string{"Umi", "Luke"}}
	err = db.Session.Repository().Insert(&r2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Repository().Remove(bson.M{"_id": r2.Name})
	b := bytes.NewBufferString(fmt.Sprintf(`{"repositories": ["%s", "%s"], "users": ["Umi"]}`, r.Name, r2.Name))
	rec, req := del("/repository/revoke", b, c)
	RevokeAccess(rec, req)
	var repos []repository.Repository
	err = db.Session.Repository().Find(bson.M{"_id": bson.M{"$in": []string{r.Name, r2.Name}}}).All(&repos)
	c.Assert(err, gocheck.IsNil)
	for _, repo := range repos {
		c.Assert(repo.Users, gocheck.DeepEquals, []string{"Luke"})
	}
}

func (s *S) TestAddKey(c *gocheck.C) {
	usr, err := user.New("Frodo", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(usr.Name)
	b := strings.NewReader(fmt.Sprintf(`{"keyname": %q}`, rawKey))
	recorder, request := post(fmt.Sprintf("/user/%s/key?:name=%s", usr.Name, usr.Name), b, c)
	AddKey(recorder, request)
	got := readBody(recorder.Body, c)
	expected := "Key(s) successfully created"
	c.Assert(got, gocheck.Equals, expected)
	c.Assert(recorder.Code, gocheck.Equals, 200)
	var k user.Key
	err = db.Session.Key().Find(bson.M{"name": "keyname", "username": usr.Name}).One(&k)
	c.Assert(err, gocheck.IsNil)
	c.Assert(k.Body, gocheck.Equals, keyBody)
	c.Assert(k.Comment, gocheck.Equals, keyComment)
}

func (s *S) TestAddKeyShouldReturnErrorWhenUserDoesNotExists(c *gocheck.C) {
	b := strings.NewReader(`{"key": "a public key"}`)
	recorder, request := post("/user/Frodo/key?:name=Frodo", b, c)
	AddKey(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 404)
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	c.Assert(string(body), gocheck.Equals, "User not found\n")
}

func (s *S) TestAddKeyShouldReturnProperStatusCodeWhenKeyAlreadyExists(c *gocheck.C) {
	usr, err := user.New("Frodo", map[string]string{"keyname": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(usr.Name)
	b := strings.NewReader(fmt.Sprintf(`{"keyname": %q}`, rawKey))
	recorder, request := post(fmt.Sprintf("/user/%s/key?:name=%s", usr.Name, usr.Name), b, c)
	AddKey(recorder, request)
	got := readBody(recorder.Body, c)
	expected := "Key already exists.\n"
	c.Assert(got, gocheck.Equals, expected)
	c.Assert(recorder.Code, gocheck.Equals, http.StatusConflict)
}

func (s *S) TestAddKeyShouldNotAcceptRepeatedKeysForDifferentUsers(c *gocheck.C) {
	usr, err := user.New("Frodo", map[string]string{"keyname": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(usr.Name)
	usr2, err := user.New("tempo", nil)
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(usr2.Name)
	b := strings.NewReader(fmt.Sprintf(`{"keyname": %q}`, rawKey))
	recorder, request := post(fmt.Sprintf("/user/%s/key?:name=%s", usr2.Name, usr2.Name), b, c)
	AddKey(recorder, request)
	got := readBody(recorder.Body, c)
	expected := "Key already exists.\n"
	c.Assert(got, gocheck.Equals, expected)
	c.Assert(recorder.Code, gocheck.Equals, http.StatusConflict)
}

func (s *S) TestAddKeyInvalidKey(c *gocheck.C) {
	u := user.User{Name: "Frodo"}
	err := db.Session.User().Insert(&u)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(bson.M{"_id": "Frodo"})
	b := strings.NewReader(`{"keyname":"invalid-rsa"}`)
	recorder, request := post(fmt.Sprintf("/user/%s/key?:name=%s", u.Name, u.Name), b, c)
	AddKey(recorder, request)
	got := readBody(recorder.Body, c)
	expected := "Invalid key\n"
	c.Assert(got, gocheck.Equals, expected)
	c.Assert(recorder.Code, gocheck.Equals, http.StatusBadRequest)
}

func (s *S) TestAddKeyShouldRequireKey(c *gocheck.C) {
	u := user.User{Name: "Frodo"}
	err := db.Session.User().Insert(&u)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().Remove(bson.M{"_id": "Frodo"})
	b := strings.NewReader(`{}`)
	recorder, request := post("/user/Frodo/key?:name=Frodo", b, c)
	AddKey(recorder, request)
	body := readBody(recorder.Body, c)
	expected := "A key is needed"
	got := strings.Replace(body, "\n", "", -1)
	c.Assert(got, gocheck.Equals, expected)
}

func (s *S) TestAddKeyShouldWriteKeyInAuthorizedKeysFile(c *gocheck.C) {
	u := user.User{Name: "Frodo"}
	err := db.Session.User().Insert(&u)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId("Frodo")
	b := strings.NewReader(fmt.Sprintf(`{"key": "%s"}`, rawKey))
	recorder, request := post("/user/Frodo/key?:name=Frodo", b, c)
	AddKey(recorder, request)
	defer db.Session.Key().Remove(bson.M{"name": "key", "username": u.Name})
	c.Assert(recorder.Code, gocheck.Equals, 200)
	content := s.authKeysContent(c)
	c.Assert(strings.HasSuffix(strings.TrimSpace(content), rawKey), gocheck.Equals, true)
}

func (s *S) TestRemoveKeyGivesExpectedSuccessResponse(c *gocheck.C) {
	u, err := user.New("Gandalf", map[string]string{"keyname": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(u.Name)
	url := "/user/Gandalf/key/keyname?:keyname=keyname&:name=Gandalf"
	recorder, request := del(url, nil, c)
	RemoveKey(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 200)
	b := readBody(recorder.Body, c)
	c.Assert(b, gocheck.Equals, `Key "keyname" successfully removed`)
}

func (s *S) TestRemoveKeyRemovesKeyFromDatabase(c *gocheck.C) {
	u, err := user.New("Gandalf", map[string]string{"keyname": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(u.Name)
	url := "/user/Gandalf/key/keyname?:keyname=keyname&:name=Gandalf"
	recorder, request := del(url, nil, c)
	RemoveKey(recorder, request)
	count, err := db.Session.Key().Find(bson.M{"name": "keyname", "username": "Gandalf"}).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(count, gocheck.Equals, 0)
}

func (s *S) TestRemoveKeyShouldRemoveKeyFromAuthorizedKeysFile(c *gocheck.C) {
	u, err := user.New("Gandalf", map[string]string{"keyname": rawKey})
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(u.Name)
	url := "/user/Gandalf/key/keyname?:keyname=keyname&:name=Gandalf"
	recorder, request := del(url, nil, c)
	RemoveKey(recorder, request)
	content := s.authKeysContent(c)
	c.Assert(content, gocheck.Equals, "")
}

func (s *S) TestRemoveKeyShouldReturnErrorWithLineBreakAtEnd(c *gocheck.C) {
	url := "/user/idiocracy/key/keyname?:keyname=keyname&:name=idiocracy"
	recorder, request := del(url, nil, c)
	RemoveKey(recorder, request)
	b := readBody(recorder.Body, c)
	c.Assert(b, gocheck.Equals, "User not found\n")
}

func (s *S) TestListKeysGivesExpectedSuccessResponse(c *gocheck.C) {
	keys := map[string]string{"key1": rawKey, "key2": otherKey}
	u, err := user.New("Gandalf", keys)
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(u.Name)
	url := "/user/Gandalf/keys?:name=Gandalf"
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	ListKeys(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 200)
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	var data map[string]string
	err = json.Unmarshal(body, &data)
	c.Assert(err, gocheck.IsNil)
	c.Assert(data, gocheck.DeepEquals, keys)
}

func (s *S) TestListKeysWithoutKeysGivesEmptyJSON(c *gocheck.C) {
	u, err := user.New("Gandalf", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	defer user.Remove(u.Name)
	url := "/user/Gandalf/keys?:name=Gandalf"
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	ListKeys(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 200)
	b := readBody(recorder.Body, c)
	c.Assert(b, gocheck.Equals, "{}")
}

func (s *S) TestListKeysWithInvalidUserReturnsNotFound(c *gocheck.C) {
	url := "/user/no-Gandalf/keys?:name=no-Gandalf"
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	ListKeys(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 404)
	b := readBody(recorder.Body, c)
	c.Assert(b, gocheck.Equals, "User not found\n")
}

func (s *S) TestRemoveUser(c *gocheck.C) {
	u, err := user.New("username", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("/user/%s/?:name=%s", u.Name, u.Name)
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RemoveUser(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 200)
	b, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	c.Assert(string(b), gocheck.Equals, "User \"username\" successfully removed\n")
}

func (s *S) TestRemoveUserShouldRemoveFromDB(c *gocheck.C) {
	u, err := user.New("anuser", map[string]string{})
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("/user/%s/?:name=%s", u.Name, u.Name)
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RemoveUser(recorder, request)
	collection := db.Session.User()
	lenght, err := collection.Find(bson.M{"_id": u.Name}).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)
}

func (s *S) TestRemoveRepository(c *gocheck.C) {
	r, err := repository.New("myRepo", []string{"pippin"}, true)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("repository/%s/?:name=%s", r.Name, r.Name)
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RemoveRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 200)
	b, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	c.Assert(string(b), gocheck.Equals, "Repository \"myRepo\" successfully removed\n")
}

func (s *S) TestRemoveRepositoryShouldRemoveFromDB(c *gocheck.C) {
	r, err := repository.New("myRepo", []string{"pippin"}, true)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("repository/%s/?:name=%s", r.Name, r.Name)
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RemoveRepository(recorder, request)
	err = db.Session.Repository().Find(bson.M{"_id": r.Name}).One(&r)
	c.Assert(err, gocheck.ErrorMatches, "^not found$")
}

func (s *S) TestRemoveRepositoryShouldReturn400OnFailure(c *gocheck.C) {
	url := fmt.Sprintf("repository/%s/?:name=%s", "foo", "foo")
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RemoveRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, 400)
}

func (s *S) TestRemoveRepositoryShouldReturnErrorMsgWhenRepoDoesNotExists(c *gocheck.C) {
	url := fmt.Sprintf("repository/%s/?:name=%s", "foo", "foo")
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RemoveRepository(recorder, request)
	b, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	c.Assert(string(b), gocheck.Equals, "Could not remove repository: not found\n")
}

func (s *S) TestRenameRepository(c *gocheck.C) {
	r, err := repository.New("raising", []string{"guardian@what.com"}, true)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("/repository/%s/?:name=%s", r.Name, r.Name)
	body := strings.NewReader(`{"name":"freedom"}`)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RenameRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, http.StatusOK)
	_, err = repository.Get("raising")
	c.Assert(err, gocheck.NotNil)
	r.Name = "freedom"
	repo, err := repository.Get("freedom")
	c.Assert(err, gocheck.IsNil)
	c.Assert(repo, gocheck.DeepEquals, *r)
}

func (s *S) TestRenameRepositoryInvalidJSON(c *gocheck.C) {
	url := "/repository/foo/?:name=foo"
	body := strings.NewReader(`{"name""`)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RenameRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, http.StatusBadRequest)
}

func (s *S) TestRenameRepositoryNotfound(c *gocheck.C) {
	url := "/repository/foo/?:name=foo"
	body := strings.NewReader(`{"name":"freedom"}`)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	recorder := httptest.NewRecorder()
	RenameRepository(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, http.StatusNotFound)
}

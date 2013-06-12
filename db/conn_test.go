// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"github.com/globocom/config"
	"labix.org/v2/mgo"
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "gandalf_tests")
	Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	Session.DB.DropDatabase()
}

func (s *S) TestSessionRepositoryShouldReturnAMongoCollection(c *gocheck.C) {
	var rep *mgo.Collection
	rep = Session.Repository()
	cRep := Session.DB.C("repository")
	c.Assert(rep, gocheck.DeepEquals, cRep)
}

func (s *S) TestSessionUserShouldReturnAMongoCollection(c *gocheck.C) {
	var usr *mgo.Collection
	usr = Session.User()
	cUsr := Session.DB.C("user")
	c.Assert(usr, gocheck.DeepEquals, cUsr)
}

func (s *S) TestSessionKeyShouldReturnKeyCollection(c *gocheck.C) {
	key := Session.Key()
	cKey := Session.DB.C("key")
	c.Assert(key, gocheck.DeepEquals, cKey)
}

func (s *S) TestSessionKeyBodyIsUnique(c *gocheck.C) {
	key := Session.Key()
	indexes, err := key.Indexes()
	c.Assert(err, gocheck.IsNil)
	c.Assert(indexes, gocheck.HasLen, 2)
	c.Assert(indexes[1].Key, gocheck.DeepEquals, []string{"body"})
	c.Assert(indexes[1].Unique, gocheck.DeepEquals, true)
}

func (s *S) TestConnect(c *gocheck.C) {
	Connect()
	c.Assert(Session.DB.Name, gocheck.Equals, "gandalf_tests")
	err := Session.DB.Session.Ping()
	c.Assert(err, gocheck.IsNil)
}

// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/user"
	"labix.org/v2/mgo/bson"
)

func getUserOr404(name string) (user.User, error) {
	var u user.User
	if err := db.Session.User().Find(bson.M{"_id": name}).One(&u); err != nil && err.Error() == "not found" {
		return u, fmt.Errorf("User %s not found", name)
	}
	return u, nil
}

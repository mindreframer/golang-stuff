/*Create app*/
package hamster

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

var aName = "apps"

//The App type
type App struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"-"`
	ParentId  bson.ObjectId `bson:"parentId" json:"parentId"`
	Name      string        `bson:"name" json:"name"`
	OS        string        `bson:"os" json:"os"`
	ApiToken  string        `bson:"apitoken" json:"apitoken"`
	ApiSecret string        `bson:"apisecret" json:"apisecret"`
	Hash      string        `bson:"hash" json:"hash"`
	Salt      string        `bson:"salt" json:"salt"`
	Created   time.Time     `bson:"created" json:"created"`
	Updated   time.Time     `bson:"updated" json:"updated"`
	Objects   []string      `bson:"objects" json:"objects"`
}

//POST: "/api/v1/developers/:developerId/apps/" handler
func (s *Server) CreateApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateApp: ")

	//get params
	did := r.URL.Query().Get(":developerId")
	if did == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+did)
	}
	developer_id := decodeToken(did)
	if developer_id == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+developer_id)
	}

	//get collection developer
	session := s.db.GetSession()
	defer session.Close()
	d := session.DB("").C(dName)

	//if developer id exists

	if err := d.FindId(bson.ObjectIdHex(developer_id)).Limit(1).One(nil); err != nil {
		s.notFound(r, w, err, developer_id+" : id not found")
		return
	}

	//parse body
	app := &App{}
	if err := s.readJson(app, r, w); err != nil {
		s.badRequest(r, w, err, "malformed app json")
		return

	}

	//set fields
	app.Id = bson.NewObjectId() //todo:make it shorter and user friendly
	app.ParentId = bson.ObjectIdHex(developer_id)
	app.Created = time.Now()
	app.Updated = time.Now()
	app.ApiToken = encodeBase64Token(app.Id.Hex())
	secret, err := genUUID(16)
	if err != nil {
		s.internalError(r, w, err, "gen uuid")
		return
	}
	app.ApiSecret = encodeBase64Token(secret)
	hash, salt, err := encryptPassword(secret)
	if err != nil {
		s.internalError(r, w, err, "encypt secret")
		return

	}
	app.Hash = hash
	app.Salt = salt

	//get apps collection
	c := session.DB("").C(aName)

	//then insert app and respond
	if insert_err := c.Insert(app); insert_err != nil {

		s.internalError(r, w, insert_err, "error inserting: "+fmt.Sprintf("%v", app))

	} else {
		response := AppResponse{ApiToken: app.ApiToken, ApiSecret: app.ApiSecret, Name: app.Name, OS: app.OS}
		s.logger.Printf("created new app: %+v, id: %v\n", response)
		s.serveJson(w, &response)
	}

}

//GET "/api/v1/developers/apps/:objectId"
func (s *Server) QueryApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryApp: ")

	//getObjectId
	object_id := s.getObjectId(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//find and serve data
	app := App{}
	if err := c.FindId(bson.ObjectIdHex(object_id)).Limit(1).One(&app); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//respond
	response := AppResponse{ApiToken: app.ApiToken, ApiSecret: app.ApiSecret, Name: app.Name, OS: app.OS}
	//s.logger.Printf("query app: %+v, id: %v\n", response)
	s.serveJson(w, &response)

}

//GET "/api/v1/developers/:developerId/apps/"
func (s *Server) QueryAllApps(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryAllApps: ")

	//get params
	did := r.URL.Query().Get(":developerId")
	if did == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+did)
	}
	developer_id := decodeToken(did)
	if developer_id == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+developer_id)
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//find apps
	var apps []App
	iter := c.Find(bson.M{"parentId": bson.ObjectIdHex(developer_id)}).Iter()
	err := iter.All(&apps)
	if err != nil {
		s.internalError(r, w, err, "error iterating app documents")
	}

	//respond
	var re []AppResponse
	for _, app := range apps {

		re = append(re, AppResponse{ApiToken: app.ApiToken, ApiSecret: app.ApiSecret, Name: app.Name, OS: app.OS})

	}

	s.serveJson(w, &re)
}

//PUT "/api/v1/developers/apps/:objectId"
func (s *Server) UpdateApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("UpdateApp: ")

	//getObjectId
	object_id := s.getObjectId(w, r)

	//parse body
	updateRequest := &UpdateAppRequest{}
	if err := s.readJson(updateRequest, r, w); err != nil {
		s.badRequest(r, w, err, "malformed update request body")
		return
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
				"name":    updateRequest.Name,
				"os":      updateRequest.OS,
			}}}

	//find and update
	app := App{}
	if _, err := c.FindId(bson.ObjectIdHex(object_id)).Apply(change, &app); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//respond
	response := AppResponse{ApiToken: app.ApiToken, ApiSecret: app.ApiSecret, Name: app.Name, OS: app.OS}
	s.serveJson(w, &response)

}

//DELETE "/api/v1/developers/apps/:objectId"
func (s *Server) DeleteApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeleteApp: ")

	//getObjectId
	object_id := s.getObjectId(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//delete
	if err := c.RemoveId(bson.ObjectIdHex(object_id)); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//respond
	response := DeleteResponse{Status: "ok"}
	s.serveJson(w, &response)

}

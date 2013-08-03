/*Manage Developer*/
package hamster

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"

	"time"
)

var (
	dName = "developers"
)

//Stores developer account info
type Developer struct {
	Id       bson.ObjectId `bson:"_id" json:"id"`
	ParentId string        `bson:"parentId" json:"parentId"` //unused. change string to bson.ObjectId
	Name     string        `bson:"name" json:"name"`
	Email    string        `bson:"email" json:"email"`
	Verified bool          `bson:"verified" json:"verified"`
	Password string        `json:"password"` //only used for parsing incoming json
	Hash     string        `bson:"hash"`
	Salt     string        `bson:"salt"`
	Created  time.Time     `bson:"created" json:"created"`
	Updated  time.Time     `bson:"updated" json:"updated"`
	UrlToken string        `bson:"urltoken" json:"urltoken"`
}

//pre-index developers collection on startup. probably should do it through command line?
func (s *Server) IndexDevelopers() {
	s.logger.SetPrefix("IndexDev: ")
	//ensure email exists and is unique
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//ensure index key exists and is unique

	if err := c.EnsureIndex(index); err != nil {

		s.logger.Printf("failed indexing developers, err: %v \n", err)
	} else {
		s.logger.Printf("developers collection indexed!\n")
	}

}

//POST: /api/v1/developers/ handler
func (s *Server) CreateDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateDev: ")
	//get the request body
	developer := &Developer{}

	if err := s.readJson(developer, r, w); err != nil {
		s.badRequest(r, w, err, "malformed developer json")
		return

	}

	//check if email is not empty
	if developer.Email == "" {
		s.internalError(r, w, nil, "empty email ")
		return
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//set fields
	developer.Id = bson.NewObjectId() //todo:make it shorter and user friendly
	developer.UrlToken = encodeBase64Token(developer.Id.Hex())
	//encrypt password
	hash, salt, err := encryptPassword(developer.Password)
	if err != nil {
		s.internalError(r, w, err, "encypt password")
		return

	}
	developer.Hash = hash
	developer.Salt = salt
	developer.Created = time.Now()
	developer.Updated = time.Now()

	//insert new document

	if insert_err := c.Insert(developer); insert_err != nil {

		s.internalError(r, w, insert_err, "error inserting: "+fmt.Sprintf("%v", developer))

	} else {
		s.logger.Printf("created new developer: %+v, id: %v\n", developer)
		//serve created developer json
		access_token, err := s.genAccessToken(developer.Email)
		if err != nil {
			s.internalError(r, w, err, "error generating access token")
		}
		response := NewDeveloperResponse{ObjectId: developer.UrlToken, AccessToken: access_token}
		s.serveJson(w, &response)
	}

	return

}

//POST: /api/v1/developers/login/ handler
func (s *Server) LoginDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("LoginDev: ")

	//get email password from request
	email, password := getUserPassword(r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	developer := Developer{}

	//find and login developer
	s.logger.Printf("find developer: %s %s", email, password)
	if email != "" && password != "" {
		if findErr := c.Find(bson.M{"email": email}).One(&developer); findErr != nil {
			s.notFound(r, w, findErr, email+" user not found")
			return
		}
		//match password and set session
		if matchPassword(password, developer.Hash, developer.Salt) {
			access_token, err := s.genAccessToken(developer.Email)
			if err != nil {
				s.internalError(r, w, err, email+" generate access token")
			}
			//respond with developer profile
			response := LoginResponse{ObjectId: developer.UrlToken, AccessToken: access_token, Status: "ok"}
			s.serveJson(w, &response)

		} else {

			s.notFound(r, w, nil, email+" password match failed")
		}

	} else {
		s.notFound(r, w, nil, "email empty")
	}

}

//POST:/api/v1/developers/logout/ handler
func (s *Server) LogoutDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("LogoutDev: ")

	//parse body
	logoutRequest := &LogoutRequest{}

	if err := s.readJson(logoutRequest, r, w); err != nil {
		s.badRequest(r, w, err, "malformed logout request")
		return
	}

	//logout
	if logout_err := s.logout(logoutRequest.Email); logout_err != nil {
		s.internalError(r, w, logout_err, logoutRequest.Email+" : could not logout")
	}

	//response
	response := LogoutResponse{Status: "ok"}
	s.serveJson(w, &response)

}

//GET:/api/v1/developers/:objectId handler
func (s *Server) QueryDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryDev: ")

	//getObjectId
	object_id := s.getObjectId(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//find and serve data
	developer := Developer{}
	if err := c.FindId(bson.ObjectIdHex(object_id)).Limit(1).One(&developer); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//respond
	response := QueryDevResponse{Name: developer.Name, Email: developer.Email}
	s.serveJson(w, &response)

}

//PUT:/api/v1/developers/:objectId handler
func (s *Server) UpdateDev(w http.ResponseWriter, r *http.Request) {

	s.logger.SetPrefix("UpdateDev: ")

	//getObjectId
	object_id := s.getObjectId(w, r)

	//parse body
	updateRequest := &UpdateRequest{}
	if err := s.readJson(updateRequest, r, w); err != nil {
		s.badRequest(r, w, err, "malformed update request body")
		return
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
				"name":    updateRequest.Name,
			}}}

	//find and update
	developer := Developer{}
	if _, err := c.FindId(bson.ObjectIdHex(object_id)).Apply(change, &developer); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//respond
	response := QueryDevResponse{Name: developer.Name, Email: developer.Email}
	s.serveJson(w, &response)

}

//DELETE:/api/v1/developers/:objectId
func (s *Server) DeleteDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeleteDev: ")

	//getObjectId
	object_id := s.getObjectId(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//delete
	if err := c.RemoveId(bson.ObjectIdHex(object_id)); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//respond
	response := DeleteResponse{Status: "ok"}
	s.serveJson(w, &response)

}

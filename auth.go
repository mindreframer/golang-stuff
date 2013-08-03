/*Utility functions and methods for authoriztion*/
package hamster

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	//"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/kr/fernet"
	"hash"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"sync"
	"time"
)

//Allow only verified ip's from config
func (s *Server) baseAuth(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")

	s.logger.SetPrefix("BaseAuth:")
	s.logger.Printf("%v  %v", r.Method, r.URL.Path)

	if !s.ipAllowed(ip[0]) {

		http.Error(w, "Unauthorized:ip not allowed", http.StatusUnauthorized)
		return

	}

}

//Allow ip
func (s *Server) ipAllowed(ip string) bool {
	for _, client := range s.config.Clients {
		if ip == client.Ip {
			return true
		}
	}
	return false
}

//Authenticate developer. Access Token is generated using a shared secret between
//Hamster and the client. The shared secret is manually configured in hamster.toml.
//TODO: find a better implementation
func (s *Server) developerAuth(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeveloperAuth:")

	//skip access token check if logging in
	if r.Method == "POST" && r.URL.Path == "/api/v1/developers/login/" {

		return
	}

	access_token := r.Header.Get("X-Access-Token")

	if access_token == "" {

		s.unauthorized(r, w, errors.New("token is empty"), "access token invalid")
		return

	}

	//check POST /api/v1/developers
	if r.Method == "POST" && r.URL.Path == "/api/v1/developers/" {

		//check auth
		if !s.validateSharedToken(access_token) {
			s.unauthorized(r, w, errors.New("shared token is old"), "access token invalid")
		}
		return
	}

	if _, ok := s.validateAccessLoginToken(access_token); !ok {
		s.unauthorized(r, w, errors.New("token is old"), "access token invalid")
		return

	}
}

//Validate access token
func (s *Server) validateAccessLoginToken(token string) (string, bool) {

	btok, err := base64.URLEncoding.DecodeString(token)

	if err != nil {
		return "", false

	}
	k := fernet.MustDecodeKeys(s.config.Clients["browser"].Secret)
	email := fernet.VerifyAndDecrypt(btok, 60*time.Second, k)

	c := s.redisConn()
	defer c.Close()

	status, err := redis.String(c.Do("GET", email))
	if err != nil {
		return "", false
	}

	if status == "loggedin" {
		return string(email), true
	} else {
		return "", false
	}

}

//Validate access token
func (s *Server) validateSharedToken(token string) bool {

	btok, err := base64.URLEncoding.DecodeString(token)

	if err != nil {
		return false

	}
	k := fernet.MustDecodeKeys(s.config.Clients["browser"].Secret)
	shared_token := fernet.VerifyAndDecrypt(btok, 60*10*time.Second, k)
	if string(shared_token) == string(s.config.Clients["browser"].Token) {
		return true
	} else {
		return false
	}

}

//Authenticates object level requests
func (s *Server) objectAuth(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("ObjectAuth:")
	api_token := r.Header.Get("X-Api-Token")
	api_secret := r.Header.Get("X-Api-Secret")

	if api_token == "" || api_secret == "" {
		s.unauthorized(r, w, errors.New("token or secret invalid"), "access token invalid")
		return
	}

	if !s.validateApiToken(api_token, api_secret) {
		s.unauthorized(r, w, errors.New("token match failed"), "access token failed")
		return
	}

}

func (s *Server) validateApiToken(token string, secret string) bool {

	if ok, hash, salt := s.getHashSalt(token); ok {
		//fmt.Println("found key in redis")
		if matchPassword(decodeToken(secret), hash, salt) {
			return true
		} else {
			return false
		}

	} else {

		if ok, hash, salt := s.getHashSaltFromDb(token); ok {
			if matchPassword(decodeToken(secret), hash, salt) {

				return true
			} else {
				//fmt.Println("api secret match failed!")
				return false
			}

		} else {
			return false
		}

	}

}

func (s *Server) getHashSalt(token string) (bool, string, string) {
	c := s.redisConn()
	defer c.Close()

	hash, err := redis.String(c.Do("GET", token+":hash"))
	if err != nil {
		return false, "", ""
	}
	salt, err := redis.String(c.Do("GET", token+":salt"))
	if err != nil {
		return false, "", ""
	}

	return true, hash, salt
}

func (s *Server) getHashSaltFromDb(token string) (bool, string, string) {
	//get collection

	app_id := decodeToken(token)
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	app := App{}
	//TODO:select fields
	if err := c.FindId(bson.ObjectIdHex(app_id)).One(&app); err != nil {
		//fmt.Println("app not found\n")
		return false, "", ""
	}

	//cache hash-salt in redis
	rc := s.redisConn()
	defer rc.Close()

	rc.Do("SET", token+":hash", app.Hash)
	rc.Do("SET", token+":salt", app.Salt)

	return true, app.Hash, app.Salt

}

//Get Basic user password
func getUserPassword(r *http.Request) (string, string) {

	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", ""
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", ""
	}

	return pair[0], pair[1]

}

type Hmac struct {
	hashFunc hash.Hash
	m        sync.Mutex
}

/*Generates hash*/
func (h *Hmac) generateHash(data []byte) []byte {
	h.m.Lock()
	defer h.m.Unlock()

	h.hashFunc.Reset()
	h.hashFunc.Write(data)
	return h.hashFunc.Sum(nil)

}

/*Encrypts the password and returns password hash*/
func (h *Hmac) encrypt(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(h.generateHash(password), cost)
}

func (h *Hmac) compare(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, h.generateHash(password))
}

/*Returns a hmac type*/
func New(hash func() hash.Hash, salt []byte) *Hmac {
	hm := &Hmac{
		hashFunc: hmac.New(hash, salt),
		m:        sync.Mutex{},
	}
	return hm
}

/*Encrypts password. Returns hash+salt*/
func encryptPassword(password string) (string, string, error) {

	salt, err0 := genUUID(16)
	if err0 != nil {
		return "", "", err0
	}

	hm := New(sha512.New, []byte(salt))
	pass := []byte(password)
	encrypted, err := hm.encrypt(pass, bcrypt.DefaultCost)

	if err != nil {

		return "", "", err
	}

	return string(encrypted), salt, nil
}

/*Match encrypted string*/
func matchPassword(password string, hash string, salt string) bool {
	p := []byte(password)
	h := []byte(hash)
	s := []byte(salt)
	hm := New(sha512.New, s)
	err := hm.compare(h, p)
	if err != nil {
		return false
	} else {
		//matched!
		return true
	}

}

/*Encode to base64*/
func encodeBase64Token(hexVal string) string {

	token := base64.URLEncoding.EncodeToString([]byte(hexVal))

	return token

}

/*Decode from base64*/
func decodeToken(token string) string {

	hexVal, err := base64.URLEncoding.DecodeString(token)
	if err != nil {

		return ""

	}

	return string(hexVal)

}

/*Generate uuid*/
func genUUID(size int) (string, error) {
	uuid := make([]byte, size)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	uuid[8] = 0x80
	uuid[4] = 0x40

	return hex.EncodeToString(uuid), nil
}

/*
 Return a random 16-byte base64 alphabet string
*/
func randomKey() string {
	k := make([]byte, 12)
	for bytes := 0; bytes < len(k); {
		n, err := rand.Read(k[bytes:])
		if err != nil {
			panic("rand.Read() failed")
		}
		bytes += n
	}
	return base64.StdEncoding.EncodeToString(k)
}

//Generate time based access token using shared secret. See fernet project
//for more details
func (s *Server) genAccessToken(email string) (string, error) {
	//encrypt token
	k := fernet.MustDecodeKeys(s.config.Clients["browser"].Secret)
	tok, err := fernet.EncryptAndSign([]byte(email), k[0])
	if err != nil {

	}
	token := base64.URLEncoding.EncodeToString(tok)

	//cache it
	c := s.redisConn()
	defer c.Close()

	c.Do("SET", email, "loggedin")

	return token, nil

}

//Logout developer
//TODO: find a better way to handle login/logout
func (s *Server) logout(email string) error {

	c := s.redisConn()
	defer c.Close()

	status, err := redis.String(c.Do("GET", email))
	if err != nil {
		return err
	}

	if status == "loggedin" {
		c.Do("SET", email, "") //log out
	} else {
		return errors.New("user is not logged in")
	}

	return nil

}

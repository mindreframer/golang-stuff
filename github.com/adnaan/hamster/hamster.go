/*
The Hamster Server. The Server type holds instances of all the components,
*effectively making it possible to collapse all the code into one file. The separation
* of code is only for readability. To use it as a package simply:
* import ("github.com/adnaan/hamster")
* server := hamster.NewServer()
* //server.Quiet()//disable logging
* server.ListenAndServe()
* Also change hamster.toml for custom configuration.
* TODO: Pass hamster.toml as argument to the server
* TODO: make handler methods local, model method exported for pkg/rpc support
*/
package hamster

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/adnaan/routes"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

//The server type holds instances of all components
type Server struct {
	listener   net.Listener
	logger     *log.Logger
	httpServer *http.Server
	route      *routes.RouteMux
	db         *Db
	config     *Config
	cookie     *sessions.CookieStore //unused
	redisConn  func() redis.Conn
}

//dbUrl:"mongodb://adnaan:pass@localhost:27017/hamster"
//serverUrl:fmt.Sprintf("%s:%d", address, port)
//creates a new server, setups logging etc.
func NewServer() *Server {
	f, err := os.OpenFile("hamster.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("hamster.log faied to open")

	}
	//log.SetOutput(f)
	//log.SetOutput(os.Stdout)
	//router
	r := routes.New()
	//toml config
	var cfg Config
	if _, err := toml.DecodeFile("hamster.toml", &cfg); err != nil {
		fmt.Println(err)
		return nil
	}
	//cookie store
	ck := sessions.NewCookieStore([]byte(cfg.Servers["local"].CookieSecret))

	//redis
	var getRedis = func() redis.Conn {

		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			panic(err)
		}

		return c

	}

	//initialize server
	s := &Server{
		httpServer: &http.Server{Addr: fmt.Sprintf(":%d", cfg.Servers["local"].Port), Handler: r},
		route:      r,
		logger:     log.New(f, "", log.LstdFlags),
		db:         &Db{Url: cfg.DB["mongo"].Host},
		config:     &cfg,
		cookie:     ck,
		redisConn:  getRedis,
	}

	s.logger.SetFlags(log.Lshortfile)
	s.addHandlers()

	return s

}

//listen and serve a fastcgi server

func (s *Server) ListenAndServe() error {

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		s.logger.Printf("error listening: %v \n", err)
		return err
	}
	s.listener = listener

	go s.httpServer.Serve(s.listener)

	s.logger.Print("********Server Startup*********\n")
	s.logger.Print("********++++++++++++++*********\n")
	s.logger.Printf("hamster is now listening on http://localhost%s\n", s.httpServer.Addr)

	//index the collections
	s.IndexDevelopers()

	return nil
}

// stops the server.
func (s *Server) Shutdown() error {

	if s.listener != nil {
		// Then stop the server.
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}

	return nil
}

// no log
func (s *Server) Quiet() {
	s.logger = log.New(ioutil.Discard, "", log.LstdFlags)
}

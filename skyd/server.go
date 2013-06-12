package skyd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmhodges/levigo"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"regexp"
	"runtime"
	"time"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Server is the front end that controls access to tables.
type Server struct {
	httpServer      *http.Server
	router          *mux.Router
	logger          *log.Logger
	path            string
	listener        net.Listener
	servlets        []*Servlet
	tables          map[string]*Table
	factors         *Factors
	shutdownChannel chan bool
}

//------------------------------------------------------------------------------
//
// Errors
//
//------------------------------------------------------------------------------

//--------------------------------------
// Alternate Content Type
//--------------------------------------

// Hackish way to return plain text for query codegen. Might fix later but
// it's rarely used.
type TextPlainContentTypeError struct {
}

func (e *TextPlainContentTypeError) Error() string {
	return ""
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewServer returns a new Server.
func NewServer(port uint, path string) *Server {
	r := mux.NewRouter()
	s := &Server{
		httpServer: &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r},
		router:     r,
		logger:     log.New(os.Stdout, "", log.LstdFlags),
		path:       path,
		tables:     make(map[string]*Table),
	}

	s.router.HandleFunc("/debug/pprof", pprof.Index)
	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	s.addHandlers()
	s.addTableHandlers()
	s.addPropertyHandlers()
	s.addEventHandlers()
	s.addQueryHandlers()

	return s
}

//------------------------------------------------------------------------------
//
// Properties
//
//------------------------------------------------------------------------------

// The root server path.
func (s *Server) Path() string {
	return s.path
}

// The path to the data directory.
func (s *Server) DataPath() string {
	return fmt.Sprintf("%v/data", s.path)
}

// The path to the table metadata directory.
func (s *Server) TablesPath() string {
	return fmt.Sprintf("%v/tables", s.path)
}

// Generates the path for a table attached to the server.
func (s *Server) TablePath(name string) string {
	return fmt.Sprintf("%v/%v", s.TablesPath(), name)
}

// The path to the factors database.
func (s *Server) FactorsPath() string {
	return fmt.Sprintf("%v/factors", s.path)
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Lifecycle
//--------------------------------------

// Runs the server.
func (s *Server) ListenAndServe(shutdownChannel chan bool) error {
	s.shutdownChannel = shutdownChannel

	err := s.open()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		s.close()
		return err
	}
	s.listener = listener
	go s.httpServer.Serve(s.listener)

	s.logger.Printf("Sky v%s is now listening on http://localhost%s\n", Version, s.httpServer.Addr)

	return nil
}

// Stops the server.
func (s *Server) Shutdown() error {
	// Close servlets.
	s.close()

	// Close socket.
	if s.listener != nil {
		// Then stop the server.
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}

	// Notify that the server is shutdown.
	if s.shutdownChannel != nil {
		s.shutdownChannel <- true
	}

	return nil
}

// Checks if the server is listening for new connections.
func (s *Server) Running() bool {
	return (s.listener != nil)
}

// Opens the data directory and servlets.
func (s *Server) open() error {
	s.close()

	// Setup the file system if it doesn't exist.
	err := s.createIfNotExists()
	if err != nil {
		return fmt.Errorf("skyd.Server: Unable to create server folders: %v", err)
	}

	// Open factors database.
	s.factors = NewFactors(s.FactorsPath())
	err = s.factors.Open()
	if err != nil {
		s.close()
		return err
	}

	// Create servlets from child directories with numeric names.
	infos, err := ioutil.ReadDir(s.DataPath())
	if err != nil {
		return err
	}
	for _, info := range infos {
		match, _ := regexp.MatchString("^\\d$", info.Name())
		if info.IsDir() && match {
			s.servlets = append(s.servlets, NewServlet(fmt.Sprintf("%s/%s", s.DataPath(), info.Name()), s.factors))
		}
	}

	// If none exist then build them based on the number of logical CPUs available.
	if len(s.servlets) == 0 {
		cpuCount := runtime.NumCPU()
		for i := 0; i < cpuCount; i++ {
			s.servlets = append(s.servlets, NewServlet(fmt.Sprintf("%s/%v", s.DataPath(), i), s.factors))
		}
	}

	// Open servlets.
	for _, servlet := range s.servlets {
		err = servlet.Open()
		if err != nil {
			s.close()
			return err
		}
	}

	return nil
}

// Closes the data directory and servlets.
func (s *Server) close() {
	// Close servlets.
	if s.servlets != nil {
		for _, servlet := range s.servlets {
			servlet.Close()
		}
		s.servlets = nil
	}

	// Close factors database.
	if s.factors != nil {
		s.factors.Close()
		s.factors = nil
	}
}

// Creates the appropriate directory structure if one does not exist.
func (s *Server) createIfNotExists() error {
	// Create root directory.
	err := os.MkdirAll(s.path, 0700)
	if err != nil {
		return err
	}

	// Create data directory and one directory for each servlet.
	err = os.MkdirAll(s.DataPath(), 0700)
	if err != nil {
		return err
	}

	// Create tables directory.
	err = os.MkdirAll(s.TablesPath(), 0700)
	if err != nil {
		return err
	}

	return nil
}

//--------------------------------------
// Logging
//--------------------------------------

// Silences the log.
func (s *Server) Silence() {
	s.logger = log.New(ioutil.Discard, "", log.LstdFlags)
}

//--------------------------------------
// Routing
//--------------------------------------

// Parses incoming JSON objects and converts outgoing responses to JSON.
func (s *Server) ApiHandleFunc(route string, handlerFunction func(http.ResponseWriter, *http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	wrappedFunction := func(w http.ResponseWriter, req *http.Request) {
		// warn("%s \"%s %s %s\"", req.RemoteAddr, req.Method, req.RequestURI, req.Proto)
		t0 := time.Now()

		var ret interface{}
		params, err := s.decodeParams(w, req)
		if err == nil {
			ret, err = handlerFunction(w, req, params)
		}

		// If we're returning plain text then just dump out what's returned.
		if _, ok := err.(*TextPlainContentTypeError); ok {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			if str, ok := ret.(string); ok {
				w.Write([]byte(str))
			}
			return
		}

		// If there is an error then replace the return value.
		if err != nil {
			ret = map[string]interface{}{"message": err.Error()}
		}

		// Write header status.
		w.Header().Set("Content-Type", "application/json")
		var status int
		if err == nil {
			status = http.StatusOK
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)

		// Write to access log.
		s.logger.Printf("%s \"%s %s %s\" %d %0.3f", req.RemoteAddr, req.Method, req.RequestURI, req.Proto, status, time.Since(t0).Seconds())
		if status != http.StatusOK {
			s.logger.Printf("ERROR %v", err)
		}

		// Encode the return value appropriately.
		if ret != nil {
			encoder := json.NewEncoder(w)
			err := encoder.Encode(ConvertToStringKeys(ret))
			if err != nil {
				fmt.Printf("skyd.Server: Encoding error: %v\n", err.Error())
				return
			}
		}
	}

	return s.router.HandleFunc(route, wrappedFunction)
}

// Decodes the body of the message into parameters.
func (s *Server) decodeParams(w http.ResponseWriter, req *http.Request) (map[string]interface{}, error) {
	// Parses body parameters.
	params := make(map[string]interface{})
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil && err != io.EOF {
		return nil, errors.New("Malformed json request.")
	}
	return params, nil
}

//--------------------------------------
// Servlet Management
//--------------------------------------

// Executes a function through a single-threaded servlet context.
func (s *Server) GetObjectContext(tableName string, objectId string) (*Table, *Servlet, error) {
	// Return an error if the table already exists.
	table, err := s.OpenTable(tableName)
	if err != nil {
		return nil, nil, err
	}

	// Determine servlet index.
	index, err := s.GetObjectServletIndex(table, objectId)
	if err != nil {
		return nil, nil, err
	}
	servlet := s.servlets[index]

	return table, servlet, nil
}

// Calculates a tablet index based on the object identifier even hash.
func (s *Server) GetObjectServletIndex(t *Table, objectId string) (uint32, error) {
	// Encode object identifier.
	encodedObjectId, err := t.EncodeObjectId(objectId)
	if err != nil {
		return 0, err
	}

	// Calculate the even bits of the FNV1a hash.
	h := fnv.New64a()
	h.Reset()
	h.Write(encodedObjectId)
	hashcode := h.Sum64()
	index := CondenseUint64Even(hashcode) % uint32(len(s.servlets))

	return index, nil
}

//--------------------------------------
// Table Management
//--------------------------------------

// Retrieves a table that has already been opened.
func (s *Server) GetTable(name string) *Table {
	return s.tables[name]
}

// Retrieves a list of all tables in the database but does not open them.
// Do not use these table references for anything but informational purposes!
func (s *Server) GetAllTables() ([]*Table, error) {
	// Create a table object for each directory in the tables path.
	infos, err := ioutil.ReadDir(s.TablesPath())
	if err != nil {
		return nil, err
	}

	tables := []*Table{}
	for _, info := range infos {
		if info.IsDir() {
			tables = append(tables, NewTable(info.Name(), s.TablePath(info.Name())))
		}
	}

	return tables, nil
}

// Opens a table and returns a reference to it.
func (s *Server) OpenTable(name string) (*Table, error) {
	// If table already exists then use it.
	table := s.GetTable(name)
	if table != nil {
		return table, nil
	}

	// Otherwise open it and save the reference.
	table = NewTable(name, s.TablePath(name))
	err := table.Open()
	if err != nil {
		table.Close()
		return nil, err
	}
	s.tables[name] = table

	return table, nil
}

// Deletes a table.
func (s *Server) DeleteTable(name string) error {
	// Return an error if the table doesn't exist.
	table := s.GetTable(name)
	if table == nil {
		table = NewTable(name, s.TablePath(name))
	}
	if !table.Exists() {
		return fmt.Errorf("Table does not exist: %s", name)
	}

	// Determine table prefix.
	prefix, err := TablePrefix(table.Name)
	if err != nil {
		return err
	}

	// Delete data from each servlet.
	for _, servlet := range s.servlets {
		servlet.Lock()
		defer servlet.Unlock()

		// Delete the data from disk.
		ro := levigo.NewReadOptions()
		defer ro.Close()
		wo := levigo.NewWriteOptions()
		defer wo.Close()
		iterator := servlet.db.NewIterator(ro)
		defer iterator.Close()

		iterator.Seek(prefix)
		for iterator = iterator; iterator.Valid(); iterator.Next() {
			key := iterator.Key()
			if bytes.HasPrefix(key, prefix) {
				err := servlet.db.Delete(wo, key)
				if err != nil {
					return err
				}
			} else {
				break
			}
		}
	}

	// Remove the table from the lookup and remove it's schema.
	delete(s.tables, name)
	return table.Delete()
}

//--------------------------------------
// Query
//--------------------------------------

// Runs a query against a table.
func (s *Server) RunQuery(table *Table, query *Query) (interface{}, error) {
	var engine *ExecutionEngine
	engines := make([]*ExecutionEngine, 0)

	// Create a channel to receive aggregate responses.
	rchannel := make(chan interface{}, len(s.servlets))

	// Generate the query source code.
	source, err := query.Codegen()
	if err != nil {
		return nil, err
	}

	// Create an engine for merging results.
	engine, err = NewExecutionEngine(table, source)
	if err != nil {
		return nil, err
	}
	defer engine.Destroy()
	//fmt.Println(engine.FullAnnotatedSource())

	// Initialize one execution engine for each servlet.
	for _, servlet := range s.servlets {
		// Create an engine for each servlet.
		e, err := NewExecutionEngine(table, source)
		if err != nil {
			return nil, err
		}

		// Initialize iterator.
		ro := levigo.NewReadOptions()
		iterator := servlet.db.NewIterator(ro)
		err = e.SetIterator(iterator)
		if err != nil {
			return nil, err
		}

		engines = append(engines, e)
	}

	// Execute servlets asynchronously and retrieve responses outside
	// of the server context.
	for index, _ := range s.servlets {
		e := engines[index]
		go func() {
			if result, err := e.Aggregate(); err != nil {
				rchannel <- err
			} else {
				rchannel <- result
			}
		}()
	}

	// Wait for each servlet to complete and then merge the results.
	var servletError error
	var result interface{}
	result = make(map[interface{}]interface{})
	for i := 0; i < len(s.servlets); i++ {
		ret := <-rchannel
		if err, ok := ret.(error); ok {
			fmt.Printf("skyd.Server: Aggregate error: %v", err)
			servletError = err
		} else {
			// Defactorize aggregate results.
			err = query.Defactorize(ret)
			if err != nil {
				return nil, err
			}

			// Merge results.
			if ret != nil {
				result, err = engine.Merge(result, ret)
				if err != nil {
					fmt.Printf("skyd.Server: Merge error: %v", err)
					servletError = err
				}
			}
		}
	}
	err = servletError

	// Clean up engines.
	// TODO: Defer clean up earlier on in case of failure.
	for _, e := range engines {
		e.Destroy()
	}

	return result, err
}

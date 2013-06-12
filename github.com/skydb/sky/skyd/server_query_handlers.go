package skyd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addQueryHandlers() {
	s.ApiHandleFunc("/tables/{name}/stats", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.statsHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/query", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.queryHandler(w, req, params)
	}).Methods("POST")
	s.ApiHandleFunc("/tables/{name}/query/codegen", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.queryCodegenHandler(w, req, params)
	}).Methods("POST")
}

// GET /tables/:name/stats
func (s *Server) statsHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	// Run a simple count query.
	query := NewQuery(table, s.factors)
	selection := NewQuerySelection(query)
	selection.Fields = append(selection.Fields, NewQuerySelectionField("count", "count()"))
	query.Steps = append(query.Steps, selection)

	return s.RunQuery(table, query)
}

// POST /tables/:name/query
func (s *Server) queryHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	// Deserialize the query.
	query := NewQuery(table, s.factors)
	err = query.Deserialize(params)
	if err != nil {
		return nil, err
	}

	return s.RunQuery(table, query)
}

// POST /tables/:name/query/codegen
func (s *Server) queryCodegenHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Retrieve table and codegen query.
	var source string
	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	// Deserialize the query.
	query := NewQuery(table, s.factors)
	err = query.Deserialize(params)
	if err != nil {
		return nil, err
	}

	// Generate the query source code.
	source, err = query.Codegen()
	//fmt.Println(source)

	return source, &TextPlainContentTypeError{}
}

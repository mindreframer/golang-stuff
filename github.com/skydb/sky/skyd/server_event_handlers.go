package skyd

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (s *Server) addEventHandlers() {
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.getEventsHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.deleteEventsHandler(w, req, params)
	}).Methods("DELETE")

	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.getEventHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.replaceEventHandler(w, req, params)
	}).Methods("PUT")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.updateEventHandler(w, req, params)
	}).Methods("PATCH")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.deleteEventHandler(w, req, params)
	}).Methods("DELETE")
}

// GET /tables/:name/objects/:objectId/events
func (s *Server) getEventsHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	// Retrieve raw events.
	events, _, err := servlet.GetEvents(table, vars["objectId"])
	if err != nil {
		return nil, err
	}

	// Denormalize events.
	output := make([]map[string]interface{}, 0)
	for _, event := range events {
		e, err := table.SerializeEvent(event)
		if err != nil {
			return nil, err
		}
		err = table.DefactorizeEvent(event, s.factors)
		if err != nil {
			return nil, err
		}
		output = append(output, e)
	}

	return output, nil
}

// DELETE /tables/:name/objects/:objectId/events
func (s *Server) deleteEventsHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	return nil, servlet.DeleteEvents(table, vars["objectId"])
}

// GET /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) getEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	// Parse timestamp.
	timestamp, err := time.Parse(time.RFC3339, vars["timestamp"])
	if err != nil {
		return nil, err
	}

	// Find event.
	event, err := servlet.GetEvent(table, vars["objectId"], timestamp)
	if err != nil {
		return nil, err
	}
	// Return an empty event if there isn't one.
	if event == nil {
		event = NewEvent(vars["timestamp"], map[int64]interface{}{})
	}

	// Convert an event to a serializable object.
	e, err := table.SerializeEvent(event)
	if err != nil {
		return nil, err
	}
	err = table.DefactorizeEvent(event, s.factors)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// PUT /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) replaceEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	params["timestamp"] = vars["timestamp"]
	event, err := table.DeserializeEvent(params)
	if err != nil {
		return nil, err
	}
	err = table.FactorizeEvent(event, s.factors, true)
	if err != nil {
		return nil, err
	}

	return nil, servlet.PutEvent(table, vars["objectId"], event, true)
}

// PATCH /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) updateEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	params["timestamp"] = vars["timestamp"]
	event, err := table.DeserializeEvent(params)
	if err != nil {
		return nil, err
	}
	err = table.FactorizeEvent(event, s.factors, true)
	if err != nil {
		return nil, err
	}
	return nil, servlet.PutEvent(table, vars["objectId"], event, false)
}

// DELETE /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) deleteEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	timestamp, err := time.Parse(time.RFC3339, vars["timestamp"])
	if err != nil {
		return nil, fmt.Errorf("Unable to parse timestamp: %v", timestamp)
	}

	return nil, servlet.DeleteEvent(table, vars["objectId"], timestamp)
}

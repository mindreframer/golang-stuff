package skyd

import (
	"fmt"
	"os"
	"testing"
)

// Ensure that we can retrieve a list of all available tables on the server.
func TestServerGetTables(t *testing.T) {
	runTestServer(func(s *Server) {
		// Make and open one table.
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables", "application/json", `{"name":"foo"}`)
		resp.Body.Close()
		// Create another one as an empty directory.
		os.MkdirAll(s.TablePath("bar"), 0700)

		resp, err := sendTestHttpRequest("GET", "http://localhost:8586/tables", "application/json", ``)
		if err != nil {
			t.Fatalf("Unable to get tables: %v", err)
		}
		assertResponse(t, resp, 200, `[{"name":"bar"},{"name":"foo"}]`+"\n", "GET /tables failed.")
	})
}

// Ensure that we can retrieve a single table on the server.
func TestServerGetTable(t *testing.T) {
	runTestServer(func(s *Server) {
		// Make and open one table.
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables", "application/json", `{"name":"foo"}`)
		resp.Body.Close()
		resp, err := sendTestHttpRequest("GET", "http://localhost:8586/tables/foo", "application/json", ``)
		if err != nil {
			t.Fatalf("Unable to get table: %v", err)
		}
		assertResponse(t, resp, 200, `{"name":"foo"}`+"\n", "GET /table failed.")
	})
}

// Ensure that we can create a new table through the server.
func TestServerCreateTable(t *testing.T) {
	runTestServer(func(s *Server) {
		resp, err := sendTestHttpRequest("POST", "http://localhost:8586/tables", "application/json", `{"name":"foo"}`)
		if err != nil {
			t.Fatalf("Unable to create table: %v", err)
		}
		assertResponse(t, resp, 200, `{"name":"foo"}`+"\n", "POST /tables failed.")
		if _, err := os.Stat(fmt.Sprintf("%v/tables/foo", s.Path())); os.IsNotExist(err) {
			t.Fatalf("POST /tables did not create table.")
		}
	})
}

// Ensure that we can delete a table through the server.
func TestServerDeleteTable(t *testing.T) {
	runTestServer(func(s *Server) {
		// Create table.
		resp, err := sendTestHttpRequest("POST", "http://localhost:8586/tables", "application/json", `{"name":"foo"}`)
		if err != nil {
			t.Fatalf("Unable to create table: %v", err)
		}
		assertResponse(t, resp, 200, `{"name":"foo"}`+"\n", "POST /tables failed.")
		if _, err := os.Stat(fmt.Sprintf("%v/tables/foo", s.Path())); os.IsNotExist(err) {
			t.Fatalf("POST /tables did not create table.")
		}

		// Delete table.
		resp, _ = sendTestHttpRequest("DELETE", "http://localhost:8586/tables/foo", "application/json", ``)
		assertResponse(t, resp, 200, "", "DELETE /tables/:name failed.")
		if _, err := os.Stat(fmt.Sprintf("%v/tables/foo", s.Path())); !os.IsNotExist(err) {
			t.Fatalf("DELETE /tables/:name did not delete table.")
		}
	})
}

package skyd

import (
	"testing"
)

// Ensure that we can ping the server.
func TestServerPing(t *testing.T) {
	runTestServer(func(s *Server) {
		resp, err := sendTestHttpRequest("GET", "http://localhost:8586/ping", "application/json", "")
		if err != nil {
			t.Fatalf("Unable to ping: %v", err)
		}
		assertResponse(t, resp, 200, `{"message":"ok"}`+"\n", "GET /ping failed.")
	})
}

func BenchmarkPing(b *testing.B) {
	runTestServer(func(s *Server) {
		for i := 0; i < b.N; i++ {
			resp, _ := sendTestHttpRequest("GET", "http://localhost:8586/ping", "application/json", "")
			resp.Body.Close()
		}
	})
}

func BenchmarkRawPing(b *testing.B) {
	runTestServer(func(s *Server) {
		for i := 0; i < b.N; i++ {
			resp, _ := sendTestHttpRequest("GET", "http://localhost:8586/rawping", "application/json", "")
			resp.Body.Close()
		}
	})
}

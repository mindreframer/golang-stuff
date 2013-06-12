package skyd

import (
	"testing"
)

// Ensure that we can query the server for a count of events.
func TestServerSimpleCountQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "fruit", true, "string")
		setupTestData(t, "foo", [][]string{
			[]string{"a0", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
			[]string{"a1", "2012-01-01T00:00:00Z", `{"data":{"fruit":"grape"}}`},
			[]string{"a1", "2012-01-01T00:00:01Z", `{}`},
			[]string{"a2", "2012-01-01T00:00:00Z", `{"data":{"fruit":"orange"}}`},
			[]string{"a3", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
		})

		setupTestTable("bar")
		setupTestProperty("bar", "fruit", true, "string")
		setupTestData(t, "bar", [][]string{
			[]string{"a0", "2012-01-01T00:00:00Z", `{"data":{"fruit":"grape"}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"selection","dimensions":[],"fields":[{"name":"count","expression":"count()"}]}
			]
		}`
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"count":5}`+"\n", "POST /tables/:name/query failed.")
		resp, _ = sendTestHttpRequest("POST", "http://localhost:8586/tables/bar/query", "application/json", query)
		assertResponse(t, resp, 200, `{"count":1}`+"\n", "POST /tables/:name/query failed.")
	})
}

// Ensure that we can query the server for a count of events with a single dimension.
func TestServerOneDimensionCountQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "fruit", true, "string")
		setupTestData(t, "foo", [][]string{
			[]string{"b0", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
			[]string{"b1", "2012-01-01T00:00:00Z", `{"data":{"fruit":"grape"}}`},
			[]string{"b1", "2012-01-01T00:00:01Z", `{}`},
			[]string{"b2", "2012-01-01T00:00:00Z", `{"data":{"fruit":"orange"}}`},
			[]string{"b3", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"selection","dimensions":["fruit"],"fields":[{"name":"count","expression":"count()"}]}
			]
		}`
		//_codegen(t, "foo", query)
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"fruit":{"":{"count":1},"apple":{"count":2},"grape":{"count":1},"orange":{"count":1}}}`+"\n", "POST /tables/:name/query failed.")
	})
}

// Ensure that we can query the server for multiple selections with multiple dimensions.
func TestServerMultiDimensionalQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "gender", false, "string")
		setupTestProperty("foo", "state", false, "factor")
		setupTestProperty("foo", "price", true, "float")
		setupTestData(t, "foo", [][]string{
			[]string{"c0", "2012-01-01T00:00:00Z", `{"data":{"gender":"m", "state":"NY", "price":100}}`},
			[]string{"c0", "2012-01-01T00:00:01Z", `{"data":{"price":200}}`},
			[]string{"c0", "2012-01-01T00:00:02Z", `{"data":{"state":"CA","price":10}}`},

			[]string{"c1", "2012-01-01T00:00:00Z", `{"data":{"gender":"m", "state":"CA", "price":20}}`},
			[]string{"c1", "2012-01-01T00:00:01Z", `{"data":{}}`},

			[]string{"c2", "2012-01-01T00:00:00Z", `{"data":{"gender":"f", "state":"NY", "price":30}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"selection","name":"s1","dimensions":["gender","state"],"fields":[
					{"name":"count","expression":"count()"},
					{"name":"sum","expression":"sum(price)"}
				]},
				{"type":"selection","dimensions":["gender","state"],"fields":[
					{"name":"minimum","expression":"min(price)"},
					{"name":"maximum","expression":"max(price)"}
				]}
			]
		}`
		//_codegen(t, "foo", query)
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"gender":{"f":{"state":{"NY":{"maximum":30,"minimum":30}}},"m":{"state":{"CA":{"maximum":20,"minimum":0},"NY":{"maximum":200,"minimum":100}}}},"s1":{"gender":{"f":{"state":{"NY":{"count":1,"sum":30}}},"m":{"state":{"CA":{"count":3,"sum":30},"NY":{"count":2,"sum":300}}}}}}`+"\n", "POST /tables/:name/query failed.")
	})
}

// Ensure that we can perform a non-sessionized funnel analysis.
func TestServerFunnelAnalysisQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "action", false, "factor")
		setupTestData(t, "foo", [][]string{
			// A0[0..0]..A1[1..2] occurs twice for this object.
			[]string{"d0", "2012-01-01T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"d0", "2012-01-01T00:00:01Z", `{"data":{"action":"A1"}}`},
			[]string{"d0", "2012-01-01T00:00:02Z", `{"data":{"action":"A2"}}`},
			[]string{"d0", "2012-01-01T12:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"d0", "2012-01-01T13:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"d0", "2012-01-01T14:00:00Z", `{"data":{"action":"A1"}}`},

			// A0[0..0]..A1[1..2] occurs once for this object. (Second time matches A1[1..3]).
			[]string{"e1", "2012-01-01T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"e1", "2012-01-01T00:00:01Z", `{"data":{"action":"A0"}}`},
			[]string{"e1", "2012-01-01T00:00:02Z", `{"data":{"action":"A1"}}`},
			[]string{"e1", "2012-01-02T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"e1", "2012-01-02T00:00:01Z", `{"data":{"action":"A0"}}`},
			[]string{"e1", "2012-01-02T00:00:02Z", `{"data":{"action":"A0"}}`},
			[]string{"e1", "2012-01-02T00:00:03Z", `{"data":{"action":"A1"}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"condition","expression":"action == 'A0'","steps":[
					{"type":"condition","expression":"action == 'A1'","within":[1,2],"steps":[
						{"type":"selection","dimensions":["action"],"fields":[{"name":"count","expression":"count()"}]}
					]}
				]}
			]
		}`
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"action":{"A1":{"count":3}}}`+"\n", "POST /tables/:name/query failed.")
	})
}

// Ensure that we can perform a sessionized funnel analysis.
func TestServerSessionizedFunnelAnalysisQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "action", false, "string")
		setupTestData(t, "foo", [][]string{
			// A0[0..0]..A1[1..1] occurs once for this object. The second one is broken across sessions.
			[]string{"f0", "2012-01-01T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"f0", "2012-01-01T01:59:59Z", `{"data":{"action":"A1"}}`},
			[]string{"f0", "2012-01-02T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"f0", "2012-01-02T02:00:00Z", `{"data":{"action":"A1"}}`},
		})

		// Run query.
		query := `{
			"sessionIdleTime":7200,
			"steps":[
				{"type":"condition","expression":"action == 'A0'","steps":[
					{"type":"condition","expression":"action == 'A1'","within":[1,1],"steps":[
						{"type":"selection","dimensions":["action"],"fields":[{"name":"count","expression":"count()"}]}
					]}
				]}
			]
		}`
		//_codegen(t, "foo", query)
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8586/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"action":{"A1":{"count":1}}}`+"\n", "POST /tables/:name/query failed.")
	})
}

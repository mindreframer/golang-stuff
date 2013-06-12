package skyd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Query is a structured way of aggregating data in the database.
type Query struct {
	table           *Table
	factors         *Factors
	sequence        int
	Steps           QueryStepList
	SessionIdleTime int
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewQuery returns a new query.
func NewQuery(table *Table, factors *Factors) *Query {
	return &Query{
		table:   table,
		factors: factors,
		Steps:   make(QueryStepList, 0),
	}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// Retrieves the table this query is associated with.
func (q *Query) Table() *Table {
	return q.table
}

// Retrieves the factors this query is associated with.
func (q *Query) Factors() *Factors {
	return q.factors
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query into an untyped map.
func (q *Query) Serialize() map[string]interface{} {
	obj := map[string]interface{}{
		"sessionIdleTime": q.SessionIdleTime,
		"steps":           q.Steps.Serialize(),
	}
	return obj
}

// Decodes a query from an untyped map.
func (q *Query) Deserialize(obj map[string]interface{}) error {
	var err error

	// Deserialize "session idle time".
	if sessionIdleTime, ok := obj["sessionIdleTime"].(float64); ok || obj["sessionIdleTime"] == nil {
		q.SessionIdleTime = int(sessionIdleTime)
	} else {
		return fmt.Errorf("Invalid 'sessionIdleTime': %v", obj["sessionIdleTime"])
	}

	q.Steps, err = DeserializeQueryStepList(obj["steps"], q)
	if err != nil {
		return err
	}
	return nil
}

//--------------------------------------
// Encoding
//--------------------------------------

// Encodes a query to JSON.
func (q *Query) Encode(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(q.Serialize())
	return err
}

// Decodes a query from JSON.
func (q *Query) Decode(reader io.Reader) error {
	// Decode into an untyped object first since we need to determine the
	// type of steps to create.
	var obj map[string]interface{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&obj)
	if err != nil {
		return err
	}

	return q.Deserialize(obj)
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the query.
func (q *Query) Codegen() (string, error) {
	buffer := new(bytes.Buffer)

	// Generate aggregation functions.
	str, err := q.Steps.CodegenAggregateFunctions()
	if err != nil {
		return "", err
	}
	buffer.WriteString(str)
	buffer.WriteString(q.CodegenAggregateFunction())

	// Generate merge functions.
	str, err = q.Steps.CodegenMergeFunctions()
	if err != nil {
		return "", err
	}
	buffer.WriteString(str)
	buffer.WriteString(q.CodegenMergeFunction())

	return buffer.String(), nil
}

// Generates the 'aggregate()' function.
func (q *Query) CodegenAggregateFunction() string {
	buffer := new(bytes.Buffer)

	// Generate the function definition.
	fmt.Fprintln(buffer, "function aggregate(cursor, data)")

	// Set the session idle if one is available.
	if q.SessionIdleTime > 0 {
		fmt.Fprintf(buffer, "  cursor:set_session_idle(%d)\n", q.SessionIdleTime)
	}

	// Begin cursor loop.
	fmt.Fprintln(buffer, "  while cursor:next_session() do")
	fmt.Fprintln(buffer, "    while cursor:next() do")

	// Call each step function.
	for _, step := range q.Steps {
		fmt.Fprintf(buffer, "      %s(cursor, data)\n", step.FunctionName())
	}

	// End cursor loop.
	fmt.Fprintln(buffer, "    end")
	fmt.Fprintln(buffer, "  end")

	// End function.
	fmt.Fprintln(buffer, "end\n")

	return buffer.String()
}

// Generates the 'merge()' function.
func (q *Query) CodegenMergeFunction() string {
	buffer := new(bytes.Buffer)

	// Generate the function definition.
	fmt.Fprintln(buffer, "function merge(results, data)")

	// Call each step function if it has a merge function.
	fmt.Fprintf(buffer, q.Steps.CodegenMergeInvoke())

	// End function.
	fmt.Fprintln(buffer, "end\n")

	return buffer.String()
}

// Returns an autoincrementing numeric identifier.
func (q *Query) NextIdentifier() int {
	q.sequence += 1
	return q.sequence
}

//--------------------------------------
// Factorization
//--------------------------------------

// Converts factorized results from the aggregate function results to use
// the appropriate strings.
func (q *Query) Defactorize(data interface{}) error {
	return q.Steps.Defactorize(data)
}

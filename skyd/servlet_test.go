package skyd

import (
	"io/ioutil"
	"os"
	"testing"
)

// Ensure that we can open and close a servlet.
func TestOpen(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)

	servlet := NewServlet(path, nil)
	defer servlet.Close()
	err = servlet.Open()
	if err != nil {
		t.Fatalf("Unable to open servlet: %v", err)
	}
}

// Ensure that we can add events and read them back.
func TestServletPutEvent(t *testing.T) {
	// Setup blank database.
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	table := NewTable("test", "/tmp/test")
	servlet := NewServlet(path, nil)
	defer servlet.Close()
	_ = servlet.Open()

	// Setup source events.
	input := make([]*Event, 3)
	input[0] = NewEvent("2012-01-02T00:00:00Z", map[int64]interface{}{-1: 20, 1: "foo", 3: "baz"})
	input[1] = NewEvent("2012-01-01T00:00:00Z", map[int64]interface{}{-1: 20, 2: "bar", 3: "baz"})
	input[2] = NewEvent("2012-01-03T00:00:00Z", map[int64]interface{}{-1: 20, 1: "foo", 3: "baz"})

	// Add events.
	for _, e := range input {
		err = servlet.PutEvent(table, "bob", e, true)
		if err != nil {
			t.Fatalf("Unable to add event: %v", err)
		}
	}

	// Setup expected events.
	expected := make([]*Event, len(input))
	expected[0] = NewEvent("2012-01-01T00:00:00Z", map[int64]interface{}{-1: 20, 2: "bar"})
	expected[1] = input[0]
	expected[2] = NewEvent("2012-01-03T00:00:00Z", map[int64]interface{}{-1: 20})
	expectedState := NewEvent("2012-01-03T00:00:00Z", map[int64]interface{}{1: "foo", 2: "bar", 3: "baz"})

	// Read events out.
	output, state, err := servlet.GetEvents(table, "bob")
	if err != nil {
		t.Fatalf("Unable to retrieve events: %v", err)
	}
	if !expectedState.Equal(state) {
		t.Fatalf("Incorrect state.\nexp: %v\ngot: %v", expectedState, state)
	}
	if len(output) != len(expected) {
		t.Fatalf("Expected %v events, received %v", len(expected), len(output))
	}
	for i := range output {
		if !expected[i].Equal(output[i]) {
			t.Fatalf("Events not equal:\n  IN:  %v\n  OUT: %v", expected[i], output[i])
		}
	}
}

package skyd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// Ensure that we can create a new table.
func TestCreate(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	path = fmt.Sprintf("%v/test", path)

	table := NewTable("test", path)
	err = table.Create()
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}
	if !table.Exists() {
		t.Fatalf("Table doesn't exist: %v", path)
	}
}

// Ensure that we can create properties on a table.
func TestTableCreateProperty(t *testing.T) {
	table := createTempTable(t)
	table.Open()
	defer table.Close()

	property, err := table.CreateProperty("name", false, "string")
	if property == nil || err != nil {
		t.Fatalf("Unable to add property to table: %v", err)
	}

	content, _ := ioutil.ReadFile(fmt.Sprintf("%v/properties", table.Path()))
	if string(content) != "[{\"id\":1,\"name\":\"name\",\"transient\":false,\"dataType\":\"string\"}]\n" {
		t.Fatalf("Invalid properties file:\n%v", string(content))
	}
}

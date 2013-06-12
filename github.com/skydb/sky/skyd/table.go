package skyd

import (
	"errors"
	"fmt"
	"github.com/ugorji/go-msgpack"
	"os"
	"path/filepath"
	"time"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Table is a collection of objects.
type Table struct {
	Name         string `json:"name"`
	path         string
	propertyFile *PropertyFile
}

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// NewTable returns a new Table that is stored at a given path.
func NewTable(name string, path string) *Table {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil
	}

	return &Table{
		Name: name,
		path: path,
	}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// Retrieves the path on the table.
func (t *Table) Path() string {
	return t.path
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a table directory structure.
func (t *Table) Create() error {
	if t.Exists() {
		return fmt.Errorf("Table already exist: %v", t.Name)
	}

	// Create root directory.
	err := os.MkdirAll(t.path, 0700)
	if err != nil {
		return err
	}

	return nil
}

// Deletes a table.
func (t *Table) Delete() error {
	if !t.Exists() {
		return fmt.Errorf("Table does not exist: %v", t.Name)
	}

	// Close everything if it's open.
	if t.IsOpen() {
		t.Close()
	}

	// Delete the whole damn directory.
	os.RemoveAll(t.path)

	return nil
}

// Opens the table.
func (t *Table) Open() error {
	if !t.Exists() {
		return errors.New("Table does not exist")
	}

	// Load property file.
	t.propertyFile = NewPropertyFile(fmt.Sprintf("%v/%v", t.path, "properties"))
	err := t.propertyFile.Open()
	if err != nil {
		t.Close()
		return err
	}

	return nil
}

// Closes the table.
func (t *Table) Close() {
	if t.propertyFile != nil {
		t.propertyFile.Close()
	}
	t.propertyFile = nil
}

// Checks if the table is currently open.
func (t *Table) IsOpen() bool {
	return t.propertyFile != nil
}

// Checks if the table exists on disk.
func (t *Table) Exists() bool {
	if _, err := os.Stat(t.path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Generates a prefix key used for iterating over the table's data.
func TablePrefix(tableName string) ([]byte, error) {
	// The table prefix should match the encoded object id syntax but without the last item.
	prefix, err := msgpack.Marshal([]interface{}{tableName, nil})
	if err != nil {
		return nil, err
	}
	return prefix[0 : len(prefix)-1], nil
}

//--------------------------------------
// Property Management
//--------------------------------------

// Adds a property to the table.
func (t *Table) CreateProperty(name string, transient bool, dataType string) (*Property, error) {
	if !t.IsOpen() {
		return nil, errors.New("Table is not open")
	}

	// Create property on property file.
	property, err := t.propertyFile.CreateProperty(name, transient, dataType)
	if err != nil {
		return nil, err
	}

	// Save the property file to disk.
	err = t.propertyFile.Save()
	if err != nil {
		return nil, err
	}

	return property, err
}

// Retrieves a list of all properties on the table.
func (t *Table) GetProperties() ([]*Property, error) {
	if !t.IsOpen() {
		return nil, errors.New("Table is not open")
	}
	return t.propertyFile.GetProperties(), nil
}

// Retrieves a single property from the table by id.
func (t *Table) GetProperty(id int64) (*Property, error) {
	if !t.IsOpen() {
		return nil, errors.New("Table is not open")
	}
	return t.propertyFile.GetProperty(id), nil
}

// Retrieves a single property from the table by name.
func (t *Table) GetPropertyByName(name string) (*Property, error) {
	if !t.IsOpen() {
		return nil, errors.New("Table is not open")
	}
	return t.propertyFile.GetPropertyByName(name), nil
}

// Deletes a single property on the table.
func (t *Table) DeleteProperty(property *Property) error {
	if !t.IsOpen() {
		return errors.New("Table is not open")
	}
	t.propertyFile.DeleteProperty(property)
	return nil
}

// Saves the property file on the table.
func (t *Table) SavePropertyFile() error {
	if !t.IsOpen() {
		return errors.New("Table is not open")
	}
	return t.propertyFile.Save()
}

// Converts a map with string keys to use property identifier keys.
func (t *Table) NormalizeMap(m map[string]interface{}) (map[int64]interface{}, error) {
	return t.propertyFile.NormalizeMap(m)
}

// Converts a map with property identifier keys to use string keys.
func (t *Table) DenormalizeMap(m map[int64]interface{}) (map[string]interface{}, error) {
	return t.propertyFile.DenormalizeMap(m)
}

//--------------------------------------
// Event Encoding
//--------------------------------------

// Encodes an object identifier for this table.
func (t *Table) EncodeObjectId(objectId string) ([]byte, error) {
	return msgpack.Marshal([]string{t.Name, objectId})
}

// Deserializes a map into a normalized event.
func (t *Table) DeserializeEvent(m map[string]interface{}) (*Event, error) {
	event := &Event{}

	// Parse timestamp.
	if timestamp, ok := m["timestamp"].(string); ok {
		ts, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse timestamp: %v", timestamp)
		}
		event.Timestamp = ts
	} else {
		return nil, errors.New("Timestamp required.")
	}

	// Convert maps to use property identifiers.
	if data, ok := m["data"].(map[string]interface{}); ok {
		normalizedData, err := t.NormalizeMap(data)
		if err != nil {
			return nil, err
		}
		event.Data = normalizedData
	}

	return event, nil
}

// Serializes a normalized event into a map.
func (t *Table) SerializeEvent(event *Event) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	// Format timestamp.
	m["timestamp"] = event.Timestamp.UTC().Format(time.RFC3339)

	// Convert data map to use property names.
	if event.Data != nil {
		denormalizedData, err := t.DenormalizeMap(event.Data)
		if err != nil {
			return nil, err
		}
		m["data"] = denormalizedData
	} else {
		m["data"] = map[string]interface{}{}
	}

	return m, nil
}

//--------------------------------------
// Factorization
//--------------------------------------

// Factorizes the values in an event.
func (t *Table) FactorizeEvent(event *Event, factors *Factors, createIfMissing bool) error {
	if event == nil {
		return nil
	}

	propertyFile := t.propertyFile
	for k, v := range event.Data {
		property := propertyFile.GetProperty(k)
		if property.DataType == FactorDataType {
			if stringValue, ok := v.(string); ok {
				sequence, err := factors.Factorize(t.Name, property.Name, stringValue, createIfMissing)
				if err != nil {
					return err
				}
				event.Data[k] = sequence
			}
		}
	}

	return nil
}

// Defactorizes the values in an event.
func (t *Table) DefactorizeEvent(event *Event, factors *Factors) error {
	if event == nil {
		return nil
	}

	propertyFile := t.propertyFile
	for k, v := range event.Data {
		property := propertyFile.GetProperty(k)
		if property.DataType == FactorDataType {
			if sequence, ok := v.(uint64); ok {
				stringValue, err := factors.Defactorize(t.Name, property.Name, sequence)
				if err != nil {
					return err
				}
				event.Data[k] = stringValue
			}
		}
	}

	return nil
}

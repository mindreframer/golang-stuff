package skyd

import (
	"fmt"
)

// A Property is a loose schema column on a Table.
type Property struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Transient bool   `json:"transient"`
	DataType  string `json:"dataType"`
}

// NewProperty returns a new Property.
func NewProperty(id int64, name string, transient bool, dataType string) (*Property, error) {
	// Validate data type.
	switch dataType {
	case FactorDataType, StringDataType, IntegerDataType, FloatDataType, BooleanDataType:
	default:
		return nil, fmt.Errorf("Invalid property data type: %v", dataType)
	}

	return &Property{
		Id:        id,
		Name:      name,
		Transient: transient,
		DataType:  dataType,
	}, nil
}

package skyd

import (
	"errors"
	"fmt"
	"regexp"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

type QuerySelectionField struct {
	Name       string
	Expression string
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// Creates a new selection field.
func NewQuerySelectionField(name string, expression string) *QuerySelectionField {
	return &QuerySelectionField{Name:name, Expression:expression}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query selection into an untyped map.
func (f *QuerySelectionField) Serialize() map[string]interface{} {
	obj := map[string]interface{}{
		"name":       f.Name,
		"expression": f.Expression,
	}
	return obj
}

// Decodes a query selection from an untyped map.
func (f *QuerySelectionField) Deserialize(obj map[string]interface{}) error {
	if obj == nil {
		return errors.New("skyd.QuerySelectionField: Unable to deserialize nil.")
	}

	// Deserialize "expression".
	if expression, ok := obj["expression"].(string); ok && len(expression) > 0 {
		f.Expression = expression
	} else {
		return fmt.Errorf("skyd.QuerySelectionField: Invalid expression: %v", obj["expression"])
	}

	// Deserialize "name".
	if name, ok := obj["name"].(string); ok && len(name) > 0 {
		f.Name = name
	} else {
		return fmt.Errorf("skyd.QuerySelectionField: Invalid name: %v", obj["name"])
	}

	return nil
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the expression.
func (f *QuerySelectionField) CodegenExpression() (string, error) {
	r, _ := regexp.Compile(`^ *(?:count\(\)|(sum|min|max)\((\w+)\)|(\w+)) *$`)
	if m := r.FindStringSubmatch(f.Expression); m != nil {
		if len(m[1]) > 0 { // sum()/min()/max()
			switch m[1] {
			case "sum":
				return fmt.Sprintf("data.%s = (data.%s or 0) + cursor.event:%s()", f.Name, f.Name, m[2]), nil
			case "min":
				return fmt.Sprintf("if(data.%s == nil or data.%s > cursor.event:%s()) then data.%s = cursor.event:%s() end", f.Name, f.Name, m[2], f.Name, m[2]), nil
			case "max":
				return fmt.Sprintf("if(data.%s == nil or data.%s < cursor.event:%s()) then data.%s = cursor.event:%s() end", f.Name, f.Name, m[2], f.Name, m[2]), nil
			}
		} else if len(m[3]) > 0 { // assignment
			return fmt.Sprintf("data.%s = cursor.event:%s()", f.Name, m[3]), nil
		} else { // count()
			return fmt.Sprintf("data.%s = (data.%s or 0) + 1", f.Name, f.Name), nil
		}
	}

	return "", fmt.Errorf("skyd.QuerySelectionField: Invalid expression: %q", f.Expression)
}

// Generates Lua code for the merge expression.
func (f *QuerySelectionField) CodegenMergeExpression() (string, error) {
	r, _ := regexp.Compile(`^ *(?:count\(\)|(sum|min|max)\((\w+)\)|(\w+)) *$`)
	if m := r.FindStringSubmatch(f.Expression); m != nil {
		if len(m[1]) > 0 { // sum()/min()/max()
			switch m[1] {
			case "sum":
				return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", f.Name, f.Name, f.Name), nil
			case "min":
				return fmt.Sprintf("if(result.%s == nil or result.%s > data.%s) then result.%s = data.%s end", f.Name, f.Name, f.Name, f.Name, f.Name), nil
			case "max":
				return fmt.Sprintf("if(result.%s == nil or result.%s < data.%s) then result.%s = data.%s end", f.Name, f.Name, f.Name, f.Name, f.Name), nil
			}
		} else if len(m[3]) > 0 { // assignment
			return fmt.Sprintf("result.%s = data.%s", f.Name, f.Name), nil
		} else { // count()
			return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", f.Name, f.Name, f.Name), nil
		}
	}

	return "", fmt.Errorf("skyd.QuerySelectionField: Invalid merge expression: %q", f.Expression)
}

package skyd

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	QueryConditionUnitSteps    = "steps"
	QueryConditionUnitSessions = "sessions"
	QueryConditionUnitSeconds  = "seconds"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A condition step made within a query.
type QueryCondition struct {
	query            *Query
	functionName     string
	Expression       string
	WithinRangeStart int
	WithinRangeEnd   int
	WithinUnits      string
	Steps            QueryStepList
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// Creates a new condition.
func NewQueryCondition(query *Query) *QueryCondition {
	id := query.NextIdentifier()
	return &QueryCondition{
		query:            query,
		functionName:     fmt.Sprintf("a%d", id),
		WithinRangeStart: 0,
		WithinRangeEnd:   0,
		WithinUnits:      QueryConditionUnitSteps,
	}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// Retrieves the query this condition is associated with.
func (c *QueryCondition) Query() *Query {
	return c.query
}

// Retrieves the function name used during codegen.
func (c *QueryCondition) FunctionName() string {
	return c.functionName
}

// Retrieves the merge function name used during codegen.
func (c *QueryCondition) MergeFunctionName() string {
	return ""
}

// Retrieves the child steps.
func (c *QueryCondition) GetSteps() QueryStepList {
	return c.Steps
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query condition into an untyped map.
func (c *QueryCondition) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":        QueryStepTypeCondition,
		"expression":  c.Expression,
		"within":      []int{c.WithinRangeStart, c.WithinRangeEnd},
		"withinUnits": c.WithinUnits,
		"steps":       c.Steps.Serialize(),
	}
}

// Decodes a query condition from an untyped map.
func (c *QueryCondition) Deserialize(obj map[string]interface{}) error {
	if obj == nil {
		return errors.New("skyd.QueryCondition: Unable to deserialize nil.")
	}
	if obj["type"] != QueryStepTypeCondition {
		return fmt.Errorf("skyd.QueryCondition: Invalid step type: %v", obj["type"])
	}

	// Deserialize "expression".
	if expression, ok := obj["expression"].(string); ok {
		c.Expression = expression
	} else {
		if obj["expression"] == nil {
			c.Expression = "true"
		} else {
			return fmt.Errorf("Invalid 'expression': %v", obj["expression"])
		}
	}

	// Deserialize "within" range.
	if withinRange, ok := obj["within"].([]interface{}); ok && len(withinRange) == 2 {
		if withinRangeStart, ok := withinRange[0].(float64); ok {
			c.WithinRangeStart = int(withinRangeStart)
		} else {
			return fmt.Errorf("skyd.QueryCondition: Invalid 'within' range start: %v", withinRange[0])
		}
		if withinRangeEnd, ok := withinRange[1].(float64); ok {
			c.WithinRangeEnd = int(withinRangeEnd)
		} else {
			return fmt.Errorf("skyd.QueryCondition: Invalid 'within' range end: %v", withinRange[1])
		}
	} else {
		if obj["within"] == nil {
			c.WithinRangeStart = 0
			c.WithinRangeEnd = 0
		} else {
			return fmt.Errorf("Invalid 'within' range: %v", obj["within"])
		}
	}

	// Deserialize "within units".
	if withinUnits, ok := obj["withinUnits"].(string); ok {
		switch withinUnits {
		case QueryConditionUnitSteps, QueryConditionUnitSessions, QueryConditionUnitSeconds:
			c.WithinUnits = withinUnits
		default:
			return fmt.Errorf("Invalid 'within units': %v", withinUnits)
		}
	} else {
		if obj["withinUnits"] == nil {
			c.WithinUnits = QueryConditionUnitSteps
		} else {
			return fmt.Errorf("Invalid 'within units': %v", obj["withinUnits"])
		}
	}

	// Deserialize steps.
	var err error
	c.Steps, err = DeserializeQueryStepList(obj["steps"], c.query)
	if err != nil {
		return err
	}

	return nil
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the query.
func (c *QueryCondition) CodegenAggregateFunction() (string, error) {
	buffer := new(bytes.Buffer)

	// Validate.
	if c.WithinRangeStart > c.WithinRangeEnd {
		return "", fmt.Errorf("skyd.QueryCondition: Invalid 'within' range: %d..%d", c.WithinRangeStart, c.WithinRangeEnd)
	}

	// Generate child step functions.
	str, err := c.Steps.CodegenAggregateFunctions()
	if err != nil {
		return "", err
	}
	buffer.WriteString(str)

	// Generate main function.
	fmt.Fprintf(buffer, "function %s(cursor, data)\n", c.FunctionName())
	if c.WithinRangeStart > 0 {
		fmt.Fprintf(buffer, "  if cursor:eos() or cursor:eof() then return false end\n")
	}
	if c.WithinUnits == QueryConditionUnitSteps {
		fmt.Fprintf(buffer, "  index = 0\n")
	}
	fmt.Fprintf(buffer, "  repeat\n")
	if c.WithinUnits == QueryConditionUnitSteps {
		fmt.Fprintf(buffer, "    if index >= %d and index <= %d then\n", c.WithinRangeStart, c.WithinRangeEnd)
	}

	// Generate conditional expression.
	expressionCode, err := c.CodegenExpression()
	if err != nil {
		return "", err
	}
	fmt.Fprintf(buffer, "      if %s then\n", expressionCode)

	// Call each step function.
	for _, step := range c.Steps {
		fmt.Fprintf(buffer, "        %s(cursor, data)\n", step.FunctionName())
	}

	fmt.Fprintf(buffer, "        return true\n")
	fmt.Fprintf(buffer, "      end\n")
	fmt.Fprintf(buffer, "    end\n")
	if c.WithinUnits == QueryConditionUnitSteps {
		fmt.Fprintf(buffer, "    if index >= %d then break end\n", c.WithinRangeEnd)
		fmt.Fprintf(buffer, "    index = index + 1\n")
	}
	fmt.Fprintf(buffer, "  until not cursor:next()\n")
	fmt.Fprintf(buffer, "  return false\n")

	// End function definition.
	fmt.Fprintln(buffer, "end")

	return buffer.String(), nil
}

// Generates Lua code for the query.
func (c *QueryCondition) CodegenMergeFunction() (string, error) {
	buffer := new(bytes.Buffer)

	// Generate child step functions.
	str, err := c.Steps.CodegenMergeFunctions()
	if err != nil {
		return "", err
	}
	buffer.WriteString(str)

	return buffer.String(), nil
}

// Generates Lua code for the expression.
func (c *QueryCondition) CodegenExpression() (string, error) {
	// Do not transform simple booleans.
	if c.Expression == "true" || c.Expression == "false" {
		return c.Expression, nil
	}

	// Full expressions should be prepended with cursor's event reference.
	r, _ := regexp.Compile(`^ *(\w+) *(==) *(?:"([^"]*)"|'([^']*)'|(\d+(?:\.\d+)?)|(true|false)) *$`)
	m := r.FindSubmatch([]byte(c.Expression))
	if m == nil {
		return "", fmt.Errorf("skyd.QueryCondition: Invalid expression: %v", c.Expression)
	}

	// Find the property.
	property := c.query.table.propertyFile.GetPropertyByName(string(m[1]))
	if property == nil {
		return "", fmt.Errorf("skyd.QueryCondition: Property not found: %v", string(m[1]))
	}

	// Validate the expression value.
	var value string
	switch property.DataType {
	case FactorDataType, StringDataType:
		// Validate string value.
		var stringValue string
		if m[3] != nil {
			stringValue = string(m[3])
		} else if m[4] != nil {
			stringValue = string(m[4])
		} else {
			return "", fmt.Errorf("skyd.QueryCondition: Expression value must be a string literal for string and factor properties: %v", c.Expression)
		}

		// Convert factors.
		if property.DataType == FactorDataType {
			sequence, err := c.query.factors.Factorize(c.query.table.Name, property.Name, stringValue, false)
			if err != nil {
				return "", err
			} else {
				value = strconv.FormatUint(sequence, 10)
			}
		} else {
			value = fmt.Sprintf(`"%s"`, stringValue)
		}

	case IntegerDataType, FloatDataType:
		if m[5] == nil {
			return "", fmt.Errorf("skyd.QueryCondition: Expression value must be a numeric literal for integer and float properties: %v", c.Expression)
		}
		value = string(m[5])

	case BooleanDataType:
		if m[6] == nil {
			return "", fmt.Errorf("skyd.QueryCondition: Expression value must be a boolean literal for boolean properties: %v", c.Expression)
		}
		value = string(m[6])
	}

	return fmt.Sprintf("cursor.event:%s() %s %s", m[1], m[2], value), nil
}

//--------------------------------------
// Factorization
//--------------------------------------

// Converts factorized fields back to their original strings.
func (c *QueryCondition) Defactorize(data interface{}) error {
	return c.Steps.Defactorize(data)
}

package hood

import (
	"bytes"
	"reflect"
	"strings"
)

func toSnake(s string) string {
	buf := bytes.NewBufferString("")
	for i, v := range s {
		if i > 0 && v >= 'A' && v <= 'Z' {
			buf.WriteRune('_')
		}
		buf.WriteRune(v)
	}
	return strings.ToLower(buf.String())
}

func interfaceToSnake(f interface{}) string {
	t := reflect.TypeOf(f)
	for {
		c := false
		switch t.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
			t = t.Elem()
			c = true
		}
		if !c {
			break
		}
	}
	return toSnake(t.Name())
}

func snakeToUpperCamel(s string) string {
	buf := bytes.NewBufferString("")
	for _, v := range strings.Split(s, "_") {
		if len(v) > 0 {
			buf.WriteString(strings.ToUpper(v[:1]))
			buf.WriteString(v[1:])
		}
	}
	return buf.String()
}

func columnsMarkersAndValuesForModel(dialect Dialect, model *Model, markerPos *int) ([]string, []string, []interface{}) {
	columns := make([]string, 0, len(model.Fields))
	markers := make([]string, 0, len(columns))
	values := make([]interface{}, 0, len(columns))
	for _, column := range model.Fields {
		if !column.PrimaryKey() {
			columns = append(columns, column.Name)
			markers = append(markers, dialect.NextMarker(markerPos))
			values = append(values, column.Value)
		}
	}
	return columns, markers, values
}

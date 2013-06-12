package hood

import (
	"fmt"
	"reflect"
	"time"
)

func init() {
	RegisterDialect("mymysql", NewMysql())
}

type mysql struct {
	base
}

func NewMysql() Dialect {
	d := &mysql{}
	d.base.Dialect = d
	return d
}

func (d *mysql) NextMarker(pos *int) string {
	return "?"
}

func (d *mysql) Quote(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func (d *mysql) ParseBool(value reflect.Value) bool {
	return value.Int() != 0
}

func (d *mysql) SqlType(f interface{}, size int) string {
	switch f.(type) {
	case Id:
		return "bigint"
	case time.Time, Created, Updated:
		return "timestamp"
	case bool:
		return "boolean"
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return "int"
	case int64, uint64:
		return "bigint"
	case float32, float64:
		return "double"
	case []byte:
		if size > 0 && size < 65532 {
			return fmt.Sprintf("varbinary(%d)", size)
		}
		return "longblob"
	case string:
		if size > 0 && size < 65532 {
			return fmt.Sprintf("varchar(%d)", size)
		}
		return "longtext"
	}
	panic("invalid sql type")
}

func (d *mysql) KeywordAutoIncrement() string {
	return "AUTO_INCREMENT"
}

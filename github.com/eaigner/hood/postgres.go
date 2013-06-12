package hood

import (
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

func init() {
	RegisterDialect("postgres", NewPostgres())
}

type postgres struct {
	base
}

func NewPostgres() Dialect {
	d := &postgres{}
	d.base.Dialect = d
	return d
}

func (d *postgres) SqlType(f interface{}, size int) string {
	switch f.(type) {
	case Id:
		return "bigserial"
	case time.Time, Created, Updated:
		return "timestamp with time zone"
	case bool:
		return "boolean"
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return "integer"
	case int64, uint64:
		return "bigint"
	case float32, float64:
		return "double precision"
	case []byte:
		return "bytea"
	case string:
		if size > 0 && size < 65532 {
			return fmt.Sprintf("varchar(%d)", size)
		}
		return "text"
	}
	panic("invalid sql type")
}

func (d *postgres) Insert(hood *Hood, model *Model) (Id, error) {
	sql, args := d.Dialect.InsertSql(model)
	var id int64
	err := hood.QueryRow(sql, args...).Scan(&id)
	return Id(id), err
}

func (d *postgres) InsertSql(model *Model) (string, []interface{}) {
	m := 0
	columns, markers, values := columnsMarkersAndValuesForModel(d.Dialect, model, &m)
	quotedColumns := make([]string, 0, len(columns))
	for _, c := range columns {
		quotedColumns = append(quotedColumns, d.Dialect.Quote(c))
	}
	sql := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v) RETURNING %v",
		d.Dialect.Quote(model.Table),
		strings.Join(quotedColumns, ", "),
		strings.Join(markers, ", "),
		d.Dialect.Quote(model.Pk.Name),
	)
	return sql, values
}

func (d *postgres) KeywordAutoIncrement() string {
	// postgres has not auto increment keyword, uses SERIAL type
	return ""
}

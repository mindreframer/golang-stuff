package hood

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type base struct {
	Dialect Dialect
}

func (d *base) NextMarker(pos *int) string {
	m := fmt.Sprintf("$%d", *pos+1)
	*pos++
	return m
}

func (d *base) Quote(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}

func (d *base) ParseBool(value reflect.Value) bool {
	return value.Bool()
}

func (d *base) SetModelValue(driverValue, fieldValue reflect.Value) error {
	// ignore zero types
	if !driverValue.Elem().IsValid() {
		return nil
	}
	fieldType := fieldValue.Type()
	switch fieldValue.Type().Kind() {
	case reflect.Bool:
		fieldValue.SetBool(d.Dialect.ParseBool(driverValue.Elem()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue.SetInt(driverValue.Elem().Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// reading uint from int value causes panic
		switch driverValue.Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValue.SetUint(uint64(driverValue.Elem().Int()))
		default:
			fieldValue.SetUint(driverValue.Elem().Uint())
		}
	case reflect.Float32, reflect.Float64:
		fieldValue.SetFloat(driverValue.Elem().Float())
	case reflect.String:
		fieldValue.SetString(string(driverValue.Elem().Bytes()))
	case reflect.Slice:
		if reflect.TypeOf(driverValue.Interface()).Elem().Kind() == reflect.Uint8 {
			fieldValue.SetBytes(driverValue.Elem().Bytes())
		}
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			fieldValue.Set(driverValue.Elem())
		} else if fieldType == reflect.TypeOf(Updated{}) {
			if time, ok := driverValue.Elem().Interface().(time.Time); ok {
				fieldValue.Set(reflect.ValueOf(Updated{time}))
			} else {
				panic(fmt.Sprintf("cannot set updated value %T", driverValue.Elem().Interface()))
			}
		} else if fieldType == reflect.TypeOf(Created{}) {
			if time, ok := driverValue.Elem().Interface().(time.Time); ok {
				fieldValue.Set(reflect.ValueOf(Created{time}))
			} else {
				panic(fmt.Sprintf("cannot set created value %T", driverValue.Elem().Interface()))
			}
		}
	}
	return nil
}

func (d *base) ConvertHoodType(f interface{}) interface{} {
	if t, ok := f.(Created); ok {
		return t.Time
	}
	if t, ok := f.(Updated); ok {
		return t.Time
	}
	return f
}

func (d *base) appendWhere(query *[]string, args *[]interface{}, hood *Hood) {
	if x := hood.where; len(x) > 0 {
		for _, v := range x {
			// TODO: could be prettier!
			var c *clause
			switch p := v.(type) {
			case *whereClause:
				*query = append(*query, "WHERE")
				c = (*clause)(p)
			case *andClause:
				*query = append(*query, "AND")
				c = (*clause)(p)
			case *orClause:
				*query = append(*query, "OR")
				c = (*clause)(p)
			}
			if c != nil {
				*query = append(*query, c.a.Quote(d.Dialect), c.op)
				if path, ok := c.b.(Path); ok {
					*query = append(*query, path.Quote(d.Dialect))
				} else {
					*query = append(*query, "?")
					*args = append(*args, c.b)
				}
			} else {
				panic(fmt.Sprintf("invalid where clause %T", v))
			}
		}
	}
}

func (d *base) QuerySql(hood *Hood) (string, []interface{}) {
	query := make([]string, 0, 20)
	args := make([]interface{}, 0, 20)
	if hood.selectTable != "" {
		selector := "*"
		if paths := hood.selectPaths; len(paths) > 0 {
			quoted := []string{}
			for _, p := range paths {
				quoted = append(quoted, p.Quote(d.Dialect))
			}
			selector = strings.Join(quoted, ", ")
		}
		query = append(query, fmt.Sprintf("SELECT %v FROM %v", selector, d.Dialect.Quote(hood.selectTable)))
	}
	for _, j := range hood.joins {
		joinType := "INNER"
		switch j.join {
		case LeftJoin:
			joinType = "LEFT"
		case RightJoin:
			joinType = "RIGHT"
		case FullJoin:
			joinType = "FULL"
		}
		query = append(query, fmt.Sprintf(
			"%s JOIN %s ON %s = %s",
			joinType,
			d.Dialect.Quote(j.table),
			j.a.Quote(d.Dialect),
			j.b.Quote(d.Dialect),
		))
	}
	d.appendWhere(&query, &args, hood)
	if x := hood.groupBy; x != "" {
		query = append(query, fmt.Sprintf("GROUP BY %v", x.Quote(d.Dialect)))
	}
	if x := hood.havingCond; x != "" {
		query = append(query, fmt.Sprintf("HAVING %v", x))
		args = append(args, hood.havingArgs...)
	}
	if x := hood.orderBy; x != "" {
		query = append(query, fmt.Sprintf("ORDER BY %v", x.Quote(d.Dialect)))

		if x := hood.order; x != "" {
			query = append(query, fmt.Sprintf("%v", x))
		}
	}
	if x := hood.limit; x > 0 {
		query = append(query, "LIMIT ?")
		args = append(args, hood.limit)
	}
	if x := hood.offset; x > 0 {
		query = append(query, "OFFSET ?")
		args = append(args, hood.offset)
	}
	return hood.substituteMarkers(strings.Join(query, " ")), args
}

func (d *base) Insert(hood *Hood, model *Model) (Id, error) {
	sql, args := d.Dialect.InsertSql(model)
	result, err := hood.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return Id(id), nil
}

func (d *base) InsertSql(model *Model) (string, []interface{}) {
	m := 0
	columns, markers, values := columnsMarkersAndValuesForModel(d.Dialect, model, &m)
	quotedColumns := make([]string, 0, len(columns))
	for _, c := range columns {
		quotedColumns = append(quotedColumns, d.Dialect.Quote(c))
	}
	sql := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v)",
		d.Dialect.Quote(model.Table),
		strings.Join(quotedColumns, ", "),
		strings.Join(markers, ", "),
	)
	return sql, values
}

func (d *base) Update(hood *Hood, model *Model) (Id, error) {
	sql, args := d.Dialect.UpdateSql(model)
	_, err := hood.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return model.Pk.Value.(Id), nil
}

func (d *base) UpdateSql(model *Model) (string, []interface{}) {
	m := 0
	columns, markers, values := columnsMarkersAndValuesForModel(d.Dialect, model, &m)
	pairs := make([]string, 0, len(columns))
	for i, column := range columns {
		pairs = append(pairs, fmt.Sprintf("%v = %v", d.Dialect.Quote(column), markers[i]))
	}
	sql := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v = %v",
		d.Dialect.Quote(model.Table),
		strings.Join(pairs, ", "),
		d.Dialect.Quote(model.Pk.Name),
		d.Dialect.NextMarker(&m),
	)
	values = append(values, model.Pk.Value)
	return sql, values
}

func (d *base) Delete(hood *Hood, model *Model) (Id, error) {
	sql, args := d.Dialect.DeleteSql(model)
	_, err := hood.Exec(sql, args...)
	return args[0].(Id), err
}

func (d *base) DeleteSql(model *Model) (string, []interface{}) {
	n := 0
	return fmt.Sprintf(
		"DELETE FROM %v WHERE %v = %v",
		d.Dialect.Quote(model.Table),
		d.Dialect.Quote(model.Pk.Name),
		d.Dialect.NextMarker(&n),
	), []interface{}{model.Pk.Value}
}

func (d *base) DeleteFrom(hood *Hood, table string) error {
	sql, args := d.Dialect.DeleteFromSql(hood, table)
	_, err := hood.Exec(sql, args...)
	return err
}

func (d *base) DeleteFromSql(hood *Hood, table string) (string, []interface{}) {
	if len(hood.where) == 0 {
		panic("no where clause specified")
	}
	query := []string{
		fmt.Sprintf("DELETE FROM %s", d.Dialect.Quote(table)),
	}
	args := []interface{}{}
	d.appendWhere(&query, &args, hood)

	return hood.substituteMarkers(strings.Join(query, " ")), args
}

func (d *base) CreateTable(hood *Hood, model *Model) error {
	_, err := hood.Exec(d.Dialect.CreateTableSql(model, false))
	return err
}

func (d *base) CreateTableIfNotExists(hood *Hood, model *Model) error {
	_, err := hood.Exec(d.Dialect.CreateTableSql(model, true))
	return err
}

func (d *base) CreateTableSql(model *Model, ifNotExists bool) string {
	a := []string{"CREATE TABLE "}
	if ifNotExists {
		a = append(a, "IF NOT EXISTS ")
	}
	a = append(a, d.Dialect.Quote(model.Table), " ( ")
	for i, field := range model.Fields {
		b := []string{
			d.Dialect.Quote(field.Name),
			d.Dialect.SqlType(field.Value, field.Size()),
		}
		if field.NotNull() {
			b = append(b, d.Dialect.KeywordNotNull())
		}
		if x := field.Default(); x != "" {
			b = append(b, d.Dialect.KeywordDefault(x))
		}
		if field.PrimaryKey() {
			b = append(b, d.Dialect.KeywordPrimaryKey())
		}
		if incKeyword := d.Dialect.KeywordAutoIncrement(); field.PrimaryKey() && incKeyword != "" {
			b = append(b, incKeyword)
		}
		a = append(a, strings.Join(b, " "))
		if i < len(model.Fields)-1 {
			a = append(a, ", ")
		}
	}
	a = append(a, " )")
	return strings.Join(a, "")
}

func (d *base) DropTable(hood *Hood, table string) error {
	_, err := hood.Exec(d.Dialect.DropTableSql(table, false))
	return err
}

func (d *base) DropTableIfExists(hood *Hood, table string) error {
	_, err := hood.Exec(d.Dialect.DropTableSql(table, true))
	return err
}

func (d *base) DropTableSql(table string, ifExists bool) string {
	a := []string{"DROP TABLE"}
	if ifExists {
		a = append(a, "IF EXISTS")
	}
	a = append(a, d.Dialect.Quote(table))
	return strings.Join(a, " ")
}

func (d *base) RenameTable(hood *Hood, from, to string) error {
	_, err := hood.Exec(d.Dialect.RenameTableSql(from, to))
	return err
}

func (d *base) RenameTableSql(from, to string) string {
	return fmt.Sprintf("ALTER TABLE %v RENAME TO %v", d.Dialect.Quote(from), d.Dialect.Quote(to))
}

func (d *base) AddColumn(hood *Hood, table, column string, typ interface{}, size int) error {
	_, err := hood.Exec(d.Dialect.AddColumnSql(table, column, typ, size))
	return err
}

func (d *base) AddColumnSql(table, column string, typ interface{}, size int) string {
	return fmt.Sprintf(
		"ALTER TABLE %v ADD COLUMN %v %v",
		d.Dialect.Quote(table),
		d.Dialect.Quote(column),
		d.Dialect.SqlType(typ, size),
	)
}

func (d *base) RenameColumn(hood *Hood, table, from, to string) error {
	_, err := hood.Exec(d.Dialect.RenameColumnSql(table, from, to))
	return err
}

func (d *base) RenameColumnSql(table, from, to string) string {
	return fmt.Sprintf(
		"ALTER TABLE %v RENAME COLUMN %v TO %v",
		d.Dialect.Quote(table),
		d.Dialect.Quote(from),
		d.Dialect.Quote(to),
	)
}

func (d *base) ChangeColumn(hood *Hood, table, column string, typ interface{}, size int) error {
	_, err := hood.Exec(d.Dialect.ChangeColumnSql(table, column, typ, size))
	return err
}

func (d *base) ChangeColumnSql(table, column string, typ interface{}, size int) string {
	return fmt.Sprintf(
		"ALTER TABLE %v ALTER COLUMN %v TYPE %v",
		d.Dialect.Quote(table),
		d.Dialect.Quote(column),
		d.Dialect.SqlType(typ, size),
	)
}

func (d *base) DropColumn(hood *Hood, table, column string) error {
	_, err := hood.Exec(d.Dialect.DropColumnSql(table, column))
	return err
}

func (d *base) DropColumnSql(table, column string) string {
	return fmt.Sprintf(
		"ALTER TABLE %v DROP COLUMN %v",
		d.Dialect.Quote(table),
		d.Dialect.Quote(column),
	)
}

func (d *base) CreateIndex(hood *Hood, name, table string, unique bool, columns ...string) error {
	_, err := hood.Exec(d.Dialect.CreateIndexSql(name, table, unique, columns...))
	return err
}

func (d *base) CreateIndexSql(name, table string, unique bool, columns ...string) string {
	a := []string{"CREATE"}
	if unique {
		a = append(a, "UNIQUE")
	}
	quotedColumns := make([]string, 0, len(columns))
	for _, c := range columns {
		quotedColumns = append(quotedColumns, d.Dialect.Quote(c))
	}
	a = append(a, fmt.Sprintf(
		"INDEX %v ON %v (%v)",
		d.Dialect.Quote(name),
		d.Dialect.Quote(table),
		strings.Join(quotedColumns, ", "),
	))
	return strings.Join(a, " ")
}

func (d *base) DropIndex(hood *Hood, name string) error {
	_, err := hood.Exec(d.Dialect.DropIndexSql(name))
	return err
}

func (d *base) DropIndexSql(name string) string {
	return fmt.Sprintf("DROP INDEX %v", d.Dialect.Quote(name))
}

func (d *base) KeywordNotNull() string {
	return "NOT NULL"
}

func (d *base) KeywordDefault(s string) string {
	return fmt.Sprintf("DEFAULT %v", s)
}

func (d *base) KeywordPrimaryKey() string {
	return "PRIMARY KEY"
}

func (d *base) KeywordAutoIncrement() string {
	return "AUTOINCREMENT"
}

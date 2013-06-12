// Package hood provides a database agnostic, transactional ORM for the sql
// package
package hood

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	// Hood is an ORM handle.
	Hood struct {
		Db           *sql.DB
		Dialect      Dialect
		Log          bool
		qo           qo     // the query object
		schema       Schema // keeping track of the schema
		dryRun       bool   // if actual sql is executed or not
		selectPaths  []Path
		selectTable  string
		where        []interface{}
		markerPos    int
		limit        int
		offset       int
		orderBy      Path
		order        string
		joins        []*join
		groupBy      Path
		havingCond   string
		havingArgs   []interface{}
		firstTxError error
		mutex        sync.Mutex
	}

	// Id represents a auto-incrementing integer primary key type.
	Id int64

	// Index represents a table index and is returned via the Indexed interface.
	Index struct {
		Name    string
		Columns []string
		Unique  bool
	}

	// Indexes represents an array of indexes.
	Indexes []*Index

	// Created denotes a timestamp field that is automatically set on insert.
	Created struct {
		time.Time
	}

	// Updated denotes a timestamp field that is automatically set on update.
	Updated struct {
		time.Time
	}

	// Model represents a parsed schema interface{}.
	Model struct {
		Pk      *ModelField
		Table   string
		Fields  []*ModelField
		Indexes Indexes
	}

	// ModelField represents a schema field of a parsed model.
	ModelField struct {
		Name         string            // Column name
		Value        interface{}       // Value
		SqlTags      map[string]string // The sql struct tags for this field
		ValidateTags map[string]string // The validate struct tags for this field
		RawTag       reflect.StructTag // The raw tag
	}

	// Schema is a collection of models
	Schema []*Model

	// Config represents an environment entry in the config.json file
	Config map[string]string

	// Environment represents a configuration map for each environment specified
	// in the config.json file
	Environment map[string]Config

	// Path denotes a combined sql identifier such as 'table.column'
	Path string

	// Indexed defines the indexes for a table. You can invoke Add on the passed instance.
	Indexed interface {
		Indexes(indexes *Indexes)
	}

	// TODO: implement aggregate function types
	//
	// // Avg denotes an average aggregate function argument
	// Avg interface{}

	// // Min denotes an minimum aggregate function argument
	// Min interface{}

	// // Max denotes an maximum aggregate function argument
	// Max interface{}

	// // Std denotes an standard derivation aggregate function argument
	// Std interface{}

	// // Sum denotes an sum aggregate function argument
	// Sum interface{}

	qo interface {
		Prepare(query string) (*sql.Stmt, error)
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryRow(query string, args ...interface{}) *sql.Row
	}

	clause struct {
		a  Path
		op string
		b  interface{}
	}

	whereClause clause
	andClause   clause
	orClause    clause

	join struct {
		join  Join
		table string
		a     Path
		b     Path
	}
)

const (
	InnerJoin = Join(iota)
	LeftJoin
	RightJoin
	FullJoin
)

type Join int

// Add adds an index
func (ix *Indexes) Add(name string, columns ...string) {
	*ix = append(*ix, &Index{Name: name, Columns: columns, Unique: false})
}

// AddUnique adds an unique index
func (ix *Indexes) AddUnique(name string, columns ...string) {
	*ix = append(*ix, &Index{Name: name, Columns: columns, Unique: true})
}

// Quote quotes the path using the given dialects Quote method
func (p Path) Quote(d Dialect) string {
	sep := "."
	a := []string{}
	c := strings.Split(string(p), sep)
	for _, v := range c {
		a = append(a, d.Quote(v))
	}
	return strings.Join(a, sep)
}

// PrimaryKey tests if the field is declared using the sql tag "pk" or is of type Id
func (field *ModelField) PrimaryKey() bool {
	_, isPk := field.SqlTags["pk"]
	_, isId := field.Value.(Id)
	return isPk || isId
}

// NotNull tests if the field is declared as NOT NULL
func (field *ModelField) NotNull() bool {
	_, ok := field.SqlTags["notnull"]
	return ok
}

// Default returns the default value for the field
func (field *ModelField) Default() string {
	return field.SqlTags["default"]
}

// Size returns the field size, e.g. for varchars
func (field *ModelField) Size() int {
	v, ok := field.SqlTags["size"]
	if ok {
		i, _ := strconv.Atoi(v)
		return i
	}
	return 0
}

// Zero tests wether or not the field is set
func (field *ModelField) Zero() bool {
	x := field.Value
	return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// String returns the field string value and a bool flag indicating if the
// conversion was successful
func (field *ModelField) String() (string, bool) {
	t, ok := field.Value.(string)
	return t, ok
}

// Int returns the field int value and a bool flag indication if the conversion
// was successful
func (field *ModelField) Int() (int64, bool) {
	switch t := field.Value.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(t).Int(), true
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(t).Uint()), true
	}
	return 0, false
}

func (field *ModelField) GoDeclaration() string {
	t := ""
	if x := field.RawTag; len(x) > 0 {
		t = fmt.Sprintf("\t`%s`", x)
	}
	return fmt.Sprintf(
		"%s\t%s%s",
		snakeToUpperCamel(field.Name),
		reflect.TypeOf(field.Value).String(),
		t,
	)
}

// Validate tests if the field conforms to it's validation constraints specified
// int the "validate" struct tag
func (field *ModelField) Validate() error {
	// length
	if tuple, ok := field.ValidateTags["len"]; ok {
		s, ok := field.String()
		if ok {
			if err := validateLen(s, tuple, field.Name); err != nil {
				return err
			}
		}
	}
	// range
	if tuple, ok := field.ValidateTags["range"]; ok {
		i, ok := field.Int()
		if ok {
			if err := validateRange(i, tuple, field.Name); err != nil {
				return err
			}
		}
	}
	// presence
	if _, ok := field.ValidateTags["presence"]; ok {
		if field.Zero() {
			return NewValidationError(ValidationErrorValueNotSet, field.Name)
		}
	}

	// regexp
	if reg, ok := field.ValidateTags["regexp"]; ok {
		s, ok := field.String()
		if ok {
			if err := validateRegexp(s, reg, field.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseTuple(tuple string) (string, string) {
	c := strings.Split(tuple, ":")
	a := c[0]
	b := c[1]
	if len(c) != 2 || (len(a) == 0 && len(b) == 0) {
		panic("invalid validation tuple")
	}
	return a, b
}

func validateLen(s, tuple, field string) error {
	a, b := parseTuple(tuple)
	if len(a) > 0 {
		min, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		if len(s) < min {
			return NewValidationError(ValidationErrorValueTooShort, field)
		}
	}
	if len(b) > 0 {
		max, err := strconv.Atoi(b)
		if err != nil {
			panic(err)
		}
		if len(s) > max {
			return NewValidationError(ValidationErrorValueTooLong, field)
		}
	}
	return nil
}

func validateRange(i int64, tuple, field string) error {
	a, b := parseTuple(tuple)
	if len(a) > 0 {
		min, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			return err
		}
		if i < min {
			return NewValidationError(ValidationErrorValueTooSmall, field)
		}
	}
	if len(b) > 0 {
		max, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			return err
		}
		if i > max {
			return NewValidationError(ValidationErrorValueTooBig, field)
		}
	}
	return nil
}

func validateRegexp(s, reg, field string) error {
	matched, err := regexp.MatchString(reg, s)
	if err != nil {
		return err
	}
	if !matched {
		return NewValidationError(ValidationErrorValueNotMatch, field)
	}
	return nil
}

func (index *Index) GoDeclaration() string {
	u := ""
	if index.Unique {
		u = "Unique"
	}
	return fmt.Sprintf(
		"indexes.Add%s(\"%s\", \"%s\")",
		u,
		index.Name,
		strings.Join(index.Columns, "\", \""),
	)
}

func (model *Model) Validate() error {
	for _, field := range model.Fields {
		err := field.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (model *Model) GoDeclaration() string {
	tableName := snakeToUpperCamel(model.Table)
	a := []string{fmt.Sprintf("type %s struct {", tableName)}
	for _, f := range model.Fields {
		a = append(a, "\t"+f.GoDeclaration())
	}
	a = append(a, "}")
	if len(model.Indexes) > 0 {
		a = append(a,
			fmt.Sprintf("\nfunc (table *%s) Indexes(indexes *hood.Indexes) {", tableName),
		)
		for _, i := range model.Indexes {
			a = append(a, "\t"+i.GoDeclaration())
		}
		a = append(a, "}")
	}
	return strings.Join(a, "\n")
}

func (schema Schema) GoDeclaration() string {
	a := []string{}
	for _, m := range schema {
		a = append(a, m.GoDeclaration())
	}
	return strings.Join(a, "\n\n")
}

var registeredDialects map[string]Dialect = make(map[string]Dialect)

// New creates a new Hood using the specified DB and dialect.
func New(database *sql.DB, dialect Dialect) *Hood {
	hood := &Hood{
		Db:      database,
		Dialect: dialect,
		qo:      database,
	}
	hood.Reset()
	return hood
}

// Dry creates a new Hood instance for schema generation.
func Dry() *Hood {
	hd := New(nil, nil)
	hd.dryRun = true
	return hd
}

// Open opens a new database connection using the specified driver and data
// source name. It matches the sql.Open method signature.
func Open(driverName, dataSourceName string) (*Hood, error) {
	database, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	dialect := registeredDialects[driverName]
	if dialect == nil {
		return nil, errors.New("no dialect registered for driver name")
	}
	return New(database, dialect), nil
}

// Load opens a new database from a config.json file with the specified
// environment, or development if none is specified.
func Load(path, env string) (*Hood, error) {
	if env == "" {
		env = "development"
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var envs Environment
	err = dec.Decode(&envs)
	if err != nil {
		return nil, err
	}
	conf, ok := envs[env]
	if !ok {
		return nil, fmt.Errorf("config entry for specified environment '%s' not found", env)
	}
	return Open(conf["driver"], conf["source"])
}

// RegisterDialect registers a new dialect using the specified name and dialect.
func RegisterDialect(name string, dialect Dialect) {
	registeredDialects[name] = dialect
}

// Reset resets the internal state.
func (hood *Hood) Reset() {
	hood.selectPaths = nil
	hood.selectTable = ""
	hood.where = []interface{}{}
	hood.markerPos = 0
	hood.limit = 0
	hood.offset = 0
	hood.orderBy = ""
	hood.order = ""
	hood.joins = []*join{}
	hood.groupBy = ""
	hood.havingCond = ""
	hood.havingArgs = make([]interface{}, 0, 20)
}

// Copy copies the hood instance for safe context manipulation.
func (hood *Hood) Copy() *Hood {
	c := new(Hood)
	*c = *hood

	return c
}

// Begin starts a new transaction and returns a copy of Hood. You have to call
// subsequent methods on the newly returned object.
func (hood *Hood) Begin() *Hood {
	if hood.IsTransaction() {
		panic("cannot start nested transaction")
	}
	c := hood.Copy()
	q, err := hood.Db.Begin()
	if err != nil {
		panic(err)
	}
	c.firstTxError = nil
	c.qo = q

	return c
}

func (hood *Hood) logSql(sql string, args ...interface{}) {
	if hood.Log {
		a := make([]interface{}, 0, len(args))
		for _, v := range args {
			switch t := v.(type) {
			case []uint8:
				a = append(a, fmt.Sprintf("<[]byte#%d>", len(t)))
			default:
				a = append(a, v)
			}
		}
		log.Printf("\x1b[35mSQL: %s ARGS: %v\x1b[0m\n", sql, a)
	}
}

func (hood *Hood) updateTxError(e error) error {
	if e != nil {
		if hood.Log {
			log.Println("ERROR:", e)
		}
		// don't shadow the first error
		if hood.firstTxError == nil {
			hood.firstTxError = e
		}
	}
	return e
}

// Commit commits a started transaction and will report the first error that
// occurred inside the transaction.
func (hood *Hood) Commit() error {
	if v, ok := hood.qo.(*sql.Tx); ok {
		err := v.Commit()
		hood.updateTxError(err)
		return hood.firstTxError
	}
	return nil
}

// Rollback rolls back a started transaction.
func (hood *Hood) Rollback() error {
	if v, ok := hood.qo.(*sql.Tx); ok {
		return v.Rollback()
	}
	return nil
}

// IsTransaction returns wether the hood object represents an active transaction or not.
func (hood *Hood) IsTransaction() bool {
	_, ok := hood.qo.(*sql.Tx)
	return ok
}

// GoSchema returns a string of the schema file in Go syntax.
func (hood *Hood) GoSchema() string {
	timeRequired := false
L:
	for _, m := range hood.schema {
		for _, f := range m.Fields {
			switch f.Value.(type) {
			case time.Time:
				timeRequired = true
				break L
			}
		}
	}
	head := []string{
		"package db",
		"",
		"import (",
		"\t\"github.com/eaigner/hood\"",
	}
	if timeRequired {
		head = append(head, "\t\"time\"")
	}
	head = append(head, []string{")\n\n", hood.schema.GoDeclaration()}...)

	return strings.Join(head, "\n")
}

// Select adds a SELECT clause to the query with the specified table and columns.
// The table can either be a string or it's name can be inferred from the passed
// interface{} type.
func (hood *Hood) Select(table interface{}, paths ...Path) *Hood {
	hood.selectPaths = paths
	switch f := table.(type) {
	case string:
		hood.selectTable = f
	case interface{}:
		hood.selectTable = interfaceToSnake(f)
	default:
		panic("invalid table")
	}
	return hood
}

// Where adds a WHERE clause to the query. You can concatenate using the
// And and Or methods.
func (hood *Hood) Where(a Path, op string, b interface{}) *Hood {
	hood.where = append(hood.where, &whereClause{
		a:  a,
		op: op,
		b:  b,
	})
	return hood
}

// Where adds a AND clause to the WHERE query. You can concatenate using the
// And and Or methods.
func (hood *Hood) And(a Path, op string, b interface{}) *Hood {
	hood.where = append(hood.where, &andClause{
		a:  a,
		op: op,
		b:  b,
	})
	return hood
}

// Where adds a OR clause to the WHERE query. You can concatenate using the
// And and Or methods.
func (hood *Hood) Or(a Path, op string, b interface{}) *Hood {
	hood.where = append(hood.where, &orClause{
		a:  a,
		op: op,
		b:  b,
	})
	return hood
}

// Limit adds a LIMIT clause to the query.
func (hood *Hood) Limit(limit int) *Hood {
	hood.limit = limit
	return hood
}

// Offset adds an OFFSET clause to the query.
func (hood *Hood) Offset(offset int) *Hood {
	hood.offset = offset
	return hood
}

// OrderBy adds an ORDER BY clause to the query.
func (hood *Hood) OrderBy(path Path) *Hood {
	hood.orderBy = path
	return hood
}

func (hood *Hood) Asc() *Hood {
	hood.order = "ASC"
	return hood
}

func (hood *Hood) Desc() *Hood {
	hood.order = "DESC"
	return hood
}

// Join performs a JOIN on tables, for example
//   Join(hood.InnerJoin, &User{}, "user.id", "order.id")
func (hood *Hood) Join(op Join, table interface{}, a Path, b Path) *Hood {
	hood.joins = append(hood.joins, &join{
		join:  op,
		table: tableName(table),
		a:     a,
		b:     b,
	})
	return hood
}

// GroupBy adds a GROUP BY clause to the query.
func (hood *Hood) GroupBy(path Path) *Hood {
	hood.groupBy = path
	return hood
}

// Having adds a HAVING clause to the query.
func (hood *Hood) Having(condition string, args ...interface{}) *Hood {
	hood.havingCond = condition
	hood.havingArgs = append(hood.havingArgs, args...)
	return hood
}

// Find performs a find using the previously specified query. If no explicit
// SELECT clause was specified earlier, the select is inferred from the passed
// interface type.
func (hood *Hood) Find(out interface{}) error {
	// infer the select statement from the type if not set
	if hood.selectTable == "" {
		hood.Select(out)
	}
	query, args := hood.Dialect.QuerySql(hood)
	return hood.FindSql(out, query, args...)
}

// FindSql performs a find using the specified custom sql query and arguments and
// writes the results to the specified out interface{}.
func (hood *Hood) FindSql(out interface{}, query string, args ...interface{}) error {
	hood.mutex.Lock()
	defer hood.mutex.Unlock()
	defer hood.Reset()

	panicMsg := errors.New("expected pointer to struct slice *[]struct")
	if x := reflect.TypeOf(out).Kind(); x != reflect.Ptr {
		panic(panicMsg)
	}
	sliceValue := reflect.Indirect(reflect.ValueOf(out))
	if x := sliceValue.Kind(); x != reflect.Slice {
		panic(panicMsg)
	}
	sliceType := sliceValue.Type().Elem()
	if x := sliceType.Kind(); x != reflect.Struct {
		panic(panicMsg)
	}
	hood.logSql(query, args...)
	stmt, err := hood.qo.Prepare(query)
	if err != nil {
		return hood.updateTxError(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return hood.updateTxError(err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return hood.updateTxError(err)
	}
	for rows.Next() {
		containers := make([]interface{}, 0, len(cols))
		for i := 0; i < cap(containers); i++ {
			var v interface{}
			containers = append(containers, &v)
		}
		err := rows.Scan(containers...)
		if err != nil {
			return err
		}
		// create a new row and fill
		rowValue := reflect.New(sliceType)
		for i, v := range containers {
			key := cols[i]
			value := reflect.Indirect(reflect.ValueOf(v))
			name := snakeToUpperCamel(key)
			field := rowValue.Elem().FieldByName(name)
			if field.IsValid() {
				err = hood.Dialect.SetModelValue(value, field)
				if err != nil {
					return err
				}
			}
		}
		// append to output
		sliceValue.Set(reflect.Append(sliceValue, rowValue.Elem()))
	}
	return nil
}

// Exec executes a raw sql query.
func (hood *Hood) Exec(query string, args ...interface{}) (sql.Result, error) {
	hood.mutex.Lock()
	defer hood.mutex.Unlock()
	defer hood.Reset()

	query = hood.substituteMarkers(query)
	hood.logSql(query, args...)
	stmt, err := hood.qo.Prepare(query + ";")
	if err != nil {
		return nil, hood.updateTxError(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(hood.convertSpecialTypes(args)...)
	if err != nil {
		return nil, hood.updateTxError(err)
	}
	return result, nil
}

// Query executes a query that returns rows, typically a SELECT
func (hood *Hood) Query(query string, args ...interface{}) (*sql.Rows, error) {
	hood.mutex.Lock()
	defer hood.mutex.Unlock()

	hood.logSql(query, args...)
	return hood.qo.Query(query, hood.convertSpecialTypes(args)...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always return a non-nil value. Errors are deferred until Row's Scan
// method is called.
func (hood *Hood) QueryRow(query string, args ...interface{}) *sql.Row {
	hood.mutex.Lock()
	defer hood.mutex.Unlock()

	hood.logSql(query, args...)
	return hood.qo.QueryRow(query, hood.convertSpecialTypes(args)...)
	// TODO: switch to this implementation, as soon as
	//
	//   https://github.com/bmizerany/pq/issues/68
	//
	// is resolved!
	//
	//
	// query = hood.substituteMarkers(query)
	// if hood.Log {
	// 	fmt.Println(query)
	// }
	// stmt, err := hood.qo.Prepare(query)
	// if err != nil {
	// 	panic(err)
	// }
	// defer stmt.Close()
	// if hood.Log {
	// 	fmt.Println(args...)
	// }
	// return stmt.QueryRow(hood.convertSpecialTypes(args)...)
}

func (hood *Hood) convertSpecialTypes(a []interface{}) []interface{} {
	args := make([]interface{}, 0, len(a))
	for _, v := range a {
		args = append(args, hood.Dialect.ConvertHoodType(v))
	}
	return args
}

// Validate validates the provided struct
func (hood *Hood) Validate(f interface{}) error {
	model, err := interfaceToModel(f)
	if err != nil {
		return err
	}
	err = model.Validate()
	if err != nil {
		return err
	}
	// call validate methods
	err = callModelMethod(f, "Validate", true)
	if err != nil {
		return err
	}
	return nil
}

func callModelMethod(f interface{}, methodName string, isPrefix bool) error {
	typ := reflect.TypeOf(f)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if (isPrefix && strings.HasPrefix(method.Name, methodName)) ||
			(!isPrefix && method.Name == methodName) {
			ft := method.Func.Type()
			if ft.NumOut() == 1 &&
				ft.NumIn() == 1 {
				v := reflect.ValueOf(f).Method(i).Call([]reflect.Value{})
				if vdErr, ok := v[0].Interface().(error); ok {
					return vdErr
				}
			}
		}
	}
	return nil
}

// Save performs an INSERT, or UPDATE if the passed structs Id is set.
func (hood *Hood) Save(f interface{}) (Id, error) {
	var (
		id  Id = -1
		err error
	)
	model, err := interfaceToModel(f)
	if err != nil {
		return id, err
	}
	err = model.Validate()
	if err != nil {
		return id, err
	}
	err = callModelMethod(f, "BeforeSave", false)
	if err != nil {
		return id, err
	}
	if model.Pk == nil {
		panic("no primary key field")
	}
	now := time.Now()
	isUpdate := model.Pk != nil && !model.Pk.Zero()
	if isUpdate {
		err = callModelMethod(f, "BeforeUpdate", false)
		if err != nil {
			return id, err
		}
		for _, f := range model.Fields {
			switch f.Value.(type) {
			case Updated:
				f.Value = now
			}
		}
		id, err = hood.Dialect.Update(hood, model)
		if err == nil {
			err = callModelMethod(f, "AfterUpdate", false)
		}
	} else {
		err = callModelMethod(f, "BeforeInsert", false)
		if err != nil {
			return id, err
		}
		for _, f := range model.Fields {
			switch f.Value.(type) {
			case Created, Updated:
				f.Value = now
			}
		}
		id, err = hood.Dialect.Insert(hood, model)
		if err == nil {
			err = callModelMethod(f, "AfterInsert", false)
		}
	}
	if err == nil {
		err = callModelMethod(f, "AfterSave", false)
	}
	if id != -1 {
		// update model id after save
		structValue := reflect.Indirect(reflect.ValueOf(f))
		for i := 0; i < structValue.NumField(); i++ {
			field := structValue.Field(i)
			switch field.Interface().(type) {
			case Id:
				field.SetInt(int64(id))
			case Updated:
				field.Set(reflect.ValueOf(Updated{now}))
			case Created:
				if !isUpdate {
					field.Set(reflect.ValueOf(Created{now}))
				}
			}
		}
	}
	return id, err
}

func (hood *Hood) doAll(f interface{}, doFunc func(f2 interface{}) (Id, error)) ([]Id, error) {
	panicMsg := "expected pointer to struct slice *[]struct"
	if reflect.TypeOf(f).Kind() != reflect.Ptr {
		panic(panicMsg)
	}
	if reflect.TypeOf(f).Elem().Kind() != reflect.Slice {
		panic(panicMsg)
	}
	sliceValue := reflect.ValueOf(f).Elem()
	sliceLen := sliceValue.Len()
	ids := make([]Id, 0, sliceLen)
	for i := 0; i < sliceLen; i++ {
		id, err := doFunc(sliceValue.Index(i).Addr().Interface())
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// SaveAll performs an INSERT or UPDATE on a slice of structs.
func (hood *Hood) SaveAll(f interface{}) ([]Id, error) {
	return hood.doAll(f, func(f2 interface{}) (Id, error) {
		return hood.Save(f2)
	})
}

// Delete deletes the row matching the specified structs Id.
func (hood *Hood) Delete(f interface{}) (Id, error) {
	model, err := interfaceToModel(f)
	if err != nil {
		return -1, err
	}
	err = callModelMethod(f, "BeforeDelete", false)
	if err != nil {
		return -1, err
	}
	if model.Pk == nil {
		panic("no primary key field")
	}
	id, err := hood.Dialect.Delete(hood, model)
	if err != nil {
		return -1, err
	}
	return id, callModelMethod(f, "AfterDelete", false)
}

// DeleteAll deletes the rows matching the specified struct slice Ids.
func (hood *Hood) DeleteAll(f interface{}) ([]Id, error) {
	return hood.doAll(f, func(f2 interface{}) (Id, error) {
		return hood.Delete(f2)
	})
}

// DeleteFrom deletes the rows matched by the previous Where clause. table can
// either be a table struct or a string.
//
// Example:
//
//    hd.Where("amount", "=", 0).DeleteFrom("stock")
//
func (hood *Hood) DeleteFrom(table interface{}) error {
	defer hood.Reset()
	return hood.Dialect.DeleteFrom(hood, tableName(table))
}

// CreateTable creates a new table based on the provided schema.
func (hood *Hood) CreateTable(table interface{}) error {
	return hood.createTable(table, false)
}

// CreateTableIfNotExists creates a new table based on the provided schema
// if it does not exist yet.
func (hood *Hood) CreateTableIfNotExists(table interface{}) error {
	return hood.createTable(table, true)
}

func (hood *Hood) createTable(table interface{}, ifNotExists bool) error {
	if !hood.dryRun && !hood.IsTransaction() {
		panic("CreateTable can only be invoked inside a transaction")
	}
	model, err := interfaceToModel(table)
	if err != nil {
		return err
	}
	hood.schema = append(hood.schema, model)
	if hood.dryRun {
		return nil
	}
	if ifNotExists {
		hood.Dialect.CreateTableIfNotExists(hood, model)
	} else {
		hood.Dialect.CreateTable(hood, model)
	}
	for _, i := range model.Indexes {
		hood.Dialect.CreateIndex(hood, i.Name, model.Table, i.Unique, i.Columns...)
	}
	return hood.firstTxError
}

// DropTable drops the table matching the provided table name.
func (hood *Hood) DropTable(table interface{}) error {
	return hood.dropTable(table, false)
}

// DropTableIfExists drops the table matching the provided table name if it exists.
func (hood *Hood) DropTableIfExists(table interface{}) error {
	return hood.dropTable(table, true)
}

func (hood *Hood) dropTable(table interface{}, ifExists bool) error {
	s := []*Model{}
	for _, m := range hood.schema {
		if m.Table != tableName(table) {
			s = append(s, m)
		}
	}
	hood.schema = s
	if hood.dryRun {
		return nil
	}
	if ifExists {
		return hood.Dialect.DropTableIfExists(hood, tableName(table))
	}
	return hood.Dialect.DropTable(hood, tableName(table))
}

// RenameTable renames a table. The arguments can either be a schema definition
// or plain strings.
func (hood *Hood) RenameTable(from, to interface{}) error {
	for _, m := range hood.schema {
		if m.Table == tableName(from) {
			m.Table = tableName(to)
		}
	}
	if hood.dryRun {
		return nil
	}
	return hood.Dialect.RenameTable(hood, tableName(from), tableName(to))
}

// AddColumns adds the columns in the specified schema to the table.
func (hood *Hood) AddColumns(table, columns interface{}) error {
	if !hood.dryRun && !hood.IsTransaction() {
		panic("AddColumns can only be invoked inside a transaction")
	}
	m, err := interfaceToModel(columns)
	if err != nil {
		return err
	}
	for _, s := range hood.schema {
		if s.Table == tableName(table) {
			if m.Pk != nil {
				panic("primary keys can only be specified on table create (for now)")
			}
			s.Fields = append(s.Fields, m.Fields...)
		}
	}
	if hood.dryRun {
		return nil
	}
	for _, column := range m.Fields {
		err = hood.Dialect.AddColumn(hood, tableName(table), column.Name, column.Value, column.Size())
		if err != nil {
			return err
		}
	}
	return hood.firstTxError
}

// RenameColumn renames the column in the specified table.
func (hood *Hood) RenameColumn(table interface{}, from, to string) error {
	for _, s := range hood.schema {
		if s.Table == tableName(table) {
			for _, f := range s.Fields {
				if f.Name == from {
					f.Name = to
				}
			}
		}
	}
	if hood.dryRun {
		return nil
	}
	return hood.Dialect.RenameColumn(hood, tableName(table), from, to)
}

// ChangeColumn changes the data type of the specified column.
func (hood *Hood) ChangeColumns(table, column interface{}) error {
	if !hood.dryRun && !hood.IsTransaction() {
		panic("ChangeColumns can only be invoked inside a transaction")
	}
	m, err := interfaceToModel(column)
	if err != nil {
		return err
	}
	for _, s := range hood.schema {
		if s.Table == tableName(table) {
			fields := []*ModelField{}
			for _, oldField := range s.Fields {
				for _, newField := range m.Fields {
					if newField.Name == oldField.Name {
						fields = append(fields, newField)
					} else {
						fields = append(fields, oldField)
					}
				}
			}
			s.Fields = fields
		}
	}
	if hood.dryRun {
		return nil
	}
	for _, column := range m.Fields {
		err = hood.Dialect.ChangeColumn(hood, tableName(table), column.Name, column.Value, column.Size())
		if err != nil {
			return err
		}
	}
	return hood.firstTxError
}

// RemoveColumns removes the specified columns from the table.
func (hood *Hood) RemoveColumns(table, columns interface{}) error {
	if !hood.dryRun && !hood.IsTransaction() {
		panic("RemoveColumns can only be invoked inside a transaction")
	}
	m, err := interfaceToModel(columns)
	if err != nil {
		return err
	}
	for _, s := range hood.schema {
		if s.Table == tableName(table) {
			fields := []*ModelField{}
			for _, field := range s.Fields {
				remove := false
				for _, fieldToRemove := range m.Fields {
					if field.Name == fieldToRemove.Name {
						remove = true
						break
					}
				}
				if !remove {
					fields = append(fields, field)
				}
			}
			s.Fields = fields
		}
	}
	if hood.dryRun {
		return nil
	}
	for _, column := range m.Fields {
		err = hood.Dialect.DropColumn(hood, tableName(table), column.Name)
		if err != nil {
			return err
		}
	}
	return hood.firstTxError
}

// CreateIndex creates the specified index on table.
func (hood *Hood) CreateIndex(table interface{}, name string, unique bool, columns ...string) error {
	if !hood.dryRun && !hood.IsTransaction() {
		panic("CreateIndex can only be invoked inside a transaction")
	}
	tn := tableName(table)
	index := &Index{Name: name, Columns: columns, Unique: unique}
	for _, s := range hood.schema {
		if s.Table == tn {
			s.Indexes = append(s.Indexes, index)
		}
	}
	if hood.dryRun {
		return nil
	}
	err := hood.Dialect.CreateIndex(hood, index.Name, tn, index.Unique, index.Columns...)
	if err != nil {
		return err
	}
	return hood.firstTxError
}

// DropIndex drops the specified index from table.
func (hood *Hood) DropIndex(table interface{}, name string) error {
	tn := tableName(table)
	for _, s := range hood.schema {
		if s.Table == tn {
			indexes := []*Index{}
			for _, i := range s.Indexes {
				if i.Name != name {
					indexes = append(indexes, i)
				}
			}
			s.Indexes = indexes
		}
	}
	if hood.dryRun {
		return nil
	}
	return hood.Dialect.DropIndex(hood, name)
}

func (hood *Hood) substituteMarkers(query string) string {
	// in order to use a uniform marker syntax, substitute
	// all question marks with the dialect marker
	chunks := make([]string, 0, len(query)*2)
	for _, v := range query {
		if v == '?' {
			chunks = append(chunks, hood.Dialect.NextMarker(&hood.markerPos))
		} else {
			chunks = append(chunks, string(v))
		}
	}
	return strings.Join(chunks, "")
}

func parseTags(s string) map[string]string {
	c := strings.Split(s, ",")
	m := make(map[string]string)
	for _, v := range c {
		c2 := strings.Split(v, "(")
		if len(c2) == 2 && len(c2[1]) > 1 {
			m[c2[0]] = c2[1][:len(c2[1])-1]
		} else {
			m[v] = ""
		}
	}
	return m
}

func addFields(m *Model, t reflect.Type, v reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sqlTag := field.Tag.Get("sql")
		if sqlTag == "-" {
			continue
		}
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			addFields(m, field.Type, v.Field(i))
			continue
		}
		parsedSqlTags := parseTags(sqlTag)
		rawValidateTag := field.Tag.Get("validate")
		parsedValidateTags := make(map[string]string)
		if len(rawValidateTag) > 0 {
			if rawValidateTag[:1] == "^" {
				parsedValidateTags["regexp"] = rawValidateTag
			} else {
				parsedValidateTags = parseTags(rawValidateTag)
			}
		}
		fd := &ModelField{
			Name:         toSnake(field.Name),
			Value:        v.FieldByName(field.Name).Interface(),
			SqlTags:      parsedSqlTags,
			ValidateTags: parsedValidateTags,
			RawTag:       field.Tag,
		}
		if fd.PrimaryKey() {
			m.Pk = fd
		}
		m.Fields = append(m.Fields, fd)
	}
}

func addIndexes(m *Model, f interface{}) {
	if t, ok := f.(Indexed); ok {
		t.Indexes(&m.Indexes)
	}
}

func interfaceToModel(f interface{}) (*Model, error) {
	v := reflect.Indirect(reflect.ValueOf(f))
	if v.Kind() != reflect.Struct {
		return nil, errors.New("model is not a struct")
	}
	t := v.Type()
	m := &Model{
		Pk:      nil,
		Table:   interfaceToSnake(f),
		Fields:  []*ModelField{},
		Indexes: Indexes{},
	}
	addFields(m, t, v)
	addIndexes(m, f)
	return m, nil
}

func tableName(f interface{}) string {
	switch t := f.(type) {
	case string:
		return t
	}
	m, _ := interfaceToModel(f)
	if m != nil {
		return m.Table
	}
	panic("invalid table name")
}

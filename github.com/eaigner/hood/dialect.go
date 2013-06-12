package hood

import (
	"reflect"
)

type Dialect interface {
	// NextMarker returns the dialect specific marker for a prepared statement,
	// for instance $1, $2, ... and increments the position by one.
	// The position starts at 0.
	NextMarker(pos *int) string

	// Quote will quote identifiers in a SQL statement.
	Quote(s string) string

	// SqlType returns the SQL type for the provided interface type. The size
	// parameter delcares the data size for the column (e.g. for VARCHARs).
	SqlType(f interface{}, size int) string

	// If database do not support boolean type this can be used to parse int
	// value to boolean value.
	ParseBool(value reflect.Value) bool

	// SetModelValue sets a model field from a db value.
	//
	// For example: time.Time objects needs to be marshalled back and forth
	// as Strings for databases that don't have a native "time" type.
	SetModelValue(value reflect.Value, field reflect.Value) error

	// ConvertHoodType converts special types such as Created or Updated to
	// values the driver can understand.
	ConvertHoodType(f interface{}) interface{}

	// QuerySql returns the resulting query sql and attributes.
	QuerySql(hood *Hood) (sql string, args []interface{})

	// Insert inserts the values in model and returns the inserted rows Id.
	Insert(hood *Hood, model *Model) (Id, error)

	// InsertSql returns the sql for inserting the passed model.
	InsertSql(model *Model) (sql string, args []interface{})

	// Update updates the values in the specified model and returns the
	// updated rows Id.
	Update(hood *Hood, model *Model) (Id, error)

	// UpdateSql returns the sql for updating the specified model.
	UpdateSql(model *Model) (string, []interface{})

	// Delete drops the row matching the primary key of model and returns the affected Id.
	Delete(hood *Hood, model *Model) (Id, error)

	// DeleteSql returns the sql for deleting the row matching model's primary key.
	DeleteSql(model *Model) (string, []interface{})

	// DeleteFrom deletes the matching rows in the specified table
	DeleteFrom(hood *Hood, table string) error

	// DeleteFromSql returns the sql for DeleteFrom
	DeleteFromSql(hood *Hood, table string) (string, []interface{})

	// CreateTable creates the table specified in model.
	CreateTable(hood *Hood, model *Model) error

	// CreateTableIfNotExists creates the table specified in model if it does not exist.
	CreateTableIfNotExists(hood *Hood, model *Model) error

	// CreateTableSql returns the sql for creating a table.
	CreateTableSql(model *Model, ifNotExists bool) string

	// DropTable drops the specified table.
	DropTable(hood *Hood, table string) error

	// DropTableIfExists drops the table if it exists.
	DropTableIfExists(hood *Hood, table string) error

	// DropTableSql returns the sql for dropping the specified table.
	DropTableSql(table string, ifExists bool) string

	// RenameTable renames the specified table.
	RenameTable(hood *Hood, from, to string) error

	// RenameTableSql returns the sql for renaming the specified table.
	RenameTableSql(from, to string) string

	// AddColumn adds the columns to the corresponding table.
	AddColumn(hood *Hood, table, column string, typ interface{}, size int) error

	// AddColumnSql returns the sql for adding the specified column in table.
	AddColumnSql(table, column string, typ interface{}, size int) string

	// RenameColumn renames a table column in the specified table.
	RenameColumn(hood *Hood, table, from, to string) error

	// RenameColumnSql returns the sql for renaming the specified column in table.
	RenameColumnSql(table, from, to string) string

	// ChangeColumn changes the data type of the specified column.
	ChangeColumn(hood *Hood, table, column string, typ interface{}, size int) error

	// ChangeColumnSql returns the sql for changing the column data type.
	ChangeColumnSql(table, column string, typ interface{}, size int) string

	// DropColumn removes the specified column.
	DropColumn(hood *Hood, table, column string) error

	// DropColumnSql returns the sql for removing the column.
	DropColumnSql(table, column string) string

	// CreateIndex creates an index on the specified column.
	CreateIndex(hood *Hood, name, table string, unique bool, columns ...string) error

	// CreateIndexSql returns the sql for creating an index on the specified column.
	CreateIndexSql(name, table string, unique bool, columns ...string) string

	// DropIndex drops the index.
	DropIndex(hood *Hood, name string) error

	// DropIndexSql returns the sql for dropping the index.
	DropIndexSql(name string) string

	// KeywordNotNull returns the dialect specific keyword for 'NOT NULL'.
	KeywordNotNull() string

	// KeywordDefault returns the dialect specific keyword for 'DEFAULT'.
	KeywordDefault(s string) string

	// KeywordPrimaryKey returns the dialect specific keyword for 'PRIMARY KEY'.
	KeywordPrimaryKey() string

	// KeywordAutoIncrement returns the dialect specific keyword for 'AUTO_INCREMENT'.
	KeywordAutoIncrement() string
}

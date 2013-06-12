package hood

import (
	"strings"
	"testing"
	"time"
)

func TestParseTags(t *testing.T) {
	m := parseTags(`pk`)
	if x, ok := m["pk"]; !ok {
		t.Fatal("wrong value", ok, x)
	}
	m = parseTags(`notnull,default('banana')`)
	if x, ok := m["notnull"]; !ok {
		t.Fatal("wrong value", ok, x)
	}
	if x, ok := m["default"]; !ok || x != "'banana'" {
		t.Fatal("wrong value", x)
	}
}

func TestFieldZero(t *testing.T) {
	field := &ModelField{}
	field.Value = nil
	if !field.Zero() {
		t.Fatal("should be zero")
	}
	field.Value = 0
	if !field.Zero() {
		t.Fatal("should be zero")
	}
	field.Value = ""
	if !field.Zero() {
		t.Fatal("should be zero")
	}
	field.Value = false
	if !field.Zero() {
		t.Fatal("should be zero")
	}
	field.Value = true
	if field.Zero() {
		t.Fatal("should not be zero")
	}
	field.Value = -1
	if field.Zero() {
		t.Fatal("should not be zero")
	}
	field.Value = 1
	if field.Zero() {
		t.Fatal("should not be zero")
	}
	field.Value = "asdf"
	if field.Zero() {
		t.Fatal("should not be zero")
	}
}

func TestFieldValidate(t *testing.T) {
	type Schema struct {
		A string `validate:"len(3:6)"`
		B int    `validate:"range(10:20)"`
		C string `validate:"len(:4),presence"`
		D string `validate:"^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"`
	}
	m, _ := interfaceToModel(&Schema{})
	a := m.Fields[0]
	if x := len(a.ValidateTags); x != 1 {
		t.Fatal("wrong len", x)
	}
	if x, ok := a.ValidateTags["len"]; !ok || x != "3:6" {
		t.Fatal("wrong value", x, ok)
	}
	if err := a.Validate(); err == nil || err.Error() != "a too short" {
		t.Fatal("should not validate")
	}
	a.Value = "abc"
	if err := a.Validate(); err != nil {
		t.Fatal("should validate", err)
	}
	a.Value = "abcdefg"
	if err := a.Validate(); err == nil || err.Error() != "a too long" {
		t.Fatal("should not validate")
	}

	b := m.Fields[1]
	if x := len(b.ValidateTags); x != 1 {
		t.Fatal("wrong len", x)
	}
	if err := b.Validate(); err == nil || err.Error() != "b too small" {
		t.Fatal("should not validate")
	}
	b.Value = 10
	if err := b.Validate(); err != nil {
		t.Fatal("should validate", err)
	}
	b.Value = 21
	if err := b.Validate(); err == nil || err.Error() != "b too big" {
		t.Fatal("should not validate")
	}

	c := m.Fields[2]
	if x := len(c.ValidateTags); x != 2 {
		t.Fatal("wrong len", x)
	}
	if err := c.Validate(); err == nil || err.Error() != "c not set" {
		t.Fatal("should not validate")
	}
	c.Value = "a"
	if err := c.Validate(); err != nil {
		t.Fatal("should validate", err)
	}
	c.Value = "abcde"
	if err := c.Validate(); err == nil || err.Error() != "c too long" {
		t.Fatal("should not validate")
	}

	d := m.Fields[3]
	if x := len(d.ValidateTags); x != 1 {
		t.Fatal("wrong len", x)
	}
	d.Value = "gggg@gmail.com"
	if err := d.Validate(); err != nil {
		t.Fatal("should validate", err)
	}
	d.Value = "www.google.com"
	if err := d.Validate(); err == nil || err.Error() != "d not match" {
		t.Fatal("should not validate", err)
	}
}

func TestFieldOmit(t *testing.T) {
	type Schema struct {
		A string `sql:"-"`
		B string
	}
	m, _ := interfaceToModel(&Schema{})
	if x := len(m.Fields); x != 1 {
		t.Fatal("wrong len", x)
	}
}

type validateSchema struct {
	A string
}

var numValidateFuncCalls = 0

func (v *validateSchema) ValidateX() error {
	numValidateFuncCalls++
	if v.A == "banana" {
		return NewValidationError(1, "value cannot be banana")
	}
	return nil
}

func (v *validateSchema) ValidateY() error {
	numValidateFuncCalls++
	return NewValidationError(2, "ValidateY failed")
}

func TestValidationMethods(t *testing.T) {
	hd := New(nil, &postgres{})
	m := &validateSchema{}
	err := hd.Validate(m)
	if err == nil || err.Error() != "ValidateY failed" {
		t.Fatal("wrong error", err)
	}
	if v, ok := err.(*ValidationError); !ok {
		t.Fatal("should be of type ValidationError", v)
	}
	if numValidateFuncCalls != 2 {
		t.Fatal("should have called validation func")
	}
	numValidateFuncCalls = 0
	m.A = "banana"
	err = hd.Validate(m)
	if err == nil || err.Error() != "value cannot be banana" {
		t.Fatal("wrong error", err)
	}
	if numValidateFuncCalls != 1 {
		t.Fatal("should have called validation func")
	}
}

func TestInterfaceToModelWithEmbedded(t *testing.T) {
	type embed struct {
		Name  string
		Value string
	}
	type table struct {
		ColPrimary Id
		embed
	}
	table1 := &table{
		6, embed{"Mrs. A", "infinite"},
	}
	m, err := interfaceToModel(table1)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	f := m.Fields[1]
	if x, ok := f.Value.(string); !ok || x != "Mrs. A" {
		t.Fatal("wrong value from embedded struct")
	}
}

type indexedTable struct {
	ColPrimary    Id
	ColAltPrimary string `sql:"pk"`
	ColNotNull    string `sql:"notnull,default('banana')"`
	ColVarChar    string `sql:"size(64)"`
	ColTime       time.Time
}

func (table *indexedTable) Indexes(indexes *Indexes) {
	indexes.Add("my_index", "col_primary", "col_time")
	indexes.AddUnique("my_unique_index", "col_var_char", "col_time")
}

func TestInterfaceToModel(t *testing.T) {
	type table struct {
		ColPrimary    Id
		ColAltPrimary string `sql:"pk"`
		ColNotNull    string `sql:"notnull,default('banana')"`
		ColVarChar    string `sql:"size(64)"`
		ColTime       time.Time
	}
	now := time.Now()
	table1 := &indexedTable{
		ColPrimary:    6,
		ColAltPrimary: "banana",
		ColVarChar:    "orange",
		ColTime:       now,
	}
	m, err := interfaceToModel(table1)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if m.Pk == nil {
		t.Fatal("pk nil")
	}
	if m.Pk.Name != "col_alt_primary" {
		t.Fatal("wrong value", m.Pk.Name)
	}
	if x := len(m.Fields); x != 5 {
		t.Fatal("wrong value", x)
	}
	if x := len(m.Indexes); x != 2 {
		t.Fatal("wrong value", x)
	}
	if x := m.Indexes[0].Name; x != "my_index" {
		t.Fatal("wrong index name", x)
	}
	if x := m.Indexes[0].Columns; strings.Join(x, ":") != "col_primary:col_time" {
		t.Fatal("wrong index columns", x)
	}
	if x := m.Indexes[0].Unique; x != false {
		t.Fatal("wrong index uniqueness", x)
	}
	if x := m.Indexes[1].Name; x != "my_unique_index" {
		t.Fatal("wrong index name", x)
	}
	if x := m.Indexes[1].Columns; strings.Join(x, ":") != "col_var_char:col_time" {
		t.Fatal("wrong index columns", x)
	}
	if x := m.Indexes[1].Unique; x != true {
		t.Fatal("wrong index uniqueness", x)
	}
	f := m.Fields[0]
	if x, ok := f.Value.(Id); !ok || x != 6 {
		t.Fatal("wrong value", x)
	}
	if !f.PrimaryKey() {
		t.Fatal("wrong value")
	}
	f = m.Fields[1]
	if x, ok := f.Value.(string); !ok || x != "banana" {
		t.Fatal("wrong value", x)
	}
	if !f.PrimaryKey() {
		t.Fatal("wrong value")
	}
	f = m.Fields[2]
	if x, ok := f.Value.(string); !ok || x != "" {
		t.Fatal("wrong value", x)
	}
	if f.Default() != "'banana'" {
		t.Fatal("should value", f.Default())
	}
	if !f.NotNull() {
		t.Fatal("wrong value")
	}
	f = m.Fields[3]
	if x, ok := f.Value.(string); !ok || x != "orange" {
		t.Fatal("wrong value", x)
	}
	if x := f.Size(); x != 64 {
		t.Fatal("wrong value", x)
	}
	f = m.Fields[4]
	if x, ok := f.Value.(time.Time); !ok || !now.Equal(x) {
		t.Fatal("wrong value", x)
	}
}

func makeWhitespaceVisible(s string) string {
	s = strings.Replace(s, "\t", "\\t", -1)
	s = strings.Replace(s, "\r\n", "\\r\\n", -1)
	s = strings.Replace(s, "\r", "\\r", -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}

type TestSchemaGenerationUserTable struct {
	Id    Id
	First string `sql:"size(30)"`
	Last  string
}

func (table *TestSchemaGenerationUserTable) Indexes(indexes *Indexes) {
	indexes.AddUnique("name_index", "first", "last")
}

func TestSchemaGeneration(t *testing.T) {
	hd := Dry()
	if x := len(hd.schema); x != 0 {
		t.Fatal("invalid schema state", x)
	}

	hd.CreateTable(&TestSchemaGenerationUserTable{})
	decl1 := "type TestSchemaGenerationUserTable struct {\n" +
		"\tId\thood.Id\n" +
		"\tFirst\tstring\t`sql:\"size(30)\"`\n" +
		"\tLast\tstring\n" +
		"}\n" +
		"\n" +
		"func (table *TestSchemaGenerationUserTable) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl1 {
		t.Fatalf("invalid schema\n%s\n---\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl1))
	}
	type DropMe struct {
		Id Id
	}
	hd.CreateTable(&DropMe{})
	decl2 := "type TestSchemaGenerationUserTable struct {\n" +
		"\tId\thood.Id\n" +
		"\tFirst\tstring\t`sql:\"size(30)\"`\n" +
		"\tLast\tstring\n" +
		"}\n" +
		"\n" +
		"func (table *TestSchemaGenerationUserTable) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"}\n" +
		"\n" +
		"type DropMe struct {\n" +
		"\tId\thood.Id\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl2 {
		t.Fatalf("invalid schema\n%s\n---\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl2))
	}
	hd.DropTable(&DropMe{})
	if x := hd.schema.GoDeclaration(); x != decl1 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl1))
	}
	hd.RenameTable(&TestSchemaGenerationUserTable{}, "customers")
	decl3 := "type Customers struct {\n" +
		"\tId\thood.Id\n" +
		"\tFirst\tstring\t`sql:\"size(30)\"`\n" +
		"\tLast\tstring\n" +
		"}\n" +
		"\n" +
		"func (table *Customers) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl3 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl3))
	}
	hd.AddColumns("customers", struct {
		Balance int
	}{})
	decl4 := "type Customers struct {\n" +
		"\tId\thood.Id\n" +
		"\tFirst\tstring\t`sql:\"size(30)\"`\n" +
		"\tLast\tstring\n" +
		"\tBalance\tint\n" +
		"}\n" +
		"\n" +
		"func (table *Customers) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl4 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl4))
	}
	hd.RenameColumn("customers", "balance", "amount")
	decl5 := "type Customers struct {\n" +
		"\tId\thood.Id\n" +
		"\tFirst\tstring\t`sql:\"size(30)\"`\n" +
		"\tLast\tstring\n" +
		"\tAmount\tint\n" +
		"}\n" +
		"\n" +
		"func (table *Customers) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl5 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl5))
	}
	hd.ChangeColumns("customers", struct {
		Amount string
	}{})
	decl6 := "type Customers struct {\n" +
		"\tId\thood.Id\n" +
		"\tFirst\tstring\t`sql:\"size(30)\"`\n" +
		"\tLast\tstring\n" +
		"\tAmount\tstring\n" +
		"}\n" +
		"\n" +
		"func (table *Customers) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl6 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl6))
	}
	hd.RemoveColumns("customers", struct {
		First string
		Last  string
	}{})
	decl7 := "type Customers struct {\n" +
		"\tId\thood.Id\n" +
		"\tAmount\tstring\n" +
		"}\n" +
		"\n" +
		"func (table *Customers) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl7 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl7))
	}
	hd.CreateIndex("customers", "amount_index", false, "amount")
	decl8 := "type Customers struct {\n" +
		"\tId\thood.Id\n" +
		"\tAmount\tstring\n" +
		"}\n" +
		"\n" +
		"func (table *Customers) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.AddUnique(\"name_index\", \"first\", \"last\")\n" +
		"\tindexes.Add(\"amount_index\", \"amount\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl8 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl8))
	}
	hd.DropIndex("customers", "name_index")
	decl9 := "type Customers struct {\n" +
		"\tId\thood.Id\n" +
		"\tAmount\tstring\n" +
		"}\n" +
		"\n" +
		"func (table *Customers) Indexes(indexes *hood.Indexes) {\n" +
		"\tindexes.Add(\"amount_index\", \"amount\")\n" +
		"}"
	if x := hd.schema.GoDeclaration(); x != decl9 {
		t.Fatalf("invalid schema\n%s\n\n%s", makeWhitespaceVisible(x), makeWhitespaceVisible(decl9))
	}
}

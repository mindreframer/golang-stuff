package hood

import (
	"database/sql"
	"testing"
	"time"
)

// THE LIVE TESTS ARE DISABLED BY DEFAULT, NOT TO INTERFERE WITH
// REAL UNIT TESTS, SINCE THEY DO REQUIRE A CERTAIN SYSTEM CONFIGURATION!
//
// ONLY ENABLE THE LIVE TESTS IF NECESSARY
//
// TO ENABLE/DISABLE LIVE TESTS UNCOMMENT/COMMENT THE CORRESPONDING DIALECT
// INFO IN THE TO_RUN ARRAY!

import (
	_ "github.com/lib/pq"
	_ "github.com/ziutek/mymysql/godrv"
)

var toRun = []dialectInfo{
// allDialectInfos[0],
// allDialectInfos[1],
}

var allDialectInfos = []dialectInfo{
	dialectInfo{
		NewPostgres(),
		setupPgDb,
		`CREATE TABLE "without_pk" ( "first" text, "last" text, "amount" integer )`,
		`CREATE TABLE IF NOT EXISTS "without_pk" ( "first" text, "last" text, "amount" integer )`,
		`CREATE TABLE "with_pk" ( "primary" bigserial PRIMARY KEY, "first" text, "last" text, "amount" integer )`,
		`INSERT INTO "sql_gen_model" ("first", "last", "amount") VALUES ($1, $2, $3) RETURNING "prim"`,
		`UPDATE "sql_gen_model" SET "first" = $1, "last" = $2, "amount" = $3 WHERE "prim" = $4`,
		`DELETE FROM "sql_gen_model" WHERE "prim" = $1`,
		`DELETE FROM "sql_del_from" WHERE "a" = $1 AND "b" > $2 OR "c" < $3`,
		`SELECT * FROM "sql_gen_model"`,
		`SELECT "col1", "col2" FROM "sql_gen_model" INNER JOIN "orders" ON "sql_gen_model"."id1" = "orders"."id2" WHERE "user"."id" = "order"."id" AND "a" > $1 OR "b" < $2 AND "c" = $3 OR "d" = $4 GROUP BY "user"."name" HAVING SUM(price) < $5 ORDER BY "user"."first_name" LIMIT $6 OFFSET $7`,
		`SELECT "col1", "col2" FROM "sql_gen_model" INNER JOIN "orders" ON "sql_gen_model"."id1" = "orders"."id2" WHERE "user"."id" = "order"."id" AND "a" > $1 OR "b" < $2 AND "c" = $3 OR "d" = $4 GROUP BY "user"."name" HAVING SUM(price) < $5 ORDER BY "user"."first_name" ASC LIMIT $6 OFFSET $7`,
		`SELECT "col1", "col2" FROM "sql_gen_model" INNER JOIN "orders" ON "sql_gen_model"."id1" = "orders"."id2" WHERE "user"."id" = "order"."id" AND "a" > $1 OR "b" < $2 AND "c" = $3 OR "d" = $4 GROUP BY "user"."name" HAVING SUM(price) < $5 ORDER BY "user"."first_name" DESC LIMIT $6 OFFSET $7`,
		`DROP TABLE "drop_table"`,
		`DROP TABLE IF EXISTS "drop_table"`,
		`ALTER TABLE "table_a" RENAME TO "table_b"`,
		`ALTER TABLE "a" ADD COLUMN "c" varchar(100)`,
		`ALTER TABLE "a" RENAME COLUMN "b" TO "c"`,
		`ALTER TABLE "a" ALTER COLUMN "b" TYPE varchar(100)`,
		`ALTER TABLE "a" DROP COLUMN "b"`,
		`CREATE UNIQUE INDEX "iname" ON "itable" ("a", "b", "c")`,
		`CREATE INDEX "iname2" ON "itable2" ("d", "e")`,
		`DROP INDEX "iname"`,
	},
	dialectInfo{
		NewMysql(),
		setupMysql,
		"CREATE TABLE `without_pk` ( `first` longtext, `last` longtext, `amount` int )",
		"CREATE TABLE IF NOT EXISTS `without_pk` ( `first` longtext, `last` longtext, `amount` int )",
		"CREATE TABLE `with_pk` ( `primary` bigint PRIMARY KEY AUTO_INCREMENT, `first` longtext, `last` longtext, `amount` int )",
		"INSERT INTO `sql_gen_model` (`first`, `last`, `amount`) VALUES (?, ?, ?)",
		"UPDATE `sql_gen_model` SET `first` = ?, `last` = ?, `amount` = ? WHERE `prim` = ?",
		"DELETE FROM `sql_gen_model` WHERE `prim` = ?",
		"DELETE FROM `sql_del_from` WHERE `a` = ? AND `b` > ? OR `c` < ?",
		"SELECT * FROM `sql_gen_model`",
		"SELECT `col1`, `col2` FROM `sql_gen_model` INNER JOIN `orders` ON `sql_gen_model`.`id1` = `orders`.`id2` WHERE `user`.`id` = `order`.`id` AND `a` > ? OR `b` < ? AND `c` = ? OR `d` = ? GROUP BY `user`.`name` HAVING SUM(price) < ? ORDER BY `user`.`first_name` LIMIT ? OFFSET ?",
		"SELECT `col1`, `col2` FROM `sql_gen_model` INNER JOIN `orders` ON `sql_gen_model`.`id1` = `orders`.`id2` WHERE `user`.`id` = `order`.`id` AND `a` > ? OR `b` < ? AND `c` = ? OR `d` = ? GROUP BY `user`.`name` HAVING SUM(price) < ? ORDER BY `user`.`first_name` ASC LIMIT ? OFFSET ?",
		"SELECT `col1`, `col2` FROM `sql_gen_model` INNER JOIN `orders` ON `sql_gen_model`.`id1` = `orders`.`id2` WHERE `user`.`id` = `order`.`id` AND `a` > ? OR `b` < ? AND `c` = ? OR `d` = ? GROUP BY `user`.`name` HAVING SUM(price) < ? ORDER BY `user`.`first_name` DESC LIMIT ? OFFSET ?",
		"DROP TABLE `drop_table`",
		"DROP TABLE IF EXISTS `drop_table`",
		"ALTER TABLE `table_a` RENAME TO `table_b`",
		"ALTER TABLE `a` ADD COLUMN `c` varchar(100)",
		"ALTER TABLE `a` RENAME COLUMN `b` TO `c`",
		"ALTER TABLE `a` ALTER COLUMN `b` TYPE varchar(100)",
		"ALTER TABLE `a` DROP COLUMN `b`",
		"CREATE UNIQUE INDEX `iname` ON `itable` (`a`, `b`, `c`)",
		"CREATE INDEX `iname2` ON `itable2` (`d`, `e`)",
		"DROP INDEX `iname`",
	},
}

type dialectInfo struct {
	dialect                         Dialect
	setupDbFunc                     func(t *testing.T) *Hood
	createTableWithoutPkSql         string
	createTableWithoutPkIfExistsSql string
	createTableWithPkSql            string
	insertSql                       string
	updateSql                       string
	deleteSql                       string
	deleteFromSql                   string
	wcQuerySql                      string
	querySql                        string
	querySqlAsc                     string
	querySqlDesc                    string
	dropTableSql                    string
	dropTableIfExistsSql            string
	renameTableSql                  string
	addColumnSql                    string
	renameColumnSql                 string
	changeColumnSql                 string
	dropColumnSql                   string
	createUniqueIndexSql            string
	createIndexSql                  string
	dropIndexSql                    string
}

func setupPgDb(t *testing.T) *Hood {
	db, err := sql.Open("postgres", "user=hood dbname=hood_test sslmode=disable")
	if err != nil {
		t.Fatal("could not open db", err)
	}
	hd := New(db, NewPostgres())
	hd.Log = true
	return hd
}

func setupMysql(t *testing.T) *Hood {
	// db, err := sql.Open("mymysql", "hood_test/hood/")
	db, err := sql.Open("mymysql", "unix:/Applications/MAMP/tmp/mysql/mysql.sock*hood_test/hood/")
	if err != nil {
		t.Fatal("could not open db", err)
	}
	hd := New(db, NewMysql())
	hd.Log = true
	return hd
}

func TestTransaction(t *testing.T) {
	for _, info := range toRun {
		DoTestTransaction(t, info)
	}
}

func DoTestTransaction(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	hd := info.setupDbFunc(t)
	type txModel struct {
		Id Id
		A  string
	}
	table := txModel{
		A: "A",
	}

	hd.DropTable(&table)
	tx := hd.Begin()
	tx.CreateTable(&table)
	err := tx.Commit()
	if err != nil {
		t.Fatal("error not nil", err)
	}

	tx = hd.Begin()
	if _, ok := hd.qo.(*sql.DB); !ok {
		t.Fatal("wrong type")
	}
	if _, ok := tx.qo.(*sql.Tx); !ok {
		t.Fatal("wrong type")
	}
	_, err = tx.Save(&table)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	err = tx.Rollback()
	if err != nil {
		t.Fatal("error not nil", err)
	}

	var out []txModel
	err = hd.Find(&out)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(out); x > 0 {
		t.Fatal("wrong length", x)
	}

	tx = hd.Begin()
	table.Id = 0 // force insert by resetting id
	_, err = tx.Save(&table)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatal("error not nil", err)
	}

	out = nil
	err = hd.Find(&out)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(out); x != 1 {
		t.Fatal("wrong length", x)
	}
}

func TestSaveAndDelete(t *testing.T) {
	for _, info := range toRun {
		DoTestSaveAndDelete(t, info)
	}
}

func DoTestSaveAndDelete(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	// time identity test
	if x := time.Now(); x.Sub(x.UTC()) != 0 {
		t.Fatal("not equal")
	}
	now := time.Now()
	hd := info.setupDbFunc(t)
	type saveModel struct {
		Id      Id
		A       string
		B       int
		Updated Updated
		Created Created
	}
	model1 := saveModel{
		A: "banana",
		B: 5,
	}
	model2 := saveModel{
		A: "orange",
		B: 4,
	}

	hd.DropTable(&model1)

	tx := hd.Begin()
	tx.CreateTable(&model1)
	err := tx.Commit()
	if err != nil {
		t.Fatal("error not nil", err)
	}
	id, err := hd.Save(&model1)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if id != 1 {
		t.Fatal("wrong id", id)
	}
	if x := model1.Created; x.Sub(now) <= 0 {
		t.Fatal("wrong timestamp", x, now)
	}
	if x := model1.Updated; x.Sub(now) <= 0 {
		t.Fatal("wrong timestamp", x, now)
	}

	// make sure created/updated values match the db
	var model1r []saveModel
	err = hd.Where("id", "=", model1.Id).Find(&model1r)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(model1r); x != 1 {
		t.Fatal("wrong result count", x)
	}
	if model1r[0].Created.Unix() != model1.Created.Unix() {
		t.Fatal("created fields do not match", model1r[0].Created, model1.Created)
	}
	if model1r[0].Updated.Unix() != model1.Updated.Unix() {
		t.Fatal("updated fields do not match", model1r[0].Updated, model1.Updated)
	}

	oldCreate := model1.Created
	oldUpdate := model1.Updated
	model1.A = "grape"
	model1.B = 9

	time.Sleep(time.Second * 1) // sleep for 1 sec

	id, err = hd.Save(&model1)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if id != 1 {
		t.Fatal("wrong id", id)
	}
	if x := model1.Created; !x.Equal(oldCreate.Time) {
		t.Fatal("wrong timestamp", x)
	}
	if x := model1.Updated; x.Sub(oldUpdate.Time) <= 0 {
		t.Fatal("wrong timestamp", x, oldUpdate)
	}

	// make sure created/updated values match the db
	var model1r2 []saveModel
	err = hd.Where("id", "=", model1.Id).Find(&model1r2)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(model1r2); x != 1 {
		t.Fatal("wrong result count", x)
	}
	if x := model1r2[0].Updated; x.Sub(model1r2[0].Created.Time) < 1 {
		t.Fatal("diff mismatch", x, model1r2[0].Created.Time)
	}
	if model1r2[0].Created.Unix() != model1.Created.Unix() {
		t.Fatal("created fields do not match", model1r2[0].Created, oldCreate)
	}
	if model1r2[0].Updated.Unix() != model1.Updated.Unix() {
		t.Fatal("updated fields do not match", model1r2[0].Updated, model1.Updated)
	}

	id, err = hd.Save(&model2)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if id != 2 {
		t.Fatal("wrong id", id)
	}
	if model2.Id != id {
		t.Fatal("id should have been copied", model2.Id)
	}

	id2, err := hd.Delete(&model2)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if id != id2 {
		t.Fatal("wrong id", id, id2)
	}
}

func TestSaveDeleteAllAndHooks(t *testing.T) {
	for _, info := range toRun {
		DoTestSaveDeleteAllAndHooks(t, info)
	}
}

type sdAllModel struct {
	Id Id
	A  string
}

var sdAllHooks []string

func (m *sdAllModel) BeforeSave() error {
	sdAllHooks = append(sdAllHooks, "bsave")
	return nil
}

func (m *sdAllModel) AfterSave() error {
	sdAllHooks = append(sdAllHooks, "asave")
	return nil
}

func (m *sdAllModel) BeforeInsert() error {
	sdAllHooks = append(sdAllHooks, "binsert")
	return nil
}

func (m *sdAllModel) AfterInsert() error {
	sdAllHooks = append(sdAllHooks, "ainsert")
	return nil
}

func (m *sdAllModel) BeforeUpdate() error {
	sdAllHooks = append(sdAllHooks, "bupdate")
	return nil
}

func (m *sdAllModel) AfterUpdate() error {
	sdAllHooks = append(sdAllHooks, "aupdate")
	return nil
}

func (m *sdAllModel) BeforeDelete() error {
	sdAllHooks = append(sdAllHooks, "bdelete")
	return nil
}

func (m *sdAllModel) AfterDelete() error {
	sdAllHooks = append(sdAllHooks, "adelete")
	return nil
}

func DoTestSaveDeleteAllAndHooks(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	hd := info.setupDbFunc(t)
	hd.DropTable(&sdAllModel{})

	models := []sdAllModel{
		sdAllModel{A: "A"},
		sdAllModel{A: "B"},
	}

	sdAllHooks = make([]string, 0, 20)
	tx := hd.Begin()
	tx.CreateTable(&sdAllModel{})
	err := tx.Commit()
	if err != nil {
		t.Fatal("error not nil", err)
	}

	ids, err := hd.SaveAll(&models)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(ids); x != 2 {
		t.Fatal("wrong id count", x)
	}
	if x := ids[0]; x != 1 {
		t.Fatal("wrong id", x)
	}
	if x := ids[1]; x != 2 {
		t.Fatal("wrong id", x)
	}
	if x := models[0].Id; x != 1 {
		t.Fatal("wrong id", x)
	}
	if x := models[1].Id; x != 2 {
		t.Fatal("wrong id", x)
	}

	hd.SaveAll(&models) // force update for hooks test 

	_, err = hd.DeleteAll(&models)
	if err != nil {
		t.Fatal("error not nil", err)
	}

	if x := len(sdAllHooks); x != 20 {
		t.Fatal("wrong hook call count", x)
	}
	hookMatch := []string{
		"bsave",
		"binsert",
		"ainsert",
		"asave",
		"bsave",
		"binsert",
		"ainsert",
		"asave",
		"bsave",
		"bupdate",
		"aupdate",
		"asave",
		"bsave",
		"bupdate",
		"aupdate",
		"asave",
		"bdelete",
		"adelete",
		"bdelete",
		"adelete",
	}
	for i, v := range hookMatch {
		if x := sdAllHooks[i]; x != v {
			t.Fatal("wrong hook sequence", x, v)
		}
	}
}

func TestFind(t *testing.T) {
	for _, info := range toRun {
		DoTestFind(t, info)
	}
}

func DoTestFind(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	hd := info.setupDbFunc(t)
	now := time.Now()

	type findModel struct {
		Id Id
		A  string
		B  int
		C  int8
		D  int16
		E  int32
		F  int64
		G  uint
		H  uint8
		I  uint16
		J  uint32
		K  uint64
		L  float32
		M  float64
		N  []byte
		P  time.Time
		Q  Created
		R  Updated
	}
	model1 := findModel{
		A: "string!",
		B: -1,
		C: -2,
		D: -3,
		E: -4,
		F: -5,
		G: 6,
		H: 7,
		I: 8,
		J: 9,
		K: 10,
		L: 11.5,
		M: 12.6,
		N: []byte("bytes!"),
		P: now,
	}

	hd.DropTable(&model1)

	tx := hd.Begin()
	tx.CreateTable(&model1)
	err := tx.Commit()
	if err != nil {
		t.Fatal("error not nil", err)
	}

	var out []findModel
	err = hd.Where("a", "=", "string!").And("j", "=", 9).Find(&out)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if out != nil {
		t.Fatal("output should be nil", out)
	}

	id, err := hd.Save(&model1)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if id != 1 {
		t.Fatal("wrong id", id)
	}

	err = hd.Where("a", "=", "string!").And("j", "=", 9).Find(&out)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if out == nil {
		t.Fatal("output should not be nil")
	}
	if x := len(out); x != 1 {
		t.Fatal("invalid output length", x)
	}
	for _, v := range out {
		if x := v.Id; x != 1 {
			t.Fatal("invalid value", x)
		}
		if x := v.A; x != "string!" {
			t.Fatal("invalid value", x)
		}
		if x := v.B; x != -1 {
			t.Fatal("invalid value", x)
		}
		if x := v.C; x != -2 {
			t.Fatal("invalid value", x)
		}
		if x := v.D; x != -3 {
			t.Fatal("invalid value", x)
		}
		if x := v.E; x != -4 {
			t.Fatal("invalid value", x)
		}
		if x := v.F; x != -5 {
			t.Fatal("invalid value", x)
		}
		if x := v.G; x != 6 {
			t.Fatal("invalid value", x)
		}
		if x := v.H; x != 7 {
			t.Fatal("invalid value", x)
		}
		if x := v.I; x != 8 {
			t.Fatal("invalid value", x)
		}
		if x := v.J; x != 9 {
			t.Fatal("invalid value", x)
		}
		if x := v.K; x != 10 {
			t.Fatal("invalid value", x)
		}
		if x := v.L; x != 11.5 {
			t.Fatal("invalid value", x)
		}
		if x := v.M; x != 12.6 {
			t.Fatal("invalid value", x)
		}
		if x := v.N; string(x) != "bytes!" {
			t.Fatal("invalid value", x)
		}
		if x := v.P; now.Unix() != x.Unix() {
			t.Fatal("invalid value", x, now)
		}
	}

	model1.Id = 0 // force insert, would update otherwise
	model1.A = "row2"

	id, err = hd.Save(&model1)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if id != 2 {
		t.Fatal("wrong id", id)
	}

	out = nil
	err = hd.Where("a", "=", "row2").And("j", "=", 9).Find(&out)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(out); x != 1 {
		t.Fatal("invalid output length", x)
	}

	out = nil
	err = hd.Where("j", "=", 9).Find(&out)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(out); x != 2 {
		t.Fatal("invalid output length", x)
	}
}

func TestCreateTable(t *testing.T) {
	for _, info := range toRun {
		DoTestCreateTable(t, info)
	}
}

type CreateTableTestModel struct {
	Prim   Id
	First  string `sql:"size(64),notnull"`
	Last   string `sql:"size(128),default('defaultValue')"`
	Amount int
}

func (table *CreateTableTestModel) Indexes(indexes *Indexes) {
	indexes.AddUnique("create_table_test_model_index", "first", "last")
}

func DoTestCreateTable(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	hd := info.setupDbFunc(t)
	table := &CreateTableTestModel{}
	err := hd.DropTableIfExists(table)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	tx := hd.Begin()
	tx.CreateTable(table)
	err = tx.Commit()
	if err != nil {
		t.Fatal("error not nil", err)
	}
	err = hd.DropTable(table)
	if err != nil {
		t.Fatal("error not nil", err)
	}
}

func TestCreateTableSql(t *testing.T) {
	for _, info := range toRun {
		DoTestCreateTableSql(t, info)
	}
}

func DoTestCreateTableSql(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	type withoutPk struct {
		First  string
		Last   string
		Amount int
	}
	table := &withoutPk{"a", "b", 5}
	model, err := interfaceToModel(table)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := info.dialect.CreateTableSql(model, false); x != info.createTableWithoutPkSql {
		t.Fatal("wrong sql", x)
	}
	if x := info.dialect.CreateTableSql(model, true); x != info.createTableWithoutPkIfExistsSql {
		t.Fatal("wrong sql", x)
	}
	type withPk struct {
		Primary Id
		First   string
		Last    string
		Amount  int
	}
	table2 := &withPk{First: "a", Last: "b", Amount: 5}
	model, err = interfaceToModel(table2)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := info.dialect.CreateTableSql(model, false); x != info.createTableWithPkSql {
		t.Fatal("wrong query", x)
	}
}

type sqlGenModel struct {
	Prim   Id
	First  string
	Last   string
	Amount int
}

var sqlGenSampleData = &sqlGenModel{3, "FirstName", "LastName", 6}

func TestInsertSQL(t *testing.T) {
	for _, info := range toRun {
		DoTestInsertSQL(t, info)
	}
}

func DoTestInsertSQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	model, _ := interfaceToModel(sqlGenSampleData)
	sql, _ := info.dialect.InsertSql(model)
	if x := info.insertSql; x != sql {
		t.Log(sql)
		t.Log(x)
		t.Fatal("invalid sql")
	}
}

func TestUpdateSQL(t *testing.T) {
	for _, info := range toRun {
		DoTestUpdateSQL(t, info)
	}
}

func DoTestUpdateSQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	model, _ := interfaceToModel(sqlGenSampleData)
	sql, _ := info.dialect.UpdateSql(model)
	if x := info.updateSql; x != sql {
		t.Log(sql)
		t.Log(x)
		t.Fatal("invalid sql")
	}
}

func TestDeleteSQL(t *testing.T) {
	for _, info := range toRun {
		DoTestDeleteSQL(t, info)
	}
}

func DoTestDeleteSQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	model, _ := interfaceToModel(sqlGenSampleData)
	sql, _ := info.dialect.DeleteSql(model)
	if x := info.deleteSql; x != sql {
		t.Log(sql)
		t.Log(x)
		t.Fatal("invalid sql")
	}
}

func TestDeleteFromSQL(t *testing.T) {
	for _, info := range toRun {
		DoTestDeleteFromSQL(t, info)
	}
}

func DoTestDeleteFromSQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	hd := info.setupDbFunc(t)
	hd.Where("a", "=", 2).And("b", ">", 3).Or("c", "<", 4)

	sql, args := info.dialect.DeleteFromSql(hd, "sql_del_from")
	if x := info.deleteFromSql; x != sql {
		t.Log(sql)
		t.Log(x)
		t.Fatal("invalid sql")
	}
	if len(args) != 3 {
		t.Log(args)
		t.Fatal("invalid args")
	}
}

func TestQuerySQL(t *testing.T) {
	for _, info := range toRun {
		DoTestQuerySQL(t, info)
	}
}

func DoTestQuerySQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	hood := New(nil, info.dialect)
	hood.Select(&sqlGenModel{})
	query, _ := hood.Dialect.QuerySql(hood)
	if x := info.wcQuerySql; x != query {
		t.Log(query)
		t.Log(x)
		t.Fatal("invalid query", query, x)
	}

	hood = New(nil, info.dialect)
	hood.Select(&sqlGenModel{}, "col1", "col2")
	hood.Where("user.id", "=", Path("order.id"))
	hood.And("a", ">", 4)
	hood.Or("b", "<", 5)
	hood.And("c", "=", 6)
	hood.Or("d", "=", 7)
	hood.Join(InnerJoin, "orders", "sql_gen_model.id1", "orders.id2")
	hood.GroupBy("user.name")
	hood.Having("SUM(price) < ?", 2000)
	hood.OrderBy("user.first_name")
	hood.Offset(3)
	hood.Limit(10)

	hoodDesc := hood.Copy()
	hoodAsc := hood.Copy()

	// TODO: verify 2nd argument ARGS
	// Without ASC/DESC
	query, _ = hood.Dialect.QuerySql(hood)
	if x := info.querySql; x != query {
		t.Fatalf("invalid query:\n%s\n---should be---\n%s\n", x, query)
	}

	// With DESC
	hoodDesc.Desc()
	query, _ = hoodDesc.Dialect.QuerySql(hoodDesc)
	if x := info.querySqlDesc; x != query {
		t.Fatalf("invalid query:\n%s\n---should be---\n%s\n", x, query)
	}

	// With ASC
	hoodAsc.Asc()
	query, _ = hoodAsc.Dialect.QuerySql(hoodAsc)
	if x := info.querySqlAsc; x != query {
		t.Fatalf("invalid query:\n%s\n---should be---\n%s\n", x, query)
	}
}

func TestDropTableSQL(t *testing.T) {
	for _, info := range toRun {
		DoTestDropTableSQL(t, info)
	}
}

func DoTestDropTableSQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.DropTableSql("drop_table", false); x != info.dropTableSql {
		t.Fatal("wrong sql", x)
	}
	if x := info.dialect.DropTableSql("drop_table", true); x != info.dropTableIfExistsSql {
		t.Fatal("wrong sql", x)
	}
}

func TestRenameTableSQL(t *testing.T) {
	for _, info := range toRun {
		DoTestRenameTableSQL(t, info)
	}
}

func DoTestRenameTableSQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.RenameTableSql("table_a", "table_b"); x != info.renameTableSql {
		t.Fatal("wrong sql", x)
	}
}

func TestAddColumSQL(t *testing.T) {
	for _, info := range toRun {
		DoTestAddColumSQL(t, info)
	}
}

func DoTestAddColumSQL(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.AddColumnSql("a", "c", "", 100); x != info.addColumnSql {
		t.Fatal("wrong sql", x)
	}
}

func TestRenameColumnSql(t *testing.T) {
	for _, info := range toRun {
		DoTestRenameColumnSql(t, info)
	}
}

func DoTestRenameColumnSql(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.RenameColumnSql("a", "b", "c"); x != info.renameColumnSql {
		t.Fatal("wrong sql", x)
	}
}

func TestChangeColumnSql(t *testing.T) {
	for _, info := range toRun {
		DoTestChangeColumnSql(t, info)
	}
}

func DoTestChangeColumnSql(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.ChangeColumnSql("a", "b", "", 100); x != info.changeColumnSql {
		t.Fatal("wrong sql", x)
	}
}

func TestRemoveColumnSql(t *testing.T) {
	for _, info := range toRun {
		DoTestRemoveColumnSql(t, info)
	}
}

func DoTestRemoveColumnSql(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.DropColumnSql("a", "b"); x != info.dropColumnSql {
		t.Fatal("wrong sql", x)
	}
}

func TestCreateIndexSql(t *testing.T) {
	for _, info := range toRun {
		DoTestCreateIndexSql(t, info)
	}
}

func DoTestCreateIndexSql(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.CreateIndexSql("iname", "itable", true, "a", "b", "c"); x != info.createUniqueIndexSql {
		t.Fatal("wrong sql", x)
	}
	if x := info.dialect.CreateIndexSql("iname2", "itable2", false, "d", "e"); x != info.createIndexSql {
		t.Fatal("wrong sql", x)
	}
}

func TestDropIndexSql(t *testing.T) {
	for _, info := range toRun {
		DoTestDropIndexSql(t, info)
	}
}

func DoTestDropIndexSql(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	if x := info.dialect.DropIndexSql("iname"); x != info.dropIndexSql {
		t.Fatal("wrong sql", x)
	}
}

func TestNullValues(t *testing.T) {
	for _, info := range toRun {
		DoTestNullValues(t, info)
	}
}

func DoTestNullValues(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	type nullModel struct {
		Id Id
		A  string
	}
	hd := info.setupDbFunc(t)
	err := hd.DropTableIfExists(&nullModel{})
	if err != nil {
		t.Fatal(err.Error())
	}
	tx := hd.Begin()
	tx.CreateTable(&nullModel{})
	err = tx.Commit()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = hd.Exec("INSERT INTO null_model (A) VALUES (NULL)")
	if err != nil {
		t.Fatal(err.Error())
	}
	var out []nullModel
	err = hd.Find(&out)
	if err != nil {
		t.Fatal(err.Error())
	}
	if x := len(out); x != 1 {
		t.Fatal("should return 1 entry, has", x)
	}
	if x := out[0].A; x != "" {
		t.Fatal("A should be empty (NULL)", x)
	}
}

func TestCommitLastError(t *testing.T) {
	for _, info := range toRun {
		DoTestCommitLastError(t, info)
	}
}

func DoTestCommitLastError(t *testing.T, info dialectInfo) {
	t.Logf("Dialect %T\n", info.dialect)
	type commitClash struct {
		Id Id
	}
	hd := info.setupDbFunc(t)
	tx := hd.Begin()
	if !tx.IsTransaction() {
		t.Fatal("should be a transaction")
	}
	tx.CreateTable(&commitClash{})
	tx.CreateTable(&commitClash{})
	err := tx.Commit()
	if tx.firstTxError == nil {
		t.Fatal("tx error should be set")
	}
	if err == nil {
		t.Fatal("should return error")
	}
}

func TestSqlTypeForPgDialect(t *testing.T) {
	d := NewPostgres()
	if x := d.SqlType(true, 0); x != "boolean" {
		t.Fatal("wrong type", x)
	}
	var indirect interface{} = true
	if x := d.SqlType(indirect, 0); x != "boolean" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(uint32(2), 0); x != "integer" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(Id(1), 0); x != "bigserial" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(int64(1), 0); x != "bigint" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(1.8, 0); x != "double precision" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType([]byte("asdf"), 0); x != "bytea" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType("astring", 0); x != "text" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType("a", 255); x != "varchar(255)" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType("b", 128); x != "varchar(128)" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(time.Now(), 0); x != "timestamp with time zone" {
		t.Fatal("wrong type", x)
	}
}

func TestSqlTypeForMysqlDialect(t *testing.T) {
	d := NewMysql()
	if x := d.SqlType(true, 0); x != "boolean" {
		t.Fatal("wrong type", x)
	}
	var indirect interface{} = true
	if x := d.SqlType(indirect, 0); x != "boolean" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(uint32(2), 0); x != "int" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(Id(1), 0); x != "bigint" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(int64(1), 0); x != "bigint" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(1.8, 0); x != "double" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType([]byte("asdf"), 0); x != "longblob" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType("astring", 0); x != "longtext" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType("a", 65536); x != "longtext" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType("b", 128); x != "varchar(128)" {
		t.Fatal("wrong type", x)
	}
	if x := d.SqlType(time.Now(), 0); x != "timestamp" {
		t.Fatal("wrong type", x)
	}
}

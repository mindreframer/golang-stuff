package beedb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var OnDebug = false
var PluralizeTableNames = false

type Model struct {
	Db              *sql.DB
	TableName       string
	LimitStr        int
	OffsetStr       int
	WhereStr        string
	ParamStr        []interface{}
	OrderStr        string
	ColumnStr       string
	PrimaryKey      string
	JoinStr         string
	GroupByStr      string
	HavingStr       string
	QuoteIdentifier string
	ParamIdentifier string
	ParamIteration  int
}

/**
 * Add New sql.DB in the future i will add ConnectionPool.Get() 
 */
func New(db *sql.DB, options ...interface{}) (m Model) {
	if len(options) == 0 {
		m = Model{Db: db, ColumnStr: "*", PrimaryKey: "Id", QuoteIdentifier: "`", ParamIdentifier: "?", ParamIteration: 1}
	} else if options[0] == "pg" {
		m = Model{Db: db, ColumnStr: "id", PrimaryKey: "id", QuoteIdentifier: "\"", ParamIdentifier: options[0].(string), ParamIteration: 1}
	} else if options[0] == "mssql" {
		m = Model{Db: db, ColumnStr: "id", PrimaryKey: "id", QuoteIdentifier: "", ParamIdentifier: options[0].(string), ParamIteration: 1}
	}
	return
}

func (orm *Model) SetTable(tbname string) *Model {
	orm.TableName = tbname
	return orm
}

func (orm *Model) SetPK(pk string) *Model {
	orm.PrimaryKey = pk
	return orm
}

func (orm *Model) Where(querystring interface{}, args ...interface{}) *Model {
	switch querystring := querystring.(type) {
	case string:
		orm.WhereStr = querystring
	case int:
		if orm.ParamIdentifier == "pg" {
			orm.WhereStr = fmt.Sprintf("%v%v%v = $%v", orm.QuoteIdentifier, orm.PrimaryKey, orm.QuoteIdentifier, orm.ParamIteration)
		} else {
			orm.WhereStr = fmt.Sprintf("%v%v%v = ?", orm.QuoteIdentifier, orm.PrimaryKey, orm.QuoteIdentifier)
			orm.ParamIteration++
		}
		args = append(args, querystring)
	}
	orm.ParamStr = args
	return orm
}

func (orm *Model) Limit(start int, size ...int) *Model {
	orm.LimitStr = start
	if len(size) > 0 {
		orm.OffsetStr = size[0]
	}
	return orm
}

func (orm *Model) Offset(offset int) *Model {
	orm.OffsetStr = offset
	return orm
}

func (orm *Model) OrderBy(order string) *Model {
	orm.OrderStr = order
	return orm
}

func (orm *Model) Select(colums string) *Model {
	orm.ColumnStr = colums
	return orm
}

func (orm *Model) ScanPK(output interface{}) *Model {
	if reflect.TypeOf(reflect.Indirect(reflect.ValueOf(output)).Interface()).Kind() == reflect.Slice {
		sliceValue := reflect.Indirect(reflect.ValueOf(output))
		sliceElementType := sliceValue.Type().Elem()
		for i := 0; i < sliceElementType.NumField(); i++ {
			bb := sliceElementType.Field(i).Tag
			if bb.Get("beedb") == "PK" || reflect.ValueOf(bb).String() == "PK" {
				orm.PrimaryKey = sliceElementType.Field(i).Name
			}
		}
	} else {
		tt := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(output)).Interface())
		for i := 0; i < tt.NumField(); i++ {
			bb := tt.Field(i).Tag
			if bb.Get("beedb") == "PK" || reflect.ValueOf(bb).String() == "PK" {
				orm.PrimaryKey = tt.Field(i).Name
			}
		}
	}
	return orm

}

//The join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (orm *Model) Join(join_operator, tablename, condition string) *Model {
	if orm.JoinStr != "" {
		orm.JoinStr = orm.JoinStr + fmt.Sprintf(" %v JOIN %v ON %v", join_operator, tablename, condition)
	} else {
		orm.JoinStr = fmt.Sprintf("%v JOIN %v ON %v", join_operator, tablename, condition)
	}

	return orm
}

func (orm *Model) GroupBy(keys string) *Model {
	orm.GroupByStr = fmt.Sprintf("GROUP BY %v", keys)
	return orm
}

func (orm *Model) Having(conditions string) *Model {
	orm.HavingStr = fmt.Sprintf("HAVING %v", conditions)
	return orm
}

func (orm *Model) Find(output interface{}) error {
	orm.ScanPK(output)
	var keys []string
	results, err := scanStructIntoMap(output)
	if err != nil {
		return err
	}

	if orm.TableName == "" {
		orm.TableName = getTableName(StructName(output))
	}
	for key, _ := range results {
		keys = append(keys, key)
	}
	orm.ColumnStr = strings.Join(keys, ", ")
	orm.Limit(1)
	resultsSlice, err := orm.FindMap()
	if err != nil {
		return err
	}
	if len(resultsSlice) == 0 {
		return errors.New("No record found")
	} else if len(resultsSlice) == 1 {
		results := resultsSlice[0]
		err := scanMapIntoStruct(output, results)
		if err != nil {
			return err
		}
	} else {
		return errors.New("More than one record")
	}
	return nil
}

func (orm *Model) FindAll(rowsSlicePtr interface{}) error {
	orm.ScanPK(rowsSlicePtr)
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if sliceValue.Kind() != reflect.Slice {
		return errors.New("needs a pointer to a slice")
	}

	sliceElementType := sliceValue.Type().Elem()
	st := reflect.New(sliceElementType)
	var keys []string
	results, err := scanStructIntoMap(st.Interface())
	if err != nil {
		return err
	}

	if orm.TableName == "" {
		orm.TableName = getTableName(getTypeName(rowsSlicePtr))
	}
	for key, _ := range results {
		keys = append(keys, key)
	}
	orm.ColumnStr = strings.Join(keys, ", ")

	resultsSlice, err := orm.FindMap()
	if err != nil {
		return err
	}

	for _, results := range resultsSlice {
		newValue := reflect.New(sliceElementType)
		err := scanMapIntoStruct(newValue.Interface(), results)
		if err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(newValue.Interface()))))
	}
	return nil
}

func (orm *Model) FindMap() (resultsSlice []map[string][]byte, err error) {
	defer orm.InitModel()
	sqls := orm.generateSql()
	if OnDebug {
		fmt.Println(sqls)
		fmt.Println(orm)
	}
	s, err := orm.Db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	res, err := s.Query(orm.ParamStr...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	fields, err := res.Columns()
	if err != nil {
		return nil, err
	}
	for res.Next() {
		result := make(map[string][]byte)
		var scanResultContainers []interface{}
		for i := 0; i < len(fields); i++ {
			var scanResultContainer interface{}
			scanResultContainers = append(scanResultContainers, &scanResultContainer)
		}
		if err := res.Scan(scanResultContainers...); err != nil {
			return nil, err
		}
		for ii, key := range fields {
			rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
			//if row is null then ignore
			if rawValue.Interface() == nil {
				continue
			}
			aa := reflect.TypeOf(rawValue.Interface())
			vv := reflect.ValueOf(rawValue.Interface())
			var str string
			switch aa.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				str = strconv.FormatInt(vv.Int(), 10)
				result[key] = []byte(str)
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				str = strconv.FormatUint(vv.Uint(), 10)
				result[key] = []byte(str)
			case reflect.Float32, reflect.Float64:
				str = strconv.FormatFloat(vv.Float(), 'f', -1, 64)
				result[key] = []byte(str)
			case reflect.Slice:
				if aa.Elem().Kind() == reflect.Uint8 {
					result[key] = rawValue.Interface().([]byte)
					break
				}
			case reflect.String:
				str = vv.String()
				result[key] = []byte(str)
			//时间类型	
			case reflect.Struct:
				str = rawValue.Interface().(time.Time).Format("2006-01-02 15:04:05.000 -0700")
				result[key] = []byte(str)
			}

		}
		resultsSlice = append(resultsSlice, result)
	}
	return resultsSlice, nil
}

func (orm *Model) generateSql() (a string) {
	if orm.ParamIdentifier == "mssql" {
		if orm.OffsetStr > 0 {
			a = fmt.Sprintf("select ROW_NUMBER() OVER(order by %v )as rownum,%v from %v",
				orm.PrimaryKey,
				orm.ColumnStr,
				orm.TableName)
			if orm.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, orm.WhereStr)
			}
			a = fmt.Sprintf("select * from (%v) "+
				"as a where rownum between %v and %v",
				a,
				orm.OffsetStr,
				orm.LimitStr)
		} else if orm.LimitStr > 0 {
			a = fmt.Sprintf("SELECT top %v %v FROM %v", orm.LimitStr, orm.ColumnStr, orm.TableName)
			if orm.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, orm.WhereStr)
			}
			if orm.GroupByStr != "" {
				a = fmt.Sprintf("%v %v", a, orm.GroupByStr)
			}
			if orm.HavingStr != "" {
				a = fmt.Sprintf("%v %v", a, orm.HavingStr)
			}
			if orm.OrderStr != "" {
				a = fmt.Sprintf("%v ORDER BY %v", a, orm.OrderStr)
			}
		} else {
			a = fmt.Sprintf("SELECT %v FROM %v", orm.ColumnStr, orm.TableName)
			if orm.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, orm.WhereStr)
			}
			if orm.GroupByStr != "" {
				a = fmt.Sprintf("%v %v", a, orm.GroupByStr)
			}
			if orm.HavingStr != "" {
				a = fmt.Sprintf("%v %v", a, orm.HavingStr)
			}
			if orm.OrderStr != "" {
				a = fmt.Sprintf("%v ORDER BY %v", a, orm.OrderStr)
			}
		}
	} else {
		a = fmt.Sprintf("SELECT %v FROM %v", orm.ColumnStr, orm.TableName)
		if orm.JoinStr != "" {
			a = fmt.Sprintf("%v %v", a, orm.JoinStr)
		}
		if orm.WhereStr != "" {
			a = fmt.Sprintf("%v WHERE %v", a, orm.WhereStr)
		}
		if orm.GroupByStr != "" {
			a = fmt.Sprintf("%v %v", a, orm.GroupByStr)
		}
		if orm.HavingStr != "" {
			a = fmt.Sprintf("%v %v", a, orm.HavingStr)
		}
		if orm.OrderStr != "" {
			a = fmt.Sprintf("%v ORDER BY %v", a, orm.OrderStr)
		}
		if orm.OffsetStr > 0 {
			a = fmt.Sprintf("%v LIMIT %v, %v", a, orm.OffsetStr, orm.LimitStr)
		} else if orm.LimitStr > 0 {
			a = fmt.Sprintf("%v LIMIT %v", a, orm.LimitStr)
		}
	}
	return
}

//Execute sql
func (orm *Model) Exec(finalQueryString string, args ...interface{}) (sql.Result, error) {
	rs, err := orm.Db.Prepare(finalQueryString)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	res, err := rs.Exec(args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//if the struct has PrimaryKey == 0 insert else update
func (orm *Model) Save(output interface{}) error {
	orm.ScanPK(output)
	results, err := scanStructIntoMap(output)
	if err != nil {
		return err
	}

	if orm.TableName == "" {
		orm.TableName = getTableName(StructName(output))
	}
	id := results[snakeCasedName(orm.PrimaryKey)]
	delete(results, snakeCasedName(orm.PrimaryKey))
	if reflect.ValueOf(id).Int() == 0 {
		structPtr := reflect.ValueOf(output)
		structVal := structPtr.Elem()
		structField := structVal.FieldByName(orm.PrimaryKey)
		id, err := orm.Insert(results)
		if err != nil {
			return err
		}
		var v interface{}
		x, err := strconv.Atoi(strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}
		v = x
		structField.Set(reflect.ValueOf(v))
		return nil
	} else {
		var condition string
		if orm.ParamIdentifier == "pg" {
			condition = fmt.Sprintf("%v%v%v=$%v", orm.QuoteIdentifier, strings.ToLower(orm.PrimaryKey), orm.QuoteIdentifier, orm.ParamIteration)
		} else {
			condition = fmt.Sprintf("%v%v%v=?", orm.QuoteIdentifier, orm.PrimaryKey, orm.QuoteIdentifier)
		}
		orm.Where(condition, id)
		_, err := orm.Update(results)
		if err != nil {
			return err
		}
	}
	return nil
}

//inert one info
func (orm *Model) Insert(properties map[string]interface{}) (int64, error) {
	defer orm.InitModel()
	var keys []string
	var placeholders []string
	var args []interface{}
	for key, val := range properties {
		keys = append(keys, key)
		if orm.ParamIdentifier == "pg" {
			ds := fmt.Sprintf("$%d", orm.ParamIteration)
			placeholders = append(placeholders, ds)
		} else {
			placeholders = append(placeholders, "?")
		}
		orm.ParamIteration++
		args = append(args, val)
	}
	ss := fmt.Sprintf("%v,%v", orm.QuoteIdentifier, orm.QuoteIdentifier)
	statement := fmt.Sprintf("INSERT INTO %v%v%v (%v%v%v) VALUES (%v)",
		orm.QuoteIdentifier,
		orm.TableName,
		orm.QuoteIdentifier,
		orm.QuoteIdentifier,
		strings.Join(keys, ss),
		orm.QuoteIdentifier,
		strings.Join(placeholders, ", "))
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(orm)
	}
	if orm.ParamIdentifier == "pg" {
		statement = fmt.Sprintf("%v RETURNING %v", statement, snakeCasedName(orm.PrimaryKey))
		var id int64
		orm.Db.QueryRow(statement, args...).Scan(&id)
		return id, nil
	} else {
		res, err := orm.Exec(statement, args...)
		if err != nil {
			return -1, err
		}

		id, err := res.LastInsertId()

		if err != nil {
			return -1, err
		}
		return id, nil
	}
	return -1, nil
}

//insert batch info
func (orm *Model) InsertBatch(rows []map[string]interface{}) ([]int64, error) {
	var ids []int64
	tablename := orm.TableName
	if len(rows) <= 0 {
		return ids, nil
	}
	for i := 0; i < len(rows); i++ {
		orm.TableName = tablename
		id, err := orm.Insert(rows[i])
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}
	return ids, nil
}

// update info
func (orm *Model) Update(properties map[string]interface{}) (int64, error) {
	defer orm.InitModel()
	var updates []string
	var args []interface{}
	for key, val := range properties {
		if orm.ParamIdentifier == "pg" {
			ds := fmt.Sprintf("$%d", orm.ParamIteration)
			updates = append(updates, fmt.Sprintf("%v%v%v = %v", orm.QuoteIdentifier, key, orm.QuoteIdentifier, ds))
		} else {
			updates = append(updates, fmt.Sprintf("%v%v%v = ?", orm.QuoteIdentifier, key, orm.QuoteIdentifier))
		}
		args = append(args, val)
		orm.ParamIteration++
	}
	args = append(args, orm.ParamStr...)
	if orm.ParamIdentifier == "pg" {
		if n := len(orm.ParamStr); n > 0 {
			for i := 1; i <= n; i++ {
				orm.WhereStr = strings.Replace(orm.WhereStr, "$"+strconv.Itoa(i), "$"+strconv.Itoa(orm.ParamIteration), 1)
			}
		}
	}
	var condition string
	if orm.WhereStr != "" {
		condition = fmt.Sprintf("WHERE %v", orm.WhereStr)
	} else {
		condition = ""
	}
	statement := fmt.Sprintf("UPDATE %v%v%v SET %v %v",
		orm.QuoteIdentifier,
		orm.TableName,
		orm.QuoteIdentifier,
		strings.Join(updates, ", "),
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(orm)
	}
	res, err := orm.Exec(statement, args...)
	if err != nil {
		return -1, err
	}
	id, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}
	return id, nil
}

func (orm *Model) Delete(output interface{}) (int64, error) {
	defer orm.InitModel()
	orm.ScanPK(output)
	results, err := scanStructIntoMap(output)
	if err != nil {
		return 0, err
	}

	if orm.TableName == "" {
		orm.TableName = getTableName(StructName(output))
	}
	id := results[strings.ToLower(orm.PrimaryKey)]
	condition := fmt.Sprintf("%v%v%v='%v'", orm.QuoteIdentifier, strings.ToLower(orm.PrimaryKey), orm.QuoteIdentifier, id)
	statement := fmt.Sprintf("DELETE FROM %v%v%v WHERE %v",
		orm.QuoteIdentifier,
		orm.TableName,
		orm.QuoteIdentifier,
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(orm)
	}
	res, err := orm.Exec(statement)
	if err != nil {
		return -1, err
	}
	Affectid, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}
	return Affectid, nil
}

func (orm *Model) DeleteAll(rowsSlicePtr interface{}) (int64, error) {
	defer orm.InitModel()
	orm.ScanPK(rowsSlicePtr)
	if orm.TableName == "" {
		orm.TableName = getTableName(getTypeName(rowsSlicePtr))
	}
	var ids []string
	val := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if val.Len() == 0 {
		return 0, nil
	}
	for i := 0; i < val.Len(); i++ {
		results, err := scanStructIntoMap(val.Index(i).Interface())
		if err != nil {
			return 0, err
		}

		id := results[strings.ToLower(orm.PrimaryKey)]
		switch id.(type) {
		case string:
			ids = append(ids, id.(string))
		case int, int64, int32:
			str := strconv.Itoa(id.(int))
			ids = append(ids, str)
		}
	}
	condition := fmt.Sprintf("%v%v%v in ('%v')", orm.QuoteIdentifier, strings.ToLower(orm.PrimaryKey), orm.QuoteIdentifier, strings.Join(ids, "','"))
	statement := fmt.Sprintf("DELETE FROM %v%v%v WHERE %v",
		orm.QuoteIdentifier,
		orm.TableName,
		orm.QuoteIdentifier,
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(orm)
	}
	res, err := orm.Exec(statement)
	if err != nil {
		return -1, err
	}
	Affectid, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}
	return Affectid, nil
}

func (orm *Model) DeleteRow() (int64, error) {
	defer orm.InitModel()
	var condition string
	if orm.WhereStr != "" {
		condition = fmt.Sprintf("WHERE %v", orm.WhereStr)
	} else {
		condition = ""
	}
	statement := fmt.Sprintf("DELETE FROM %v%v%v %v",
		orm.QuoteIdentifier,
		orm.TableName,
		orm.QuoteIdentifier,
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(orm)
	}
	res, err := orm.Exec(statement, orm.ParamStr...)
	if err != nil {
		return -1, err
	}
	Affectid, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}
	return Affectid, nil
}

func (orm *Model) InitModel() {
	orm.TableName = ""
	orm.LimitStr = 0
	orm.OffsetStr = 0
	orm.WhereStr = ""
	orm.ParamStr = make([]interface{}, 0)
	orm.OrderStr = ""
	orm.ColumnStr = "*"
	orm.PrimaryKey = "id"
	orm.JoinStr = ""
	orm.GroupByStr = ""
	orm.HavingStr = ""
	orm.ParamIteration = 1
}

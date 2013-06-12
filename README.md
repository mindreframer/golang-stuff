beedb
=====

beedb is an ORM for Go. It lets you map Go structs to tables in a database. It's intended to be very lightweight, doing very little beyond what you really want. For example, when fetching data, instead of re-inventing a query syntax, we just delegate your query to the underlying database, so you can write the "where" clause of your SQL statements directly. This allows you to have more flexibility while giving you a convenience layer. But beedb also has some smart defaults, for those times when complex queries aren't necessary.

Right now, it interfaces with Mysql/SQLite/PostgreSQL/DB2/MS ADODB/ODBC/Oracle. The goal however is to add support for other databases in the future, including maybe MongoDb or NoSQL? 

Also, at the moment, relationship-support is in the works, but not yet implemented.

All in all, it's not entirely ready for advanced use yet, but it's getting there.

Drivers for Go's sql package which support database/sql includes:

Mysql:[github.com/ziutek/mymysql/godrv](https://github.com/ziutek/mymysql/godrv)`[*]`

Mysql:[github.com/Go-SQL-Driver/MySQL](https://github.com/Go-SQL-Driver/MySQL)`[*]`

PostgreSQL:[github.com/bmizerany/pq](https://github.com/bmizerany/pq)`[*]`

SQLite:[github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)`[*]`

DB2: [bitbucket.org/phiggins/go-db2-cli](https://bitbucket.org/phiggins/go-db2-cli)

MS ADODB: [github.com/mattn/go-adodb](https://github.com/mattn/go-adodb)`[*]`

ODBC: [bitbucket.org/miquella/mgodbc](https://bitbucket.org/miquella/mgodbc)`[*]`

Oracle: [github.com/mattn/go-oci8](https://github.com/mattn/go-oci8)

Drivers marked with a `[*]` are tested with beedb

### API Interface 
[wiki/API-Interface](https://github.com/astaxie/beedb/wiki/API-Interface)

### Installing beedb
    go get github.com/astaxie/beedb

### How do we use it?

Open a database link(may be will support ConnectionPool in the future)

```go
db, err := sql.Open("mymysql", "test/xiemengjun/123456")
if err != nil {
	panic(err)
}
orm := beedb.New(db)
```

with PostgreSQL,

```go
orm := beedb.New(db, "pg")
```
	
Open Debug log, turn on the debug
  
```go
beedb.OnDebug=true
```

Model a struct after a table in the db

```go
type Userinfo struct {
	Uid		int	`beedb:"PK"` //if the table's PrimaryKey is not "Id", use this tag
	Username	string
	Departname	string
	Created		time.Time
}
```

###***Caution***
The structs Name 'UserInfo' will turn into the table name 'user_info', the same as the keyname.	
If the keyname is 'UserName' will turn into the select colum 'user_name'	
If you want table names to be pluralized so that 'UserInfo' struct was treated as 'user_infos' table, just set following option:
```go
beedb.PluralizeTableNames=true
```
	

Create an object and save it

```go
var saveone Userinfo
saveone.Username = "Test Add User"
saveone.Departname = "Test Add Departname"
saveone.Created = time.Now()
orm.Save(&saveone)
```

Saving new and existing objects

```go
saveone.Username = "Update Username"  
saveone.Departname = "Update Departname"
saveone.Created = time.Now()
orm.Save(&saveone)  //now saveone has the primarykey value it will update
```

Fetch a single object

```go
var user Userinfo
orm.Where("uid=?", 27).Find(&user)

var user2 Userinfo
orm.Where(3).Find(&user2) // this is shorthand for the version above

var user3 Userinfo
orm.Where("name = ?", "john").Find(&user3) // more complex query

var user4 Userinfo
orm.Where("name = ? and age < ?", "john", 88).Find(&user4) // even more complex
```

Fetch multiple objects

```go
var allusers []Userinfo
err := orm.Where("id > ?", "3").Limit(10,20).FindAll(&allusers) //Get id>3 limit 10 offset 20

var tenusers []Userinfo
err := orm.Where("id > ?", "3").Limit(10).FindAll(&tenusers) //Get id>3 limit 10  if omit offset the default is 0

var everyone []Userinfo
err := orm.FindAll(&everyone)
```

Find result as Map

```go
//Original SQL Backinfo resultsSlice []map[string][]byte 
//default PrimaryKey id
a, _ := orm.SetTable("userinfo").SetPK("uid").Where(2).Select("uid,username").FindMap()
```

Update with Map

```go
t := make(map[string]interface{})
var j interface{}
j = "astaxie"
t["username"] = j
//update one
orm.SetTable("userinfo").SetPK("uid").Where(2).Update(t)
```

Update batch with Map
```go
orm.SetTable("userinfo").Where("uid>?", 3).Update(t)
```

Insert data with Map	

```go
add := make(map[string]interface{})
j = "astaxie"
add["username"] = j
j = "cloud develop"
add["departname"] = j
j = "2012-12-02"
add["created"] = j
orm.SetTable("userinfo").Insert(add)
```

Insert batch with map

```go
addslice := make([]map[string]interface{})
add:=make(map[string]interface{})
add2:=make(map[string]interface{})
j = "astaxie"
add["username"] = j
j = "cloud develop"
add["departname"] = j
j = "2012-12-02"
add["created"] = j
j = "astaxie2"
add2["username"] = j
j = "cloud develop2"
add2["departname"] = j
j = "2012-12-02"
add2["created"] = j
addslice =append(addslice, add, add2)
orm.SetTable("userinfo").Insert(addslice)
```

Join Table

```go
a, _ := orm.SetTable("userinfo").Join("LEFT", "userdeatail", "userinfo.uid=userdeatail.uid").Where("userinfo.uid=?", 1).Select("userinfo.uid,userinfo.username,userdeatail.profile").FindMap()
```

Group By And Having

```go
a, _ := orm.SetTable("userinfo").GroupBy("username").Having("username='astaxie'").FindMap()
```

## LICENSE

 BSD License
 [http://creativecommons.org/licenses/BSD/](http://creativecommons.org/licenses/BSD/)

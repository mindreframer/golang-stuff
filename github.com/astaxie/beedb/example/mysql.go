package main

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/ziutek/mymysql/godrv"
	"strconv"
	"time"
)

/*
CREATE TABLE `userinfo` (
	`uid` INT(10) NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NULL,
	`departname` VARCHAR(64) NULL,
	`created` DATE NULL,
	PRIMARY KEY (`uid`)
);
CREATE TABLE `userdeatail` (
	`uid` INT(10) NULL,
	`intro` TEXT NULL,
	`profile` TEXT NULL,
	PRIMARY KEY (`uid`)
);
*/

var orm beedb.Model

type Userinfo struct {
	Uid        int `beedb:"PK"`
	Username   string
	Departname string
	Created    time.Time
}

func main() {
	db, err := sql.Open("mymysql", "test/xiemengjun/123456")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db)
	//insertbatch()
	// insert()
	// insertsql()
	// a := selectone()
	// fmt.Println(a)
	// b := selectall()
	// fmt.Println(b)
	// update()
	//updatesql()
	findmap()
	//groupby()
	//jointable()
	//delete()
	//deleteall()
	//deletesql()
}

func insert() {
	//save data
	var saveone Userinfo
	saveone.Username = "Test Add User"
	saveone.Departname = "Test Add Departname"
	saveone.Created = time.Now()
	orm.Save(&saveone)
	fmt.Println(saveone)
}

func insertsql() {
	// add one
	add := make(map[string]interface{})
	add["username"] = "astaxie"
	add["departname"] = "cloud develop"
	add["created"] = "2012-12-02"
	orm.SetTable("userinfo").Insert(add)
}

func insertbatch() {
	rows := make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		add := make(map[string]interface{})
		name := "person" + strconv.Itoa(i)
		add["username"] = name
		add["departname"] = "IT"
		add["created"] = "2012-06-13"
		rows[i] = add
	}
	fmt.Println(rows)
	orm.SetTable("userinfo").InsertBatch(rows)
}

func selectone() Userinfo {
	//get one info
	var one Userinfo
	orm.Where("uid=?", 1).Find(&one)
	return one
}

func selectall() []Userinfo {
	//get all data
	var alluser []Userinfo
	orm.Limit(10).Where("uid>?", 1).FindAll(&alluser)
	return alluser
}
func update() {
	// //update data
	var saveone Userinfo
	saveone.Uid = 1
	saveone.Username = "Update Username"
	saveone.Departname = "Update Departname"
	saveone.Created = time.Now()
	orm.Save(&saveone)
	fmt.Println(saveone)
}

func updatesql() {
	//original SQL update 
	t := make(map[string]interface{})
	t["username"] = "updateastaxie"
	//update one
	orm.SetTable("userinfo").SetPK("uid").Where(2).Update(t)
	//update batch
	//orm.SetTable("userinfo").Where("uid>?", 3).Update(t)
}

func findmap() {
	//Original SQL Backinfo resultsSlice []map[string][]byte 
	//default PrimaryKey id
	c, _ := orm.
		SetTable("userinfo").
		SetPK("uid").
		Where("username like ?", "%per%").
		Select("uid,username").
		FindMap()
	fmt.Println(c)
}

func groupby() {
	//Original SQL Group By 
	b, _ := orm.SetTable("userinfo").GroupBy("username").Having("username='updateastaxie'").FindMap()
	fmt.Println(b)
}

func jointable() {
	//Original SQL Join Table
	a, _ := orm.SetTable("userinfo").Join("LEFT", "userdeatail", "userinfo.uid=userdeatail.uid").Where("userinfo.uid=?", 1).Select("userinfo.uid,userinfo.username,userdeatail.profile").FindMap()
	fmt.Println(a)
}

func delete() {
	// // //delete one data
	saveone := selectone()
	orm.Delete(&saveone)
}

func deletesql() {
	//original SQL delete
	orm.SetTable("userinfo").Where("uid>?", 3).DeleteRow()
}

func deleteall() {
	// //delete all data
	alluser := selectall()
	orm.DeleteAll(&alluser)
}

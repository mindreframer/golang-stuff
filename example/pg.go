package main

import (
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/bmizerany/pq"
	//"time"
	"database/sql"
)

/*
CREATE TABLE userinfo
(
  uid serial NOT NULL,
  username character varying(100) NOT NULL,
  departname character varying(500) NOT NULL,
  Created date,
  CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE userdeatail
(
  uid integer,
  intro character varying(100),
  profile character varying(100)
)
WITH(OIDS=FALSE);
*/

var orm beedb.Model

type Userinfo struct {
	Uid        int `beedb:"PK"`
	Username   string
	Departname string
	Created    string
}

func main() {
	db, err := sql.Open("postgres", "user=asta password=123456 dbname=test sslmode=disable")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db, "pg")
	insert()
	//insertsql()
	// a := selectone()
	// fmt.Println(a)
	// b := selectall()
	// fmt.Println(b)
	//update()
	//updatesql()
	//findmap()
	//groupby()
	//jointable()
	//delete()
	//deleteall()
	//deletesql()
}

func insert() {
	//save data
	var saveone Userinfo
	saveone.Username = "Test_Add_User"
	saveone.Departname = "Test_Add_Departname"
	saveone.Created = "2011-12-12"
	err := orm.Save(&saveone)
	if err != nil {
		fmt.Println(err)
	}
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

func selectone() Userinfo {
	//get one info
	var one Userinfo
	orm.Where("uid=$1", 1).Find(&one)
	return one
}

func selectall() []Userinfo {
	//get all data
	var alluser []Userinfo
	orm.Limit(10).Where("uid>$1", 1).FindAll(&alluser)
	return alluser
}
func update() {
	// //update data
	var saveone Userinfo
	saveone.Uid = 1
	saveone.Username = "Update Username"
	saveone.Departname = "Update Departname"
	saveone.Created = "2012-12-02"
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
	orm.SetTable("userinfo").Where("uid>$1", 2).Update(t)
}

func findmap() {
	//Original SQL Backinfo resultsSlice []map[string][]byte 
	//default PrimaryKey id
	c, _ := orm.SetTable("userinfo").SetPK("uid").Where(2).Select("uid,username").FindMap()
	fmt.Println(c)
}

func groupby() {
	//Original SQL Group By 
	b, _ := orm.SetTable("userinfo").GroupBy("username").Having("username='updateastaxie'").Select("username").FindMap()
	fmt.Println(b)
}

func jointable() {
	//Original SQL Join Table
	a, _ := orm.SetTable("userinfo").Join("LEFT", "userdeatail", "userinfo.uid=userdeatail.uid").Where("userinfo.uid=$1", 1).Select("userinfo.uid,userinfo.username,userdeatail.profile").FindMap()
	fmt.Println(a)
}

func delete() {
	// // //delete one data
	saveone := selectone()
	orm.Delete(&saveone)
}

func deletesql() {
	//original SQL delete
	orm.SetTable("userinfo").Where("uid>$1", 3).DeleteRow()
}

func deleteall() {
	// //delete all data
	alluser := selectall()
	orm.DeleteAll(&alluser)
}

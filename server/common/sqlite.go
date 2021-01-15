package common

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func init()  {
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

//	createtb, err := db.Prepare(`CREATE TABLE userinfo(
//   uid INTEGER PRIMARY KEY AUTOINCREMENT,
//   username           TEXT    NOT NULL,
//   departname         TEXT    NOT NULL,
//   created        TEXT
//);`)
//	restab, err := createtb.Exec()
//	fmt.Println(restab,err)

	rows1, err := db.Query("SELECT name,type FROM sqlite_master ORDER BY name")
	checkErr(err)
	fmt.Println(rows1)
	for rows1.Next() {
		var name string
		var tp string
		err = rows1.Scan(&name,&tp)
		checkErr(err)
		fmt.Println(name,tp)
	}
	//插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)
	res, err := stmt.Exec("361way", "研发部", "2019-03-06")
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
	//更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)
	res, err = stmt.Exec("361wayupdate", id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	for rows.Next() {
		//fmt.Println(rows)
		list:=make([]interface{},0)
		columns,_ :=rows.Columns()
		for k,_:=range columns{
			fmt.Println(k)
			var key string
			list = append(list,&key)
		}
		for k,_:=range columns{
			var key string
			list[k]=&key
			err = rows.Scan(list...)
			fmt.Println(key)
		}

		//ColumnTypes,_ :=rows.ColumnTypes()
		//fmt.Println(columns,ColumnTypes[0].Name())
		//var uid int
		//var username string
		//var department string
		//var created string
		//list = append(list, &uid)
		//list = append(list, &username)
		//list = append(list, &department)
		//list = append(list, &created)
		//err = rows.Scan(&uid,&department,&uid,&uid)
		//err = rows.Scan(list...)
		//fmt.Println(&list[0],list[1],&list[2],&list[3])
		//fmt.Println(uid,department)
		////fmt.Println(username)
		//fmt.Println(department)
		//fmt.Println(created)
	}
	//删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)
	res, err = stmt.Exec(id)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	db.Close()
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
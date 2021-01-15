package sqllite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"strings"
	"sync"
)

type DataBase struct {
	DB    *sql.DB
	LogId int64
	sync.RWMutex
}

func NewDataBase(path string) *DataBase {
	list := strings.Split(path, "/")
	dir := strings.Join(list[0:len(list)-1], "/")
	os.Mkdir(dir, os.ModePerm)
	db, err := sql.Open("sqlite3", path)
	checkErr(err)
	c, err := db.Prepare(`CREATE TABLE binlog(
id INT PRIMARY KEY NOT NULL
    )`)
	if err == nil {
		_, err1 := c.Exec()
		fmt.Println(err1)
	}
	var logID int64 = 1
	rows, err := db.Query(`select max(id) as id from binlog`)
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&logID)
		}
	}
	fmt.Println(logID)
	return &DataBase{
		DB:    db,
		LogId: logID,
	}
}
func (d *DataBase) AddId() int64 {
	d.RLock()
	defer d.RUnlock()
	d.LogId++
	return d.LogId
}

func (d *DataBase) Query(sql string) [][]string {
	d.RLock()
	rows, err := d.DB.Query(sql)
	//list := make([]string, 0)
	result := make([][]string, 0)
	if err != nil {
		fmt.Println(err)
		return result
	}
	titleBool := true
	for rows.Next() {
		item := make([]string, 0)
		//fmt.Println(rows)
		list := make([]interface{}, 0)
		columns, _ := rows.ColumnTypes()
		title := make([]string, 0)
		for k, v := range columns {
			fmt.Println(k)
			var key string
			list = append(list, &key)
			title = append(title, v.Name())
		}
		if titleBool {
			result = append(result, title)
			titleBool = false
		}

		for k, _ := range columns {
			var key string
			list[k] = &key
			err = rows.Scan(list...)
			fmt.Println(key)
			item = append(item, key)
		}
		result = append(result, item)
	}
	d.RUnlock()
	return result
}

func (d *DataBase) Check(id int64) error {
	sql := "insert into binlog (id) values (" + strconv.FormatInt(id, 10) + ")"
	fmt.Println(sql)
	stmt, err := d.DB.Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (d *DataBase) Prepare(sql string, id int64) error {
	d.Lock()
	defer d.Unlock()
	err := d.Check(id)
	if err != nil {
		return nil
	}
	stmt, err := d.DB.Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("sqlerr")
		//panic(err)
	}
}

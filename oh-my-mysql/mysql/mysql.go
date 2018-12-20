package main

import "os"
import "encoding/json"
import "fmt"
import "reflect"
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func ChechErr(err error) {
	if err != nil {
		panic(err)
	}
}

type SSHMySQL struct {
	dbName string
}

func readConfigFile(fileName string, dbName string) map[string]string {
	f, err := os.Open(fileName)
	ChechErr(err)
	fileInfo, err := os.Stat(fileName)
	ChechErr(err)
	fileSize := fileInfo.Size()
	container := make([]byte, fileSize)
	f.Read(container)
	m := make(map[string]map[string]string)
	json.Unmarshal(container, &m)
	return m[dbName]
}

func (sm SSHMySQL) GetDB() *sql.DB {
	s := reflect.ValueOf(sm)
	dbName := s.FieldByName("dbName").String()
	dbConfig := readConfigFile("/home/zy/go/src/oh-my-mysql/mysql/DBConfig.txt", dbName)
	user := dbConfig["dbUser"]
	passwd := dbConfig["dbPwd"]
	host := dbConfig["dbHost"]
	port := dbConfig["dbPort"]
	db_name := dbConfig["dbName"]
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&autocommit=true", user, passwd, host, port, db_name))
	ChechErr(err)

	return db
}

func (sm SSHMySQL) SQL(db *sql.DB, sql string) []map[string]string {
	rows, err := db.Query(sql)
	ChechErr(err)
	column, _ := rows.Columns()
	values := make([][]byte, len(column))
	scans := make([]interface{}, len(column))
	for i := range values {
		scans[i] = &values[i]
	}

	result := make([]map[string]string, 0)
	for rows.Next() {
		err := rows.Scan(scans...)
		ChechErr(err)
		row := make(map[string]string)
		for k, v := range values {
			key := column[k]
			row[key] = string(v)
		}
		result = append(result, row)
	}
	return result
}

func main() {
	sm := SSHMySQL{dbName: "oh_my_mysql"}
	rows := sm.SQL(sm.GetDB(), "select * from test")
	for _, v := range rows {
		fmt.Println(v)
	}

}

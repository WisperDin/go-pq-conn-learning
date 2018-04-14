package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB
var dBDriverName string

func InitDB(host, port, user, pwd, dbName, driverName string, maxConns int) {
	dBDriverName = driverName
	//构建连接字符串
	dateSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, dbName)
	db, err := sql.Open(dBDriverName, dateSource)
	if err != nil {
		panic(err)
	}
	Db = db
	err = Db.Ping()
	if err != nil {
		log.Println("InitDB failed at Ping " + err.Error())
		panic(err)
	}
	db.SetMaxOpenConns(maxConns)

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS foo(
		id SERIAL,
		name text NOT NULL,
		pwd text NOT NULL,
		PRIMARY KEY ("id")
	)`)
	if err != nil {
		panic(err)
	}
}

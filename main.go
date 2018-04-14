package main

import (
	"./db"
	"log"
)

//Summarize several methods of operating the database

//1. stmt :: once prepare multiple exec
func PrepareStmt() {
	//prepare 从连接池那一条连接,并且追踪记录prepare的这条连接,返回一个stmt
	//下次真正执行时,只需要传语句中的参数即可执行sql,并且可以多次执行不同参数的这条prepare好的语句
	//每次exec时,优先尝试使用prepare的时候使用的那条连接,如果找不到,就需要re-prepared,当exec后,将连接归还连接池
	stmt, err := db.Db.Prepare(`INSERT INTO foo(name,pwd) values($1,$2)`)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("lzy", "123456")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec("hyx", "789123")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("ok")
}

//2. exec
func DirectExec() {
	//Exec 从连接池获取一条连接 执行之后直接归还连接
	_, err := db.Db.Exec(`INSERT INTO foo(name,pwd) values($1,$2)`, "lzy", "123456")
	if err != nil {
		log.Println(err)
		return
	}
}

//3. query/queryrow
func Query() {
	//query() 从连接池获取一条连接后,如果执行过程中无错误,将释放连接的函数传给rows结构返回
	// 直至rows.Close() 释放连接
	rows, err := db.Db.Query(`SELECT name from foo`)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var tName string
		err = rows.Scan(&tName)
		if err != nil {
			log.Println(err)
			return
		}
		names = append(names, tName)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(names)
}

//4. transaction :: atom
func Transaction() {
	//一个事务与一条连接唯一绑定
	//事务开始后将释放连接的方法传给了tx结构,而释放连接当且仅当在Commit 或者 Rollback时才会调用
	//即事务所占用的连接当且仅当Commit 或者 Rollback时才会释放
	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = tx.Exec(`INSERT INTO foo(name,pwd) values($1,$2)`, "who", "123456")
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, err = tx.Exec(`UPDATE foo SET name='fyz' WHERE name=$1`, "who")
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	db.InitDB("127.0.0.1", "6543", "liziyi", "", "liziyi", "postgres", 1)
	PrepareStmt()
	DirectExec()
	Query()
	Transaction()
}

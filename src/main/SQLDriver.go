package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDBDriver() {
	fmt.Println("Initializing DB Driver...")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)

	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
}

func testWrite() {

}

func testRead() {

}

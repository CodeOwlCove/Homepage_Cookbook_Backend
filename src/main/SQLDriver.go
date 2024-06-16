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

	fmt.Println("Creating connection with parameters: ")
	fmt.Println("DB_USER: ", DB_USER)
	fmt.Println("DB_PASS: ", DB_PASS)
	fmt.Println("DB_NAME: ", DB_NAME)
	fmt.Println("DB_HOST: ", DB_HOST)
	fmt.Println("DB_PORT: ", DB_PORT)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)

	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Error while connecting to database!")
		panic(err.Error())
	}
}

func testWrite() {

}

func testRead() {

}

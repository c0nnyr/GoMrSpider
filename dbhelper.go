package main

import (
	"database/sql"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main2() {
	_, err := sql.Open("sqlite3", "./data/house_cd.sqlite")
	if err != nil {
		log.Fatal("cannot open db")
		return
	}
	//stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	//checkErr(err)
}

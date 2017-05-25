package main

import (
	gin "gopkg.in/gin-gonic/gin.v1"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func init(){
	db, err := sql.Open("sqlite3", "data/xx.sqlite")
}

func rootHandler(c *gin.Context){
	c.JSON(200, gin.H{
		"msg":"ok",
	})
}

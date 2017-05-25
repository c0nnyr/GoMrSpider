package main

import (
	//_ "gopkg.in/gin-gonic/gin.v1"
	"fmt"
	"runtime"
)

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	fmt.Println(runtime.NumCPU())
	dispatcher := &Dispatcher{}
	net := &NetService{}
	dispatcher.SetNetService(net)
	dispatcher.Dispatch(NewTestSpider())
	//router := gin.Default()

	//router.GET("/", rootHandler)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	//router.Run()
	// router.Run(":3000") for a hard coded port
}

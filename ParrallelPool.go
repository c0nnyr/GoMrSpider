package main

import (
	"runtime"
)


type JobFunc func(params... interface{})

func ParallelRun(jobFunc JobFunc, params ...interface{}){
	ParrallelChan<- 1
	go pack(jobFunc, params...)
}

func pack(jobFunc JobFunc, params ...interface{}){
	jobFunc(params...)
	<-ParrallelChan
}


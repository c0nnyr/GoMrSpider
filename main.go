package main

import (
)

func main() {
	dispatcher := &Dispatcher{}
	net := &NetService{}
	dispatcher.SetNetService(net)
	dispatcher.Dispatch(NewTestSpider())
}

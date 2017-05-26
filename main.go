package main

import (
)

func main() {
	dispatcher := NewDispatcher()
	dispatcher.SetNetService(&NetService{})
	dispatcher.Dispatch(NewTestSpider())
}

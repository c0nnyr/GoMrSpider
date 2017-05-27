package main
import "fmt"

func main() {
	dispatcher := NewDispatcher()
	dispatcher.SetNetService(&NetService{})
	dispatcher.Dispatch(NewTestSpider())
}

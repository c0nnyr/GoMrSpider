package main

func main() {
	dispatcher := NewDispatcher()
	dispatcher.SetNetService(&NetService{})
	//dispatcher.Dispatch(NewTestSpider())
	dispatcher.Dispatch(NewProxySpider())
}

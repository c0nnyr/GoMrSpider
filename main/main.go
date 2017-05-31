package main

import (
	_ "mrspider"
	"mrspider/proxy"
)

func main() {
	//session := mrspider.CreateMongoSession()
	//if session == nil {
	//	return
	//}
	//dispatcher := mrspider.NewDispatcher()
	//dispatcher.SetNetService(&mrspider.NetService{})
	//dispatcher.AddItemMidware(mrspider.ItemMidwareSaveToDB)
	//dispatcher.AddResponseMidware(mrspider.ResponseMidwareWait)
	//dispatcher.SetConfigFile("dispatcher.conf")
	//dispatcher.Dispatch(proxy.NewKuaidailiProxySpider())
	proxy.TestProxy(nil)
}

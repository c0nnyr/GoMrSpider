package main

import (
	"mrspider"
	"mrspider/proxy"
)

func save_to_db(item mrspider.IBaseItem) bool{
	if item, ok := item.(mrspider.IDBItem); ok{
		mrspider.DBUpsert(item)
	}
	return true
}

func main() {
	session := mrspider.CreateMongoSession()
	if session == nil {
		return
	}
	dispatcher := mrspider.NewDispatcher()
	dispatcher.SetNetService(&mrspider.NetService{})
	//dispatcher.Dispatch(mrspider.NewTestSpider())
	dispatcher.AddItemMidware(save_to_db)
	dispatcher.Dispatch(proxy.NewKuaidailiProxySpider())
}

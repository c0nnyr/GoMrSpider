package main

import (
	"fmt"
	"time"
)

type TestSpider struct{
	BaseSpider
}

type TestItem struct{
	BaseItem
	url string
}

func (self *TestItem)String() string {
	return self.url
}

func NewTestSpider() *TestSpider{
	self := &TestSpider{}
	self.start_urls = append(self.start_urls, "http://1.com", "http://2.com")
	self.defaultParser = self.Parse
	return self
}

func (self *TestSpider) Parse(response *Response) (requestOrItems RequestOrItems){
	fmt.Println("parse in test spider")
	requestOrItems = append(requestOrItems, &TestItem{url:response.url})
	requestOrItems = append(requestOrItems, &TestItem{url:response.url + "100"})
	time.Sleep(time.Second)
	return
}

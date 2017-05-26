package main

import (
	_ "fmt"
	"time"
	"fmt"
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
	for ind := 0; ind < 100; ind++ {
		self.start_urls = append(self.start_urls, fmt.Sprintf("http://%v", ind))
	}
	//self.start_urls = []string{"http://1.com", "http://2.com"}
	self.defaultCallback = self.Parse
	return self
}

func (self *TestSpider) Parse(response *Response) (requestOrItems RequestOrItems){
	requestOrItems = append(requestOrItems, &TestItem{url:response.url})
	requestOrItems = append(requestOrItems, &Request{url:"2" + response.url, callback:self.ParseDone})
	time.Sleep(time.Second)
	return
}

func (self *TestSpider) ParseDone(response *Response) (requestOrItems RequestOrItems){
	requestOrItems = append(requestOrItems, &TestItem{url:response.url})
	return
}

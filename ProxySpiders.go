package main

import (
	_ "fmt"
	_ "time"
)

type ProxySpider struct{
	BaseSpider
}

type ProxyItem struct{
	url string
}

func (self *ProxyItem)String() string {
	return self.url
}

func NewProxySpider() *ProxySpider{
	self := &ProxySpider{}
	const MAX_PAGE = 2
	//self.start_urls = make([]string, 0, 2 * MAX_PAGE)

	self.start_urls = append(self.start_urls, "http://www.kuaidaili.com/free/outha/1/")
	//self.start_urls = append(self.start_urls, "http://www.xicidaili.com/nn")
	//for i := 2; i < MAX_PAGE; i++{
	//	self.start_urls = append(self.start_urls, fmt.Sprintf("http://www.xicidaili.com/nn/%v", i))
	//}
	//self.start_urls = append(self.start_urls, "http://www.xicidaili.com/nt")
	//for i := 2; i < MAX_PAGE; i++{
	//	self.start_urls = append(self.start_urls, fmt.Sprintf("http://www.xicidaili.com/nt/%v", i))
	//}
	self.defaultCallback = self.Parse
	return self
}

func (self *ProxySpider) Parse(response *Response) (requests []*Request, items []IBaseItem){
	//fmt.Println(response.body)
	return
}

func (self *ProxySpider)IsValidResponse(_ *Response)bool{
	return true
}

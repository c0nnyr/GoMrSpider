package main

type TestSpider struct{
	BaseSpider
}

type TestItem struct{
	BaseItem
	url string
}

func NewTestSpider() *TestSpider{
	self := &TestSpider{}
	append(self.start_urls, "http://1.com", "http://2.com")
	return self
}

func (self *TestSpider) Parse(response *Response) (requestOrItems RequestOrItems){
	append(requestOrItems, &TestItem{url:response.url})
	append(requestOrItems, &TestItem{url:response.url})
	return
}

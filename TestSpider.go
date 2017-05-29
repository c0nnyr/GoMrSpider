package mrspider

import (
	_ "fmt"
	_ "time"
	"fmt"
)

type TestSpider struct{
	BaseSpider
}

type TestItem struct{
	url string
}

func (self *TestItem)String() string {
	return self.url
}

func NewTestSpider() *TestSpider{
	self := &TestSpider{}
	for ind := 0; ind < 50; ind++ {
		self.StartUrls = append(self.StartUrls, fmt.Sprintf("http://%v", ind))
	}
	//self.start_urls = []string{"http://1.com", "http://2.com"}
	self.DefaultCallback = self.Parse
	return self
}

func (self *TestSpider) Parse(response *Response) (requests []*Request, items []IBaseItem){
	items = append(items, &TestItem{url:response.url})
	requests = append(requests, &Request{url:"2" + response.url, callback:self.ParseDone})
	return
}

func (self *TestSpider) ParseDone(response *Response) (requests []*Request, items []IBaseItem){
	items = append(items, &TestItem{url:response.url})
	return
}

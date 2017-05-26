package main

import (
	"fmt"
)

type IBaseSpider  interface {
	GetStartRequests(_ *Response) (requestOrItems RequestOrItems)
	Parse(response *Response) (requestOrItems RequestOrItems)
}

type BaseSpider struct {
	start_urls []string
	metas [][]interface{}
	defaultParser RequestCallback
}

func (self *BaseSpider)GetStartRequests(_ *Response) (requestOrItems RequestOrItems){
	if self.metas == nil {
		for _, start_url := range self.start_urls {
			requestOrItems = append(requestOrItems, &Request{
				method:"GET",
				url:start_url,
				callback:self.defaultParser,//不能直接用self.Parser，否则就只是绑定到这里的Parser了
			})
		}
	} else {
		for ind, start_url := range self.start_urls {
			requestOrItems = append(requestOrItems, &Request{
				method:"GET",
				url:fmt.Sprintf(start_url, self.metas[ind]...),
				callback:self.defaultParser,
			})
		}
	}
	return
}

func (self *BaseSpider)Parse(response *Response) (requestOrItems RequestOrItems){
	fmt.Println("parse in base spider")
	return
}

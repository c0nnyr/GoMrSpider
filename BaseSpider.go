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
}

func (self *BaseSpider)GetStartRequests(_ *Response) (requestOrItems RequestOrItems){
	if self.metas == nil {
		for _, start_url := range self.start_urls {
			requestOrItems = append(requestOrItems, &Request{
				method:"GET",
				url:start_url,
				callback:self.Parse,
			})
		}
	} else {
		for ind, start_url := range self.start_urls {
			requestOrItems = append(requestOrItems, &Request{
				method:"GET",
				url:fmt.Sprintf(start_url, self.metas[ind]...),
				callback:self.Parse,
			})
		}
	}
	return
}

func (self *BaseSpider)Parse(response *Response) (requestOrItems RequestOrItems){
	return
}

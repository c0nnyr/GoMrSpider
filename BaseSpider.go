package main

import (
	"fmt"
)

type BaseSpider struct {
	start_urls []string
	metas [][]string
}

func (self *BaseSpider)GetStartRequests(_ *Response) (requestOrItems RequestOrItems){
	if self.metas == nil {
		for ind := range self.start_urls {
			requestOrItems = append(requestOrItems, &Request{
				method:"GET",
				url:self.start_urls[ind],
				callback:self.Parse,
			})
		}
	} else {
		for ind := range self.start_urls {
			requestOrItems = append(requestOrItems, &Request{
				method:"GET",
				url:fmt.Sprintf(self.start_urls[ind], self.metas[ind]...),
				callback:self.Parse,
			})
		}
	}
	return
}

func (self *BaseSpider)Parse(response *Response) (requestOrItems RequestOrItems){
	return
}

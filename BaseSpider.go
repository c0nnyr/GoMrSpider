package main

import (
	"fmt"
)

type IBaseSpider  interface {
	GetStartRequests(_ *Response) (requests []*Request)
}

type BaseSpider struct {
	start_urls []string
	metas [][]interface{}
	defaultCallback RequestCallback
}

func (self *BaseSpider)GetStartRequests(_ *Response) (requests []*Request){
	if self.metas == nil {
		for _, start_url := range self.start_urls {
			requests = append(requests, &Request{
				method:"GET",
				url:start_url,
				callback:self.defaultCallback,//不能直接用self.Parser，否则就只是绑定到这里的Parser了
			})
		}
	} else {
		for ind, start_url := range self.start_urls {
			requests = append(requests, &Request{
				method:"GET",
				url:fmt.Sprintf(start_url, self.metas[ind]...),
				callback:self.defaultCallback,
			})
		}
	}
	return
}


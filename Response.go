package main

import (
	"fmt"
)

type Response struct {
	body string
	url string
}
func (self *Response)String() string{
	return fmt.Sprintf("%v, %v", self.url, self.body)
}

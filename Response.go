package mrspider

import (
	"fmt"
	"gopkg.in/xmlpath.v2"
)

type Response struct {
	body []byte
	url string
	htmlRoot *xmlpath.Node
}
func (self *Response)String() string{
	//return fmt.Sprintf("%v, %v", self.url, self.body[:100])
	return fmt.Sprintf("%v", self.url)
}

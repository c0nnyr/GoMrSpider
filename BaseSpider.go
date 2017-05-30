package mrspider

import (
	"fmt"
	"gopkg.in/xmlpath.v2"
	"log"
)

type IBaseSpider  interface {
	GetStartRequests(*Response) ([]*Request)
	IsValidResponse(*Response) bool
}

type BaseSpider struct {
	StartUrls []string
	Metas [][]interface{}
	DefaultCallback RequestCallback
	ValidXpath string
	ValidXpathCompiled *xmlpath.Path
}

func (self *BaseSpider)GetStartRequests(_ *Response) (requests []*Request){
	if self.Metas == nil {
		for _, startUrl := range self.StartUrls {
			requests = append(requests, &Request{
				method:"GET",
				url:startUrl,
				callback:self.DefaultCallback,//不能直接用self.Parser，否则就只是绑定到这里的Parser了
			})
		}
	} else {
		for ind, start_url := range self.StartUrls {
			requests = append(requests, &Request{
				method:"GET",
				url:fmt.Sprintf(start_url, self.Metas[ind]...),
				callback:self.DefaultCallback,
			})
		}
	}
	return
}

func (self *BaseSpider)IsValidResponse(response *Response)bool{
	if self.ValidXpathCompiled == nil{
		if self.ValidXpath != ""{
			var err error
			self.ValidXpathCompiled, err = xmlpath.Compile(self.ValidXpath)
			if err != nil{
				log.Fatalln("valid xpath compile failed", err)
				return false
			}
		} else {
			return true//如果没有设置就返回true了
		}
	}
	if response.htmlRoot != nil{
		return false//不是html
	}
	iter := self.ValidXpathCompiled.Iter(response.htmlRoot)
	if iter.Next(){
		return true
	}
	return false
}

type NewItemFunc func (map[string]string)IBaseItem
func (self *BaseSpider)ParseItems(response *Response, attrMap map[string]string, newItemFunc NewItemFunc, listXpath string) (items []IBaseItem){
	listXpathCompiled, err := xmlpath.Compile(listXpath)//这里需要有tbody,因为了thea
	if err != nil{
		log.Fatalln("list_xpath", err)
		return
	}
	attrXpathMap := map[string]*xmlpath.Path{}
	for attr, xpath := range attrMap{
		if xpath == ""{
			continue
		}
		attrXpathMap[attr], err = xmlpath.Compile(xpath)
		if err != nil{
			log.Fatalln("attr_xpath", attr, xpath, err)
			return
		}
	}

	if response.htmlRoot != nil{
		listIter := listXpathCompiled.Iter(response.htmlRoot)
		for listIter.Next(){
			itemData := map[string]string{}
			for attr, xpath := range attrXpathMap{
				iter := xpath.Iter(listIter.Node())
				if iter.Next(){
					itemData[attr] = iter.Node().String()
				}
			}
			items = append(items, newItemFunc(itemData))
		}
	}
	return
}

package proxy

import (
	_ "fmt"
	"log"
	"mrspider"
)

type KuaidailiProxySpider struct{
	mrspider.BaseSpider
}

func NewKuaidailiProxySpider() *KuaidailiProxySpider{
	self := &KuaidailiProxySpider{}
	self.StartUrls = append(self.StartUrls, "http://www.kuaidaili.com/free/outha/1/")

	self.DefaultCallback = self.Parse
	self.ValidXpath = `//*[@id="list"]/table`
	return self
}

func (self *KuaidailiProxySpider) Parse(response *mrspider.Response) (requests []*mrspider.Request, items []mrspider.IBaseItem){
	listXpath := `//*[@id="list"]/table/tbody/tr`//这里需要有tbody,因为了thea
	attrMap := map[string]string{
		//attr xpath, re_filter
		"Country":"",
		"IP":"td[1]/text()",
		"Port":"td[2]/text()",
		"AnonymouseType":"td[3]/text()",
		"LinkType":"td[4]/text()",
		}

	items = append(items, self.ParseItems(response, attrMap, NewProxyItem, listXpath)...)
	return
}

//const MAX_PAGE = 2
//self.start_urls = make([]string, 0, 2 * MAX_PAGE)

//self.start_urls = append(self.start_urls, "http://www.xicidaili.com/nn")
//for i := 2; i < MAX_PAGE; i++{
//	self.start_urls = append(self.start_urls, fmt.Sprintf("http://www.xicidaili.com/nn/%v", i))
//}
//self.start_urls = append(self.start_urls, "http://www.xicidaili.com/nt")
//for i := 2; i < MAX_PAGE; i++{
//	self.start_urls = append(self.start_urls, fmt.Sprintf("http://www.xicidaili.com/nt/%v", i))
//}

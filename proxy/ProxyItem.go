package proxy

import (
	"fmt"
	"mrspider"
	"time"
	"labix.org/v2/mgo/bson"
	"net/url"
	"net/http"
	"strings"
)

type ProxyItem struct{
	Country string `bson:"country`
	IP string `bson:"ip`
	Port string `bson:"port`
	AnonymouseType string `bson:"anonymouse_type`
	LinkType string `bson:"link_type`
	CrawlTime int `bson:"crawl_time"`
}

func (self *ProxyItem)String() string {
	return fmt.Sprintf("%v:%v at time %v", self.IP, self.Port, self.CrawlTime)
}

func NewProxyItem(m map[string]string) mrspider.IBaseItem{
	self := &ProxyItem{}
	v, ok := m["Country"]
	if ok{
		self.Country = v
	}
	v, ok = m["IP"]
	if ok{
		self.IP = v
	}
	v, ok = m["Port"]
	if ok{
		self.Port = v
	}
	v, ok = m["AnonymouseType"]
	if ok{
		self.AnonymouseType = v
	}
	v, ok = m["LinkType"]
	if ok{
		self.LinkType = v
	}
	self.CrawlTime = time.Now().Nanosecond()
	return self
}

func (self *ProxyItem)GetMgoID() *bson.M{
	return &bson.M{
		"ip":self.IP,
		"port":self.Port,
	}
}

func (self *ProxyItem)GetCollectionName() string{
	return "proxy_item"
}

func (self *ProxyItem)CreateProxyClient(timeOut time.Duration) *http.Client{
	proxyFunc := func (_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("%v://%v:%v", strings.ToLower(self.LinkType), self.IP, self.Port))
	}
	transport := &http.Transport{Proxy:proxyFunc}
	client := &http.Client{Transport:transport, Timeout:timeOut}
	return client
}

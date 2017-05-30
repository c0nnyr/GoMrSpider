package mrspider

import (
	"labix.org/v2/mgo/bson"
)

type ResponseCache struct {
	RequestFullUrl string `bson:"request_full_url"`//首字母必须大写,才能被mgo访问,reflect中structInfo中的PkgPath才是空...
	ResponseBody  []byte  `bson:"response_body"`
}

func (self *ResponseCache)GetMgoID() *bson.M{
	return &bson.M{
		"request_full_url":self.RequestFullUrl,
	}
}

func (self *ResponseCache)GetCollectionName() string{
	return MONGO_RESPONSE_CACHE_COLLECTION
}

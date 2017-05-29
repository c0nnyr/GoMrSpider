package mrspider

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"errors"
)

var mongoBasicSesstion *mgo.Session = nil

func CreateMongoSession() *mgo.Session {
	if mongoBasicSesstion == nil {
		var err error
		mongoBasicSesstion, err = mgo.Dial(MONGO_URL)
		if err != nil {
			log.Fatal("mongo cannot be connected")
		}
	}
	return mongoBasicSesstion
}

type IDBItem interface{
	GetMgoID() *bson.M
	GetCollectionName() string
}
func DBUpsert(item IDBItem) {
	if mongoBasicSesstion == nil {
		return
	}
	session := mongoBasicSesstion.Clone()//
	defer session.Close()
	c := session.DB(MONGO_DB).C(item.GetCollectionName())
	_, err := c.Upsert(item.GetMgoID(), item)
	if err != nil {
		log.Fatalf("upsert response cache error with id%v", item.GetMgoID())
		return
	}
}

func DBFind(item IDBItem) error{
	if mongoBasicSesstion == nil {
		return errors.New("mongoBasicSesstion is nil")
	}
	session := mongoBasicSesstion.Clone()
	defer session.Close()
	c := session.DB(MONGO_DB).C(item.GetCollectionName())
	err := c.Find(item.GetMgoID()).One(item)
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////

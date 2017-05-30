package mrspider
import (
	"time"
	"log"
)

func ItemMidwareSaveToDB(item IBaseItem) bool{
	if item, ok := item.(IDBItem); ok{
		DBUpsert(item)
	}
	return true
}

func ResponseMidwareWait(_ *Response) bool{
	time.Sleep(2 * time.Second)
	log.Println("sleep done")
	return true
}

package mrspider
import (
	"testing"
	"math/rand"
	"strconv"
	"time"
	"log"
)

func TestResponseCache(t *testing.T){
	session := CreateMongoSession()
	if session == nil {
		return
	}
	requestFullUrl := "hello"
	rand.Seed(int64(time.Now().Second()))
	responseBody := []byte("world" + strconv.Itoa(rand.Int()))
	responseCache := &ResponseCache{RequestFullUrl:requestFullUrl, ResponseBody:responseBody}
	DBUpsert(responseCache)
	newResponseCahce := &ResponseCache{RequestFullUrl:requestFullUrl}
	err := DBFind(newResponseCahce)
	log.Println(newResponseCahce)
	if err != nil {
		t.Error("Cannot find the response cache item inserted now. err:", err, "url:", requestFullUrl, "body:", responseBody)
		return
	}
	if string(newResponseCahce.ResponseBody) != string(responseBody) {
		t.Error("The item found is not equal to inserted. url:", requestFullUrl, "body", responseBody, "founded body",
			newResponseCahce.ResponseBody, "founded url", newResponseCahce.RequestFullUrl)
	}
}

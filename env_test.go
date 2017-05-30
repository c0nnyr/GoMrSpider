package mrspider
import (
	"testing"
	"log"
	"os"
)

func TestMgoEnv(t *testing.T){
	session := CreateMongoSession()
	if session == nil {
		t.Error("Cannot dial mgo with", MONGO_URL)
	}
}

func TestDispatchConfig(t *testing.T){
	dispatcher := &Dispatcher{}
	log.Println("pre", dispatcher.config)
	dispatcher.SetConfigFile("/Users/conny/self_project/web/bin/dispatcher.conf")
	log.Println("post", dispatcher.config)
	log.Println(os.Args[0])
	t.Fail()
}


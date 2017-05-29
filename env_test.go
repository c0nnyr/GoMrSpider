package mrspider
import (
	"testing"
)

func TestMgoEnv(t *testing.T){
	session := CreateMongoSession()
	if session == nil {
		t.Error("Cannot dial mgo with", MONGO_URL)
	}
}


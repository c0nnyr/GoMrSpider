package proxy
import (
	"testing"
	"log"
	"io/ioutil"
	"mrspider"
	"time"
)
type timeoutError interface {
	Timeout() bool
}
func TestProxy(t *testing.T){
	mrspider.CreateMongoSession()
	proxy := &ProxyItem{}
	iter := mrspider.DBFindAll(proxy.GetCollectionName())
	for iter.Next(proxy){
		func (){
			log.Println(proxy)
			if proxy.LinkType != "HTTP"{
				return
			}
			client := proxy.CreateProxyClient(3 * time.Second)
			log.Println("trying connecting")
			res, err := client.Get("http://www.whatismyip.com.tw")
			if err != nil {
				log.Printf("%T %v\n", err, err)
				return
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			log.Println("body", string(body))
			time.Sleep(2 * time.Second)
		}()
	}
	//t.FailNow()
}

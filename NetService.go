package main
import (
	"net/url"
	"net/http"
	"io/ioutil"
	_ "time"
	"log"
	"os"
)

type NetService struct {

}

func (self *NetService) SendRequest(request *Request) *Response{
	var res *http.Response
	var err error
	if request.method == "GET" {
		target_url := request.url
		if request.data != nil {
			params := url.Values{}
			for k, v := range request.data{
				params[k] = []string{v}
			}
			target_url = request.url + params.Encode()
		}
		res, err = http.Get(target_url)
		if err != nil {
			log.Printf("HTTP GET ERROR %v", err)
			return nil
		}
	} else if request.method == "POST"{
		log.Printf("Unsupported method POST")
		return nil
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return &Response{
		body:string(body),
		url:request.url,
	}
	//time.Sleep(2 * time.Second)
	//return &Response{
	//	body:"noting",
	//	url:request.url,
	//}
}

package main
import (
	"net/url"
	"net/http"
	"io/ioutil"
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
			panic("HTTP GET ERROR %s")
		}
	} else if request.method == "POST"{
		panic("NOT SUPPORT POST")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return &Response{
		body:string(body),
		url:request.url,
	}
}

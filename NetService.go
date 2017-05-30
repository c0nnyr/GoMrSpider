package mrspider
import (
	"net/url"
	"net/http"
	"io/ioutil"
	"log"
	"gopkg.in/xmlpath.v2"
	"bytes"
)

type NetService struct {

}

func (self *NetService) SendRequest(request *Request) *Response{
	var res *http.Response
	var err error
	var target_url string
	response := &Response{}
	if request.method == "GET" {
		target_url = request.url
		if request.data != nil {
			params := url.Values{}
			for k, v := range request.data{
				params[k] = []string{v}
			}
			target_url = request.url + params.Encode()
		}
		if ENABLE_USE_RESPONSE_CACHE{
			response_cache := &ResponseCache{RequestFullUrl:target_url}
			if DBFind(response_cache) == nil {
				response.url = target_url
				response.body = response_cache.ResponseBody
				log.Println("using response cache")
				response.htmlRoot, err = xmlpath.ParseHTML(bytes.NewBuffer(response.body))//有可能是其他类型的
				if err != nil{
					log.Fatalln("cannot parse html", err)
				}
				return response
			}
		}
		log.Println("requesting", target_url)
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

	if ENABLE_SAVE_RESPONSE_CACHE{
		response_cache := &ResponseCache{
			RequestFullUrl:target_url,
			ResponseBody:body,
		}
		DBUpsert(response_cache)
	}

	log.Printf("received body len %v", len(body))

	response.body = body
	response.url = target_url
	response.htmlRoot, err = xmlpath.ParseHTML(bytes.NewBuffer(response.body))//有可能是其他类型的

	return response
}

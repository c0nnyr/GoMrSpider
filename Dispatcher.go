package main
import (
	"fmt"
	_ "runtime"
)

const (
	DISPATCH_MODE_DEPTH = iota
	DISPATCH_MODE_WIDTH
)

type ResponsePack struct {
	response *Response
	callback RequestCallback
}

type Dispatcher struct {
	config map[string]int
	net *NetService
	requests []*Request

	itemChan chan *BaseItem

	requestChan chan *Request
	requestHeadChan chan *Request

	responseChan chan *ResponsePack
}

func NewDispatcher() *Dispatcher{
	//const MAX_PROCS = runtime.GOMAXPROCS(0)
	const MAX_PROCS = 1
	fmt.Println("max procs %s", MAX_PROCS)
	self := &Dispatcher{
		config:{
			"mode":DISPATCH_MODE_DEPTH,//遍历模式
		},
		net:nil,

		itemChan:make(chan *BaseItem, MAX_PROCS),
		requestChan:make(chan *Request, MAX_PROCS),
		requestHeadChan:make(chan *Request, MAX_PROCS),
		responseChan:make(chan *Response, MAX_PROCS),
	}
	return self
}

func (self *Dispatcher)SetNetService(net *NetService){
	self.net = net
}

func (self *Dispatcher)Dispatch(spiders... *BaseSpider){
	for i := 0; i < cap(self.requestChan); i++ {
		go self.handleRequest(i)
	}
	for i := 0; i < cap(self.itemChan); i++ {
		go self.handleItem(i)
	}
	for i := 0; i < cap(self.responseChan); i++ {
		go self.handleResponse(i)
	}

	fmt.Println("Dispatching spiders")
	for _, spider := range spiders{
		for _, request := range spider.GetStartRequests(nil) {
			self.requestChan <- request
		}
	}
}

func (self *Dispatcher)handleRequest(ind int){
	fmt.Println("handling request with go %s", ind)
	sendRequest := func (request *Request){
		fmt.Println("handleRequest send request %s with ind %s", request.url, ind)
		response := self.net.SendRequest(request)
		self.responseChan<- &ResponsePack{response, request.callback}
	}
	for {
		select {
		case request := <-self.requestHeadChan://先看head里面有没有
			sendRequest(request)
		default:
				select {
				case request := <-self.requestHeadChan:
					sendRequest(request)
				case request := <-self.requestChan:
					sendRequest(request)
				}
		}
	}
}

func (self *Dispatcher)handleResponse(ind int){
	fmt.Println("handling response with go %s", ind)
	for {
		responsePack := <-self.responseChan
		fmt.Println("handleResponse receive response %s with ind %s", responsePack.response.url, ind)
		requestOrItems := responsePack.callback(responsePack.response)
		for _, requestOrItem := range requestOrItems {
			switch requestOrItem := requestOrItem.(type) {
			case *Request:
				if mode := self.config["mode"]; mode == DISPATCH_MODE_DEPTH {
					self.requestHeadChan <- requestOrItem
				} else if mode == DISPATCH_MODE_WIDTH{
					self.requestChan<- requestOrItem
				}
			case *BaseItem:
				self.itemChan<- requestOrItem
			default:
			}
		}

	}
}

func (self *Dispatcher)handleItem(ind int){
	fmt.Println("handling item with go %s", ind)
	for {
		item := <-self.itemChan
		fmt.Println("handleItem receive item %s with ind %s", item, ind)
	}
}

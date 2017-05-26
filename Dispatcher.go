package main
import (
	"fmt"
	_ "runtime"
	"time"
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

	itemChan chan IBaseItem

	requestCacheChan chan *Request
	requestChan chan *Request

	responseChan chan *ResponsePack

	heartChan chan string
}

func NewDispatcher() *Dispatcher{
	//const MAX_PROCS = runtime.GOMAXPROCS(0)
	const MAX_PROCS = 0
	fmt.Printf("max procs %v\n", MAX_PROCS)
	self := &Dispatcher{
		config:map[string]int{
			"mode":DISPATCH_MODE_WIDTH,//遍历模式
		},
		net:nil,

		itemChan:make(chan IBaseItem, MAX_PROCS),
		requestCacheChan:make(chan *Request),
		requestChan:make(chan *Request, MAX_PROCS),
		responseChan:make(chan *ResponsePack, MAX_PROCS),
		heartChan:make(chan string, MAX_PROCS),
	}
	return self
}

func (self *Dispatcher)SetNetService(net *NetService){
	self.net = net
}

func (self *Dispatcher)Dispatch(spiders... IBaseSpider){
	for i := 0; i < Max(cap(self.requestChan), 1); i++ {
		go self.handleRequest(i)
	}
	for i := 0; i < Max(cap(self.itemChan), 1); i++ {
		go self.handleItem(i)
	}
	for i := 0; i < Max(cap(self.responseChan), 1); i++ {
		go self.handleResponse(i)
	}

	go self.dispatchRequests()

	fmt.Printf("Dispatching spiders\n")
	for _, spider := range spiders {
		requests := spider.GetStartRequests(nil)
		if mode := self.config["mode"]; mode == DISPATCH_MODE_DEPTH {
			//倒着插进去
			for i := len(requests); i > 0; i--{
				self.requestCacheChan <- requests[i - 1]
			}
		} else {
			for i := 0; i < len(requests); i++ {
				self.requestCacheChan <- requests[i]
			}
		}
	}
	self.WaitAllDone()
	fmt.Printf("All Done!!\n")
}

func (self *Dispatcher)WaitAllDone(){
	const DURATION = 10 * time.Second
	timer := time.NewTimer(DURATION)

	OUTER_LOOP:
	for {
		timer.Reset(DURATION)
		select {
		case <-self.heartChan:
		//case v := <-self.heartChan:
			//fmt.Printf("heart beating....%v\n", v)
		case <-timer.C:
			break OUTER_LOOP
		}
	}
}

func (self *Dispatcher)dispatchRequests(){
	var requests []*Request
	insertRequest := func (request *Request){
		requests = append(requests, request)
		if mode := self.config["mode"]; mode == DISPATCH_MODE_DEPTH {
			copy(requests[1:], requests[:len(requests) - 1]) //头部插入
			requests[0] = request
		}
	}
	handledRequestCount := 0
	for {
		if len(requests) != 0 {
			select {
			case request := <-self.requestCacheChan:
				insertRequest(request)
			case self.requestChan <- requests[0]:
				//未测试
				requests = requests[1:]//这个操作导致数组越来越大,要适当时候处理下,怎么处理呢,背后的数组都找不到头了,只有拷贝一下
				handledRequestCount += 1
				if handledRequestCount > 100 {
					newRequests := make([]*Request, len(requests))
					copy(newRequests, requests)
					requests = newRequests
				}

			}
		} else{
			request := <-self.requestCacheChan
			insertRequest(request)
		}
	}

}

func (self *Dispatcher)handleRequest(ind int){
	sendRequest := func (request *Request){
		self.heartChan<- "handleRequest"
		fmt.Printf("handleRequest send request %v with go %v\n", request.url, ind)
		response := self.net.SendRequest(request)
		self.responseChan<- &ResponsePack{response, request.callback}
	}
	for {
		request := <-self.requestChan
		sendRequest(request)
	}
}

func (self *Dispatcher)handleResponse(ind int){
	for {
		responsePack := <-self.responseChan
		self.heartChan<- "handleResponse"
		fmt.Printf("handleResponse receive response %v with go %v\n", responsePack.response.url, ind)
		requests, items := responsePack.callback(responsePack.response)
		for _, request := range requests {
			self.requestCacheChan<- request
		}
		for _, item := range items {
			self.itemChan<- item
		}
	}
}

func (self *Dispatcher)handleItem(ind int){
	for {
		item := <-self.itemChan
		self.heartChan<- "handleItem"
		fmt.Printf("handleItem receive item %v with go %v\n", item, ind)
	}
}

package main
import (
	"fmt"
	"runtime"
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
	responseChan chan *ResponsePack
}

func NewDispatcher() *Dispatcher{
	const MAX_PROCS = runtime.GOMAXPROCS(0)
	self := &Dispatcher{
		config:{
			"mode":DISPATCH_MODE_DEPTH,//遍历模式
		},
		net:nil,

		itemChan:make(chan *BaseItem, MAX_PROCS / 2),
		requestChan:make(chan *Request, MAX_PROCS / 2),
		responseChan:make(chan *Response, MAX_PROCS / 2),
	}
	return self
}

func (self *Dispatcher)SetNetService(net *NetService){
	self.net = net
}

func (self *Dispatcher)Dispatch(spiders... *BaseSpider){
	for i := 0; i < cap(self.requestChan); i++ {
		go self.handleRequest()
	}
	for i := 0; i < cap(self.itemChan); i++ {
		go self.handleItem()
	}
	for i := 0; i < cap(self.responseChan); i++ {
		go self.handleResponse()
	}

	for _, spider := range spiders{
		for _, request := range spider.GetStartRequests(nil) {
			self.requestChan <- request
		}
	}

	self.do_dispatch_job()
}

func (self *Dispatcher)handleRequest(){
	for {
		request := <-self.requestChan
		//append(self.requests, request)//先扩容
		//if cur_mode := self.config["mode"]; cur_mode == DISPATCH_MODE_DEPTH{
		//	//队列头部插入
		//	copy_len := len(self.requests)
		//	copy(self.requests[1:copy_len], self.requests[0:copy_len - 1])
		//	self.requests[0] = request
		//}
		response := self.net.SendRequest(request)
		self.responseChan<- &ResponsePack{response, request.callback}
	}
}

func (self *Dispatcher)handleResponse(){
	for {
		responsePack := <-self.responseChan
		request_or_items := responsePack.callback(responsePack.response)
		for _, request_or_item := range request_or_items {
			switch request_or_item := request_or_item.(type) {
			case *Request:
				self.requestChan<- request_or_item//只能宽度优先了....
			case *BaseItem:
				self.itemChan<- request_or_item
			default:
			}
		}

	}
}

func (self *Dispatcher)do_dispatch_job(){
	for {
		<-self.startJobChan
		for {
			if len(self.requests) != 0 {
				request := self.requests[0]
				self.requests = self.requests[1:]
				self.jobChan <- 1
				go self.runJob(request)
			} else {
				<-self.stopJobChan
				break
			}
		}
	}
}


func (self *Dispatcher) runJob(request *Request){
	response := self.net.SendRequest(request)
	new_request_or_items := request.callback(response)
	self.dispatch_request_or_items(new_request_or_items)
	<-self.jobChan
}

func (self *Dispatcher)do_dispatch_item(){
	for {
		item := <-self.itemChan
		fmt.Println(item)
	}
}

func (self *Dispatcher)dispatch_request_or_items(request_or_items RequestOrItems){
	for _, request_or_item := range request_or_items {
		switch request_or_item := request_or_item.(type) {
		case *Request:
			self.requestChan<- request_or_item
		case *BaseItem:
			self.itemChan<- request_or_item
		default:
		}
	}
}

func (self *Dispatcher)ClearBuffer(request_or_items RequestOrItems, original_slice RequestOrItems) RequestOrItems{
	if cap(original_slice) > cap(request_or_items) * 2 {
		//前面没用的太多了,清理一下
		original_slice = original_slice[:len(request_or_items)]//保证拷贝全
		copy(original_slice, request_or_items)
		request_or_items = original_slice
	}
	return request_or_items
}

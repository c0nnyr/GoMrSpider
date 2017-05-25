package main
import (
	"fmt"
	"runtime"
)

const (
	DISPATCH_MODE_DEPTH = iota
	DISPATCH_MODE_WIDTH
)

type Dispatcher struct {
	config map[string]int
	net *NetService
	requests []*Request

	itemChan chan *BaseItem
	requestChan chan *Request
	jobChan chan int

	startJobChan chan int
	stopJobChan chan int
}

func New() *Dispatcher{
	dispatcher := &Dispatcher{
		config:{
			"mode":DISPATCH_MODE_DEPTH,//遍历模式
		},
		net:nil,

		itemChan:make(chan *BaseItem),
		requestChan:make(chan *Request),
		startJobChan:make(chan int, 1),
		stopJobChan:make(chan int, 1),
	}
	dispatcher.jobChan = make(chan int, runtime.GOMAXPROCS(0))
	dispatcher.stopJobChan<-1//表示已经停止了,没有任务待分发
	return dispatcher
}

func (self *Dispatcher)SetNetService(net *NetService){
	self.net = net
}

func (self *Dispatcher)Dispatch(spiders... *BaseSpider){
	go self.do_dispatch_request()
	go self.do_dispatch_item()

	for _, spider := range spiders{
		self.dispatch_request_or_items(spider.GetStartRequests(nil))
	}

	self.do_dispatch_job()
}

func (self *Dispatcher)do_dispatch_request(){
	for {
		request := <-self.requestChan
		append(self.requests, request)//先扩容
		if cur_mode := self.config["mode"]; cur_mode == DISPATCH_MODE_DEPTH{
			//队列头部插入
			copy_len := len(self.requests)
			copy(self.requests[1:copy_len], self.requests[0:copy_len - 1])
			self.requests[0] = request
		}
		select {
		case <-self.stopJobChan://必须要有一个缓冲才好
			self.startJobChan<-1
		default:
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

package main
import (
	"fmt"
	"time"
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

	itemChan chan IBaseItem

	requestCacheChan chan *Request
	requestChan chan *Request

	responseChan chan *ResponsePack

	heartChan chan int
}

func NewDispatcher() *Dispatcher{
	MAX_PROCS := runtime.GOMAXPROCS(0)
	fmt.Printf("max procs %v\n", MAX_PROCS)
	self := &Dispatcher{
		config:map[string]int{
			"mode":DISPATCH_MODE_WIDTH,//遍历模式
			"max_procs":MAX_PROCS,
		},
		net:nil,

		itemChan:make(chan IBaseItem, MAX_PROCS),
		requestCacheChan:make(chan *Request),//这个只要一路就好了,总管分发
		requestChan:make(chan *Request, MAX_PROCS),
		responseChan:make(chan *ResponsePack, MAX_PROCS),
		heartChan:make(chan int, MAX_PROCS),
	}
	return self
}

func (self *Dispatcher)SetConfig(config map[string]int){
	self.config = config
}

func (self *Dispatcher)SetNetService(net *NetService){
	self.net = net
}

func (self *Dispatcher)Dispatch(spiders... IBaseSpider){
	go self.dispatchRequests(spiders...)

	for i := 0; i < Max(cap(self.requestChan), 1); i++ {
		go self.handleRequest(i)
	}
	for i := 0; i < Max(cap(self.itemChan), 1); i++ {
		go self.handleItem(i)
	}
	for i := 0; i < Max(cap(self.responseChan), 1); i++ {
		go self.handleResponse(i)
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
		case v := <-self.heartChan:
			fmt.Printf("heart beating....%v\n", v)
		case <-timer.C:
			break OUTER_LOOP
		}
	}
}

func (self *Dispatcher)dispatchRequests(spiders ...IBaseSpider){
	const MAX_BUFFER_COUNT = 200
	var requests []*Request = make([]*Request, 0, MAX_BUFFER_COUNT)
	insertRequest := func (request *Request){
		requests = append(requests, request)//应该是拷贝有用的那一块来扩充吧
		if mode := self.config["mode"]; mode == DISPATCH_MODE_DEPTH &&
		len(requests) > cap(self.requestChan){//考虑到并行度,这种情况下拷贝才有价值
			copy(requests[1:], requests[:len(requests) - 1]) //头部插入
			requests[0] = request
		}
	}

	for _, spider := range spiders {
		requests = append(requests, spider.GetStartRequests(nil)...)
	}

	//handledRequestCount := 0
	for {
		if len(requests) != 0 {
			select {
			case request := <-self.requestCacheChan://优先这个,好大概保持个顺序
				insertRequest(request)
			default:
				select{
					case request := <-self.requestCacheChan:
						insertRequest(request)
					case self.requestChan <- requests[0]:
					//未测试
						requests = requests[1:]
					//这个操作导致数组越来越大,要适当时候处理下,怎么处理呢,背后的数组都找不到头了,只有拷贝一下
					//上面的理解可能有误,不会拷贝数组的无用部分的.具体还是得看看源代码了
					//handledRequestCount += 1
					//if handledRequestCount > MAX_BUFFER_COUNT / 2 {
					//newRequests := make([]*Request, len(requests), MAX_BUFFER_COUNT)
					//copy(newRequests, requests)
					//requests = newRequests
					//}

					}
			}
		} else{
			request := <-self.requestCacheChan
			insertRequest(request)
		}
		//self.heartChan<- 0#不能有这个,会卡主主线程的
	}

}

func (self *Dispatcher)handleRequest(ind int){
	sendRequest := func (request *Request){
		self.heartChan<- 1
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
		self.heartChan<- 2
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
		self.heartChan<- 3
		fmt.Printf("handleItem receive item %v with go %v\n", item, ind)
	}
}

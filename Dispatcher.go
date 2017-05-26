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

	requestChan chan *Request
	requestHeadChan chan *Request

	responseChan chan *ResponsePack

	heartChan chan string
}

func NewDispatcher() *Dispatcher{
	//const MAX_PROCS = runtime.GOMAXPROCS(0)
	const MAX_PROCS = 0
	fmt.Printf("max procs %v\n", MAX_PROCS)
	self := &Dispatcher{
		config:map[string]int{
			"mode":DISPATCH_MODE_DEPTH,//遍历模式
		},
		net:nil,

		itemChan:make(chan IBaseItem, MAX_PROCS),
		requestChan:make(chan *Request, MAX_PROCS),
		requestHeadChan:make(chan *Request, MAX_PROCS),
		responseChan:make(chan *ResponsePack, MAX_PROCS),
		//heartChan:make(chan string, MAX_PROCS),
		heartChan:make(chan string),
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

	fmt.Printf("Dispatching spiders\n")
	go func() {
		for _, spider := range spiders {
			for _, request := range spider.GetStartRequests(nil) {
				if request, ok := request.(*Request); ok {
					self.requestChan <- request
				}
			}
		}
	}()
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

func (self *Dispatcher)handleRequest(ind int){
	sendRequest := func (request *Request){
		self.heartChan<- "handleRequest"
		fmt.Printf("handleRequest send request %v with go %v\n", request.url, ind)
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
	for {
		responsePack := <-self.responseChan
		self.heartChan<- "handleResponse"
		fmt.Printf("handleResponse receive response %v with go %v\n", responsePack.response.url, ind)
		requestOrItems := responsePack.callback(responsePack.response)
		for _, requestOrItem := range requestOrItems {
			switch requestOrItem := requestOrItem.(type) {
			case *Request://先判断这个，因为IBaseItem太多了
				go func(request *Request){//Request和Response容易造成相互等待，看看如何更轻量化， 避免太多协程了
					if mode := self.config["mode"]; mode == DISPATCH_MODE_DEPTH {
						self.requestHeadChan <- requestOrItem
					} else if mode == DISPATCH_MODE_WIDTH {
						self.requestChan<- requestOrItem
					}
				}(requestOrItem)
			case IBaseItem:
				self.itemChan<- requestOrItem
			default:
				fmt.Printf("received unhandled item %T\n", requestOrItem)
			}
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

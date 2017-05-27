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
type RoutineStatus struct {
	running bool
	//registerTime time.Time
	routineInd int
}

func (self *RoutineStatus)String()string{
	return fmt.Sprintf("running:%v, routineInd:%v", self.running, self.routineInd)
}

type ItemMidwareFunc func (IBaseItem) bool
type RequestMidwareFunc func (*Request) bool
type ResponseMidwareFunc func (*Response) bool

type Dispatcher struct {
	config map[string]int
	net *NetService
	requests []*Request

	requestMidware []RequestMidwareFunc
	responseMidware []ResponseMidwareFunc
	itemMidware []ItemMidwareFunc

	itemChan chan IBaseItem
	requestCacheChan chan *Request
	requestChan chan *Request
	responseChan chan *ResponsePack
	heartChan chan *RoutineStatus
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
		heartChan:make(chan *RoutineStatus, MAX_PROCS),
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
	totalRoutineCount := 0
	go self.dispatchRequests(totalRoutineCount, spiders...)
	totalRoutineCount += 1

	routineCount := Max(cap(self.requestChan), 1)
	for i := 0; i < routineCount; i++ {
		go self.handleRequest(totalRoutineCount + i)
	}
	totalRoutineCount += routineCount
	routineCount = Max(cap(self.itemChan), 1)
	for i := 0; i < routineCount; i++ {
		go self.handleItem(totalRoutineCount + i)
	}
	totalRoutineCount += routineCount
	routineCount = Max(cap(self.itemChan), 1)
	for i := 0; i < routineCount; i++ {
		go self.handleResponse(totalRoutineCount + i)
	}
	totalRoutineCount += routineCount
	self.waitAllDone(totalRoutineCount)
	fmt.Printf("All Done!!\n")
}

func (self *Dispatcher)waitAllDone(totalRoutineCount int){
	const DURATION = 2 * time.Second
	timer := time.NewTimer(DURATION)

	var subRoutinesStatus [2][]*RoutineStatus
	for ind := range subRoutinesStatus{
		subRoutinesStatus[ind] = make([]*RoutineStatus, totalRoutineCount)
	}
	curInd := 0

	OUTER_LOOP:
	for {
		timer.Reset(DURATION)
		select {
		case v := <-self.heartChan:
			subRoutinesStatus[curInd][v.routineInd] = v
			fmt.Printf("heart beating....%v\n", v)
		case <-timer.C:
			fmt.Printf("time out \n")
			isAllNil := true
			for _, status := range subRoutinesStatus[curInd]{
				if status != nil && status.running{//有活着的超时了
					fmt.Printf("has some routine working\n")
					continue OUTER_LOOP
				} else if status != nil{
					isAllNil = false
				}
			}
			if isAllNil{
				fmt.Printf("all nil, get out\n")
				break OUTER_LOOP//真正全部结束了
			}

			fmt.Printf("give another chance to wait\n")
			curInd = (curInd + 1) % 2
			//等下一次超时过来，看看你们是不是真的全死了
			for ind := range subRoutinesStatus[curInd]{
				subRoutinesStatus[curInd][ind] = nil
			}
		}
	}
}

func (self *Dispatcher)registerStatus(ind int, running bool){
	//self.heartChan<- &RoutineStatus{running:running, registerTime:time.Now(), routineInd:ind}
	self.heartChan<- &RoutineStatus{running:running, routineInd:ind}
}

func (self *Dispatcher)dispatchRequests(ind int, spiders ...IBaseSpider){
	self.registerStatus(ind, true)
	const MAX_BUFFER_COUNT = 200
	var requests []*Request = make([]*Request, 0, MAX_BUFFER_COUNT)
	handleMidware := func(request *Request){
		for _, midware := range self.requestMidware{
			if ! midware(request){//提供终止服务
				break
			}
		}
	}
	insertRequest := func (request *Request){
		self.registerStatus(ind, true)
		handleMidware(request)
		if request == nil{
			return
		}
		requests = append(requests, request)//应该是拷贝有用的那一块来扩充吧
		if mode := self.config["mode"]; mode == DISPATCH_MODE_DEPTH &&
		len(requests) > cap(self.requestChan){//考虑到并行度,这种情况下拷贝才有价值
			copy(requests[1:], requests[:len(requests) - 1]) //头部插入
			requests[0] = request
		}
	}


	for _, spider := range spiders {
		for _, request := range spider.GetStartRequests(nil){
			handleMidware(request)
			if request == nil{
				return
			}
			requests = append(requests, request)
		}
	}

	//handledRequestCount := 0
	for {
		self.registerStatus(ind, false)
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
						self.registerStatus(ind, true)
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
	for {
		self.registerStatus(ind, false)
		request := <-self.requestChan
		self.registerStatus(ind, true)
		fmt.Printf("handleRequest send request %v with go %v\n", request, ind)
		response := self.net.SendRequest(request)
		for _, midware := range self.responseMidware{
			if !midware(response){
				break
			}
		}
		if response == nil{
			return
		}
		self.responseChan<- &ResponsePack{response, request.callback}
	}
}

func (self *Dispatcher)handleResponse(ind int){
	for {
		self.registerStatus(ind, false)
		responsePack := <-self.responseChan
		self.registerStatus(ind, true)
		fmt.Printf("handleResponse receive response %v with go %v\n", responsePack.response, ind)
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
		self.registerStatus(ind, false)
		item := <-self.itemChan
		self.registerStatus(ind, true)
		for _, midware := range self.itemMidware{
			if ! midware(item){//提供终止服务
				break
			}
		}
		fmt.Printf("handleItem receive item %v with go %v\n", item, ind)
	}
}

func (self *Dispatcher)AddRequestMidware(midware RequestMidwareFunc){
	self.requestMidware = append(self.requestMidware, midware)
}
func (self *Dispatcher)AddResponseMidware(midware ResponseMidwareFunc){
	self.responseMidware = append(self.responseMidware, midware)
}
func (self *Dispatcher)AddItemMidware(midware ItemMidwareFunc){
	self.itemMidware = append(self.itemMidware, midware)
}

package main

type RequestOrItems []interface{}//会产生重新拷贝吗?
type RequestCallback func (*Response)RequestOrItems

type Request struct {
	method string
	url string
	data map[string]string
	callback RequestCallback
}

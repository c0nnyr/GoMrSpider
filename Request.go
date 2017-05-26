package main

type RequestOrItems []interface{}//会产生重新拷贝吗?
type RequestCallback func (*Response)([]*Request, []IBaseItem)

type Request struct {
	method string
	url string
	data map[string]string
	callback RequestCallback
}

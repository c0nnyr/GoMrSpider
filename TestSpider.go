package main

type TestSpider struct{
	BaseSpider
}

func NewTestSpider() *TestSpider{
	spider := &TestSpider{
	}
	append(spider.start_urls, "http://baidu.com")
	return spider
}

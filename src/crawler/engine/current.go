package engine

import (
	"crawler/model"
)

type Processor func(Request) (ParserResult,error)
type CurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemServer chan model.Item
	Processor Processor
}

type Scheduler interface {
	Submit(Request)
	WorkMaskChan() chan Request
	WorkerReady(chan Request)
	Run()
}

//注意代码顺序问题  worker 必须是在创建完之后才能使用  所以先走Run
func (c *CurrentEngine) Run(seeds ...Request) {
	c.Scheduler.Run()

	out := make(chan ParserResult)
	for i := 0; i < c.WorkerCount; i++ {
		c.createWorker(c.Scheduler.WorkMaskChan(), out, c.Scheduler)
	}

	for _, r := range seeds {
		if isDuplicate(r.Url){
			continue
		}
		c.Scheduler.Submit(r)
	}

	for {
		result := <-out
		for _, item := range result.Items {
			go func() { c.ItemServer <- item }()
		}

		//这里去重
		for _, p := range result.Requests {
			if isDuplicate(p.Url){
				continue
			}
			c.Scheduler.Submit(p)
		}
	}

}

func (c *CurrentEngine) createWorker(in chan Request, out chan ParserResult, s Scheduler) {
	//一般网络请求是并发的  所有用go

	go func() {

		for {
			s.WorkerReady(in)
			result := <-in
			parserResult, err := c.Processor(result)
			if err != nil {
				continue
			}

			out <- parserResult
		}
	}()
}

var isDupUrls = make(map[string]bool)
func isDuplicate(url string) bool {
	if isDupUrls[url] {
		return true
	}

	isDupUrls[url] = true
	return false
}

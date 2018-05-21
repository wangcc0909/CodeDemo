package simple

import "crawler/engine"

type SimpleSheduler struct {
	WorkChan chan engine.Request
}

func (s *SimpleSheduler) WorkMaskChan() chan engine.Request {
	return s.WorkChan
}

func (s *SimpleSheduler) WorkerReady(chan engine.Request) {

}

func (s *SimpleSheduler) Run() {
	s.WorkChan = make(chan engine.Request)
}


func (s *SimpleSheduler) Submit(r engine.Request) {
	go func() {
		s.WorkChan <- r
	}()
}




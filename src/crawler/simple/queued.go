package simple

import (
	"crawler/engine"
)

type QueuedScheduler struct {
	RequestChan chan engine.Request

	WorkerChan chan chan engine.Request
}

func (s *QueuedScheduler) Submit(r engine.Request) {
	s.RequestChan <- r
}

func (s *QueuedScheduler) WorkMaskChan() chan engine.Request {
	return make(chan engine.Request)
}

func (s *QueuedScheduler) WorkerReady(c chan engine.Request) {
	s.WorkerChan <- c
}

func (s *QueuedScheduler) Run() {
	s.RequestChan = make(chan engine.Request)
	s.WorkerChan = make(chan chan engine.Request)
	go func() {

		var WorkerQ []chan engine.Request
		var RequestQ []engine.Request

		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(RequestQ) > 0 && len(WorkerQ) > 0 {
				activeRequest = RequestQ[0]
				activeWorker = WorkerQ[0]
			}

			select {
			case r := <-s.RequestChan:
				RequestQ = append(RequestQ,r)
			case w := <-s.WorkerChan:
				WorkerQ = append(WorkerQ,w)
			case activeWorker <- activeRequest:
				 WorkerQ = WorkerQ[1:]
				 RequestQ = RequestQ[1:]
			}
		}
	}()
}

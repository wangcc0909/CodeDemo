package engine

import (
	"log"
)

type SimpleEngine struct {
}

func (s SimpleEngine) Run(seeds ...Request) {
	var requests []Request

	for _,r := range seeds{
		requests = append(requests, r)
	}

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		result,err := Worker(r)
		if err != nil {
			continue
		}

		requests = append(requests,result.Requests...)

		for _,item := range result.Items{
			log.Printf("get item : %v\n",item)
		}

	}

}



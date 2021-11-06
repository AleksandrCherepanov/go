package main

import (
	"fmt"
	"time"
)

func main() {
	inputData := []int{0, 1, 3, 4, 5, 6, 7, 8, 9, 10}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data := dataRaw.(string)
			fmt.Println(data)
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)

	fmt.Println(end)
}

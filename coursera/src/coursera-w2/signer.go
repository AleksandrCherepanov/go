package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func crc32Async(input string, in chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	in <- DataSignerCrc32(input)
}

func closeChannel(wg *sync.WaitGroup, ch chan string) {
	wg.Wait()
	close(ch)
}

// SingleHash function
func SingleHash(in, out chan interface{}) {
	crc32Channel := make(chan string)
	crc32Wg := &sync.WaitGroup{}

	crc32Md5Channel := make(chan string)
	crc32Md5Wg := &sync.WaitGroup{}

	for value := range in {
		intValue := value.(int)

		stringData := strconv.Itoa(intValue)

		crc32Wg.Add(1)
		go crc32Async(stringData, crc32Channel, crc32Wg)

		crc32Md5Wg.Add(1)
		go crc32Async(DataSignerMd5(stringData), crc32Md5Channel, crc32Md5Wg)
	}

	go closeChannel(crc32Wg, crc32Channel)
	go closeChannel(crc32Md5Wg, crc32Md5Channel)

	var a []string
	for result := range crc32Channel {
		a = append(a, result)
	}

	var b []string
	for result := range crc32Md5Channel {
		b = append(b, result)
	}

	for i, v := range a {
		out <- strings.Join([]string{v, b[i]}, "~")
	}
}

func asyncFor(value interface{}, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	channels := make([]chan string, 6, 6)
	hashParts := make([]string, 6, 6)

	for i := 0; i < 6; i++ {
		channels[i] = make(chan string)
		go func(input string, in chan<- string) {
			in <- DataSignerCrc32(input)
		}(strconv.Itoa(i)+value.(string), channels[i])
	}

	hashParts[0] = <-channels[0]
	hashParts[1] = <-channels[1]
	hashParts[2] = <-channels[2]
	hashParts[3] = <-channels[3]
	hashParts[4] = <-channels[4]
	hashParts[5] = <-channels[5]

	out <- strings.Join(hashParts, "")
}

// MultiHash function
func MultiHash(in, out chan interface{}) {
	forOut := make(chan string)
	wg := &sync.WaitGroup{}

	for input := range in {
		wg.Add(1)
		go asyncFor(input, forOut, wg)
	}

	go closeChannel(wg, forOut)
	for result := range forOut {
		out <- result
	}
}

// CombineResults function
func CombineResults(in, out chan interface{}) {
	var result []string

	for value := range in {
		result = append(result, value.(string))
	}

	sort.Strings(result)
	out <- strings.Join(result, "_")
}

func workerExec(worker job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	worker(in, out)
	close(out)
}

// ExecutePipeline run functions as a one flow
func ExecutePipeline(input ...job) {
	var channels []chan interface{}

	for i := 0; i <= len(input); i++ {
		channels = append(channels, make(chan interface{}, 100))
	}

	wg := &sync.WaitGroup{}
	for i, worker := range input {
		wg.Add(1)
		go workerExec(worker, channels[i], channels[i+1], wg)
	}
	wg.Wait()
}

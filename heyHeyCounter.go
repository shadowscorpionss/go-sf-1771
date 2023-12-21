package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type HeyHeyCounterInterface interface {
	Start(routineCount, endCount int) <-chan int
}

type heyHeyCounterClass struct {
	send        chan int
	endCount    int
	counter     int
	countLock   *sync.RWMutex
	startLock   *sync.Mutex
	started     bool
	checkSignal *sync.Cond
	exitChan    chan int
}

func (h *heyHeyCounterClass) check() {
	for {
		//wait till any count
		h.checkSignal.L.Lock()
		h.checkSignal.Wait()
		h.checkSignal.L.Unlock()

		h.countLock.RLock()
		if h.counter >= h.endCount {
			h.exitChan <- h.counter
			close(h.exitChan)
			break
		}
		h.countLock.RUnlock()
	}

	h.startLock.Lock()
	h.started = false
	h.startLock.Unlock()

}

func (h *heyHeyCounterClass) count() {
	immer := true

	go func() {
		<-h.exitChan
		immer = false

	}()

	for immer {
		h.countLock.Lock()
		h.counter += <-h.send
		fmt.Printf("Расчет (%d рутин): %d\n", runtime.NumGoroutine(), h.counter)
		h.countLock.Unlock()
	}
	close(h.send)

}

func (h *heyHeyCounterClass) sender(cnt int, send chan<- int) {
	for i := 0; i < cnt; i++ {
		send <- 1
		h.checkSignal.Signal()
		//искусственная случайная задержка
		p := rand.Intn(100)
		time.Sleep(time.Duration(p) * time.Millisecond)
	}

}

func (h *heyHeyCounterClass) Start(routineCount, endCount int) <-chan int {

	h.startLock.Lock()
	if h.started {
		defer h.startLock.Unlock()
		return h.exitChan
	}
	h.started = true
	h.startLock.Unlock()

	h.countLock.Lock()
	h.endCount = endCount
	h.countLock.Unlock()

	h.exitChan = make(chan int)

	perRoutine := endCount / routineCount
	capacity := perRoutine + 10 //parallel writing
	rests := endCount - perRoutine*routineCount

	h.send = make(chan int, capacity)

	for i := 0; i < routineCount; i++ {

		go h.sender(perRoutine, h.send)
	}	

	for rests > perRoutine {
		go h.sender(perRoutine, h.send)
		rests-=perRoutine
	}

	if rests>0 {
		go h.sender(rests, h.send)
	}


	go h.check()

	go h.count()

	return h.exitChan
}

func NewHeyHeyCounterClass() HeyHeyCounterInterface {
	res := &heyHeyCounterClass{
		countLock:   &sync.RWMutex{},
		startLock:   &sync.Mutex{},
		checkSignal: &sync.Cond{L: &sync.Mutex{}},
	}

	return res
}

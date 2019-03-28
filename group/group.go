package group

import "sync"

var wg = new(sync.WaitGroup)

func Add(c int) {
	wg.Add(c)
}

func Done() {
	wg.Done()
}

func Wait() {
	wg.Wait()
}

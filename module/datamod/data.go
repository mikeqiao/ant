package datamod

import "sync"

type Data interface {
	Close() bool
	Start() bool
	Run()
	GetKey() interface{}
	Update()
}

type BaseData struct {
	Data  interface{}
	Key   interface{}
	state bool
	mutex sync.Mutex // 锁
}

//example
func (b *BaseData) Start() bool {
	return true
}

func (b *BaseData) Close() bool {
	return true
}

func (b *BaseData) GetKey() interface{} {
	return b.Key
}

func (b *BaseData) Run() {

}

func (b *BaseData) Update() {
	if !b.state {
		return
	}
	b.mutex.Lock()
	b.mutex.Unlock()

}

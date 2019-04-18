package mod

import (
	"sync"
	"time"

	"github.com/mikeqiao/ant/group"
)

type DataMod struct {
	Dlist    map[interface{}]Data
	mutex    sync.Mutex // 锁
	closeSig chan bool
}

func (d *DataMod) Init() {
	d.Dlist = make(map[interface{}]Data)
	d.closeSig = make(chan bool, 1)
}

func (d *DataMod) Start() {

	group.Add(1)
	go d.Run()
}

func (d *DataMod) Close() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	for k, v := range d.Dlist {
		if nil != v {
			v.Close()
			delete(d.Dlist, k)
		}
	}
}

func (d *DataMod) Add(data Data) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if nil != data {
		if data.Start() {
			d.Dlist[data.GetKey()] = data
		}
	}
}

func (d *DataMod) Del(key interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if v, ok := d.Dlist[key]; ok {
		if v.Close() {
			delete(d.Dlist, key)
		}
	}
}

func (d *DataMod) Get(key interface{}) interface{} {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if v, ok := d.Dlist[key]; ok {
		return v
	}
	return nil
}

func (d *DataMod) Run() {
	t1 := time.NewTimer(time.Second * 1)
	for {
		select {
		case <-d.closeSig:
			d.Close()
			goto Loop
		case <-t1.C:
			d.Update()
			t1.Reset(time.Second * 1)
		}
	}
Loop:
	group.Done()
}

func (d *DataMod) Update() {
	for _, v := range d.Dlist {
		if nil != v {
			v.Update()
		}
	}
}

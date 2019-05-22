package mod

import (
	"sync"

	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

var ModuleControl *ModuleManager

func Init() {
	ModuleControl = NewModuleManager()
}

func NewModuleManager() *ModuleManager {
	m := new(ModuleManager)
	m.Init()
	return m
}

type ModuleManager struct {
	ml       map[int64]*Module
	createid int64
	mutex    sync.RWMutex
}

func (m *ModuleManager) Init() {
	m.ml = make(map[int64]*Module)
	m.createid = 0
}

func (m *ModuleManager) Run() {
	for k, v := range m.ml {
		if nil != v && v.NeedStart() {
			v.Start()
		} else {
			log.Debug("this mod  is nil:%v", k)
			delete(m.ml, k)
		}
	}
}

func (m *ModuleManager) Close() {
	m.mutex.RLock()
	for _, v := range m.ml {
		if nil != v {
			v.Close()
		}
	}
	m.mutex.RUnlock()
}

func (m *ModuleManager) GetNewID() int64 {
	m.mutex.Lock()
	m.createid += 1
	m.mutex.Unlock()
	return m.createid
}

func (m *ModuleManager) RegisterMod(mod *Module) {
	if nil == mod {
		return
	}
	m.mutex.Lock()
	if _, ok := m.ml[mod.mid]; ok {
		log.Debug("already have this mod:%v", mod)
	} else {
		m.ml[mod.mid] = mod
		//	log.Debug("add mod:%v", mod)
	}
	m.mutex.Unlock()
}

func (m *ModuleManager) DelMod(id int64) {
	m.mutex.Lock()
	if v, ok := m.ml[id]; ok {
		log.Debug("del mod:%v", v)
		v.Destroy()
		delete(m.ml, id)

	} else {
		log.Debug("not have this mod id:%v", id)
	}
	m.mutex.Unlock()
}

func (m *ModuleManager) CloseMod(id int64) {
	m.mutex.Lock()
	if v, ok := m.ml[id]; ok {
		if nil != v {
			v.Close()
		}
	} else {
		log.Debug("not have this mod id:%v", id)
	}
	m.mutex.Unlock()
}

func (m *ModuleManager) StartMod(id int64) {
	m.mutex.Lock()
	if v, ok := m.ml[id]; ok {
		if nil != v {
			v.Start()
		}
	} else {
		log.Debug("not have this mod id:%v", id)
	}
	m.mutex.Unlock()
}

func (m *ModuleManager) UpdatMod(id, state int64) {
	m.mutex.Lock()
	if v, ok := m.ml[id]; ok {
		if nil != v {
			v.Update(state)
		}
	} else {
		log.Debug("not have this mod id:%v", id)
	}
	m.mutex.Unlock()
}

func (m *ModuleManager) GetModule(uid int64) *Module {
	m.mutex.RLock()
	if v, ok := m.ml[uid]; ok {
		if nil != v {
			m.mutex.RUnlock()
			return v
		}
	}
	m.mutex.RUnlock()
	return nil
}

func Route(mid int64, fid, did uint32, cb interface{}, in interface{}, data *net.UserData) {
	m := ModuleControl.GetModule(mid)
	if nil != m {
		m.Route(fid, did, cb, in, data)
	} else {
		log.Error("not have this mod id :%v", mid)
		ExecCb(cb)
	}
}

func Run() {
	ModuleControl.Run()
}

func Close() {
	ModuleControl.Close()
}

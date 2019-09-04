package mod

import (
	"fmt"
	"runtime"
	"sync"

	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

var RPC = NewRPC()

func NewRPC() *RPCServer {
	r := new(RPCServer)
	r.Init()
	return r
}

type RServer struct {
	ml    map[int64]*Module //服务提供者
	fid   uint32            //所提供的服务(type)
	did   uint32
	ftype int32 // 0 本地服务 1 远程服务
}

func (r *RServer) Init(id uint32) {
	r.ml = make(map[int64]*Module)
	r.fid = id
}

func (r *RServer) SetType(t int32) {
	r.ftype = t
}

func (r *RServer) Register(s *Module) {
	if nil == s {
		return
	}
	if _, ok := r.ml[s.mid]; ok {
		log.Debug("already have this module:%v, fid:%v", s, r.fid)
	} else {
		r.ml[s.mid] = s
	}
}

func (r *RServer) DelModule(s *Module) {
	if nil == s {
		return
	}
	if _, ok := r.ml[s.mid]; ok {
		delete(r.ml, s.mid)
	} else {
		log.Debug("not have this module:%v, fid:%v", s, r.fid)
	}
}

func (r *RServer) Route(fid uint32, cb interface{}, in interface{}, data *net.UserData) {
	//	log.Debug(" start find call time:%v", time.Now().String())
	did := r.fid
	var m *Module
	for k, v := range r.ml {
		if nil != v {
			if nil == m || m.v < v.v || m.s > v.s {
				m = v
			}
		} else {
			delete(r.ml, k)
		}
	}
	if nil != m {
//		log.Debug("mod id:%v", m.mid)
		m.Route(fid, did, cb, in, data)
	} else {
		log.Debug("not find working module, fid:%v", did)
		ExecCb(cb)
	}
}

//服务内容管理
type RPCServer struct {
	Functions map[uint32]*RServer
	mutex     sync.Mutex // 读写锁
}

func (r *RPCServer) Init() {
	r.Functions = make(map[uint32]*RServer)
}

func (r *RPCServer) Register(id uint32, s *Module) {
	r.mutex.Lock()
	if v, ok := r.Functions[id]; ok {
		if nil != v {
			v.Register(s)
			v.SetType(0)
		} else {
			delete(r.Functions, id)
		}
	} else {
		g := new(RServer)
		g.Init(id)
		g.Register(s)
		g.SetType(0)
		r.Functions[id] = g
	}
	r.mutex.Unlock()
}

func (r *RPCServer) RegisterRemote(id uint32, s *Module) {
	r.mutex.Lock()
	if v, ok := r.Functions[id]; ok {
		if nil != v {
			v.Register(s)
		} else {
			delete(r.Functions, id)
		}
	} else {
		g := new(RServer)
		g.Init(id)
		g.Register(s)
		g.SetType(1)
		r.Functions[id] = g
	}
	r.mutex.Unlock()
}

func (r *RPCServer) DelFunction(id uint32, s *Module) {
	r.mutex.Lock()
	if v, ok := r.Functions[id]; ok {
		if nil != v {
			v.DelModule(s)
		} else {
			delete(r.Functions, id)
		}
	}
	r.mutex.Unlock()
}

func (r *RPCServer) Route(fid, did uint32, cb interface{}, in interface{}, data *net.UserData) {
	if nil == in {
		log.Debug("in is nil ")
		return
	}
	if v, ok := r.Functions[did]; ok {
		if nil != v {
		//	log.Debug("have did:%v ", did)
			v.Route(fid, cb, in, data)
		}
	} else {
		if fid != HandleForwardMsg {
			if fv, ok := r.Functions[HandleForwardMsg]; ok && nil != fv && 1 == fv.ftype {
				fv.Route(fid, cb, in, data)
			} else {
				ExecCb(cb)
			}
		} else {
			ExecCb(cb)
		}

	}
}

func ExecCb(cb interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if conf.Config.LenStackBuf > 0 {
				buf := make([]byte, conf.Config.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Debug("%v: %s", r, buf[:l])
			} else {
				log.Debug("%v", r)
			}
		}
	}()
	if nil == cb {
		return
	}
	err := fmt.Errorf("no this func")
	cb.(func(interface{}, error))(nil, err)
	return
}

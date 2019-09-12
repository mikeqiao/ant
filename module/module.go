package mod

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mikeqiao/ant/group"
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
	"github.com/mikeqiao/ant/net/proto"
	"github.com/mikeqiao/ant/rpc"
)

type Module struct {
	mid      int64
	s        int64 //繁忙状态值（越小的越优先选择）
	v        int32 //版本号（越大的越优先选择）
	chanlen  int32
	state    int32 //模块状态值 1 开机启动， 0 触发启动（服务器启动时stop状态)
	working  bool
	closeSig chan bool
	server   *rpc.Server
	client   *rpc.Client
	Agent    *net.TcpAgent
	waitback map[string]*rpc.CallInfo
	lock     sync.Mutex
	wg       sync.WaitGroup
	fl       map[uint32]bool //服务列表
}

func NewModule(id, state int64, version, start int32) *Module {
	if 0 == id {
		id = ModuleControl.GetNewID()
	}
	m := new(Module)
	m.mid = id
	m.s = state
	m.v = version
	m.state = start
	m.working = false
	m.chanlen = 1024
	m.Init()
	ModuleControl.RegisterMod(m)
	return m
}

func (m *Module) SetAgent(a *net.TcpAgent) {
	if nil != a {
		a.SetRemotUID(m.mid)
		a.SetLogin()
		m.Agent = a
	}
}

func (m *Module) Init() {
	m.closeSig = make(chan bool, 1)
	m.fl = make(map[uint32]bool)
	m.waitback = make(map[string]*rpc.CallInfo)
	m.server = rpc.NewServer(m.chanlen)
	m.client = rpc.NewClient(m.chanlen, m.server)
	if nil != m.Agent {
		m.Agent.SetRemotUID(m.mid)
	}
}

func (m *Module) NeedStart() bool {
	if 1 == m.state {
		return true
	}
	return false
}

func (m *Module) Start() {
	log.Debug("modul start id:%v", m.mid)
	m.working = true
	group.Add(1)
	go m.Run()
}

func (m *Module) Close() {
	err := errors.New("Module colsed")
	for _, v := range m.waitback {
		if nil != v {

			v.SetResult(nil, err)
		}
	}

	//	log.Debug(" end find  callback time:%v", time.Now().String())

	m.closeSig <- true
}

func (m *Module) Update(state int64) {
	m.s = state
}

func (m *Module) Run() {
	for {
		select {
		case <-m.closeSig:
			m.server.Close()
			goto Loop
		case ri := <-m.client.ChanAsynRet:
			m.client.CallBack(ri)
		case ci := <-m.server.ChanCall:
			m.server.Exec(ci)
		}
	}
Loop:
	log.Debug("modul close id:%v", m.mid)
	m.working = false
	group.Done()
}

func (m *Module) Register(id uint32, f interface{}) {
	if nil == m.server {
		m.server = rpc.NewServer(m.chanlen)
	}
	m.fl[id] = true
	m.server.Register(id, f)
	RPC.Register(id, m)
}

func (m *Module) RegisterRemote(id uint32, f interface{}) {
	if nil == m.server {
		m.server = rpc.NewServer(m.chanlen)
	}
	m.fl[id] = true
	m.server.Register(id, f)
	RPC.RegisterRemote(id, m)
}

func (m *Module) Route(fid, did uint32, cb interface{}, in interface{}, data *net.UserData) {

	if !m.working {
		log.Error("module not working")
		ExecCb(cb)
		return
	}
	if nil != m.client {

		m.client.CallAsyn(m.mid, fid, did, cb, in, data)
	} else {
		log.Error("client not working")
	}

	//	log.Debug(" end find call time:%v", time.Now().String())
}

func (m *Module) DelFunc(id uint32) {
	if nil == m.server {
		return
	}
	if _, ok := m.fl[id]; ok {
		delete(m.fl, id)
	}
	m.server.DelFunc(id)
	RPC.DelFunction(id, m)
}

func (m *Module) Destroy() {
	if m.working {
		m.Close()
	}
	if nil == m.server {
		return
	}
	for k, _ := range m.fl {
		m.server.DelFunc(k)
		RPC.DelFunction(k, m)
	}
	m.fl = make(map[uint32]bool)
}

func (m *Module) GetVersion() int32 {
	return m.v
}

func (m *Module) GetFunc() (fl []uint32) {
	if nil != m.server {
		fl = m.server.GetFunc()
	}
	return
}

func (m *Module) SendFunc(t int32, uid int64) {
	msg := new(proto.ServerRegister)
	msg.Type = t
	msg.Uid = uid
	for k, v := range RPC.Functions {
		if nil != v && 0 == v.ftype {
			msg.Fuid = append(msg.Fuid, k)
		}
	}
	if nil != m.Agent {
		m.Agent.WriteMsg(msg)
	}

}

func (m *Module) AddWaitCall(c *rpc.CallInfo) (key string) {
	if nil != c && nil != c.Cb {
		key = fmt.Sprintf("%p", c)
		m.lock.Lock()
		m.waitback[key] = c
		m.lock.Unlock()
	}
	return
}

func (m *Module) ExecBack(key string, msg interface{}, e error) {

	m.lock.Lock()
	c, ok := m.waitback[key]
	if ok {
		delete(m.waitback, key)
	}
	m.lock.Unlock()

	//	log.Debug(" end find  callback time:%v", time.Now().String())
	if nil != c {
		c.SetResult(msg, e)
	}
}

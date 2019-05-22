package network

import (
	"sync"
	"time"

	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

var CM *ConnManager

func NewConnManager() {
	CM = new(ConnManager)
	CM.Init()
}

type ConnManager struct {
	userConn    map[int64]*net.TcpAgent
	userNewConn map[int64]*net.TcpAgent
	mutexConns  sync.Mutex // 锁
	isClose     bool
}

func (this *ConnManager) Init() {
	this.userConn = make(map[int64]*net.TcpAgent)
	this.userNewConn = make(map[int64]*net.TcpAgent)
	this.isClose = false
}

func (this *ConnManager) AddUserConn(id int64, uid int64) bool {
	if this.isClose == true {
		return false
	}
	this.mutexConns.Lock()
	relink := false
	if v, ok := this.userNewConn[id]; ok {
		if _, ok := this.userConn[uid]; ok {
			log.Debug("userConn agent have add,  stop old add new, uid：%v,  id:%v", uid, id)
			this.userConn[uid].Close()
			//触发断线重连
			relink = true
		}
		this.userConn[uid] = v
		this.userConn[uid].SetLogin()
		this.userConn[uid].SetRemotUID(uid)

	} else {
		log.Debug("no this userNewConn, id:%v", id)
	}
	this.mutexConns.Unlock()
	this.DelNewUserConn(id)
	return relink
}

func (this *ConnManager) AddNewUserConn(id int64, agent *net.TcpAgent) {
	if this.isClose == true {
		return
	}
	this.mutexConns.Lock()
	if _, ok := this.userNewConn[id]; ok {
		log.Debug("userNewConn agent have add,  stop old add new,  id:%v", id)
		this.userNewConn[id].Close()
		this.userNewConn[id] = agent
	} else {
		this.userNewConn[id] = agent
	}
	this.mutexConns.Unlock()
}

func (this *ConnManager) DelNewUserConn(id int64) {
	if this.isClose == true {
		return
	}
	this.mutexConns.Lock()
	if _, ok := this.userNewConn[id]; ok {
		delete(this.userNewConn, id)
	} else {
		log.Debug("no this userNewConn, id:%v", id)
	}
	this.mutexConns.Unlock()

}

func (this *ConnManager) DelUserConn(id int64) {
	if this.isClose == true {
		return
	}
	this.mutexConns.Lock()
	if _, ok := this.userNewConn[id]; ok {
		delete(this.userNewConn, id)
	} else {
		log.Debug("no this userNewConn, id:%v", id)
	}
	if _, ok := this.userConn[id]; ok {
		delete(this.userConn, id)
	} else {
		log.Debug("no this userConn, id:%v", id)
	}
	this.mutexConns.Unlock()
}

func (this *ConnManager) Run() {
	go this.Updata()
}

func (this *ConnManager) Updata() {

	t1 := time.NewTimer(time.Second * 1)
	//	t2 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			if this.isClose == true {
				goto Loop
			}
			this.mutexConns.Lock()
			for k, v := range this.userConn {
				//v.agent.OnClose()

				if v.IsClose() == true {
					log.Debug("userConn close, id:%v", k)
					delete(this.userConn, k)
				}
			}
			for k, v := range this.userNewConn {
				//v.agent.OnClose()

				if v.IsClose() == true {
					log.Debug("userNewConn close, id:%v", k)
					delete(this.userNewConn, k)
				}
			}
			this.mutexConns.Unlock()

			t1.Reset(time.Second * 1)
		}

	}

Loop:
}

func (this *ConnManager) Close() {
	this.isClose = true
	this.mutexConns.Lock()
	for k, v := range this.userConn {
		if nil != v {
			v.Close()
		}
		delete(this.userConn, k)
	}
	for k, v := range this.userNewConn {
		if nil != v {
			v.Close()
		}
		delete(this.userNewConn, k)
	}
	this.mutexConns.Unlock()
}

func (this *ConnManager) ForwardClientMsg(id int64, uid int64, msg interface{}) {

	this.mutexConns.Lock()
	if _, ok := this.userNewConn[id]; ok {

		this.userNewConn[id].WriteMsg(msg)
	} else if _, ok := this.userConn[uid]; ok {
		this.userConn[uid].WriteMsg(msg)
	} else {
		log.Debug("no this userNewConn, id:%v", id)

	}
	this.mutexConns.Unlock()
}

func (this *ConnManager) ForwardClientNoticeMsg(uid []int64, msg interface{}) {
	this.mutexConns.Lock()
	for _, v := range uid {

		if _, ok := this.userConn[v]; ok {
			this.userConn[v].WriteMsg(msg)
		} else {
			log.Debug("no this userConn, id:%v", v)
		}
	}
	this.mutexConns.Unlock()
}

func (this *ConnManager) ForwardDBMsg(id, uid int64, rtype int32, column string, key, value interface{}, f func(m, a interface{})) {
	/*
		var c *ConnInfo
		this.mutexConns.Lock()
		if v, ok := this.serverConn[id]; ok {
			c = v
		} else {
			log.Debug("no this serverConn, id:%v", id)
		}
		this.mutexConns.Unlock()
		if nil != c {
			var cid int64
			if nil != f {
				cid = Call.Register(id, f)
			}
			kdata := this.serverConn[id].agent.GetForwardMsg(key)
			vdata := this.serverConn[id].agent.GetForwardMsg(value)
			if kdata != nil && vdata != nil {
				this.serverConn[id].agent.WriteMsg(
					&tproto.ServerDBMsgReq{
						Userid: uid,
						Fuid:   cid,
						Rtype:  rtype,
						Key:    kdata[:],
						Value:  vdata[:],
						Column: column,
					})
			} else {
				log.Debug("package err key:%v, value:%v", key, value)
			}
		}
	*/
}

func startManager() {
	CM.Run()
}

func closeManager() {
	CM.Close()
}

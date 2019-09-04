package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net/msgtype"
	"github.com/mikeqiao/ant/net/proto"
)

var CreatID = int64(1)

type UserData struct {
	UId     int64
	UsersId []int64
	Agent   *TcpAgent
}

type TcpAgent struct {
	LUId      int64 //本端local uid
	RUId      int64 //对方端remote uid
	conn      Conn
	Processor Processor
	lifetime  int64 //链接未验证有效期时间（秒）
	starttime int64 //链接开始的时间戳
	tick      int64 //上次心跳的时间戳
	islogin   bool  //是否登陆验证过
	isClose   bool  //是否关闭
	Closed    bool
	isUpdate  bool           //是否验证心跳
	ctype     int32          //连接类型 1 server  2 client
	wg        sync.WaitGroup // 链接wait
	userData  interface{}
	Version   int32
	Type      int32
}

func NewAgent(conn *TCPConn, tp Processor) *TcpAgent {
	a := new(TcpAgent)
	CreatID += 1
	a.RUId = CreatID
	a.conn = conn
	a.Processor = tp
	a.lifetime = 10
	a.starttime = time.Now().Unix()
	a.tick = time.Now().Unix()
	a.islogin = false
	a.isClose = false
	a.isUpdate = true
	a.Type = 1
	a.Closed = false
	return a
}

func NewAgentWeb(conn *WSConn, tp Processor) *TcpAgent {
	a := new(TcpAgent)
	CreatID += 1
	a.RUId = CreatID
	a.conn = conn
	a.Processor = tp
	a.lifetime = 10
	a.starttime = time.Now().Unix()
	a.tick = time.Now().Unix()
	a.islogin = false
	a.isClose = false
	a.isUpdate = false
	a.Type = 2
	a.Closed = false
	return a
}

func (a *TcpAgent) SetLocalUID(uid int64) {
	a.LUId = uid
}

func (a *TcpAgent) SetRemotUID(uid int64) {
	a.RUId = uid
}

func (a *TcpAgent) Start(name string) {
	if name != "Out" {
		a.ctype = 1
		return
	}
	msg := &proto.NewConnect{
		Id: a.RUId,
	}
	ta := &UserData{
		Agent: a,
		UId:   a.RUId,
	}
	err := a.Processor.Route(msgtype.NewConnect, msg, ta)
	if err != nil {
		log.Debug("Start Add Agent error: %v", err)
	}
	a.ctype = 2
}

func (a *TcpAgent) Run() {

	if a.isUpdate == true {
		a.wg.Add(1)
		go a.Update()
	}

	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Error("agent:localUID:%v, RemoteUID:%v, ,read message err: %v", a.LUId, a.RUId, err)
			goto Loop
		}
		if a.Processor != nil {
			id, msg, err := a.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
			}
			if msg != nil {
				ta := &UserData{
					Agent: a,
					UId:   a.RUId,
				}
				err = a.Processor.Route(id, msg, ta)
				if err != nil {
					log.Debug("route message error: %v", err)
				}
			}
		}
	}
Loop:
	a.isClose = true
	if a.isUpdate != true {
		a.Close()
	}
}

func (a *TcpAgent) Update() {
	t1 := time.NewTimer(time.Second * 1)
	t2 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			if a.isClose == true {
				log.Debug("agent closed")
				goto Loop
			}

			if a.islogin != true {
				nowtime := time.Now().Unix()
				if (a.starttime + a.lifetime) < nowtime {
					log.Debug("outtime to not login: %v", a.conn.RemoteAddr())
					goto Loop
				}
			} else {
				nowtime := time.Now().Unix()
				if (a.tick + a.lifetime*3) < nowtime {
					log.Debug("outtime to no tick: %v", a.conn.RemoteAddr())
					goto Loop
				}
			}
			t1.Reset(time.Second * 1)

		case <-t2.C:
			if a.isUpdate == true {
				nowtime := time.Now().Unix()
				a.WriteMsg(&proto.ServerTick{
					Time: nowtime,
				})
			}
			t2.Reset(time.Second * 10)
		}
	}
Loop:
	a.wg.Done()
	a.Close()
}

func (a *TcpAgent) IsClose() bool {
	if a.isClose == true {
		return true
	}
	return false
}

func (a *TcpAgent) GetProcessor() Processor {
	return a.Processor
}

func (a *TcpAgent) SetTick(time int64) {
	a.tick = time
}

func (a *TcpAgent) SetLogin() {
	a.islogin = true
}

func (a *TcpAgent) WriteMsg(msg interface{}) {
	if a.Processor != nil {
		id, data, err := a.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message id:%v error: %v", id, err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message id:%v error: %v", id, err)
		}
	}
}

func (a *TcpAgent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *TcpAgent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *TcpAgent) Close() {
	if a.Closed {
		return
	}
	a.Closed = true
	var msg interface{}
	var id uint16
	if a.ctype == 2 {
		msg = &proto.DelConnect{
			Id: a.RUId,
		}
		id = msgtype.DelConnect
	} else {
		msg = &proto.ServerDelConnect{
			Id: a.RUId,
		}
		id = msgtype.ServerDelConnect
	}
	ta := &UserData{
		Agent: a,
		UId:   a.RUId,
	}
	err := a.Processor.Route(id, msg, ta)
	if err != nil {
		log.Debug("Agent error: %v", err)

	}
	a.conn.Close()
	a.isClose = true
	a.wg.Wait()
}

func (a *TcpAgent) OnlyClose() {
	if a.Closed {
		return
	}
	a.Closed = true
	a.conn.Close()
	a.isClose = true
	a.wg.Wait()
}

func (a *TcpAgent) Destroy() {
	a.conn.Destroy()
}

func (a *TcpAgent) UserData() interface{} {
	return a.userData
}

func (a *TcpAgent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *TcpAgent) GetForwardMsg(msg interface{}) []byte {
	if a.Processor != nil {
		id, data, err := a.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message id:%v error: %v", id, err)
			return nil
		}

		tdata, err := a.conn.GetForwardMsg(data...)
		if err != nil {
			log.Error("write message id:%v error: %v", id, err)
			return nil
		}
		return tdata
	}

	return nil
}

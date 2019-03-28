package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/ant/log"
<<<<<<< HEAD
	"github.com/mikeqiao/ant/net/msgtype"
	"github.com/mikeqiao/ant/net/proto"
=======
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
)

var CreatID = int64(1)

type UserData struct {
	UId     int64
	UsersId []int64
	Agent   *TcpAgent
}

type TcpAgent struct {
	UId       int64
	conn      *TCPConn
	Processor Processor
	lifetime  int64          //链接未验证有效期时间（秒）
	starttime int64          //链接开始的时间戳
	tick      int64          //上次心跳的时间戳
	islogin   bool           //是否登陆验证过
	isClose   bool           //是否关闭
	isUpdate  bool           //是否验证心跳
	ctype     int32          //连接类型 1 server  2 client
	wg        sync.WaitGroup // 链接wait
	userData  interface{}
<<<<<<< HEAD
	Version   int32
	SUID      int64
=======
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func NewAgent(conn *TCPConn, tp Processor) *TcpAgent {
	a := new(TcpAgent)
	CreatID += 1
	a.UId = CreatID
	a.conn = conn
	a.Processor = tp
	a.lifetime = 10
	a.starttime = time.Now().Unix()
	a.tick = time.Now().Unix()
	a.islogin = false
	a.isClose = false
<<<<<<< HEAD
	a.isUpdate = true
=======
	a.isUpdate = false
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
	return a
}

func (a *TcpAgent) SetUID(uid int64) {
	a.UId = uid
}

func (a *TcpAgent) Start(name string) {
<<<<<<< HEAD
	if name != "Out" {
		a.ctype = 1
		return
	}
	msg := &proto.NewConnect{
		Id: a.UId,
	}
	ta := UserData{
		Agent: a,
		UId:   a.UId,
	}
	err := a.Processor.Route(msgtype.NewConnect, msg, ta)
	if err != nil {
		log.Debug("Start Add Agent error: %v", err)
	}
	a.ctype = 2
=======

>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (a *TcpAgent) Run() {

	if a.isUpdate == true {
		a.wg.Add(1)
		go a.Update()
	}

	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			a.Close()
			break
		}
		if a.Processor != nil {
			id, msg, err := a.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
			}
			if msg != nil {
				ta := UserData{
					Agent: a,
					UId:   a.UId,
				}
				err = a.Processor.Route(id, msg, ta)
				if err != nil {
					log.Debug("route message error: %v", err)
				}
			}
		}
	}
	a.isClose = true
}

func (a *TcpAgent) Update() {
	t1 := time.NewTimer(time.Second * 1)
	t2 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			if a.isClose == true {
				log.Debug("agent closed")
				a.Close()
				goto Loop
			}

			if a.islogin != true {
				nowtime := time.Now().Unix()
				if (a.starttime + a.lifetime) < nowtime {
					log.Debug("outtime to not login: %v", a.conn.RemoteAddr())
					a.Close()
					goto Loop
				}
			} else {
				nowtime := time.Now().Unix()
				if (a.tick + a.lifetime*3) < nowtime {
					log.Debug("outtime to no tick: %v", a.conn.RemoteAddr())
					a.Close()
					goto Loop
				}
			}
			t1.Reset(time.Second * 1)

		case <-t2.C:
			if a.isUpdate == true {
<<<<<<< HEAD
				nowtime := time.Now().Unix()
				a.WriteMsg(&proto.ServerTick{
					Time: nowtime,
				})

=======
				//		nowtime := time.Now().Unix()
				//		a.WriteMsg(&tproto.ServerTick{
				//			Time: nowtime,
				//		})
				//	log.Debug("send tick: %v, %v", nowtime, a.conn.RemoteAddr())
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
			}
			t2.Reset(time.Second * 10)
		}
	}
Loop:
	a.wg.Done()
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
<<<<<<< HEAD
	var msg interface{}
	var id uint16
	if a.ctype == 2 {
		msg = &proto.DelConnect{
			Id: a.UId,
		}
		id = msgtype.DelConnect
	} else {
		msg = &proto.ServerDelConnect{
			Id: a.UId,
		}
		id = msgtype.ServerDelConnect
	}
	ta := UserData{
		Agent: a,
		UId:   a.UId,
	}
	err := a.Processor.Route(id, msg, ta)
	if err != nil {
		log.Debug("OutTime Add Agent error: %v", err)

	}

=======
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
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

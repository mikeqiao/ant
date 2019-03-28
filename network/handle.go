package network

import (
	"fmt"

	"github.com/mikeqiao/ant/log"
	mod "github.com/mikeqiao/ant/module"
	"github.com/mikeqiao/ant/net"
	"github.com/mikeqiao/ant/net/proto"
	"github.com/mikeqiao/ant/rpc"
)

func HandleNewConnect(msg interface{}, data *net.UserData) {
	m := msg.(*proto.NewConnect)
	if nil != data && nil != data.Agent {
		CM.AddNewUserConn(m.GetId(), data.Agent)
	}
}

func HandleDelConnect(msg interface{}, data *net.UserData) {
	m := msg.(*proto.DelConnect)
	uid := m.GetId()
	CM.DelUserConn(uid)
}

func HandleServerDelConnect(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerDelConnect)
	uid := m.GetId()
	module := mod.ModuleControl.GetModule(uid)
	if nil == module {
		log.Debug(" module is nil uid:%v", uid)
		return
	}
	module.Close()
}

func HandleServerTick(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerTick)
	// 消息的发送者
	if nil != data && nil != data.Agent {
		data.Agent.SetTick(m.GetTime())
	}
}

func HandleServerLoginRQ(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerLoginRQ)
	// 消息的发送者
	if nil == data || nil == data.Agent {
		log.Debug(" wrong agent ,msg: %v", msg)
		return
	}
	s := m.GetServerinfo()
	if nil == s {
		log.Debug(" wrong serverinfo msg: %v", m)
		return
	}
	uid := s.GetServerId()
	module := mod.ModuleControl.GetModule(uid)
	if nil != module {
		mod.ModuleControl.DelMod(uid)
	}
	module = mod.NewModule(uid, s.GetState(), s.GetServerVersion(), 0)
	module.SetAgent(data.Agent)

	data.Agent.WriteMsg(&proto.ServerLoginRS{
		Result: 0,
		Serverinfo: &proto.ServerInfo{
			ServerId:      data.Agent.SUID,
			ServerVersion: data.Agent.Version,
			State:         1,
		},
	})
}

func HandleServerLoginRS(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerLoginRS)
	// 消息的发送者
	if nil == data || nil == data.Agent {
		log.Debug(" wrong agent ,msg: %v", msg)
		return
	}
	s := m.GetServerinfo()
	if nil == s {
		log.Debug(" wrong serverinfo msg: %v", m)
		return
	}
	uid := s.GetServerId()
	module := mod.ModuleControl.GetModule(uid)
	if nil != module {
		mod.ModuleControl.DelMod(uid)
	}
	module = mod.NewModule(uid, s.GetState(), s.GetServerVersion(), 0)
	module.SetAgent(data.Agent)
	module.SendFunc(1)
}

func HandleServerRegister(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerRegister)
	// 消息的发送者
	uid := m.GetUid()
	module := mod.ModuleControl.GetModule(uid)
	if nil == module {
		log.Debug(" module is nil uid:%v", uid)
		return
	}
	f := func(c *rpc.CallInfo) {
		if nil == c || nil == module {
			return
		}
		key := module.AddWaitCall(c)
		if nil != module.Agent {
			tdata := module.Agent.GetForwardMsg(c.Args)
			if tdata != nil {
				u := new(proto.UserInfo)
				if nil != c.Data {
					u.UId = c.Data.UId
					u.UsersId = c.Data.UsersId[:]
				}
				msg := &proto.ServerCall{
					FromMId: c.Mid,
					ToMId:   c.Mid,
					FUId:    c.Fid,
					CUId:    key,
					Msginfo: tdata[:],
					User:    u,
				}
				module.Agent.WriteMsg(msg)
			}
		} else if "" != key {
			e := fmt.Errorf("too shout args")
			module.ExecBack(key, nil, e)
		}
	}
	for _, v := range m.GetFuid() {
		module.RegisterRemote(v, f)
	}
	module.Start()
	if 1 == m.GetType() {
		module.SendFunc(2)
	}
}

func HandleServerDelFunc(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerDelFunc)
	uid := m.GetUid()
	module := mod.ModuleControl.GetModule(uid)
	if nil == module {
		log.Debug(" module is nil uid:%v", uid)
		return
	}
	for _, v := range m.GetFuid() {
		module.DelFunc(v)
	}
}

func HandleServerCall(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerCall)
	// 消息的发送者
	info := m.GetMsginfo()
	lenMsgLen := 4
	if nil != m.GetUser() {
		data.UId = m.GetUser().GetUId()
		data.UsersId = m.GetUser().GetUsersId()[:]
	}
	fid := m.GetFUId()
	Cid := m.GetCUId()
	//data
	msgData := info[lenMsgLen:]
	if data.Agent.Processor != nil {
		_, msg, err := data.Agent.Processor.Unmarshal(msgData)
		if err != nil {
			log.Debug("HandleServerForwardMsg message error: %v", err)
		}
		if "" != Cid {
			cb := func(in interface{}, e error) {
				if nil != data.Agent {
					if nil != e || nil == in {
						log.Debug("err:%v", e)
						data.Agent.WriteMsg(
							&proto.ServerCallBack{
								FromMId: m.GetFromMId(),
								ToMId:   m.GetToMId(),
								CUId:    m.GetCUId(),
								User:    m.GetUser(),
							})
					} else {
						tdata := data.Agent.GetForwardMsg(in)
						if tdata != nil {
							data.Agent.WriteMsg(
								&proto.ServerCallBack{
									FromMId: m.GetFromMId(),
									ToMId:   m.GetToMId(),
									CUId:    m.GetCUId(),
									Msginfo: tdata[:],
									User:    m.GetUser(),
								})
						}
					}
				}
			}
			mod.RPC.Route(fid, cb, msg, data)
		} else {
			mod.RPC.Route(fid, nil, msg, data)
		}
	}

}

func HandleServerCallBack(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerCallBack)
	info := m.GetMsginfo()
	lenMsgLen := 4
	if nil != m.GetUser() {
		data.UId = m.GetUser().GetUId()
		data.UsersId = m.GetUser().GetUsersId()[:]
	}
	//data
	module := mod.ModuleControl.GetModule(m.GetFromMId())
	if nil == module {
		log.Debug(" module is nil uid:%v", m.GetFromMId())
		return
	}
	key := m.GetCUId()
	if len(info) > lenMsgLen {
		msgData := info[lenMsgLen:]
		if data.Agent.Processor != nil {
			_, msg, err := data.Agent.Processor.Unmarshal(msgData)
			if err != nil {
				log.Debug("HandleServerForwardMsg message error: %v", err)
			}
			module.ExecBack(key, msg, nil)
		}
	} else {
		e := fmt.Errorf("no back msg")
		module.ExecBack(key, nil, e)
	}
}

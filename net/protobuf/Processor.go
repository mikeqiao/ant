package pb

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/mikeqiao/ant/db/mysql"
	"github.com/mikeqiao/ant/log"
	mod "github.com/mikeqiao/ant/module"
	"github.com/mikeqiao/ant/net"
	"github.com/mikeqiao/ant/net/msgtype"
	msg "github.com/mikeqiao/ant/net/proto"
	"github.com/mikeqiao/ant/network"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	msgInfo      map[uint16]*MsgInfo
	msgID        map[reflect.Type]uint16
}

type MsgInfo struct {
	msgType    reflect.Type
	msgHandler MsgHandler
	fid        uint32
}

type MsgHandler func(msg interface{}, data *net.UserData)

func NewProcessor() *Processor {
	p := new(Processor)
	p.littleEndian = false
	p.msgInfo = make(map[uint16]*MsgInfo)
	p.msgID = make(map[reflect.Type]uint16)
	p.baseMsg()
	return p
}

func (p *Processor) baseMsg() {

	p.Register(&msg.NewConnect{}, msgtype.NewConnect)
	p.Register(&msg.DelConnect{}, msgtype.DelConnect)
	p.Register(&msg.ServerDelConnect{}, msgtype.ServerDelConnect)
	p.Register(&msg.ServerTick{}, msgtype.ServerTick)
	p.Register(&msg.ServerLoginRQ{}, msgtype.ServerLoginRQ)
	p.Register(&msg.ServerLoginRS{}, msgtype.ServerLoginRS)
	p.Register(&msg.ServerRegister{}, msgtype.ServerRegister)
	p.Register(&msg.ServerDelFunc{}, msgtype.ServerDelFunc)
	p.Register(&msg.ServerCall{}, msgtype.ServerCall)
	p.Register(&msg.ServerCallBack{}, msgtype.ServerCallBack)
	p.Register(&msg.DBServerRQ{}, msgtype.DBServerRQ)
	p.Register(&msg.DBServerRS{}, msgtype.DBServerRS)

	p.SetHandler(msgtype.NewConnect, network.HandleNewConnect)
	//	p.SetHandler(msgtype.DelConnect, network.HandleDelConnect)
	p.SetHandler(msgtype.ServerDelConnect, network.HandleServerDelConnect)
	p.SetHandler(msgtype.ServerTick, network.HandleServerTick)
	p.SetHandler(msgtype.ServerLoginRQ, network.HandleServerLoginRQ)
	p.SetHandler(msgtype.ServerLoginRS, network.HandleServerLoginRS)
	p.SetHandler(msgtype.ServerRegister, network.HandleServerRegister)
	p.SetHandler(msgtype.ServerDelFunc, network.HandleServerDelFunc)
	p.SetHandler(msgtype.ServerCall, network.HandleServerCall)
	p.SetHandler(msgtype.ServerCallBack, network.HandleServerCallBack)

}

func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

func (p *Processor) Register(msg interface{}, id uint16) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	if _, ok := p.msgInfo[id]; ok {
		log.Fatal("message id%v type %s is already registered", id, msgType)
	}
	if _, ok := p.msgID[msgType]; ok {
		log.Fatal("message %s is already registered", msgType)
	}
	if id >= math.MaxUint16 {
		log.Fatal("too many protobuf messages (max = %v)", math.MaxUint16)
	}
	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[id] = i
	p.msgID[msgType] = id
	if nil != mysql.DB {
		mysql.DB.Register(msg)
	}
	return id
}

func (p *Processor) SetRouter(id uint16, fid uint32) {
	if v, ok := p.msgInfo[id]; !ok || nil == v {
		log.Fatal("message %s not registered", id)
	} else {
		v.fid = fid
	}
}

func (p *Processor) SetHandler(id uint16, msgHandler MsgHandler) {
	if v, ok := p.msgInfo[id]; !ok || nil == v {
		log.Fatal("message %s not registered", id)
	} else {
		v.msgHandler = msgHandler
	}
}

func (p *Processor) Route(id uint16, msg interface{}, data *net.UserData) error {

	i, ok := p.msgInfo[id]
	if !ok || nil == i {
		return fmt.Errorf("message id:%v %s not registered", id)
	}
	if i.msgHandler != nil {
		i.msgHandler(msg, data)
	} else {
		if 0 != i.fid {
			cb := func(in interface{}, e error) {
				if nil != in && nil != data && nil != data.Agent {
					data.Agent.WriteMsg(in)
				}
			}
			log.Debug("fid:%v, did:%v, msg:%v, data:%v", i.fid, i.fid, msg, data)
			mod.RPC.Route(i.fid, i.fid, cb, msg, data)
		} else {
			return fmt.Errorf(" msgid:%v, handler is nil :%v", id, i)
		}
	}
	return nil
}

func (p *Processor) Unmarshal(data []byte) (uint16, interface{}, error) {
	if len(data) < 2 {
		return 0, nil, errors.New("protobuf data too short")
	}
	// id
	var id uint16
	if p.littleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}

	i, ok := p.msgInfo[id]
	if !ok {
		return id, nil, fmt.Errorf("message id %v not registered", id)
	}
	// msg
	msg := reflect.New(i.msgType.Elem()).Interface()
	return id, msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) (uint16, [][]byte, error) {
	msgType := reflect.TypeOf(msg)
	// id
	_id, ok := p.msgID[msgType]
	if !ok {
		err := fmt.Errorf("message %s not registered", msgType)
		return _id, nil, err
	}

	id := make([]byte, 2)
	if p.littleEndian {
		binary.LittleEndian.PutUint16(id, _id)
	} else {
		binary.BigEndian.PutUint16(id, _id)
	}
	// data
	data, err := proto.Marshal(msg.(proto.Message))
	return _id, [][]byte{id, data}, err
}

func (p *Processor) Range(f func(id uint16, t reflect.Type)) {
	for id, i := range p.msgInfo {
		f(uint16(id), i.msgType)
	}
}

func (p *Processor) GetMsg(id uint16) interface{} {
	i, ok := p.msgInfo[id]
	if !ok {
		log.Error("message id %v not registered", id)
		return nil
	}
	msg := reflect.New(i.msgType.Elem()).Interface()
	return msg
}

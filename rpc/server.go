package rpc

import (
	"fmt"
	"runtime"

	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

type Server struct {
	id        int64
	v         int32 //服务者版本
	s         int64 //服务者状态(服务数量)
	functions map[uint32]interface{}
	ChanCall  chan *CallInfo
}

type CallInfo struct {
	Mid     int64 //module uid
	Fid     uint32
	Did     uint32        //执行id(转发功能使用)
	f       interface{}   //执行function
	Cb      interface{}   //callback
	Args    interface{}   //参数
	Data    *net.UserData //附加信息
	chanRet chan *Return  //

}

func (c *CallInfo) SetResult(data interface{}, e error) {

	if c.chanRet == nil || nil == c.Cb {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Debug("%v", err)
		}
	}()
	r := &Return{
		ret: data,
		err: e,
		cb:  c.Cb,
	}
	c.chanRet <- r
}

func NewServer(l int32) *Server {
	s := new(Server)
	s.functions = make(map[uint32]interface{})
	s.ChanCall = make(chan *CallInfo, l)
	return s
}

func (s *Server) GetFunc() (fl []uint32) {
	for k, v := range s.functions {
		if nil != v {
			fl = append(fl, k)
		}
	}
	return
}

func (s *Server) Close() {

}

func (s *Server) Exec(ci *CallInfo) {
	if nil == ci {
		log.Error("nil call")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			if conf.Config.LenStackBuf > 0 {
				buf := make([]byte, int32(conf.Config.LenStackBuf))
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
			s.ret(ci, &Return{err: fmt.Errorf("%v", r)})
		}
	}()
	f, ok := ci.f.(func(*CallInfo))
	if ok {
		//	log.Debug(" start handle:%v", ci)
		f(ci)
		//	log.Debug(" end handle:%v", ci)
	} else if nil != ci.Cb {
		s.ret(ci, &Return{err: fmt.Errorf("err func format")})
		log.Error("err func format")
	}
}

func (s *Server) ret(ci *CallInfo, ri *Return) {
	if ci.chanRet == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Debug("%v", err)
		}
	}()

	ri.cb = ci.Cb
	ci.chanRet <- ri
}

func (s *Server) Register(id uint32, f interface{}) {
	_, ok := f.(func(*CallInfo))
	if !ok {
		log.Error("err func format, id:%v, func :%v", id, f)
		return
	}
	if _, ok := s.functions[id]; ok {
		log.Debug("already have this id:%v, func :%v", id, f)
	} else {
		s.functions[id] = f
	}
}

func (s *Server) DelFunc(id uint32) {
	if _, ok := s.functions[id]; ok {
		delete(s.functions, id)
	} else {
		log.Debug("not have this id:%v, func", id)
	}
}

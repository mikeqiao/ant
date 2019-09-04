package network

import (
	"sync"
	"time"

	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/group"
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

var SM *SeverManager

type NetServer struct {
	UID             int64
	Name            string
	Nettype         int //网络通信类型 1 tcp        2 websocket
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Pid             int32 //processor ID
	Version         int32
	Processor       net.Processor
	ListenAddr      string //监听地址
	// tcp
	LenMsgLen    int //消息长度字段占用字节数
	LittleEndian bool

	//websocket
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	//	clients   []*network.TCPClient
	tcpServer *net.TCPServer
	wsServer  *net.WSServer
}

func (this *NetServer) Run() {
	if this.Nettype == 1 {
		this.runTcp()
	} else if this.Nettype == 2 {
		this.runWebsocket()
	}
}

func (this *NetServer) runTcp() {
	if this.ListenAddr != "" {
		this.tcpServer = new(net.TCPServer)
		this.tcpServer.UId = this.UID
		this.tcpServer.Name = this.Name
		this.tcpServer.Addr = this.ListenAddr
		this.tcpServer.MaxConnNum = this.MaxConnNum
		this.tcpServer.PendingWriteNum = this.PendingWriteNum
		this.tcpServer.LenMsgLen = this.LenMsgLen
		this.tcpServer.MaxMsgLen = this.MaxMsgLen
		this.tcpServer.LittleEndian = this.LittleEndian
		this.tcpServer.Processor = this.Processor
		this.tcpServer.NewAgent = net.NewAgent
		this.tcpServer.Version = this.Version
	}
	if this.tcpServer != nil {
		log.Debug("server.ListenAddr: %s", this.ListenAddr)
		this.tcpServer.Start()
		group.Done()
	}
}

func (this *NetServer) runWebsocket() {
	if this.ListenAddr != "" {
		this.wsServer = new(net.WSServer)
		this.wsServer.Name = this.Name
		this.wsServer.Addr = this.ListenAddr
		this.wsServer.MaxConnNum = this.MaxConnNum
		this.wsServer.PendingWriteNum = this.PendingWriteNum
		this.wsServer.MaxMsgLen = this.MaxMsgLen
		this.wsServer.HTTPTimeout = this.HTTPTimeout
		this.wsServer.CertFile = this.CertFile
		this.wsServer.KeyFile = this.KeyFile
		this.wsServer.Processor = this.Processor
		this.wsServer.NewAgent = net.NewAgentWeb

	}
	if this.wsServer != nil {
		log.Debug("server.ListenAddr: %s", this.ListenAddr)
		this.wsServer.Start()
		group.Done()
	}
}

func (this *NetServer) Close() {
	if this.tcpServer != nil {
		this.tcpServer.Close()
	}

	if this.wsServer != nil {
		this.wsServer.Close()
	}

}

type SeverManager struct {
	Sl   map[int64]*NetServer
	lock sync.Mutex
}

func (s *SeverManager) Init() {
	s.Sl = make(map[int64]*NetServer)
}

func (s *SeverManager) AddData(d *NetServer) {
	if nil == d {
		return
	}
	s.lock.Lock()
	if _, ok := s.Sl[d.UID]; ok {
		log.Debug("server already have, data:%v", d)
	} else {
		s.Sl[d.UID] = d
	}
	s.lock.Unlock()
}

func (s *SeverManager) DelData(id int64) {
	s.lock.Lock()
	if v, ok := s.Sl[id]; ok {
		if nil != v {
			v.Close()
		}
		delete(s.Sl, id)
	} else {
		log.Debug("server not have, data:%v", id)
	}
	s.lock.Unlock()
}

func (s *SeverManager) Close() {
	s.lock.Lock()
	for k, v := range s.Sl {
		if nil != v {
			v.Close()
		}
		delete(s.Sl, k)
	}
	s.lock.Unlock()
}

func NewSeverManager() {
	SM = new(SeverManager)
	SM.Init()
}

func startServer() {
	for _, c := range conf.Config.ServerConf {
		server := new(NetServer)
		server.UID = c.UId
		server.Name = c.Name
		server.Nettype = c.Nettype
		server.MaxConnNum = c.MaxConnNum
		server.PendingWriteNum = c.PendingWriteNum
		server.MaxMsgLen = c.MaxMsgLen
		server.LenMsgLen = c.LenMsgLen
		server.LittleEndian = conf.Config.LittleEndian
		server.ListenAddr = c.ListenAddr //"192.168.31.222:6001"
		server.Pid = c.Pid
		server.Version = c.Version
		server.HTTPTimeout = 30 * time.Second
		server.CertFile = c.CertFile
		server.KeyFile = c.KeyFile
		if P != nil {
			v := P.GetData(server.Pid)
			if nil != v {
				server.Processor = v
			} else {
				log.Fatal("Processor must not be nil, server uid%v, pid:%s", c.UId, c.Pid)
			}

		} else {
			log.Fatal("P is nil")
		}
		group.Add(1)
		go server.Run()
		if nil != SM {
			SM.AddData(server)
		}
	}
}

func closeServer() {
	SM.Close()
}

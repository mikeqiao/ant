package network

import (
	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/group"
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

type NetServer struct {
	UID             int64
	Name            string
	Nettype         int //网络通信类型 1 tcp        2 websocket
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Pid             int32 //processor ID
	Processor       net.Processor
	ListenAddr      string //监听地址
	// tcp
	LenMsgLen    int //消息长度字段占用字节数
	LittleEndian bool
	//	clients   []*network.TCPClient
	tcpServer *net.TCPServer
	//	wsServer *network.WSServer
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
	}
	if this.tcpServer != nil {
		log.Debug("server.ListenAddr: %s", this.ListenAddr)
		this.tcpServer.Start()
		group.Done()
	}
}

func (this *NetServer) runWebsocket() {
	/*	if this.ListenAddr != "" {
			this.wsServer = new(network.WSServer)
			this.wsServer.Name = this.Name
			this.wsServer.Addr = this.ListenAddr
			this.wsServer.MaxConnNum = this.MaxConnNum
			this.wsServer.PendingWriteNum = this.PendingWriteNum
			this.wsServer.MaxMsgLen = this.MaxMsgLen
			this.wsServer.HTTPTimeout = this.HTTPTimeout
			this.wsServer.CertFile = this.CertFile
			this.wsServer.KeyFile = this.KeyFile
			this.wsServer.Processor = this.Processor
			this.wsServer.NewAgent = network.NewAgentWeb

		}
		if this.wsServer != nil {
			log.Debug("server.ListenAddr: %s", this.ListenAddr)
			this.wsServer.Start()
		}

		<-closeSig
		if this.wsServer != nil {
			this.wsServer.Close()
		}
	*/
}

func (this *NetServer) Close() {
	if this.tcpServer != nil {
		this.tcpServer.Close()
	}
	/*
		if this.wsServer != nil {
			this.wsServer.Close()
		}
	*/
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
	}
}

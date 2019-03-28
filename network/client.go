package network
<<<<<<< HEAD

import (
	"time"

	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/group"
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

type NetClient struct {
	UID             int64
	Nettype         int //网络通信类型 1 tcp        2 websocket
	PendingWriteNum int
	ConnectInterval time.Duration
	MaxMsgLen       uint32
	Pid             int32 //processor ID
	Version         int32
	Processor       net.Processor
	ConnAddr        string
	Client          *net.TCPClient
	// tcp
	LenMsgLen    int //消息长度字段占用字节数
	LittleEndian bool
}

func (this *NetClient) Run() {

	if this.Nettype == 1 {
		this.tcpConnect()
	} else if this.Nettype == 2 {
		this.websocketConnect()
	}
}

func (this *NetClient) tcpConnect() {
	if this.ConnAddr != "" {
		client := new(net.TCPClient)
		this.Client = client
		client.UId = this.UID
		client.Addr = this.ConnAddr
		client.ConnNum = 1
		client.ConnectInterval = this.ConnectInterval
		client.PendingWriteNum = this.PendingWriteNum
		client.LenMsgLen = this.LenMsgLen
		client.MaxMsgLen = this.MaxMsgLen
		client.Processor = this.Processor
		client.NewAgent = net.NewAgent
		client.Version = this.Version
	}
	if this.Client != nil {
		log.Debug("server.ListenAddr: %s", this.ConnAddr)
		this.Client.Start()
		group.Done()
	}

}

func (this *NetClient) websocketConnect() {
	/*
		this.WSClients = make(map[string]*network.WSClient)
		for _, t_addr := range addr {
			client := new(network.WSClient)
			client.Addr = t_addr
			client.ConnNum = 1
			client.ConnectInterval = this.ConnectInterval
			client.PendingWriteNum = this.PendingWriteNum

			client.Processor = this.Processor
			client.NewAgent = network.NewAgentWeb

			client.Start()
			if _, ok := this.WSClients[t_addr]; ok {
				log.Debug("client have start, conf.addr:%s", t_addr)
			} else {
				this.WSClients[t_addr] = client
				log.Debug("add Conn and the Addrs: %v", t_addr)
			}
			log.Debug("Conn to Addrs: %v", t_addr)
		}*/
}

func (this *NetClient) OnDestroy() {

}

func (this *NetClient) OnStop(addr string) {
}

func (this *NetClient) SendMsg(addr string, msg interface{}) {

}

func startClient() {
	for _, c := range conf.Config.ConnConf {
		client := new(NetClient)
		//	client.UId = c.UId
		client.Nettype = c.Nettype
		client.PendingWriteNum = c.PendingWriteNum
		client.MaxMsgLen = c.MaxMsgLen
		client.LenMsgLen = c.LenMsgLen
		client.ConnectInterval = conf.Config.ConnectInterval
		client.LittleEndian = conf.Config.LittleEndian
		client.ConnAddr = c.ConnAddr //[]
		client.Pid = c.Pid
		client.Version = c.Version
		if P != nil {
			v := P.GetData(client.Pid)
			if nil != v {
				client.Processor = v
			} else {
				log.Fatal("Processor must not be nil, server uid%v, pid:%s", c.UId, c.Pid)
			}

		} else {
			log.Fatal("P is nil")
		}
		group.Add(1)
		go client.Run()
	}
}
=======
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729

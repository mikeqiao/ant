package network

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/group"
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

var ClientM *ClientManager

func NewClientManager() {
	ClientM = new(ClientManager)
	ClientM.Init()
}

type ClientManager struct {
	Cl      map[int64]*NetClient
	lock    sync.Mutex // 锁
	isClose bool
}

func (c *ClientManager) Init() {
	c.Cl = make(map[int64]*NetClient)
	c.isClose = false
}

func (c *ClientManager) AddClient(d *NetClient) {
	if nil == d {
		return
	}
	c.lock.Lock()
	if _, ok := c.Cl[d.UID]; ok {
		log.Debug("clent already have, data:%v", d)
	} else {
		c.Cl[d.UID] = d
	}
	c.lock.Unlock()
}

func (c *ClientManager) DelClient(id int64) {
	c.lock.Lock()
	if v, ok := c.Cl[id]; ok {
		if nil != v {
			v.Close()
		}
		delete(c.Cl, id)
	} else {
		log.Debug("clent not have, data:%v", id)
	}
	c.lock.Unlock()
}

func (c *ClientManager) Close() {
	c.lock.Lock()
	for k, v := range c.Cl {
		if nil != v {
			v.Close()
		}
		delete(c.Cl, k)
	}
	c.lock.Unlock()
}

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

	//websocket
	WebClient        *net.WSClient
	HandshakeTimeout time.Duration
	dialer           websocket.Dialer
}

func (this *NetClient) Run() {

	if this.Nettype == 1 {
		this.tcpConnect()
	} else if this.Nettype == 2 {
		this.websocketConnect()
	}
	group.Done()
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
		log.Debug("client.ConnAddr: %s", this.ConnAddr)
		this.Client.Start()

	}

}

func (this *NetClient) websocketConnect() {

	if this.ConnAddr != "" {
		client := new(net.WSClient)
		this.WebClient = client
		client.Addr = this.ConnAddr
		client.UId = this.UID
		client.ConnNum = 1
		client.ConnectInterval = this.ConnectInterval
		client.PendingWriteNum = this.PendingWriteNum
		client.Processor = this.Processor
		client.NewAgent = net.NewAgentWeb
		client.Version = this.Version
		client.MaxMsgLen = this.MaxMsgLen
		client.HandshakeTimeout = this.HandshakeTimeout
	}

	if this.WebClient != nil {
		log.Debug("client.ConnAddr: %s", this.ConnAddr)
		this.WebClient.Start()
	}
}

func (this *NetClient) OnDestroy() {

}

func (this *NetClient) OnStop(addr string) {
}

func (this *NetClient) SendMsg(addr string, msg interface{}) {

}

func (this *NetClient) Close() {
	if nil != this.Client {
		this.Client.Close()
	}
	if nil != this.WebClient {
		this.WebClient.Close()
	}
}
func startClient() {
	for _, c := range conf.Config.ConnConf {
		client := new(NetClient)
		client.UID = c.UId
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
		if nil != ClientM {
			ClientM.AddClient(client)
		}
	}
}

func closeClient() {
	if nil != ClientM {
		ClientM.Close()
	}
}

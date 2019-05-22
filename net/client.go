package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/ant/log"

	"github.com/mikeqiao/ant/net/proto"
)

type TCPClient struct {
	UId             int64
	Version         int32
	sync.Mutex                                          // 锁
	Addr            string                              // 地址
	ConnNum         int                                 // 链接数量
	ConnectInterval time.Duration                       // 请求链接的间隔
	PendingWriteNum int                                 // 等待 阻塞 数量
	AutoReconnect   bool                                // 是否重新连接
	NewAgent        func(*TCPConn, Processor) *TcpAgent // 代理
	conns           ConnList                            // 链接list （map）
	wg              sync.WaitGroup                      // 等待
	closeFlag       bool                                // 关闭标识符
	Processor       Processor
	//msg parser
	LenMsgLen    int        // 消息长度字段占用字节数
	MinMsgLen    uint32     // 消息最小长度
	MaxMsgLen    uint32     // 消息最大长度
	LittleEndian bool       // 小端字节序
	msgParser    *MsgParser //消息解析器

	Agent *TcpAgent
}

func (this *TCPClient) init() {
	this.Lock()
	defer this.Unlock()

	if this.ConnNum <= 0 {
		this.ConnNum = 1
		log.Release("invalid ConnNum, reset to %v", this.ConnNum)
	}
	if this.ConnectInterval <= 0 {
		this.ConnectInterval = 3 * time.Second
		log.Release("invalid ConnectInterval, reset to %v", this.ConnectInterval)
	}
	if this.PendingWriteNum <= 0 {
		this.PendingWriteNum = 100
		log.Release("invalid PendingWriteNum, reset to %v", this.PendingWriteNum)
	}
	if this.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}
	if this.conns != nil {
		log.Fatal("client is running")
	}

	this.AutoReconnect = true
	this.conns = make(ConnList)
	this.closeFlag = false

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(this.LenMsgLen, this.MinMsgLen, this.MaxMsgLen)
	msgParser.SetByteOrder(this.LittleEndian)
	this.msgParser = msgParser
}

func (this *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", this.Addr)
		if err == nil || this.closeFlag {
			return conn
		}
		log.Release("connect to %v error: %v", this.Addr, err)
		time.Sleep(this.ConnectInterval)
		continue
	}
}

func (this *TCPClient) connect() {
	defer this.wg.Done()
reconnect:
	conn := this.dial()
	if conn == nil {
		log.Error("conn is nil")
		return
	}
	this.Lock()
	if this.closeFlag {
		this.Unlock()
		conn.Close()
		log.Debug("this is close")
		return
	}
	this.conns[conn] = struct{}{}
	this.Unlock()
	tcpConn := newTCPConn(conn, this.PendingWriteNum, this.msgParser)
	agent := this.NewAgent(tcpConn, this.Processor)
	this.Agent = agent
	agent.SetLocalUID(this.UId)
	agent.WriteMsg(&proto.ServerLoginRQ{
		Serverinfo: &proto.ServerInfo{
			ServerId:      this.UId,
			ServerVersion: this.Version,
			State:         1,
		},
	})
	agent.Run()
	//cleanup
	tcpConn.Close()
	this.Lock()
	delete(this.conns, conn)
	this.Unlock()
	agent.Close()
	if this.AutoReconnect {
		time.Sleep(this.ConnectInterval)
		log.Debug("agent reconnect")
		goto reconnect
	}
}

func (this *TCPClient) Start() {
	this.init()

	for i := 0; i < this.ConnNum; i++ {
		this.wg.Add(1)
		go this.connect()
	}
}

func (this *TCPClient) Close() {
	this.Lock()
	this.closeFlag = true
	for conn := range this.conns {
		conn.Close()
	}
	this.conns = nil
	this.Unlock()
	this.wg.Wait()
}

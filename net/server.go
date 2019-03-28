package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/ant/log"
)

func NewServer() {

}

type TCPServer struct {
	UId             int64
	Name            string
	Addr            string                              // 监听的地址端口
	MaxConnNum      int                                 // 最大链接数
	PendingWriteNum int                                 // (发送消息)阻塞等待数量
	NewAgent        func(*TCPConn, Processor) *TcpAgent // 代理
	ln              net.Listener                        // 监听
	conns           ConnList                            // 链接list (map)
	mutexConns      sync.Mutex                          // 锁
	wgln            sync.WaitGroup                      // 监听wait
	wgConns         sync.WaitGroup                      // 链接wait

	Processor Processor
	// msg parser
	LenMsgLen    int        // 消息长度字段占用字节数
	MinMsgLen    uint32     // 最小消息长度
	MaxMsgLen    uint32     // 最大消息长度
	LittleEndian bool       // 小端字节序
	msgParser    *MsgParser //消息解析器
}

func (this *TCPServer) init() {
	ln, err := net.Listen("tcp", this.Addr)
	if err != nil {
		log.Fatal("%v", err)
	}

	if this.MaxConnNum <= 0 {
		this.MaxConnNum = 100
		log.Release("invalid MaxConnNum, reset to %v", this.MaxConnNum)
	}

	if this.PendingWriteNum <= 0 {
		this.PendingWriteNum = 100
		log.Release("invalid PendingWriteNum, reset to %v", this.PendingWriteNum)
	}
	if this.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}

	this.ln = ln
	this.conns = make(ConnList)

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(this.LenMsgLen, this.MinMsgLen, this.MaxMsgLen)
	msgParser.SetByteOrder(this.LittleEndian)
	this.msgParser = msgParser
}

func (this *TCPServer) run() {
	this.wgln.Add(1)

	var tempDelay time.Duration
	for {
		conn, err := this.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Release("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			this.wgln.Done()
			return
		}
		tempDelay = 0

		this.mutexConns.Lock()
		if len(this.conns) >= this.MaxConnNum {
			this.mutexConns.Unlock()
			conn.Close()
			log.Debug("too many connections")
			continue
		}
		this.conns[conn] = struct{}{}
		this.mutexConns.Unlock()
		this.wgConns.Add(1)

		tcpConn := newTCPConn(conn, this.PendingWriteNum, this.msgParser)
		agent := this.NewAgent(tcpConn, this.Processor)

		go func() {

			agent.Start(this.Name)
			agent.Run()
			// cleanup
			tcpConn.Close()
			this.mutexConns.Lock()
			delete(this.conns, conn)
			this.mutexConns.Unlock()
			agent.Close()

			this.wgConns.Done()
		}()
	}
	this.wgln.Done()
}

func (this *TCPServer) Start() {
	this.init()
	this.run()
}

func (this *TCPServer) Close() {
	this.ln.Close()
	this.wgln.Wait()

	this.mutexConns.Lock()
	for conn := range this.conns {
		conn.Close()
	}

	this.conns = nil
	this.mutexConns.Unlock()
	this.wgConns.Wait()
}

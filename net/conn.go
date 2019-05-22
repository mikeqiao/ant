package net

import (
	"net"
	"sync"

	"github.com/mikeqiao/ant/log"
)

type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg(args ...[]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	GetForwardMsg(args ...[]byte) ([]byte, error)
}

type ConnList map[net.Conn]struct{}

type TCPConn struct {
	sync.Mutex
	conn      net.Conn
	writeChan chan []byte
	closeFlag bool
	msgParser *MsgParser
}

func newTCPConn(conn net.Conn, pendingWriteNum int, msgParser *MsgParser) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeChan = make(chan []byte, pendingWriteNum)
	tcpConn.msgParser = msgParser

	go func() {
		for b := range tcpConn.writeChan {
			if b == nil {
				break
			}
			_, err := conn.Write(b)
			if err != nil {
				break
			}
		}
		conn.Close()
		tcpConn.Lock()
		tcpConn.closeFlag = true
		tcpConn.Unlock()
	}()

	return tcpConn
}

func (this *TCPConn) doDestroy() {
	this.conn.(*net.TCPConn).SetLinger(0)
	this.conn.Close()
	if !this.closeFlag {
		close(this.writeChan)
		this.closeFlag = true
	}
}

func (this *TCPConn) Destroy() {
	this.Lock()
	this.doDestroy()
	this.Unlock()
}

func (this *TCPConn) Close() {
	this.Lock()
	if this.closeFlag {
		this.Unlock()
		return
	}
	this.doWrite(nil)
	this.closeFlag = true
	this.Unlock()
}

func (this *TCPConn) doWrite(b []byte) {
	if len(this.writeChan) == cap(this.writeChan) {
		log.Debug("close conn: channel full")
		this.doDestroy()
		return
	}
	this.writeChan <- b
}

func (this *TCPConn) Write(b []byte) {
	this.Lock()
	defer this.Unlock()

	if this.closeFlag || b == nil {
		log.Debug("write err: closeflag:%v", this.closeFlag)
		return
	}
	this.doWrite(b)
}

func (this *TCPConn) Read(b []byte) (int, error) {
	return this.conn.Read(b)
}

func (this *TCPConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *TCPConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *TCPConn) ReadMsg() ([]byte, error) {
	return this.msgParser.Read(this)
}

func (this *TCPConn) WriteMsg(args ...[]byte) error {
	return this.msgParser.Write(this, args...)
}

func (this *TCPConn) GetForwardMsg(args ...[]byte) ([]byte, error) {
	return this.msgParser.GetForwardMsg(args...)
}

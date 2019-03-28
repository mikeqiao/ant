package net

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type MsgParser struct {
	lenMsgLen    int
	minMsgLen    uint32
	maxMsgLen    uint32
	littleEndian bool
}

func NewMsgParser() *MsgParser {
	p := &MsgParser{
		lenMsgLen:    2,
		minMsgLen:    1,
		maxMsgLen:    4096,
		littleEndian: false,
	}
	return p
}

//It's dangerous to call the method on reading or writing
func (p *MsgParser) SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32) {
	if lenMsgLen == 1 || lenMsgLen == 2 || lenMsgLen == 4 {
		p.lenMsgLen = lenMsgLen
	}
	if minMsgLen != 0 {
		p.minMsgLen = minMsgLen
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}

	var max uint32
	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max {
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// goroutine safe
func (p *MsgParser) Read(conn *TCPConn) ([]byte, error) {
	var b [4]byte
	bufMsgLen := b[:p.lenMsgLen]

	//read len (io.ReadFull // 读取指定长度的字节)
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}
	//parse len
	var msglen uint32
	switch p.lenMsgLen {
	case 1:
		msglen = uint32(bufMsgLen[0])
	case 2:
		if p.littleEndian {
			msglen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msglen = uint32(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if p.littleEndian {
			msglen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msglen = binary.BigEndian.Uint32(bufMsgLen)
		}

	}

	//check len
	if msglen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msglen < p.minMsgLen {
		return nil, errors.New("message too short")
	}

	//data
	msgData := make([]byte, msglen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}

	return msgData, nil

}

// goroutine safe
func (p *MsgParser) Write(conn *TCPConn, args ...[]byte) error {
	//get len
	var msglen uint32
	for i := 0; i < len(args); i++ {
		msglen += uint32(len(args[i]))
	}
	//check len
	if msglen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msglen < p.minMsgLen {
		return errors.New("message too short")
	}

	msg := make([]byte, uint32(p.lenMsgLen)+msglen)

	// write len
	switch p.lenMsgLen {
	case 1:
		msg[0] = byte(msglen)
	case 2:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msglen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msglen))
		}
	case 4:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(msg, uint32(msglen))
		} else {
			binary.BigEndian.PutUint32(msg, uint32(msglen))
		}
	}

	// write data
	l := p.lenMsgLen
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}
	conn.Write(msg)

	return nil
}

func (p *MsgParser) GetForwardMsg(args ...[]byte) ([]byte, error) {
	//get len
	var msglen uint32
	for i := 0; i < len(args); i++ {
		msglen += uint32(len(args[i]))
	}
	//check len
	if msglen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msglen < p.minMsgLen {
		return nil, errors.New("message too short")
	}

	msg := make([]byte, uint32(p.lenMsgLen)+msglen)

	// write len
	switch p.lenMsgLen {
	case 1:
		msg[0] = byte(msglen)
	case 2:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msglen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msglen))
		}
	case 4:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(msg, uint32(msglen))
		} else {
			binary.BigEndian.PutUint32(msg, uint32(msglen))
		}
	}

	// write data
	l := p.lenMsgLen
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}

	return msg, nil
}

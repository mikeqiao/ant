package network

import (
	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net"
)

var P *PL

type PL struct {
	pl map[int32]net.Processor
}

func (p *PL) Init() {
	p.pl = make(map[int32]net.Processor)
}

func (p *PL) AddData(id int32, d net.Processor) {
	if _, ok := p.pl[id]; ok {
		log.Fatal("already have this id:%v", id)
	} else {
		p.pl[id] = d
	}
}

func (p *PL) GetData(id int32) net.Processor {
	if v, ok := p.pl[id]; ok {
		return v
	} else {
		log.Fatal("not have this id:%v", id)
		return nil
	}
}

func NewPl() {
	P = new(PL)
	P.Init()
}

func Init() {
	NewPl()
	NewConnManager()
	NewSeverManager()

}

func Run() {
	startManager()
	startServer()
	startClient()
}

func Close() {
	closeClient()
	closeServer()
	closeManager()
}

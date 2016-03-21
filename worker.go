package remote

import (
	"github.com/gansidui/gotcp"
	"time"
)

type Worker struct {
	fns map[int]goAble
}

func NewWorker() *Worker {
	fns := make(map[int]goAble)
	return &Worker{
		fns: fns,
	}
}

func (this *Worker) Add(id int, fn goAble) {
	this.fns[id] = fn
}

func (this *Worker) OnConnect(c *gotcp.Conn) bool {
	return true
}

func (this *Worker) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*Packet)
	id, body := packet.getId(), packet.getBody()
	response := this.fns[id](body)
	c.AsyncWritePacket(newPacket(id, response), 5*time.Second)
	return true
}

func (this *Worker) OnClose(c *gotcp.Conn) {
}

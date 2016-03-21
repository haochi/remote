package remote

import (
	"github.com/gansidui/gotcp"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type goAble func([]byte) []byte

type Remote struct {
	isWorker      bool
	protocol      *Protocol
	worker        *Worker
	conn          *net.TCPConn
	protocolMutex *sync.Mutex
	port          string
}

func New(port string, isWorker bool) *Remote {
	var r = &Remote{
		isWorker: isWorker,
		port:     port,
	}
	r.protocolMutex = &sync.Mutex{}
	r.worker = NewWorker()

	if !isWorker {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", r.port)
		if err != nil {
			return nil
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return nil
		}
		r.conn = conn
		r.protocol = &Protocol{}
	}

	return r
}

func (this *Remote) Run(fn func()) {
	if this.isWorker {
		this.listen()
	} else {
		fn()
	}
}

func (this *Remote) Register(id int, fn goAble) {
	this.worker.Add(id, fn)
}

func (this *Remote) Go(id int, request []byte, ch chan []byte) {
	requestPacket := newPacket(id, request)
	requestPayload := requestPacket.Serialize()
	this.conn.Write(requestPayload)
	this.protocolMutex.Lock()
	response, err := this.protocol.ReadPacket(this.conn)
	this.protocolMutex.Unlock()
	if err != nil {
		ch <- nil
		return
	}
	responsePacket := response.(*Packet)
	ch <- responsePacket.getBody()
}

func (this *Remote) listen() {
	if this.isWorker {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", this.port)
		if err != nil {
			log.Println(err)
			return
		}
		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Println(err)
			return
		}

		config := &gotcp.Config{
			PacketSendChanLimit:    20,
			PacketReceiveChanLimit: 20,
		}
		srv := gotcp.NewServer(config, this.worker, this.protocol)

		go srv.Start(listener, time.Second)

		chSig := make(chan os.Signal)
		signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
		<-chSig // release ...

		srv.Stop()
	}
}

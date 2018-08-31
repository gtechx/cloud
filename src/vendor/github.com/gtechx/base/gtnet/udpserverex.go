package gtnet

import (
	"net"

	"github.com/gtechx/base/pool"
)

type UdpServerEx struct {
	*UdpServer
	onNewConn func(net.Conn)

	clients map[string]*UdpConn
	delChan chan *net.UDPAddr
}

func NewUdpServerEx() *UdpServerEx {
	serex := &UdpServerEx{clients: make(map[string]*UdpConn), delChan: make(chan *net.UDPAddr, 1024)}

	return serex
}

func (this *UdpServerEx) Start(addr string, connhandler func(net.Conn)) error {
	var err error
	this.UdpServer = NewUdpServer(addr)
	this.UdpServer.Listener = this
	this.onNewConn = connhandler
	err = this.UdpServer.Start()

	if err != nil {
		return err
	}

	return nil
}

// func (this *UdpServerEx) Start() error {
// 	var err error
// 	err = this.userver.Start()

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (this *UdpServerEx) remove(raddr *net.UDPAddr) {
	this.delChan <- raddr
}

// func (this *UdpServerEx) Send(buff []byte, raddr *net.UDPAddr) {
// 	this.server.Send(buff, raddr)
// }

// func (this *UdpServerEx) Stop() error {
// 	var err error

// 	if this.userver != nil {
// 		err = this.userver.Stop()
// 		this.userver = nil
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

//var headerbuf = make([]byte, MsgHeaderSize)

func (this *UdpServerEx) OnRecv(buff []byte, raddr *net.UDPAddr) {
	select {
	case addr := <-this.delChan:
		delete(this.clients, addr.String())
	default:
	}

	conn, ok := this.clients[raddr.String()]
	if !ok {
		if this.onNewConn != nil {
			conn = newUdpConn(raddr, this) //&UdpConn{IP: raddr.IP.String(), Port: raddr.Port, raddr: raddr, server: this}
			this.clients[raddr.String()] = conn
			go this.onNewConn(conn)
		}
	}

	num := len(buff)
	newbuf := pool.ByteGet(num) //make([]byte, num)
	copy(newbuf, buff)
	conn.recvChan <- newbuf

	// datasize := 0
	// if conn.Parser != nil {
	// 	datasize = conn.Parser.ParseHeader(buff[:MsgHeaderSize])
	// }

	// if datasize > 0 {
	// 	if conn.Parser != nil {
	// 		if conn.Parser != nil {
	// 			conn.Parser.ParseMsg(buff[MsgHeaderSize : MsgHeaderSize+datasize])
	// 		}
	// 	}
	// }
}

func (this *UdpServerEx) OnPreSend(buff []byte, raddr *net.UDPAddr) {
	conn, ok := this.clients[raddr.String()]
	if ok {
		if conn.listener != nil {
			conn.listener.OnPreSend(buff)
		}
	}
}

func (this *UdpServerEx) OnPostSend(buff []byte, raddr *net.UDPAddr, num int) {
	conn, ok := this.clients[raddr.String()]
	if ok {
		if conn.listener != nil {
			conn.listener.OnPostSend(buff, num)
		}
	}
}

func (this *UdpServerEx) OnRecvBusy(buff []byte, raddr *net.UDPAddr) {
	conn, ok := this.clients[raddr.String()]
	if ok {
		if conn.listener != nil {
			conn.listener.OnRecvBusy(buff)
		}
	}
}

func (this *UdpServerEx) OnSendBusy(buff []byte, raddr *net.UDPAddr) {
	conn, ok := this.clients[raddr.String()]
	if ok {
		if conn.listener != nil {
			conn.listener.OnSendBusy(buff)
		}
	}
}

func (this *UdpServerEx) OnError(errcode int, msg string) {
	// conn, ok := this.clients[raddr]
	// if ok {
	// 	if conn.listener != nil {
	// 		conn.listener.OnError(errcode, msg)
	// 	}
	// }
}

func (this *UdpServerEx) OnStop() {
	// conn, ok := this.clients[raddr]
	// if ok {
	// 	if conn.listener != nil {
	// 		conn.listener.OnStop()
	// 	}
	// }
}

package gtnet

import (
	"net"
	"time"
)

type UdpConn struct {
	parser   IDataParser
	listener IConnListener

	addr string

	raddr    *net.UDPAddr
	server   *UdpServerEx
	recvChan chan []byte
	curBuff  []byte
}

func newUdpConn(addr *net.UDPAddr, server *UdpServerEx) *UdpConn {
	conn := &UdpConn{addr: addr.String(), raddr: addr, server: server, recvChan: make(chan []byte, UdpServerRecvChanSize)}
	//go conn.startRecv()
	return conn //&UdpConn{Addr: addr.String(), raddr: addr, server: server}
}

func (this *UdpConn) Read(p []byte) (n int, err error) {
	//println("UdpConn read p len ", len(p))
	if this.curBuff == nil {
		this.curBuff = <-this.recvChan
	}

	readcount := len(p)
	n = 0
	//println("udpconn readcount ", readcount, " curbuff len ", len(this.curBuff))
	if readcount <= len(this.curBuff) {
		n += copy(p, this.curBuff[:readcount])
	} else {
		for readcount > len(this.curBuff) {
			n += copy(p[n:], this.curBuff)

			readcount = readcount - n
			this.curBuff = <-this.recvChan
		}
		n += copy(p[n:], this.curBuff[:readcount])
	}
	if readcount == len(this.curBuff) {
		this.curBuff = nil
	} else {
		this.curBuff = this.curBuff[readcount:]
	}
	return
}

func (this *UdpConn) Write(b []byte) (n int, err error) {
	//this.Send(b)
	this.server.Send(b, this.raddr)
	return len(b), nil
}

func (this *UdpConn) Close() error {
	this.server.remove(this.raddr)
	close(this.recvChan)
	return nil
}

func (this *UdpConn) LocalAddr() net.Addr {
	return this.server.conn.LocalAddr()
}

func (this *UdpConn) RemoteAddr() net.Addr {
	return this.raddr
}

func (this *UdpConn) SetDeadline(t time.Time) error {
	return nil
}

func (this *UdpConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (this *UdpConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (this *UdpConn) Send(buff []byte) {
	this.server.Send(buff, this.raddr)
}

func (this *UdpConn) startRecv() {
	// for buff := range this.recvChan {
	// 	if buff == nil {
	// 		break
	// 	}
	// 	datasize := 0
	// 	if this.Parser != nil {
	// 		datasize = this.Parser.ParseHeader(buff[:MsgHeaderSize])
	// 	}

	// 	if datasize > 0 {
	// 		if this.Parser != nil {
	// 			if this.Parser != nil {
	// 				this.Parser.ParseMsg(buff[MsgHeaderSize : MsgHeaderSize+datasize])
	// 			}
	// 		}
	// 	}
	// 	pool.BytePut(buff)
	// }
	for {
		if this.parser != nil {
			err := this.parser.Parse(this)
			if err != nil {
				if this.listener != nil {
					this.listener.OnError(1, "Read error:"+err.Error())
				}
				// if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
				// 	//time.Sleep(tempDelay)
				// 	continue
				// }
				break
			}
		}
	}
}

// func (this *UdpConn) RemoteAddr() string {
// 	return this.raddr.String()
// }

// func (this *UdpConn) LocalAddr() string {
// 	return this.addr
// }

func (this *UdpConn) SetDataParser(parser IDataParser) {
	this.parser = parser
}

func (this *UdpConn) SetListener(listener IConnListener) {
	this.listener = listener
}

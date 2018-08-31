package gtnet

import (
	"errors"
	"net"
	"strconv"
)

type ServerEx struct {
	servers []IServer
}

func NewServerEx() *ServerEx {
	return &ServerEx{}
}

func (this *ServerEx) Start(net, ip string, startport, endport int, connhandler func(net.Conn)) error {
	var err error
	portcount := endport - startport + 1
	this.servers = make([]IServer, portcount)
	if net == "tcp" {
		for i := 0; i < portcount; i++ {
			this.servers[i] = NewTcpServer()
		}
	} else if net == "ws" {
		for i := 0; i < portcount; i++ {
			this.servers[i] = NewWsServer()
		}
	} else if net == "kcp" {
		for i := 0; i < portcount; i++ {
			this.servers[i] = NewKcpServer()
		}
	} else if net == "udp" {
		for i := 0; i < portcount; i++ {
			this.servers[i] = NewUdpServerEx()
		}
	} else {
		return errors.New("invalid network:" + net)
	}

	for i := 0; i < portcount; i++ {
		addr := ip + ":" + strconv.Itoa(startport+i)
		err = this.servers[i].Start(addr, connhandler)
		if err != nil {
			this.Stop()
			return err
		}
	}

	return nil
}

func (this *ServerEx) Stop() error {
	var err error
	for i, _ := range this.servers {
		if this.servers[i] != nil {
			serr := this.servers[i].Stop()
			if serr != nil && err == nil {
				err = serr
			}
			this.servers[i] = nil
		}
	}
	return err
}

// func (this *ServerEx) onNewConn(conn IConn) {
// 	if this.OnNewConn != nil {
// 		//sure only one conn give upper one time
// 		this.connmutex.Lock()
// 		this.OnNewConn(conn)
// 		this.connmutex.Unlock()
// 	}
// }

// func (this *ServerEx) onNewTcpConn(conn *TcpConn) {
// 	if this.OnNewConn != nil {
// 		//sure only one conn give upper one time
// 		this.connmutex.Lock()
// 		this.OnNewConn(conn)
// 		this.connmutex.Unlock()
// 	}
// }

// func (this *ServerEx) onNewUdpConn(conn *UdpConn) {
// 	if this.OnNewConn != nil {
// 		this.connmutex.Lock()
// 		this.OnNewConn(conn)
// 		this.connmutex.Unlock()
// 	}
// }

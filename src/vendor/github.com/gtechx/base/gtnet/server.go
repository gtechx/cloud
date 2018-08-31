package gtnet

import (
	"errors"
	"net"
)

type Server struct {
	server IServer
}

func NewServer() *Server {
	return &Server{}
}

func (this *Server) Start(net, addr string, connhandler func(net.Conn)) error {
	var err error
	if net == "tcp" {
		this.server = NewTcpServer()
	} else if net == "ws" {
		this.server = NewWsServer()
	} else if net == "kcp" {
		this.server = NewKcpServer()
	} else if net == "udp" {
		this.server = NewUdpServerEx()
	} else {
		return errors.New("invalid network:" + net)
	}

	err = this.server.Start(addr, connhandler)
	if err != nil {
		return err
	}

	return nil
}

func (this *Server) Stop() error {
	var err error
	if this.server != nil {
		err = this.server.Stop()
		this.server = nil
	}
	return err
}

// func (this *Server) onNewConn(conn IConn) {
// 	if this.OnNewConn != nil {
// 		this.OnNewConn(conn)
// 	}
// }

// func (this *Server) onNewTcpConn(conn *TcpConn) {
// 	if this.OnNewConn != nil {
// 		this.OnNewConn(conn)
// 	}
// }

// func (this *Server) onNewUdpConn(conn *UdpConn) {
// 	if this.OnNewConn != nil {
// 		this.OnNewConn(conn)
// 	}
// }

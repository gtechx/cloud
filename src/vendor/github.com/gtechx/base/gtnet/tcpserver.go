package gtnet

import (
	"net"
	"time"
)

type TcpServer struct {
	onNewConn func(net.Conn)
	listener  *net.TCPListener
}

var KeepAlivePeriod time.Duration = 3000

func NewTcpServer() *TcpServer {
	return &TcpServer{}
}

func (this *TcpServer) Start(addr string, connhandler func(net.Conn)) error {
	var err error

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	this.onNewConn = connhandler
	this.listener = ln.(*net.TCPListener)
	go this.startListen()
	return nil
}

func (this *TcpServer) Stop() error {
	var err error
	if this.listener != nil {
		err = this.listener.Close()
		this.listener = nil
	}
	return err
}

func (this *TcpServer) startListen() {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := this.listener.AcceptTCP()
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
				//srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				tempDelay = 0
				continue
			}
			break
		}

		//conn.SetKeepAlive(true)
		//conn.SetKeepAlivePeriod(KeepAlivePeriod * time.Millisecond)
		conn.SetNoDelay(true)

		if this.onNewConn != nil {
			go this.onNewConn(conn)
		}
		// tcpconn := newConn(conn)
		// //tcpconn.serve()
		// if this.OnNewConn != nil {
		// 	this.OnNewConn(tcpconn)
		// }
		// //tcpconn.serve()
		// go tcpconn.startSend()
		// go tcpconn.startRecv()
	}
}

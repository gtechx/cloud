package gtnet

import (
	"net"
	"time"

	"github.com/xtaci/kcp-go"
)

type KcpServer struct {
	onNewConn func(net.Conn)
	listener  *kcp.Listener
}

func NewKcpServer() *KcpServer {
	return &KcpServer{}
}

func (this *KcpServer) Start(addr string, connhandler func(net.Conn)) error {
	var err error

	ln, err := kcp.Listen(addr)
	if err != nil {
		return err
	}

	this.onNewConn = connhandler
	this.listener = ln.(*kcp.Listener)
	go this.startListen()
	return nil
}

func (this *KcpServer) Stop() error {
	var err error
	if this.listener != nil {
		err = this.listener.Close()
		this.listener = nil
	}
	return err
}

func (this *KcpServer) startListen() {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := this.listener.AcceptKCP()
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
		conn.SetNoDelay(1, 10, 2, 1)
		conn.SetACKNoDelay(true)

		if this.onNewConn != nil {
			go this.onNewConn(conn)
		}
		// kcpconn := newConn(conn)
		// //KcpConn.serve()
		// if this.OnNewConn != nil {
		// 	this.OnNewConn(kcpconn)
		// }
		// //kcpconn.serve()
		// go kcpconn.startSend()
		// go kcpconn.startRecv()
	}
}

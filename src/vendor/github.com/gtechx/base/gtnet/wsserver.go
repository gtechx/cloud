package gtnet

import (
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

type WsServer struct {
	onNewConn func(net.Conn)
	listener  net.Listener
}

func NewWsServer() *WsServer {
	return &WsServer{}
}

func (this *WsServer) Start(addr string, connhandler func(net.Conn)) error {
	var err error

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	this.onNewConn = connhandler
	this.listener = ln
	http.Handle("/", websocket.Handler(this.accept))
	go this.startListen()
	return nil
}

type wsAddr string

func (wa wsAddr) Network() string {
	return "ws"
}

func (wa wsAddr) String() string {
	return string(wa)
}

type wsConn struct {
	*websocket.Conn
	remoteAddr net.Addr
}

func (wc *wsConn) RemoteAddr() net.Addr {
	//override RemoteAddr() of websocket.Conn
	return wc.remoteAddr
}

func (this *WsServer) accept(conn *websocket.Conn) {
	conn.PayloadType = websocket.BinaryFrame

	wsconn := &wsConn{Conn: conn, remoteAddr: wsAddr(conn.Request().RemoteAddr)}
	if this.onNewConn != nil {
		this.onNewConn(wsconn)
	}
	//WsConn.serve()
	// go wsconn.startSend()
	// wsconn.startRecv()
}

func (this *WsServer) Stop() error {
	var err error
	if this.listener != nil {
		err = this.listener.Close()
		this.listener = nil
	}
	return err
}

func (this *WsServer) startListen() {
	http.Serve(this.listener, nil)
}

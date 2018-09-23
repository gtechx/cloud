package gtnet

import (
	"io"
	"net"
)

type IDataParser interface {
	Parse(io.Reader) error
}

type IConnListener interface {
	OnClose()
	OnError(int, string)
	OnPreSend([]byte)
	OnPostSend([]byte, int)
	OnRecvBusy([]byte)
	OnSendBusy([]byte)
}

type IConn interface {
	Send([]byte)
	Close() error
	LocalAddr() string
	RemoteAddr() string
	SetDataParser(func(io.Reader) error)
	SetListener(IConnListener)
}

type IClient interface {
	IConn
	Connect(addr string) error
}

type IServer interface {
	Start(addr string, connhandler func(net.Conn)) error
	Stop() error
}

var MsgHeaderSize int = 2

var TcpServerSendChanSize int = 4
var TcpServerRecvChanSize int = 4
var TcpClientSendChanSize int = 4
var TcpClientRecvChanSize int = 4

var UdpPacketSize int = 10240
var UdpServerSendChanSize int = 1024
var UdpServerRecvChanSize int = 1024
var UdpClientSendChanSize int = 4
var UdpClientRecvChanSize int = 4

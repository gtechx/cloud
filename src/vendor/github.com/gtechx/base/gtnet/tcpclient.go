package gtnet

import (
	"net"
)

type TcpClient struct {
	*Conn
}

func NewTcpClient() *TcpClient {
	return &TcpClient{}
}

func (this *TcpClient) Connect(addr string) error {
	var err error

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	this.Conn = newConn(conn)

	//this.serve()
	go this.startSend()
	go this.startRecv()
	return nil
}

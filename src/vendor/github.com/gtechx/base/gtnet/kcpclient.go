package gtnet

import (
	"github.com/xtaci/kcp-go"
)

type KcpClient struct {
	*Conn
}

func NewKcpClient() *KcpClient {
	return &KcpClient{}
}

func (this *KcpClient) Connect(addr string) error {
	var err error

	conn, err := kcp.Dial(addr)
	if err != nil {
		return err
	}
	this.Conn = newConn(conn)

	//this.serve()
	go this.startSend()
	go this.startRecv()
	return nil
}

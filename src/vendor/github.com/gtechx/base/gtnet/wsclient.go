package gtnet

import (
	"golang.org/x/net/websocket"
)

type WsClient struct {
	*Conn
}

func NewWsClient() *WsClient {
	return &WsClient{}
}

func (this *WsClient) Connect(addr string) error {
	var err error

	// 	origin := "http://localhost/"
	// url := "ws://localhost:12345/ws"
	// ws, err := websocket.Dial(url, "", origin)
	// if err != nil {
	//     log.Fatal(err)
	// }

	conn, err := websocket.Dial("ws://"+addr, "", "http://"+addr)
	if err != nil {
		return err
	}
	this.Conn = newConn(conn)

	//this.serve()
	go this.startSend()
	go this.startRecv()
	return nil
}

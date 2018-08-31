package gtnet

import (
	"errors"
	"io"
)

type Client struct {
	parser   IDataParser
	listener IConnListener

	addr string
	net  string

	client IClient
}

func NewClient(net, addr string) *Client {
	return &Client{addr: addr, net: net}
}

func (this *Client) Connect() error {
	var err error
	if this.net == "tcp" {
		this.client = NewTcpClient()
	} else if this.net == "ws" {
		this.client = NewWsClient()
	} else if this.net == "kcp" {
		this.client = NewKcpClient()
	} else if this.net == "udp" {
		this.client = NewUdpClient()
	} else {
		return errors.New("invalid network:" + this.net)
	}

	err = this.client.Connect(this.addr)
	if err != nil {
		return err
	}

	this.client.SetDataParser(this)
	this.client.SetListener(this)

	return nil
}

func (this *Client) RemoteAddr() string {
	return this.client.RemoteAddr()
}

func (this *Client) LocalAddr() string {
	return this.client.LocalAddr()
}

func (this *Client) SetDataParser(parser IDataParser) {
	this.parser = parser
}

func (this *Client) SetListener(listener IConnListener) {
	this.listener = listener
}

func (this *Client) Close() error {
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
	return nil
}

func (this *Client) Send(buff []byte) {
	if this.client != nil {
		this.client.Send(buff)
	}
}

func (this *Client) Parse(reader io.Reader) error {
	if this.parser != nil {
		return this.parser.Parse(reader)
	}
	return nil
}

func (this *Client) OnClose() {
	if this.listener != nil {
		this.listener.OnClose()
	}
}

func (this *Client) OnError(errcode int, msg string) {
	if this.listener != nil {
		this.listener.OnError(errcode, msg)
	}
}

func (this *Client) OnPreSend(buff []byte) {
	if this.listener != nil {
		this.listener.OnPreSend(buff)
	}
}

func (this *Client) OnPostSend(buff []byte, num int) {
	if this.listener != nil {
		this.listener.OnPostSend(buff, num)
	}
}

func (this *Client) OnRecvBusy(buff []byte) {
	if this.listener != nil {
		this.listener.OnRecvBusy(buff)
	}
}

func (this *Client) OnSendBusy(buff []byte) {
	if this.listener != nil {
		this.listener.OnSendBusy(buff)
	}
}

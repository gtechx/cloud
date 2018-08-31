package gtnet

import (
	"net"

	"github.com/gtechx/base/pool"
)

// type iUdpClient interface {
// 	Connect() error
// 	Close() error
// }

// type IUdpClientListener interface {
// 	udpListener
// 	OnClose()
// }

type UdpClient struct {
	parser   IDataParser
	listener IConnListener

	recvChan chan []byte
	sendChan chan []byte
	conn     *net.UDPConn
	raddr    *net.UDPAddr
	quitChan chan int
	curBuff  []byte
}

func NewUdpClient() *UdpClient {
	return &UdpClient{recvChan: make(chan []byte, UdpClientRecvChanSize), sendChan: make(chan []byte, UdpClientSendChanSize), quitChan: make(chan int, 1)}
}

func (this *UdpClient) Connect(addr string) error {
	var err error
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	this.raddr = uaddr

	this.conn, err = net.DialUDP("udp", nil, uaddr)

	if err != nil {
		return err
	}

	go this.startUDPProcess()
	go this.startUDPRecv()
	go this.startUDPSend()

	return nil
}

func (this *UdpClient) RemoteAddr() string {
	return this.conn.RemoteAddr().String()
}

func (this *UdpClient) LocalAddr() string {
	return this.conn.LocalAddr().String()
}

func (this *UdpClient) SetDataParser(parser IDataParser) {
	this.parser = parser
}

func (this *UdpClient) SetListener(listener IConnListener) {
	this.listener = listener
}

func (this *UdpClient) Close() error {
	this.quitChan <- 1
	return nil
}

func (this *UdpClient) Send(buff []byte) {
	select {
	case this.sendChan <- buff:
	default:
		if this.listener != nil {
			this.listener.OnSendBusy(buff)
		}
	}
}

func (this *UdpClient) startUDPSend() {
	for {
		select {
		case <-this.quitChan:
			goto end
		case buff := <-this.sendChan:
			if this.listener != nil {
				this.listener.OnPreSend(buff)
			}

			num, err := this.conn.Write(buff)
			if err != nil {
				if this.listener != nil {
					this.listener.OnError(1, "Send error:"+err.Error())
				}
				if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
					continue
				}
				break
			}

			if this.listener != nil {
				this.listener.OnPostSend(buff, num)
			}
		}
	}
end:
	this.conn.Close()
	this.conn = nil
}

func (this *UdpClient) startUDPRecv() {
	buff := make([]byte, UdpPacketSize)
	for {
		num, err := this.conn.Read(buff)
		if err != nil {
			if this.listener != nil {
				this.listener.OnError(1, "Recv error:"+err.Error())
			}
			if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
				continue
			}
			break
		}

		if num > 0 {
			newbuf := pool.ByteGet(num) //make([]byte, num)
			copy(newbuf, buff[0:num])
			//newbuf = append(newbuf, buffer[0:num]...)
			select {
			case this.recvChan <- newbuf:
			default:
				if this.listener != nil {
					this.listener.OnRecvBusy(newbuf)
				}
			}
		}
	}
	if this.listener != nil {
		//fmt.Println("conn close")
		this.listener.OnClose()
	}
	close(this.recvChan)
}

func (this *UdpClient) startUDPProcess() {
	// for buff := range this.recvChan {
	// 	if buff == nil {
	// 		break
	// 	}
	// 	datasize := 0
	// 	if this.parser != nil {
	// 		datasize = this.parser.ParseHeader(buff[:MsgHeaderSize])
	// 	}

	// 	if datasize > 0 {
	// 		if this.parser != nil {
	// 			if this.parser != nil {
	// 				this.parser.ParseMsg(buff[MsgHeaderSize : MsgHeaderSize+datasize])
	// 			}
	// 		}
	// 	}
	// 	pool.BytePut(buff)
	// }
	for {
		if this.parser != nil {
			err := this.parser.Parse(this)
			if err != nil {
				if this.listener != nil {
					this.listener.OnError(1, "Read error:"+err.Error())
				}
				// if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
				// 	//time.Sleep(tempDelay)
				// 	continue
				// }
				break
			}
		}
	}
}

func (this *UdpClient) Read(p []byte) (n int, err error) {
	if this.curBuff == nil {
		this.curBuff = <-this.recvChan
	}

	readcount := len(p)
	n = 0
	if readcount <= len(this.curBuff) {
		n += copy(p, this.curBuff[:readcount])
	} else {
		for readcount > len(this.curBuff) {
			n += copy(p[n:], this.curBuff)

			readcount = readcount - n
			this.curBuff = <-this.recvChan
		}
		n += copy(p[n:], this.curBuff[:readcount])
	}
	if readcount == len(this.curBuff) {
		this.curBuff = nil
	} else {
		this.curBuff = this.curBuff[readcount:]
	}
	return
}

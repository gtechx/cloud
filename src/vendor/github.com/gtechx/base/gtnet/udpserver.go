package gtnet

import (
	"net"
	"sync"

	"github.com/gtechx/base/pool"
)

// type iUdpServer interface {
// 	Start() error
// 	Stop() error
// }

// type IUdpServerListener interface {
// 	udpListener
// 	OnStop()
// }
type udpPacket struct {
	buff  []byte
	raddr *net.UDPAddr
}

type udpServerListener interface {
	OnRecv([]byte, *net.UDPAddr)
	OnPreSend([]byte, *net.UDPAddr)
	OnPostSend([]byte, *net.UDPAddr, int)
	OnRecvBusy([]byte, *net.UDPAddr)
	OnSendBusy([]byte, *net.UDPAddr)
	OnError(int, string)
	OnStop()
}

type UdpServer struct {
	Listener udpServerListener
	Addr     string

	recvChan chan *udpPacket
	sendChan chan *udpPacket
	conn     *net.UDPConn
	quitChan chan int
}

func NewUdpServer(addr string) *UdpServer {
	return &UdpServer{recvChan: make(chan *udpPacket, UdpServerRecvChanSize), sendChan: make(chan *udpPacket, UdpServerSendChanSize), quitChan: make(chan int, 1), Addr: addr}
}

func (this *UdpServer) Start() error {
	var err error
	uaddr, err := net.ResolveUDPAddr("udp", this.Addr)
	if err != nil {
		return err
	}

	this.conn, err = net.ListenUDP("udp", uaddr)

	if err != nil {
		return err
	}

	go this.startUDPProcess()
	go this.startUDPRecv()
	go this.startUDPSend()

	return nil
}

func (this *UdpServer) Stop() error {
	this.quitChan <- 1
	return nil
}

func (this *UdpServer) Send(buff []byte, raddr *net.UDPAddr) {
	packet := newUdpPacket(buff, raddr)
	select {
	case this.sendChan <- packet:
	default:
		if this.Listener != nil {
			this.Listener.OnSendBusy(buff, raddr)
		}
		putUdpPacket(packet)
	}
}

func (this *UdpServer) startUDPSend() {
	for {
		select {
		case <-this.quitChan:
			goto end
		case packet := <-this.sendChan:
			if this.Listener != nil {
				this.Listener.OnPreSend(packet.buff, packet.raddr)
			}

			num, err := this.conn.WriteToUDP(packet.buff, packet.raddr)
			if err != nil {
				if this.Listener != nil {
					this.Listener.OnError(1, "Send error:"+err.Error())
				}
				if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
					continue
				}
				goto end
			}

			if this.Listener != nil {
				this.Listener.OnPostSend(packet.buff, packet.raddr, num)
			}
			putUdpPacket(packet)
		}
	}
end:
	this.conn.Close()
	this.conn = nil
}

func (this *UdpServer) startUDPRecv() {
	buff := make([]byte, UdpPacketSize)
	for {
		num, raddr, err := this.conn.ReadFromUDP(buff)

		if err != nil {
			if this.Listener != nil {
				this.Listener.OnError(1, "Recv error:"+err.Error())
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
			packet := newUdpPacket(newbuf, raddr) //&UdpPacket{newbuf, raddr}
			select {
			case this.recvChan <- packet:
			default:
				if this.Listener != nil {
					this.Listener.OnRecvBusy(newbuf, raddr)
				}
			}
		}
	}
	if this.Listener != nil {
		this.Listener.OnStop()
	}
	close(this.recvChan)
}

func (this *UdpServer) startUDPProcess() {
	for packet := range this.recvChan {
		if packet == nil {
			break
		}
		if this.Listener != nil {
			this.Listener.OnRecv(packet.buff, packet.raddr)
		}
		pool.BytePut(packet.buff)
		putUdpPacket(packet)
	}
}

var udpPacketPool sync.Pool

func newUdpPacket(buff []byte, raddr *net.UDPAddr) *udpPacket {
	if v := udpPacketPool.Get(); v != nil {
		packet := v.(*udpPacket)
		packet.buff = buff
		packet.raddr = raddr
		return packet
	}

	return &udpPacket{buff, raddr}
}

func putUdpPacket(packet *udpPacket) {
	udpPacketPool.Put(packet)
}

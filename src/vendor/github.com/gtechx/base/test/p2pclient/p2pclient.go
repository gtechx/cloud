package main

import (
	"flag"
	"fmt"
	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
	"net"
	"strings"
)

var client *gtnet.UdpServer
var quit chan int
var clientaddr *net.UDPAddr
var serveraddr *net.UDPAddr

type Processer struct {
}

func newProcesser(client *gtnet.UdpServer) *Processer {
	pro := &Processer{}
	client.Listener = pro
	return pro
}

func (p *Processer) OnRecv(buff []byte, raddr *net.UDPAddr) {
	fmt.Println("recv:" + string(buff[2:]))
}

func (p *Processer) OnStop() {
	fmt.Println("tcpclient stoped")
	quit <- 1
}

func (p *Processer) OnError(errorcode int, msg string) {
	fmt.Println("tcpclient error, errorcode:", errorcode, "msg:", msg)
}

func (p *Processer) OnPreSend(buff []byte, raddr *net.UDPAddr) {

}

func (p *Processer) OnPostSend(buff []byte, raddr *net.UDPAddr, num int) {

}

func (p *Processer) OnRecvBusy(buff []byte, raddr *net.UDPAddr) {
	fmt.Println("client is busy for recv, msg size is ", len(buff))
}

func (p *Processer) OnSendBusy(buff []byte, raddr *net.UDPAddr) {
	fmt.Println("client is busy for send, msg size is ", len(buff))
}

var addr string = "127.0.0.1:9090"

func main() {
	var err error

	paddr := flag.String("addr", "127.0.0.1:9090", "-addr=")

	flag.Parse()

	addr = *paddr

	quit = make(chan int, 1)
	//if addr port is 0, golang will choose an avialable port.
	//for muiti instance apps, this is useful
	client = gtnet.NewUdpServer(addr)

	err = client.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Stop()
	newProcesser(client)
	go startSend()
	<-quit
}

//first, start a udpserver under test folder
//second, start this client, use cmd 3#ip:port connect to udpserver.
//then you can see the client's ip:port on udpserver console
//third, start another client, with the same step of second
//forth, use cmd 2#clientip:clientport to hole to another client
//then another client can send msg to the client. then two client can send msg to each other.
//if you test this in lan, two client always can send msg to each other if you know ip:port of each without server.
//if you test on the internet, then you need the hole.
func startSend() {
	var str string
	for {
		str = ""
		fmt.Scanln(&str)

		if str != "" {
			//fmt.Println(str)
			strarr := strings.Split(str, "#")
			//fmt.Println(strarr)
			if len(strarr) == 1 {
				bytes := Bytes(int16(len(str)))
				//fmt.Println(bytes)
				//client.Send(bytes)
				client.Send(append(bytes, []byte(str)...), serveraddr)
			} else if len(strarr) == 2 {
				cmd := Int(strarr[0])
				//fmt.Println(cmd)
				if cmd == 1 {
					var err error
					//fmt.Println("cmd is 1")
					clientaddr, err = net.ResolveUDPAddr("udp", strarr[1])
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					bytes := Bytes(int16(len(strarr[1])))
					client.Send(append(bytes, []byte(strarr[1])...), clientaddr)
				} else if cmd == 2 {
					//fmt.Println("cmd is 2")
					bytes := Bytes(int16(len(strarr[1])))
					client.Send(append(bytes, []byte(strarr[1])...), clientaddr)
				} else if cmd == 3 {
					var err error
					//first
					serveraddr, err = net.ResolveUDPAddr("udp", strarr[1])
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					bytes := Bytes(int16(len(strarr[1])))
					client.Send(append(bytes, []byte(strarr[1])...), serveraddr)
				} else {
					bytes := Bytes(int16(len(strarr[1])))
					client.Send(append(bytes, []byte(strarr[1])...), serveraddr)
				}
			}
		}

	}
}

package main

import (
	"flag"
	"fmt"
	"io"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

var client *gtnet.Client
var quit chan int

type Processer struct {
}
type Handler func(reader io.Reader) error

func newProcesser(client *gtnet.Client) *Processer {
	pro := &Processer{}
	client.SetDataParser(Handler(pro.Parse))
	client.SetListener(pro)
	return pro
}
func (h Handler) Parse(reader io.Reader) error {
	return h(reader)
}
func (p *Processer) Parse(reader io.Reader) error {
	typebuff := make([]byte, 1)
	sizebuff := make([]byte, 2)

	fmt.Println()
	fmt.Println("****start new read****")
	_, err := reader.Read(typebuff)
	if err != nil {
		return err
	}

	fmt.Println("data type:", typebuff[0])

	_, err = reader.Read(sizebuff)
	if err != nil {
		return err
	}
	size := Int(sizebuff)

	fmt.Println("data size:", size)

	databuff := make([]byte, size)

	_, err = reader.Read(databuff)
	if err != nil {
		return err
	}

	fmt.Println("recv data:", String(databuff))
	return nil
}

func (p *Processer) OnError(errorcode int, desc string) {
	fmt.Println("conn error, errorcode:", errorcode, "desc:", desc)
}

func (p *Processer) OnPreSend([]byte) {

}

func (p *Processer) OnPostSend([]byte, int) {

}

func (p *Processer) OnClose() {
	fmt.Println("client closed")
	quit <- 1
}

func (p *Processer) OnRecvBusy(buff []byte) {
	fmt.Println("client is busy for recv, msg size is ", len(buff))
}

func (p *Processer) OnSendBusy(buff []byte) {
	fmt.Println("client is busy for send, msg size is ", len(buff))
}

var nettype string = "kcp"
var addr string = "127.0.0.1:9090"

func main() {
	var err error

	pnet := flag.String("net", "kcp", "-net=")
	paddr := flag.String("addr", "127.0.0.1:9090", "-addr=")

	flag.Parse()

	nettype = *pnet
	addr = *paddr

	quit = make(chan int, 1)
	client = gtnet.NewClient(nettype, addr)

	err = client.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Close()
	newProcesser(client)
	go startSend()
	<-quit
}

func startSend() {
	var str string
	for {
		str = ""
		fmt.Scanln(&str)
		if str != "" {
			bytes := Bytes(int16(len(str)))
			//fmt.Println(bytes)
			client.Send(append([]byte{1}, append(bytes, []byte(str)...)...))
		}
	}
}

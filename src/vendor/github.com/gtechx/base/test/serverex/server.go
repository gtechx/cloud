package main

import (
	"flag"
	"fmt"
	"io"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

var server *gtnet.ServerEx
var promap map[*Processer]*Processer
var quit chan int

type Processer struct {
	conn gtnet.IConn
}

func newProcesser(conn gtnet.IConn) *Processer {
	pro := &Processer{conn}
	conn.SetDataParser(pro)
	conn.SetListener(pro)
	return pro
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

	fmt.Println("recv data:", String(databuff), " from "+p.conn.RemoteAddr())

	p.conn.Send(append([]byte{typebuff[0]}, append(Bytes(int16(len(databuff))), databuff...)...))
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
	fmt.Println("close client connection:", p.conn.RemoteAddr())
	//remove client conn
	delete(promap, p)
}

func (p *Processer) OnRecvBusy([]byte) {
	str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	p.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (p *Processer) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}

var nettype string = "tcp"
var ip string = "127.0.0.1"
var startport int = 9000
var endport int = 9005

func main() {
	var err error

	pnet := flag.String("net", "tcp", "-net=")
	paddr := flag.String("addr", "127.0.0.1", "-addr=")
	pstartport := flag.Int("startport", 9000, "-startport=")
	pendport := flag.Int("endport", 9005, "-endport=")

	flag.Parse()

	nettype = *pnet
	ip = *paddr
	startport = *pstartport
	endport = *pendport

	quit = make(chan int, 1)
	promap = make(map[*Processer]*Processer, 0)
	server = gtnet.NewServerEx(nettype, ip, startport, endport, onNewConn)

	err = server.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("server start ok...")
	defer server.Stop()
	<-quit
}

func onNewConn(conn gtnet.IConn) {
	fmt.Println("new conn:", conn.RemoteAddr())
	pro := newProcesser(conn)
	promap[pro] = pro
	// conn.Send([]byte("hello\n"))
	// time.Sleep(5000 * time.Millisecond)
	// conn.Send([]byte("hello\n"))
}

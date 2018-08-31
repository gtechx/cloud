package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"time"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

var server *gtnet.Server
var promap map[*Processer]*Processer
var quit chan int

type Processer struct {
	conn     net.Conn
	sendChan chan []byte
	quitChan chan int
}

func newProcesser(conn net.Conn) *Processer {
	pro := &Processer{conn: conn, quitChan: make(chan int, 1), sendChan: make(chan []byte, 512)}
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

	// fmt.Println("recv data:", String(databuff), " from "+p.conn.RemoteAddr())

	// p.conn.Send(append([]byte{typebuff[0]}, append(Bytes(int16(len(databuff))), databuff...)...))
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
	//str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (p *Processer) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
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
	promap = make(map[*Processer]*Processer, 0)
	server = gtnet.NewServer()

	err = server.Start(nettype, addr, onNewConn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(nettype + " server start on addr " + addr + " ok...")
	defer server.Stop()
	<-quit
}

func onNewConn(conn net.Conn) {
	fmt.Println("new conn:", conn.RemoteAddr())
	//you can read verify msg here
	isok := false
	time.AfterFunc(5*time.Second, func() {
		if !isok {
			conn.Close()
		}
	})

	typebuff := make([]byte, 1)
	sizebuff := make([]byte, 2)

	_, err := conn.Read(typebuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	fmt.Println("data type:", typebuff[0])

	_, err = conn.Read(sizebuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	size := Int(sizebuff)

	fmt.Println("data size:", size)

	databuff := make([]byte, size)

	_, err = conn.Read(databuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	fmt.Println("recv data:", String(databuff), " from "+conn.RemoteAddr().String())

	if String(databuff) != "wyq" {
		conn.Close()
		return
	}

	isok = true
	pro := newProcesser(conn)
	promap[pro] = pro
	go pro.startSend()
	// conn.Send([]byte("hello\n"))
	// time.Sleep(5000 * time.Millisecond)
	// conn.Send([]byte("hello\n"))
	quitChan := pro.quitChan

	for {
		fmt.Println()
		fmt.Println("****start new read****")
		_, err := conn.Read(typebuff)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		fmt.Println("data type:", typebuff[0])

		_, err = conn.Read(sizebuff)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		size := Int(sizebuff)

		fmt.Println("data size:", size)

		databuff := make([]byte, size)

		_, err = conn.Read(databuff)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		fmt.Println("recv data:", String(databuff), " from "+pro.conn.RemoteAddr().String())

		pro.sendChan <- append([]byte{typebuff[0]}, append(Bytes(int16(len(databuff))), databuff...)...)
	}

	quitChan <- 1
}

func (this *Processer) startSend() {
	quitChan := this.quitChan
	sendChan := this.sendChan
	conn := this.conn
	for {
		select {
		case <-quitChan:
			count := len(sendChan)
			for i := 0; i < count; i++ {
				databuf := <-sendChan
				conn.Write(databuf)
			}
			goto end
		case databuf := <-sendChan:
			_, err := conn.Write(databuf)
			if err != nil {
				fmt.Println("err Send:" + err.Error())
				// if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
				// 	//srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				// 	//time.Sleep(tempDelay)
				// 	continue
				// }
				goto end
			}
		}
	}

end:
	conn.Close()
}

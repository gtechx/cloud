package gtnet

import (
	"fmt"
	"io"
	"net"
)

type Conn struct {
	parser   func(io.Reader) error
	listener IConnListener

	conn net.Conn

	sendChan chan []byte
	quitChan chan int
}

func newConn(conn net.Conn) *Conn {
	//sendchan is 2, for one is sending ,another can be put into the chan
	return &Conn{conn: conn, quitChan: make(chan int, 1), sendChan: make(chan []byte, TcpServerSendChanSize)}
}

func (this *Conn) RemoteAddr() string {
	return this.conn.RemoteAddr().String()
}

func (this *Conn) LocalAddr() string {
	return this.conn.LocalAddr().String()
}

func (this *Conn) SetDataParser(parser func(io.Reader) error) {
	this.parser = parser
}

func (this *Conn) SetListener(listener IConnListener) {
	this.listener = listener
}

func (this *Conn) Close() error {
	this.quitChan <- 1
	this.conn = nil
	return nil
}

func (this *Conn) Send(buff []byte) {
	if this.conn == nil {
		if this.listener != nil {
			this.listener.OnError(1, "conn is nil")
		}
		return
	}

	select {
	case this.sendChan <- buff:
	default:
		if this.listener != nil {
			this.listener.OnSendBusy(buff)
		}
	}
}

func (this *Conn) startRecv() {
	quitChan := this.quitChan
	conn := this.conn
	for {
		if this.parser != nil {
			err := this.parser(conn)
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

	if this.listener != nil {
		this.listener.OnClose()
	}
	quitChan <- 1
}

func (this *Conn) doSend(conn net.Conn, buff []byte) error {
	if this.listener != nil {
		this.listener.OnPreSend(buff)
	}

	num, err := conn.Write(buff)

	if this.listener != nil && err == nil {
		this.listener.OnPostSend(buff, num)
	}
	return err
}

func (this *Conn) startSend() {
	quitChan := this.quitChan
	sendChan := this.sendChan
	conn := this.conn
	for {
		select {
		case <-quitChan:
			count := len(sendChan)
			for i := 0; i < count; i++ {
				databuf := <-sendChan
				this.doSend(conn, databuf)
			}
			goto end
		case databuf := <-sendChan:
			err := this.doSend(conn, databuf)
			if err != nil {
				fmt.Println("err Send:" + err.Error())
				if this.listener != nil {
					this.listener.OnError(2, "Send error:"+err.Error())
				}
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

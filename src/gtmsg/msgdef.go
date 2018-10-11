package gtmsg

import (
	"errors"
	"io"

	. "github.com/gtechx/base/common"
)

//客户端在发送一个请求以后，需要启动一个定时器来检测该请求是否得到反馈：
//1.如果在规定时间没有反馈，则说明网络超时。
//2.如果收到反馈，则可以根据反馈进行弹框告知用户。
//在收到反馈或者超时以后需要将该检测移除检测队列。

//处理函数中需要根据发送给谁的id进行发送，这就需要用到session管理器，根据id查找到对于的session
//并且需要发送返回的消息， 所以可能需要传session进来
//所以消息处理模块和session和db模块有交互
//const define need uppercase for first word or all uppercase with "_" connected
const (
	ReqFrame byte = iota
	RetFrame
	//RpcFrame
	//IqFrame
	//PresenceFrame
	//JsonFrame
	//BinaryFrame
	//PingFrame
	//PongFrame
	TickFrame
	//CloseFrame

	ErrorFrame
	EchoFrame
)

func PackageMsg(msgtype uint8, id uint16, msgid uint16, data interface{}) []byte {
	ret := []byte{}
	databuff := Bytes(data)
	datalen := uint16(len(databuff))
	ret = append(ret, byte(msgtype))
	ret = append(ret, Bytes(id)...)
	ret = append(ret, Bytes(datalen)...)
	ret = append(ret, Bytes(msgid)...)

	if datalen > 0 {
		ret = append(ret, databuff...)
	}
	return ret
}

func ReadMsgHeader(reader io.Reader) (byte, uint16, uint16, uint16, []byte, error) {
	typebuff := make([]byte, 1)
	idbuff := make([]byte, 2)
	sizebuff := make([]byte, 2)
	msgidbuff := make([]byte, 2)
	var id uint16
	var size uint16
	var msgid uint16
	var databuff []byte

	_, err := io.ReadFull(reader, typebuff)
	if err != nil {
		goto end
	}

	//fmt.Println("data type:", typebuff[0])

	if typebuff[0] == TickFrame {
		goto end
	}

	_, err = io.ReadFull(reader, idbuff)
	if err != nil {
		goto end
	}
	id = Uint16(idbuff)

	//fmt.Println("id:", id)

	_, err = io.ReadFull(reader, sizebuff)
	if err != nil {
		goto end
	}
	size = Uint16(sizebuff)

	//fmt.Println("data size:", size)

	if size > 65535 {
		err = errors.New("too long data size")
		goto end
	}

	_, err = io.ReadFull(reader, msgidbuff)
	if err != nil {
		goto end
	}
	msgid = Uint16(msgidbuff)

	//fmt.Println("msgid:", msgid)

	if size == 0 {
		goto end
	}

	databuff = make([]byte, size)

	_, err = io.ReadFull(reader, databuff)
	if err != nil {
		goto end
	}
end:
	return typebuff[0], id, size, msgid, databuff, err
}

//server internal msg

//*********event*************
const SMsgId_UserOnline uint16 = 10000

type SMsgUserOnline struct {
	Uids       []uint64 `json:"uidarr"`
	Platforms  []string `json:"platform"`
	ServerAddr string   `json:"serveraddr"` //use # connect platform and serveraddr
}

const SMsgId_UserOffline uint16 = 10005

type SMsgUserOffline struct {
	Uids       []uint64 `json:"uidarr"`
	Platforms  string   `json:"platform"`
	ServerAddr string   `json:"serveraddr"` //use # connect platform and serveraddr
}

const SMsgId_RoomPresence uint16 = 10010

type SMsgRoomPresence struct {
	PresenceType uint8  `json:"presencetype"`
	Rid          uint64 `json:"rid"`
	Uid          uint64 `json:"uid"`
	Data         []byte `json:"data"`
}

const SMsgId_UserPresence uint16 = 10015

type SMsgUserPresence struct {
	To   uint64 `json:"to"`
	Data []byte `json:"data"`
}

// const SMsgId_RoomAddUser uint16 = 10010

// type SMsgRoomAddUser struct {
// 	Rid uint64 `json:"rid"`
// 	Uid uint64 `json:"uid"`
// }

// const SMsgId_RoomRemoveUser uint16 = 10015

// type SMsgRoomRemoveUser struct {
// 	Rid uint64 `json:"rid"`
// 	Uid uint64 `json:"uid"`
// }

// const SMsgId_RoomDimiss uint16 = 10020

// type SMsgRoomDimiss struct {
// 	Rid uint64 `json:"rid"`
// }

const SMsgId_ServerQuit uint16 = 10025

type SMsgServerQuit struct {
}

const SMsgId_ReqChatServerList uint16 = 10030

type SMsgReqChatServerList struct {
}

type SMsgRetChatServerList struct {
	Json string
}

//*********message*************
const SMsgId_UserMessage uint16 = 20000

type SMsgUserMessage struct {
	ServerAddr string `json:"serveraddr"`
	From       uint64 `json:"from"`
	To         uint64 `json:"to"`
	Data       []byte `json:"data"`
}

// type UserMessageData struct {
// 	Uid  uint64 `json:"uid"`
// 	Data []byte `json:"data"`
// }

const SMsgId_RoomMessage uint16 = 20005

type SMsgRoomMessage struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
	Data []byte `json:"data"`
}

const SMsgId_RoomAdminMessage uint16 = 20006

type SMsgRoomAdminMessage struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
	Data []byte `json:"data"`
}

const SMsgId_ZonePublicMessage uint16 = 20010

type SMsgZonePublicMessage struct {
	Appname  string `json:"appname"`
	Zonename string `json:"zonename"`
	Data     []byte `json:"data"`
}

const SMsgId_AppPublicMessage uint16 = 20015

type SMsgAppPublicMessage struct {
	Appname string `json:"appname"`
	Data    []byte `json:"data"`
}

const SMsgId_ServerPublicMessage uint16 = 20020

type SMsgServerPublicMessage struct {
	Data []byte `json:"data,string"`
}

//server internal msg end

//user event
const EventId_UserJoinRoom uint16 = 30000

type EventUserJoinRoom struct {
	Eventid uint16 `json:"eventid"`
	Rid     uint64 `json:"rid"`
}

const EventId_UserLeaveRoom uint16 = 30005

type EventUserLeaveRoom struct {
	Eventid uint16 `json:"eventid"`
	Rid     uint64 `json:"rid"`
}

const EventId_UserRoomAdmin uint16 = 30010
const EventId_UserRoomUnAdmin uint16 = 30015

//user event end

//room event
const EventId_RoomDismiss uint16 = 31000

//room event end

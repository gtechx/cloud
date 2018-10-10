package main

import (
	"fmt"
	"gtmsg"
	"time"

	"gtdb"

	. "github.com/gtechx/base/common"
)

func init() {
	RegisterRoomMsg()
}

func RegisterRoomMsg() {
	registerMsgHandler(MsgId_ReqCreateRoom, HandlerReqCreateRoom)
	registerMsgHandler(MsgId_ReqDeleteRoom, HandlerReqDeleteRoom)
	registerMsgHandler(MsgId_RoomPresence, HandlerReqRoomPresence)
	registerMsgHandler(MsgId_ReqUpdateRoomSetting, HandlerReqUpdateRoomSetting)
	registerMsgHandler(MsgId_ReqBanRoomUser, HandlerReqBanRoomUser)
	registerMsgHandler(MsgId_ReqJinyanRoomUser, HandlerReqJinyanRoomUser)
	registerMsgHandler(MsgId_ReqUnJinyanRoomUser, HandlerReqUnJinyanRoomUser)
	registerMsgHandler(MsgId_ReqAddRoomAdmin, HandlerReqAddRoomAdmin)
	registerMsgHandler(MsgId_ReqRemoveRoomAdmin, HandlerReqRemoveRoomAdmin)
	registerMsgHandler(MsgId_RoomMessage, HandlerRoomMessage)

	registerMsgHandler(MsgId_ReqRoomList, HandlerReqRoomList)
	registerMsgHandler(MsgId_ReqRoomPresenceList, HandlerReqRoomPresenceList)
	registerMsgHandler(MsgId_ReqRoomUserList, HandlerReqRoomUserList)
}

func HandlerReqCreateRoom(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE
	var roommsg *MsgReqCreateRoom = &MsgReqCreateRoom{}
	if !jsonUnMarshal(data, roommsg, &errcode) {
		return errcode, errcode
	}

	createRoom(sess.ID(), roommsg, &errcode)
	return errcode, errcode
}

func HandlerReqDeleteRoom(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isRoomOwner(rid, sess.ID(), &errcode) {
		return errcode, errcode
	}

	// var uselist []*gtdb.RoomUser
	// if !getRoomUserIds(rid, &uselist, &errcode) {
	// 	return errcode, errcode
	// }

	deleteRoom(sess.ID(), rid, &errcode)

	return errcode, errcode
}

func HandlerReqRoomPresence(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE
	var presence *MsgRoomPresence = &MsgRoomPresence{}
	if !jsonUnMarshal(data, presence, &errcode) {
		return errcode, errcode
	}

	if !isRoomExists(presence.Rid, &errcode) {
		return errcode, errcode
	}

	presencetype := presence.PresenceType
	rid := presence.Rid
	who := presence.Who

	presence.TimeStamp = time.Now().Unix()

	switch presencetype {
	case PresenceType_Subscribe:
		presence.Nickname = sess.NickName()
		presence.Who = sess.ID()

		if !isNotRoomUser(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isRoomNotFull(rid, &errcode) {
			return errcode, errcode
		}

		var roomtype byte
		if !getRoomType(rid, &roomtype, &errcode) {
			return errcode, errcode
		}

		if roomtype == RoomType_Apply {
			var presencebytes []byte
			if !jsonMarshal(presence, &presencebytes, &errcode) {
				return errcode, errcode
			}

			if !addRoomPresence(rid, sess.ID(), &errcode) {
				return errcode, errcode
			}

			// if !sendPresenceToRoomAdmin(rid, presencebytes, &errcode) {
			// 	return errcode, errcode
			// }
			// err := dbMgr.AddUserToRoomApplyList(rid, sess.ID(), time.Now().Unix())
			// if err != nil {
			// 	return ERR_DB, ERR_DB
			// }
			SendPresenceToRoomAdmin(sess.ID(), rid, presencebytes)
		} else if roomtype == RoomType_Everyone {
			if !addRoomUser(rid, sess.ID(), presence, &errcode) {
				return errcode, errcode
			}
		} else if roomtype == RoomType_Password {
			if !isRoomPassword(rid, presence.Password, &errcode) {
				return errcode, errcode
			}

			if !addRoomUser(rid, sess.ID(), presence, &errcode) {
				return errcode, errcode
			}
		}
	case PresenceType_Subscribed:
		if !isRoomAdmin(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isAppDataExists(who, &errcode) {
			return errcode, errcode
		}

		if !isRoomPresenceExists(rid, who, &errcode) {
			return errcode, errcode
		}

		if !addRoomUser(rid, who, presence, &errcode) {
			return errcode, errcode
		}

		var presencebytes []byte
		if !jsonMarshal(presence, &presencebytes, &errcode) {
			return errcode, errcode
		}

		senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
		errcode = SendPresenceToUser(who, senddata)

		if !removeRoomPresence(rid, who, &errcode) {
			return errcode, errcode
		}
	case PresenceType_UnSubscribe:
		if !isRoomUser(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isNotRoomOwner(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !removeRoomUser(rid, sess.ID(), presence, &errcode) {
			return errcode, errcode
		}

		if !removeRoomPresence(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}
	case PresenceType_UnSubscribed:
		if !isRoomAdmin(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isAppDataExists(who, &errcode) {
			return errcode, errcode
		}

		if !isRoomPresenceExists(rid, who, &errcode) {
			return errcode, errcode
		}

		var presencebytes []byte
		if !jsonMarshal(presence, &presencebytes, &errcode) {
			return errcode, errcode
		}

		senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
		errcode = SendPresenceToUser(who, senddata)

		if !removeRoomPresence(rid, who, &errcode) {
			return errcode, errcode
		}
	case PresenceType_Available, PresenceType_UnAvailable, PresenceType_Invisible:
		//send to my friend online
		// presencebytes, err := json.Marshal(presence)
		// if err != nil {
		// 	errcode = ERR_INVALID_JSON
		// } else {
		// 	senddata := packageMsg(RetFrame, 0, MsgId_Presence, presencebytes)
		// 	SendMessageToFriendsOnline(sess.ID(), senddata)
		// }
	}

	//ret := &MsgRetUserData{errcode, jsonbytes}
	return errcode, errcode
}

func HandlerReqUpdateRoomSetting(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE
	var roomsetting *MsgReqUpdateRoomSetting = &MsgReqUpdateRoomSetting{}
	if !jsonUnMarshal(data, roomsetting, &errcode) {
		return errcode, errcode
	}

	if roomsetting.Bit&RoomSetting_RoomName != 0 {
		err := dbMgr.SetRoomName(roomsetting.Rid, roomsetting.RoomName)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_RoomType != 0 {
		err := dbMgr.SetRoomType(roomsetting.Rid, roomsetting.RoomType)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Jieshao != 0 {
		err := dbMgr.SetRoomJieshao(roomsetting.Rid, roomsetting.Jieshao)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Notice != 0 {
		err := dbMgr.SetRoomNotice(roomsetting.Rid, roomsetting.Notice)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Password != 0 {
		err := dbMgr.SetRoomPassword(roomsetting.Rid, roomsetting.Password)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	return errcode, errcode
}

func HandlerReqBanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	banRoomUser(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqJinyanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	jinyanRoomUser(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqUnJinyanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	unjinyanRoomUser(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqAddRoomAdmin(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	addRoomAdmin(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqRemoveRoomAdmin(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	removeRoomAdmin(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerRoomMessage(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE
	var roommsg *MsgRoomMessage = &MsgRoomMessage{}
	if !jsonUnMarshal(data, roommsg, &errcode) {
		return errcode, errcode
	}

	if !isRoomExists(roommsg.Rid, &errcode) {
		return errcode, errcode
	}

	if !isRoomUser(roommsg.Rid, sess.ID(), &errcode) {
		return errcode, errcode
	}

	roommsg.TimeStamp = time.Now().Unix()
	roommsg.Who = sess.ID()
	roommsg.Nickname = sess.NickName()

	var msgbytes []byte
	if !jsonMarshal(roommsg, &msgbytes, &errcode) {
		return errcode, errcode
	}

	//最后一个if判断不需要
	// if !sendMessageToRoomUser(roommsg.Rid, msgbytes, &errcode) {
	// 	return errcode, errcode
	// }
	//sendMessageToRoomUser(roommsg.Rid, msgbytes, &errcode)

	// senddata = append(senddata, Bytes(platform)...)
	// senddata = append(senddata, msgbytes...)
	timstamp := time.Now().UnixNano()
	err := dbMgr.AddRoomMsgHistory(roommsg.Rid, msgbytes, timstamp)
	if err != nil {
		return ERR_DB, ERR_DB
	}

	err = dbMgr.SetRoomLastMsgTime(roommsg.Rid, timstamp)
	if err != nil {
		return ERR_DB, ERR_DB
	}

	return SendMsgToRoom(sess.ID(), roommsg.Rid, MsgId_RoomMessage, gtmsg.SMsgId_RoomMessage, msgbytes)
}

func SendMsgToRoom(from, rid uint64, msgid, exmsgid uint16, msgbytes []byte) (uint16, uint16) {
	senddata := packageMsg(RetFrame, 0, msgid, msgbytes)
	// serverlist, err := dbMgr.GetChatServerList()
	// if err != nil {
	// 	return ERR_DB, ERR_DB
	// }

	// msg := &gtmsg.SMsgRoomMessage{From: from, To: rid, Data: senddata}
	// msgdata, _ := json.Marshal(msg)
	// sendMsgToExchangeServer(exmsgid, msgdata)

	err := dbMgr.PubRoomMsg(rid, senddata)

	if err != nil {
		return ERR_DB, ERR_DB
	}

	//send to use on local server and offline users in room
	//SendMsgToLocalRoom(rid, senddata)
	return ERR_NONE, ERR_NONE
}

func SendPresenceToRoom(rid, uid uint64, presencetype uint8, msgbytes []byte) (uint16, uint16) {
	senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, msgbytes)
	// serverlist, err := dbMgr.GetChatServerList()
	// if err != nil {
	// 	return ERR_DB, ERR_DB
	// }

	// msg := &gtmsg.SMsgRoomPresence{PresenceType: presencetype, Rid: rid, Uid: uid, Data: senddata}
	// msgdata, _ := json.Marshal(msg)
	// sendMsgToExchangeServer(gtmsg.SMsgId_RoomPresence, msgdata)
	err := dbMgr.PubRoomMsg(rid, senddata)

	if err != nil {
		return ERR_DB, ERR_DB
	}

	//send to use on local server and offline users in room
	//SendMsgToLocalRoom(rid, senddata)
	return ERR_NONE, ERR_NONE
}

func SendPresenceToRoomAdmin(from, rid uint64, msgbytes []byte) (uint16, uint16) {
	senddata := packageMsg(RetFrame, 0, MsgId_Presence, msgbytes)

	// msg := &gtmsg.SMsgRoomAdminMessage{From: from, To: rid, Data: senddata}
	// msgdata, _ := json.Marshal(msg)
	// sendMsgToExchangeServer(gtmsg.SMsgId_RoomAdminMessage, msgdata)

	err := dbMgr.PubRoomAdminMsg(rid, senddata)

	if err != nil {
		return ERR_DB, ERR_DB
	}

	//send to use on local server and offline users in room
	//SendMsgToLocalRoomAdmin(rid, senddata)
	return ERR_NONE, ERR_NONE
}

func HandlerReqRoomList(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE

	var roomlist []*gtdb.Room
	if !getRoomList(sess.ID(), &roomlist, &errcode) {
		return errcode, errcode
	}

	var msgbytes []byte
	if !jsonMarshal(roomlist, &msgbytes, &errcode) {
		return errcode, errcode
	}

	ret := &MsgRetRoomList{}
	ret.Json = msgbytes
	ret.ErrorCode = errcode
	return errcode, ret
}

func HandlerReqRoomPresenceList(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	errcode := ERR_NONE

	ret := &MsgRetRoomPresenceList{}

	if !isRoomExists(rid, &ret.ErrorCode) {
		return errcode, ret
	}
	if !isRoomAdmin(rid, sess.ID(), &ret.ErrorCode) {
		return errcode, ret
	}
	var uids []uint64
	if !getRoomPresenceList(rid, &uids, &ret.ErrorCode) {
		return errcode, ret
	}

	if len(uids) > 0 {
		fmt.Println(uids)
		datalist, err := dbMgr.GetAppDatas(uids)
		if err != nil {
			ret.ErrorCode = ERR_DB
			return ERR_DB, ret
		}

		var msgbytes []byte
		if !jsonMarshal(datalist, &msgbytes, &ret.ErrorCode) {
			return errcode, ret
		}
		ret.Json = msgbytes
	}

	return errcode, ret
}

func HandlerReqRoomUserList(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	errcode := ERR_NONE

	if !isRoomUser(rid, sess.ID(), &errcode) {
		return errcode, errcode
	}

	var datalist []*gtdb.RoomUser
	if !getRoomUserList(rid, &datalist, &errcode) {
		return errcode, errcode
	}

	var msgbytes []byte
	if !jsonMarshal(datalist, &msgbytes, &errcode) {
		return errcode, errcode
	}

	ret := &MsgRetRoomUserList{}
	ret.Json = msgbytes
	ret.ErrorCode = errcode
	return errcode, ret
}

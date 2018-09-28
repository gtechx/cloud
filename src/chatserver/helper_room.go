package main

import (
	"encoding/json"
	"time"

	"gtdb"
	//. "github.com/gtechx/base/common"
)

//所有这类函数，返回false表示出错
func createRoom(appdataid uint64, roommsg *MsgReqCreateRoom, perrcode *uint16) bool {
	tbl_room := &gtdb.Room{Ownerid: appdataid, Roomname: roommsg.RoomName, Roomtype: roommsg.RoomType, Jieshao: roommsg.Jieshao, Notice: roommsg.Notice, Password: roommsg.Password}

	err := dbMgr.CreateRoom(tbl_room)

	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	return true
}

func deleteRoom(who, rid uint64, perrcode *uint16) bool {
	err := dbMgr.DeleteRoom(rid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	presence := &MsgRoomPresence{PresenceType: PresenceType_Dismiss, Rid: rid, Who: who, TimeStamp: time.Now().Unix()}
	presencebytes, err := json.Marshal(presence)
	if err != nil {
		*perrcode = ERR_INVALID_JSON
	} else {
		SendPresenceToRoom(rid, who, PresenceType_Dismiss, presencebytes)
		return true
	}

	return false
}

func getRoomList(appdataid uint64, proomlist *[]*gtdb.Room, perrcode *uint16) bool {
	roomlist, err := dbMgr.GetRoomListByJoined(appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}
	*proomlist = roomlist
	return true
}

func getRoomPresenceList(rid uint64, pdatalist *[]uint64, perrcode *uint16) bool {
	datalist, err := dbMgr.GetAllRoomPresence(rid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}
	*pdatalist = datalist
	return true
}

func isRoomFull(rid uint64, perrcode *uint16) bool {

	usercount, err := dbMgr.GetRoomUserCount(rid)
	if err != nil {
		*perrcode = ERR_DB
	} else {
		maxusercount, err := dbMgr.GetRoomMaxUser(rid)
		if err != nil {
			*perrcode = ERR_DB
		} else {
			if usercount != maxusercount {
				*perrcode = ERR_ROOM_NOT_FULL
			} else {
				return true
			}
		}
	}

	return false
}

func isRoomNotFull(rid uint64, perrcode *uint16) bool {

	usercount, err := dbMgr.GetRoomUserCount(rid)
	if err != nil {
		*perrcode = ERR_DB
	} else {
		maxusercount, err := dbMgr.GetRoomMaxUser(rid)
		if err != nil {
			*perrcode = ERR_DB
		} else {
			if usercount == maxusercount {
				*perrcode = ERR_ROOM_FULL
			} else {
				return true
			}
		}
	}

	return false
}

func addRoomUser(rid, appdataid uint64, presence *MsgRoomPresence, perrcode *uint16) bool {

	tbl_roomuser := &gtdb.RoomUser{Rid: rid, Dataid: appdataid}
	err := dbMgr.AddRoomUser(tbl_roomuser)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		presencebytes, err := json.Marshal(presence)
		if err != nil {
			*perrcode = ERR_INVALID_JSON
		} else {
			SendPresenceToRoom(rid, appdataid, PresenceType_Subscribed, presencebytes)
			return true
		}

		// msg := &gtmsg.SMsgRoomAddUser{Rid: rid, Uid: appdataid}
		// msgdata, _ := json.Marshal(msg)
		// sendMsgToExchangeServer(gtmsg.SMsgId_RoomAddUser, msgdata)

		// msg := &SMsgRoomAddUser{Rid: rid, Uid: appdataid}
		// msg.MsgId = SMsgId_RoomAddUser
		// msgbytes := Bytes(msg)

		// if broadcastServerEvent(msgbytes) != nil {
		// 	*perrcode = ERR_DB
		// 	return false
		// }

		// presencebytes, err := json.Marshal(presence)
		// if err != nil {
		// 	*perrcode = ERR_INVALID_JSON
		// } else {
		// 	senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
		// 	userlist, err := dbMgr.GetRoomUserIds(rid)

		// 	if err != nil {
		// 		*perrcode = ERR_DB
		// 	} else {
		// 		for _, user := range userlist {
		// 			//broadcast to user in room
		// 			errcode := SendMessageToUser(user.Dataid, senddata)

		// 			if errcode != ERR_NONE {
		// 				*perrcode = errcode
		// 				return false
		// 			}
		// 		}

		// 		return true
		// 	}
		// }
	}

	return false
}

// func sendMessageToRoomUser(rid uint64, msgbytes []byte, perrcode *uint16) bool {

// 	senddata := packageMsg(RetFrame, 0, MsgId_RoomMessage, msgbytes)
// 	userlist, err := dbMgr.GetRoomUserIds(rid)

// 	if err != nil {
// 		*perrcode = ERR_DB
// 	} else {
// 		for _, user := range userlist {
// 			//broadcast to user in room
// 			errcode := SendMessageToUser(user.Dataid, senddata)

// 			if errcode != ERR_NONE {
// 				*perrcode = errcode
// 				return false
// 			}
// 		}

// 		return true
// 	}

// 	return false
// }

// func sendPresenceToRoomUser(rid uint64, presencebytes []byte, perrcode *uint16) bool {

// 	senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
// 	userlist, err := dbMgr.GetRoomUserIds(rid)

// 	if err != nil {
// 		*perrcode = ERR_DB
// 	} else {
// 		for _, user := range userlist {
// 			//broadcast to user in room
// 			errcode := SendMessageToUser(user.Dataid, senddata)

// 			if errcode != ERR_NONE {
// 				*perrcode = errcode
// 				return false
// 			}
// 		}

// 		return true
// 	}

// 	return false
// }

// func sendPresenceToRoomAdmin(rid uint64, presencebytes []byte, perrcode *uint16) bool {

// 	senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
// 	userlist, err := dbMgr.GetRoomAdminIds(rid)

// 	if err != nil {
// 		*perrcode = ERR_DB
// 	} else {
// 		for _, user := range userlist {
// 			//broadcast to user in room
// 			errcode := SendMessageToUser(user.Dataid, senddata)

// 			if errcode != ERR_NONE {
// 				*perrcode = errcode
// 				return false
// 			}
// 		}

// 		return true
// 	}

// 	return false
// }

func addRoomPresence(rid, appdataid uint64, perrcode *uint16) bool {
	err := dbMgr.AddRoomPresence(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	// err = dbMgr.AddUserToRoomApplyList(rid, appdataid, time.Now().Unix())
	// if err != nil {
	// 	*perrcode = ERR_DB
	// 	return false
	// }
	return true
}

func removeRoomPresence(rid, appdataid uint64, perrcode *uint16) bool {
	err := dbMgr.RemoveRoomPresence(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	// err = dbMgr.RemoveUserFromRoomApplyList(rid, appdataid)
	// if err != nil {
	// 	*perrcode = ERR_DB
	// 	return false
	// }
	return true
}

func addRoomAdmin(rid, appdataid uint64, perrcode *uint16) bool {
	err := dbMgr.AddRoomAdmin(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}
	return true
}

func removeRoomAdmin(rid, appdataid uint64, perrcode *uint16) bool {
	err := dbMgr.RemoveRoomAdmin(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}
	return true
}

func jinyanRoomUser(rid, appdataid uint64, perrcode *uint16) bool {
	err := dbMgr.JinyanRoomUser(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	presence := &MsgRoomPresence{PresenceType: PresenceType_Jinyan, Rid: rid, Who: appdataid, TimeStamp: time.Now().Unix()}
	presencebytes, err := json.Marshal(presence)
	if err != nil {
		*perrcode = ERR_INVALID_JSON
	} else {
		SendPresenceToRoom(rid, appdataid, PresenceType_Jinyan, presencebytes)
		return true
	}

	return false
}

func unjinyanRoomUser(rid, appdataid uint64, perrcode *uint16) bool {
	err := dbMgr.UnJinyanRoomUser(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	presence := &MsgRoomPresence{PresenceType: PresenceType_UnJinyan, Rid: rid, Who: appdataid, TimeStamp: time.Now().Unix()}
	presencebytes, err := json.Marshal(presence)
	if err != nil {
		*perrcode = ERR_INVALID_JSON
	} else {
		SendPresenceToRoom(rid, appdataid, PresenceType_UnJinyan, presencebytes)
		return true
	}
	return false
}

func banRoomUser(rid, appdataid uint64, perrcode *uint16) bool {
	err := dbMgr.RemoveRoomUser(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	presence := &MsgRoomPresence{PresenceType: PresenceType_Ban, Rid: rid, Who: appdataid, TimeStamp: time.Now().Unix()}
	presencebytes, err := json.Marshal(presence)
	if err != nil {
		*perrcode = ERR_INVALID_JSON
	} else {
		SendPresenceToRoom(rid, appdataid, PresenceType_UnSubscribe, presencebytes)
		return true
	}

	return false
}

func removeRoomUser(rid, appdataid uint64, presence *MsgRoomPresence, perrcode *uint16) bool {
	err := dbMgr.RemoveRoomUser(rid, appdataid)
	if err != nil {
		*perrcode = ERR_DB
		return false
	}

	presencebytes, err := json.Marshal(presence)
	if err != nil {
		*perrcode = ERR_INVALID_JSON
	} else {
		SendPresenceToRoom(rid, appdataid, PresenceType_UnSubscribe, presencebytes)
		return true
	}

	// msg := &SMsgRoomRemoveUser{Rid: rid, Uid: appdataid}
	// msg.MsgId = SMsgId_RoomRemoveUser
	// msgbytes := Bytes(msg)

	// if broadcastServerEvent(msgbytes) != nil {
	// 	*perrcode = ERR_DB
	// 	return false
	// }
	return false
}

func isRoomPassword(rid uint64, password string, perrcode *uint16) bool {
	roompassword, err := dbMgr.GetRoomPassword(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if password != roompassword {
			*perrcode = ERR_ROOM_PASSWORD_INVALID
		} else {
			return true
		}
	}

	return false
}

func getRoomType(rid uint64, proomtype *byte, perrcode *uint16) bool {
	roomtype, err := dbMgr.GetRoomType(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		*proomtype = roomtype
		return true
	}

	return false
}

func isRoomExists(rid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomExists(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_ROOM_NOT_EXISTS
		} else {
			return true
		}
	}

	return false
}

func isRoomNotExists(rid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomExists(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if flag {
			*perrcode = ERR_ROOM_EXISTS
		} else {
			return true
		}
	}

	return false
}

func isRoomUser(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomUser(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_ROOM_USER_INVALID
		} else {
			return true
		}
	}

	return false
}

func isNotRoomUser(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomUser(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if flag {
			*perrcode = ERR_ROOM_USER_EXISTS
		} else {
			return true
		}
	}

	return false
}

func isRoomOwner(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomOwner(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_ROOM_NOT_OWNER
		} else {
			return true
		}
	}

	return false
}

func isNotRoomOwner(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomOwner(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if flag {
			*perrcode = ERR_ROOM_OWNER
		} else {
			return true
		}
	}

	return false
}

func isRoomAdmin(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomAdmin(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_ROOM_ADMIN_INVALID
		} else {
			return true
		}
	}

	return false
}

func isNotRoomAdmin(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomAdmin(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if flag {
			*perrcode = ERR_ROOM_ADMIN_EXISTS
		} else {
			return true
		}
	}

	return false
}

func getRoomUserIds(rid uint64, ids *[]*gtdb.RoomUser, perrcode *uint16) bool {
	userlist, err := dbMgr.GetRoomUserIds(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		*ids = userlist
		return true
	}

	return false
}

func getRoomAdminIds(rid uint64, ids *[]*gtdb.RoomUser, perrcode *uint16) bool {
	userlist, err := dbMgr.GetRoomAdminIds(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		*ids = userlist
		return true
	}

	return false
}

func isRoomPresenceExists(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := dbMgr.IsRoomPresenceExists(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_ROOM_PRESENCE_NOT_EXISTS
		} else {
			return true
		}
	}

	return false
}

func getRoomUserList(rid uint64, puserlist *[]*gtdb.RoomUser, perrcode *uint16) bool {
	userlist, err := dbMgr.GetRoomUserList(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		*puserlist = userlist
		return true
	}

	return false
}

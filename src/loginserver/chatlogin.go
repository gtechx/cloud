package main

import (
	"encoding/json"
	"gtdb"
	"io"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type ChatLoginRetMsg struct {
	ErrorDesc  string `json:"errordesc,omitempty"`
	ErrorCode  uint16 `json:"errcode"`
	Token      string `json:"token,omitempty"`
	ServerAddr string `json:"serveraddr,omitempty"`
	UserData   string `json:"userdata,omitempty"`
}

func getChatLoginToken(databytes []byte) (string, error) {
	uu, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	token := uu.String()

	err = dbMgr.SaveChatLoginToken(token, databytes, srvconfig.TokenTimeout)

	if err != nil {
		return "", err
	}

	return token, nil
}

func checkAppname(appname string) (uint16, string) {
	flag, err := dbMgr.IsAppExists(appname)

	if err != nil {
		return 1, err.Error()
	}

	if !flag {
		return 1, "appname not exists"
	}

	return 0, ""
}

func chatlogin(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		account := req.PostFormValue("account")
		password := req.PostFormValue("password")
		appname := req.PostFormValue("appname")
		//platform := req.PostFormValue("platform")
		ret := &ChatLoginRetMsg{}
		var err error
		var appdata *gtdb.AppData
		var databytes []byte

		ret.ErrorCode, ret.ErrorDesc = checkLogin(account, password)
		if ret.ErrorCode != 0 {
			goto end
		}

		appdata, err = dbMgr.GetAppDataByAccount(account, appname)
		if err != nil {
			ret.ErrorCode = 1
			ret.ErrorDesc = err.Error()
		} else if appdata == nil {
			//tell client no appdata
			ret.ErrorCode = 2
			ret.ErrorDesc = "no user data"
		} else {
			databytes, err = json.Marshal(appdata)
			if err != nil {
				ret.ErrorCode = 3
				ret.ErrorDesc = err.Error()
			} else {
				//get chat server addr
				ret.ServerAddr = minUserServer
				ret.Token, err = getChatLoginToken(databytes)
				if err != nil {
					ret.ErrorCode = 1
					ret.ErrorDesc = err.Error()
				} else {
					ret.UserData = string(databytes)
				}
			}
		}
	end:
		data, _ := json.Marshal(&ret)
		io.WriteString(rw, string(data))
	}
}

func chatcreateuser(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		account := req.PostFormValue("account")
		password := req.PostFormValue("password")
		appname := req.PostFormValue("appname")
		nickname := req.PostFormValue("nickname")
		var err error
		var tbl_appdata *gtdb.AppData
		var databytes []byte
		ret := &ChatLoginRetMsg{}

		ret.ErrorCode, ret.ErrorDesc = checkLogin(account, password)
		if ret.ErrorCode != 0 {
			goto end
		}

		ret.ErrorCode, ret.ErrorDesc = checkAppname(appname)
		if ret.ErrorCode != 0 {
			goto end
		}

		dbMgr = dbMgr
		tbl_appdata, err = dbMgr.GetAppDataByNickname(nickname, appname)
		if err != nil {
			ret.ErrorCode = 1
			ret.ErrorDesc = err.Error()
			goto end
		}

		if tbl_appdata != nil {
			ret.ErrorCode = 1
			ret.ErrorDesc = "nickname already exists"
			goto end
		}

		tbl_appdata = &gtdb.AppData{Appname: appname, Zonename: "default", Account: account, Nickname: nickname, Regip: req.RemoteAddr}
		err = dbMgr.CreateAppData(tbl_appdata)

		if err != nil {
			ret.ErrorCode = 1
			ret.ErrorDesc = err.Error()
		} else {
			databytes, err = json.Marshal(tbl_appdata)
			if err != nil {
				ret.ErrorCode = 3
				ret.ErrorDesc = err.Error()
			} else {
				//get chat server addr
				ret.ServerAddr = minUserServer
				ret.Token, err = getChatLoginToken(databytes)
				if err != nil {
					ret.ErrorCode = 1
					ret.ErrorDesc = err.Error()
				} else {
					ret.UserData = string(databytes)
				}
			}
		}
	end:
		data, _ := json.Marshal(&ret)
		io.WriteString(rw, string(data))
	}
}

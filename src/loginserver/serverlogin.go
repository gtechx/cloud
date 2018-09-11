package main

import "net/http"

func serverlogin(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		//check ip

		//check account password appname
		// account := req.PostFormValue("account")
		// password := req.PostFormValue("password")
		// appname := req.PostFormValue("appname")

		//get internal server addr
	}
}

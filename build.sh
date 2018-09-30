#!/bin/sh

export GOPATH=${PWD}
export GOBIN="${PWD}/bin"

go install -tags debug loginserver 
go install -tags debug chatserver 
go install -tags debug exchangeserver


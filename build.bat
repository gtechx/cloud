call killserver.bat
md bin
set GOPATH=%~dp0
set GOBIN=%~dp0bin
REM set GOOS=windows

go install -tags debug loginserver 
go install -tags debug chatserver 

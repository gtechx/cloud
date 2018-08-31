md bin
set GOPATH=%~dp0
set GOBIN=%~dp0bin
REM set GOOS=windows

go install -tags debug -race loginserver 
go install -tags debug -race chatserver 

pause
cd bin

echo starting loginserver
start loginserver.exe

echo starting chatserver
start chatserver.exe -config="../res/config/chatserver.config"

echo starting chatserver2
start chatserver.exe -config="../res/config/chatserver2.config"

REM ping -n 1 127.1 >nul
REM echo start all done

cd ..
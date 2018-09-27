cd bin

echo starting exchangeserver
start exchangeserver.exe

ping -n 2 127.1 >nul

echo starting chatserver
start chatserver.exe -config="../res/config/chatserver.config"

echo starting chatserver2
start chatserver.exe -config="../res/config/chatserver2.config"

echo starting loginserver
start loginserver.exe

REM ping -n 1 127.1 >nul
REM echo start all done

cd ..
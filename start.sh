#!/bin/sh

cd bin

nohup ./exchangeserver &
sleep 1

nohup ./chatserver &
sleep 1

nohup ./loginserver &
sleep 1

cd ..


#!/bin/sh

ps ux | grep -i -E 'loginserver|chatserver|exchangeserver' |awk '{print $2}'| xargs  kill


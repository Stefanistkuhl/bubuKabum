#!/bin/sh
cd /bubuKabum/bot
node deploy-commands.js
node . &
cd /bubuKabum/converter
go run . &
wait

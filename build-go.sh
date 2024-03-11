#!/bin/sh

cd go-backend

mkdir -p ../go-output

go build -o ../go-output/chatroom .

cd ..


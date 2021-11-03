#!/bin/zsh
source /etc/profile
cp -r ./envs/$1-conf/* ./

go mod tidy
go build -o msgCenter

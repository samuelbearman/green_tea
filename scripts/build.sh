#!/bin/bash

# Linpeas
curl -L https://github.com/carlospolop/PEASS-ng/releases/latest/download/linpeas.sh -o ./tools/linpeas.sh
# Winpeas
curl -L https://github.com/carlospolop/PEASS-ng/releases/download/20220828/winPEASany.exe -o ./tools/winPEASany.exe
#Netcat
curl -L https://github.com/yunchih/static-binaries/raw/master/nc -o ./tools/nc
# Socat
curl -L https://github.com/andrew-d/static-binaries/raw/master/binaries/linux/x86_64/socat -o ./tools/socat

# Go build
go build -ldflags "-s -w" -o ./bin/gt 

# UPX pack
upx ./bin/gt
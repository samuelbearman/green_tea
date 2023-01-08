#!/bin/bash

# Linpeas
curl -L https://github.com/carlospolop/PEASS-ng/releases/latest/download/linpeas.sh -o ./tools/linpeas.sh
# Winpeas
curl -L https://github.com/carlospolop/PEASS-ng/releases/download/20220828/winPEASany.exe -o ./tools/winPEASany.exe
#Netcat
curl -L https://github.com/yunchih/static-binaries/raw/master/nc -o ./tools/nc
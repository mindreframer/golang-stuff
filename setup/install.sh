#!/bin/bash -e

# copy gandalf command to a place in path
go build -o gandalf bin/gandalf.go
sudo mv gandalf /usr/local/bin/

# copy default config file
sudo cp etc/gandalf.conf /etc/

# starts gandalf api web server
go build -o gandalf-webserver webserver/main.go
./gandalf-webserver > $HOME/gandalf-webserver.out 2>&1 &
git daemon --base-path=/var/repositories --syslog --export-all &

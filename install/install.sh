#!/bin/bash

DIR_SCRIPT=$(dirname -- "${BASH_SOURCE[0]}")
logname
echo $DIR_SCRIPT

sudo cp $DIR_SCRIPT/lyncser.service /etc/systemd/system/lyncser.service
sudo cp $DIR_SCRIPT/lyncser.timer /etc/systemd/system/lyncser.timer

sudo sed -i "s/###user###/$(logname)/g" /etc/systemd/system/lyncser.service

sudo systemctl start lyncser.service
sudo systemctl enable lyncser.service
sudo systemctl start lyncser.timer
sudo systemctl enable lyncser.timer

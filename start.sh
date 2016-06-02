#!/bin/sh
raspi-gpio set 17 pu
raspi-gpio set 27 pu
raspi-gpio set 22 pu
xset -dpms
xset s off
matchbox-window-manager -use_titlebar no -use_cursor no &
sleep 15
iceweasel http://localhost:8080 &
sleep 15
xdotool key --clearmodifiers F11
#!/bin/bash

 # the only argument $1 is the mac address
 # use this command to get the mac address for this command:
 # ./send 57 3E 02 01 00 00
 if [ "$1" == "" ]; then
   echo "usage: cd /bin & ./test_device.sh '<mac address>' "; exit 1
 fi
 echo ***MOTOR***
 ./send 57 3E 08 0D 10 00 $1 00 FF 06 F9 00 FF 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00
 sleep 1
 ./send 57 3E 08 0D 10 00 $1 00 FF 00 FF 00 FF 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00
 sleep 1
 echo ***ALARM***
 ./send 57 3E 08 0D 10 00 $1 00 FF 10 EF 00 FF 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00
 sleep 1
 ./send 57 3E 08 0D 10 00 $1 00 FF 00 FF 00 FF 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00
 echo ***SENSOR***
 ./send 57 3E 08 0D 10 00 $1 00 FF 00 FF 01 FE 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00
 echo ***TEMP/HUMID***
 ./send 57 3E 08 64 10 00 $1 00 FF 00 FF 00 FF 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00
 sleep 1
 ./send 57 3E 83 01 00 00
 echo '***BUTTON(hard to test, please ignore it)***'
 ./send 57 3E 08 0D 10 00 $1 01 FE 00 FF 00 FF 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00
 sleep 1
 ./send 57 3E 08 0D 10 00 $1 00 FF 00 FF 00 FF 00 FF
 sleep 1
 ./send 57 3E 83 01 00 00

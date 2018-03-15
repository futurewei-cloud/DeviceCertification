# WRTNode-2P

## Operation
## setwifi
The command to set WIFI to connect to the router

## Communicate with zigbee
* Allow the connection
```
root@LEDE:/bin# ./send 57 3E 01 01 00 00
 57 3e 01 01 00 00
write 6
read 7
buf: 57 3c 01 01 01 00 00
```
通知2P，CC2530已经开启60s入网模式

* Get online device list
```
root@LEDE:/bin# ./send 57 3E 02 01 00 00
 57 3e 02 01 00 00
write 6
read 15
buf: 57 3c 02 01 09 00 0a ea d6 a6 16 00 4b 12 00
```
|Return Code|Meanning|
|---|---|
|57 3c|Fixed header|
|02|Online devices number flag|
|01|There is one device online right now|
|09|The length of the data|
|00|Checking data, always be 00 for this testing board|
|0a|The first device's type, always be 0a in this testing board|
|ea d6 a6 16 00 4b 12 00|data section, the MAC of the first device|

Note that there will be more 9 characters sets if there are more devices online. Support 8 devices maximum right now.

## Tools
### send
The command line tool to send message to serial port on board

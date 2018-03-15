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

* Set the alarm
Turn on the alarm
```
root@LEDE:/bin# send 57 3E 08 0D 10 00 ea d6 a6 16 00 4b 12 00 00 FF 10 EF 00 FF 00 FF
 57 3e 08 0d 10 00 ffffffea ffffffd6 ffffffa6 16 00 4b 12 00 00 ffffffff 10 ffffffef 00 ffffffff 00 ffffffff
write 22
read 7
buf: 57 3c 08 0d 01 00 00
```
|Command|Meanning|
|---|---|
|57 3E|Fixed header|
|08|Read/Write IO|
|0d|CID|
|10|The length of the data|
|00|Checking data, always be 00 for this testing board|
|ea d6 a6 16 00 4b 12 00|data section, the MAC of the first device|
|00 FF 10 EF 00 FF 00 FF|Command to set alarm|

Turn off the alarm
```
root@LEDE:/bin# send 57 3E 08 0D 10 00 ea d6 a6 16 00 4b 12 00 00 FF 00 FF 00 FF 00 FF
 57 3e 08 0d 10 00 ffffffea ffffffd6 ffffffa6 16 00 4b 12 00 00 ffffffff 00 ffffffff 00 ffffffff 00 ffffffff
write 22
read 7
buf: 57 3c 08 0d 01 00 00
```

* Get output from DTH11
Send command to DTH11 to ask for data
```
root@LEDE:~# send 57 3E 08 64 10 00 ea d6 a6 16 00 4b 12 00 00 FF 00 FF 00 FF 00 FF
 57 3e 08 64 10 00 ffffffea ffffffd6 ffffffa6 16 00 4b 12 00 00 ffffffff 00 ffffffff 00 ffffffff 00 ffffffff
write 22
read 7
buf: 57 3c 08 64 01 00 00
```
Retrive data from the DTH11 sensor
```
root@LEDE:~# send 57 3E 83 01 00 00
 57 3e ffffff83 01 00 00
write 6
read 19
buf: 57 3c 83 64 0d 38 01 ea d6 a6 16 00 4b 12 00 1c 00 1b 00
```
|Return Code|Meanning|
|---|---|
|57 3c|Fixed header|
|83|Data returned from device flag|
|64|The CID of the returned data, 64 is DTH11|
|0d|The length of the data|
|38|Checking data|
|01|There is 1 data totally, which is the return value of this command|
|ea d6 a6 16 00 4b 12 00|Data section, the MAC of the first device|
|1c 00 1b 00|Data section, 1c00 is the temperature, 1b00 is the humidity|



## Device accessing
|Function|Pin|Memo|
|Motor|P1.1, P1.2||
|DTH11|P1.5||
|Alarm|P1.4|Active High|
|Button|P0.0|Long press for reset|
|Cover sensor|P2.0||
For example, the data to set alarm on is `00 FF 10 EF 00 FF 00 FF` and to switch it off is `00 FF 00 FF 00 FF 00 FF`

## Tools
### send
The command line tool to send message to serial port on board

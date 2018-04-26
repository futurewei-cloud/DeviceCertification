# The wrtnodedriver

Supports multiple demoboards for single Wrtnode Edge Box. Currently the example of demoboardmapper.go and demoboardapp.go are only supporting 2 demoboards and limited sensor/device functions. 

## 1. Build the driver mappers and function application:
Before running build, please set environment variables:

	export GOOS=linux
	export GOARCH=arm
	export GOARM=7

	go build demoboardmapper.go
	go build demoboardapp.go

## 2. then copy built binaries to wrtnode-2p box

> NOTE: Please create schemas with following contents before start application/driver mappers:

	[
	    {
	        "deviceid": "demoboard/cover1",
	        "direction": "source",
	        "description": "first cover sensor",
	        "valuetype": "Integer:0:1"
	    },
	    {
	        "deviceid": "demoboard/cover2",
	        "direction": "source",
	        "description": "second cover sensor",
	        "valuetype": "Integer:0:1"
	    },
	    {
	        "deviceid": "demoboard/motor1",
	        "direction": "target",
	        "description": "first motor",
	        "valuetype": "Integer:0:1"
	    },
	    {
	        "deviceid": "demoboard/motor2",
	        "direction": "target",
	        "description": "second motor",
	        "valuetype": "Integer:0:1"
	    }
	]

## 3. start the drivers in different remote consoles

	./demoboardmapper e1 motor
	./demoboardmapper e1 cover
	./application e1 e1 motor1

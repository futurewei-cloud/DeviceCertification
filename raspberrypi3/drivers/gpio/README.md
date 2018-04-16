# The raspberrypi3/gpio drivers

Supports multiple devices for single Raspberry Pi 3 B Box. Refer to <https://github.com/stianeikeland/go-rpio>

Currently the example of lightmapper.go, switchmapper.go and application.go are only supporting one pair of light/switch device functionalities. 

## 1. Build the driver mappers and function application:
Before running build, please set environment variables:

	export GOOS=linux
	export GOARCH=arm
	export GOARM=7

	go build lightmapper.go
	go build switchmapper.go
	go build application.go

## 2. then copy to Raspberry Pi box

> NOTE: Please create schemas with following contents before start application/driver mappers:

	[
	    {
	        "deviceid": "switch",
	        "direction": "source",
	        "description": "The switch device",
	        "valuetype": "Integer:0:1"
	    },
	    {
	        "deviceid": "light",
	        "direction": "source",
	        "description": "The light device",
	        "valuetype": "Integer:0:1"
	    }
	]

## 3. start the drivers in different pi consoles:

	./lightmapper
	./switchmapper
	./application

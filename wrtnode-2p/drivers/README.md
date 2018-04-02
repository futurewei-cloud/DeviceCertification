The wrtnodedriver supports multiple demoboards for single Wrtnode Edge Box. Currently the example of demomapper.go and demoapp.go are only supports 2 demoboards and limited sensor/device functions. 
Please create schemas with following contents(note: revision, createTime, modifiedTime are system generated):

[
    {
        "deviceid": "demoboard/coversensor1",
        "revision": 18,
        "createTime": 1522611513,
        "modifiedTime": 1522611513,
        "direction": "source",
        "description": "first cover sensor",
        "valuetype": "Integer:0:1"
    },
    {
        "deviceid": "demoboard/coversensor2",
        "revision": 18,
        "createTime": 1522611513,
        "modifiedTime": 1522611513,
        "direction": "source",
        "description": "second cover sensor",
        "valuetype": "Integer:0:1"
    },
    {
        "deviceid": "demoboard/motor1",
        "revision": 18,
        "createTime": 1522611513,
        "modifiedTime": 1522611513,
        "direction": "target",
        "description": "first motor",
        "valuetype": "Integer:0:1"
    },
    {
        "deviceid": "demoboard/motor2",
        "revision": 18,
        "createTime": 1522611513,
        "modifiedTime": 1522611513,
        "direction": "target",
        "description": "second motor",
        "valuetype": "Integer:0:1"
    }
]

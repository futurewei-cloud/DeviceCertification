---
title: 'Edge Cloud V1.0 Developer Manual'
---

Edge Cloud Service components
=============================

Hardware
--------

1)  Cloud VM (e.g. Linux Ubuntu)

2)  Edge node (e.g. Wrtnode running Mips or Rasp Berry Pi running ARM
    v7)

3)  Devices (Embedded OS burned by device manufacturer, can only
    communicate with through specific protocol, such as Zigbee)

![](https://farm1.staticflickr.com/886/41150238182_182d9ecaa0_n.jpg){width="3.7105172790901135in"
height="3.8849431321084866in"}

> Note: the black box with two antennas is the Edge Node, running Linux
> Mips OS with 128 MB memory, and users only get 61 MB for their apps.
> The other device is the demo board mount with multiple small devices:
> 1 infra sensor with red light, 1 motor, 1 temp/humid detector and 1
> push button.

![cid:8877414293326061511900973](https://farm1.staticflickr.com/785/39384760460_21d66a3f48.jpg){width="6.5in"
height="4.875429790026247in"}

Raspberry Pi 3 with 1 GB memory which can connection devices directly
through GPIO

Software
--------

1)  Center Cloud Service at Public Cloud VM

-   Service Bus with ETCD as KV store (golang)

2)  Edge Cloud Service (for each Edge Node)

    -   Logical Device Repository Service (LDRS) with Service Bus/ETCD
        (edge command in golang)

    -   Application (currently running as plugin app, in future will be
        managed by Serverless framework) (nodejs/python/golang)

    -   Device Driver (running as plugin) (nodejs/python/golang)

Relations of components
-----------------------

Scenarios
=========

1)  Manage edge local logical devices from cloud

-   Create logical device for specified edge node (POST)

-   Read device state (GET)

-   Update expected device state for target device, and expect the
    linked device takes action (PUT)

2)  Autonomous on Edge node

    -   Cover the sensor with finger, the LED is light up. The local
        LDRS states are updated by deviceDriverCover and ControlMotorApp
        and local target Motor1 starts running.

    -   Release the finger from the sensor, the states are reverted back
        to original and the motor should stop running.

3)  Control devices between edge nodes

    -   In 2), if the target device is located inside different edge
        node, the state update is sent to the other edge node through
        the Service Bus, with central cloud's help. Repeat the test
        in 2) will cause the remote Motor2 spins.

Config center and edge cloud
----------------------------

1)  Start etcd server on central cloud (./etcd)

2)  Start service bus server (go run server.go)

3)  Load configuration to etcd server.

    ./etcdctl put
    Root/System/Configure/Public/Common/EdgeList/HuaweiProject1/center
    '{\"Key\":
    \"Root/System/Configure/Public/Common/EdgeList/HuaweiProject1/center\",\
    \"Revision\": 1,\
    \"Value\": {\
    \"ClusterID\": 1,\
    \"EdgeID\": 1,\
    \"EdgeName\": \"center\",\
    \"ProjectID\": \"HuaweiProject1\"\
    }\
    }'

4)  Start central metaDB service (run edge command)

    a.  Using Config.yaml (put to the folder where all applications
        stays)

        i.  serverAddress: \"34.209.89.1:10000\"\
            clusterName: \"HuaweiProject1\"\
            edgeName: \"center\"\
            localTcpPort: 8080

5)  Create initial configurations (POST
    [http://hostname:port/v1.0/{ProjectID}/edgecloud/edges/center/metadata/configure?operation=batch](http://hostname:port/v1.0/%7bProjectID%7d/edgecloud/edges/center/metadata/configure?operation=batch)
    )

+-----------------------------------+-----------------------------------+
| Method                            | POST                              |
+===================================+===================================+
| URL                               | <http://34.209.89.1:8080/v1.0/Hua |
|                                   | weiProject1/edgecloud/edges/cente |
|                                   | r/metadata/configure?operation=ba |
|                                   | tch>                              |
+-----------------------------------+-----------------------------------+
| BODY                              | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | \"Key\":                          |
|                                   | \"Root/System/Configure/Edges/Hua |
|                                   | weiProject1/e1/Services/MetadataD |
|                                   | B/SyncToLocal\",                  |
|                                   |                                   |
|                                   | \"Value\": \[                     |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | \"EndKey\": \"\",                 |
|                                   |                                   |
|                                   | \"StartKey\":                     |
|                                   | \"Root/System/Configure/Public/Co |
|                                   | mmon/EdgeList/\"                  |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
|                                   |                                   |
|                                   | },                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | \"Key\":                          |
|                                   | \"Root/System/Configure/Edges/Hua |
|                                   | weiProject1/e2/Services/MetadataD |
|                                   | B/SyncToLocal\",                  |
|                                   |                                   |
|                                   | \"Value\": \[                     |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | \"EndKey\": \"\",                 |
|                                   |                                   |
|                                   | \"StartKey\":                     |
|                                   | \"Root/System/Configure/Public/Co |
|                                   | mmon/EdgeList/\"                  |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
|                                   |                                   |
|                                   | },                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | \"Key\":                          |
|                                   | \"Root/System/Configure/Public/Co |
|                                   | mmon/EdgeList/HuaweiProject1/e1\" |
|                                   | ,                                 |
|                                   |                                   |
|                                   | \"Value\": {                      |
|                                   |                                   |
|                                   | \"ClusterID\": 1,                 |
|                                   |                                   |
|                                   | \"EdgeID\": 2,                    |
|                                   |                                   |
|                                   | \"EdgeName\": \"e1\",             |
|                                   |                                   |
|                                   | \"ProjectID\": \"HuaweiProject1\" |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | },                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | \"Key\":                          |
|                                   | \"Root/System/Configure/Public/Co |
|                                   | mmon/EdgeList/HuaweiProject1/e2\" |
|                                   | ,                                 |
|                                   |                                   |
|                                   | \"Value\": {                      |
|                                   |                                   |
|                                   | \"ClusterID\": 1,                 |
|                                   |                                   |
|                                   | \"EdgeID\": 3,                    |
|                                   |                                   |
|                                   | \"EdgeName\": \"e2\",             |
|                                   |                                   |
|                                   | \"ProjectID\": \"HuaweiProject1\" |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK "Succeed"                  |
|                                   |                                   |
|                                   | 200 OK "Key already exist" --     |
|                                   | etcd error message                |
|                                   |                                   |
|                                   | 4xx ERROR HTTP error messages     |
+-----------------------------------+-----------------------------------+

6)  Or add new configurations (POST
    [http://hostname:port/v1.0/{ProjectID}/edgecloud/edges/center/metadata/configure?operation=batch](http://hostname:port/v1.0/%7bProjectID%7d/edgecloud/edges/center/metadata/configure?operation=batch)
    )

+-----------------------------------+-----------------------------------+
| Method                            | POST                              |
+===================================+===================================+
| URL                               | <http://34.209.89.1:8080/v1.0/Hua |
|                                   | weiProject1/edgecloud/edges/cente |
|                                   | r/metadata/configure?operation=ba |
|                                   | tch>                              |
+-----------------------------------+-----------------------------------+
| BODY                              | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | \"Key\":                          |
|                                   | \"Root/System/Configure/Public/Co |
|                                   | mmon/EdgeList/HuaweiProject1/e3\" |
|                                   | ,                                 |
|                                   |                                   |
|                                   | \"Value\": {                      |
|                                   |                                   |
|                                   | \"ClusterID\": 1,                 |
|                                   |                                   |
|                                   | \"EdgeID\": 4,                    |
|                                   |                                   |
|                                   | \"EdgeName\": \"e3\",             |
|                                   |                                   |
|                                   | \"ProjectID\": \"HuaweiProject1\" |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK "Succeed"                  |
|                                   |                                   |
|                                   | 200 OK "Key already exist" --     |
|                                   | etcd error message                |
|                                   |                                   |
|                                   | 4xx ERROR HTTP error messages     |
+-----------------------------------+-----------------------------------+

Manage local logic edge devices from cloud
------------------------------------------

On cloud VM, customer uses Postman of Chrome, or any Restful client
application, send following requests (Upon success of this step,
framework will automatically create key entries for actual/expected
states):

### Create logical device for specified edge node

+-----------------------------------+-----------------------------------+
| Method                            | POST                              |
+===================================+===================================+
| URL                               | http://34.209.89.1:8080/v1.0/Huaw |
|                                   | eiProject1/edgecloud/edges/e1/ldr |
|                                   | s/schema/?recursive=true          |
+-----------------------------------+-----------------------------------+
| BODY                              | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "DeviceId": "light1",             |
|                                   |                                   |
|                                   | "ValueType": "Integer:0:1",       |
|                                   |                                   |
|                                   | "Direction": "target",            |
|                                   |                                   |
|                                   | "Description": "light sensor"     |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK "Succeed"                  |
|                                   |                                   |
|                                   | 200 OK "Key already exist" --     |
|                                   | etcd error message                |
|                                   |                                   |
|                                   | 4xx ERROR HTTP error messages     |
+-----------------------------------+-----------------------------------+

### Read device state

+-----------------------------------+-----------------------------------+
| Method                            | GET                               |
+===================================+===================================+
| URL                               | <http://34.209.89.1:8080/v1.0/Hua |
|                                   | weiProject1/edgecloud/edges/e1/ld |
|                                   | rs/actual/motor1/>                |
|                                   | \[?recursive=true\]               |
+-----------------------------------+-----------------------------------+
| BODY                              | N/A                               |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 1.  Value of motor1 in JSON       |
|                                   |                                   |
|                                   | 2.  (?recursive=true) values of   |
|                                   |     all keys in Motor1 group in   |
|                                   |     JSON                          |
|                                   |                                   |
|                                   | 3.  "key not exist" if key is not |
|                                   |     defined                       |
|                                   |                                   |
|                                   | BODY:                             |
|                                   |                                   |
|                                   | e.x:                              |
|                                   |                                   |
|                                   | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "key": "motor1",                  |
|                                   |                                   |
|                                   | "Value": 1,                       |
|                                   |                                   |
|                                   | "Revisions": 43,                  |
|                                   |                                   |
|                                   | "CreateTime": 1522435664          |
|                                   |                                   |
|                                   | "ModifiedTime": 1522435664        |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+

### Update expected device state for target device

+-----------------------------------+-----------------------------------+
| Method                            | PUT                               |
+===================================+===================================+
| URL                               | <http://34.209.89.1:8080/v1.0/Hua |
|                                   | weiProject1/edgecloud/edges/e1/ld |
|                                   | rs/expected/?recursive=true>      |
+-----------------------------------+-----------------------------------+
| BODY                              | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "key": "light1",                  |
|                                   |                                   |
|                                   | "Value": "1",                     |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK "Succeed"                  |
|                                   |                                   |
|                                   | 200 OK "Key already exist" --     |
|                                   | etcd error message                |
|                                   |                                   |
|                                   | 4xx ERROR HTTP error messages     |
+-----------------------------------+-----------------------------------+

Autonomous on Edge node
-----------------------

Device driver for cover sensor will use following HTTP request to update
local LDRS:

+-----------------------------------+-----------------------------------+
| Method                            | PUT                               |
+===================================+===================================+
| URL                               | <http://localhost:8080/v1.0/Huawe |
|                                   | iProject1/edgecloud/edges/e1/ldrs |
|                                   | /actual/>                         |
+-----------------------------------+-----------------------------------+
| BODY                              | 1.  Value of motor1 in JSON       |
|                                   |                                   |
|                                   | 2.  (?recursive=true) values of   |
|                                   |     all keys in motor1 group in   |
|                                   |     JSON                          |
|                                   |                                   |
|                                   | 3.  "key not exist" if key is not |
|                                   |     defined                       |
|                                   |                                   |
|                                   | BODY:                             |
|                                   |                                   |
|                                   | e.x:                              |
|                                   |                                   |
|                                   | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "key": "motor1",                  |
|                                   |                                   |
|                                   | "Value": 1                        |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "key": "light1",                  |
|                                   |                                   |
|                                   | "Value": 1,                       |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK "Succeed"                  |
|                                   |                                   |
|                                   | 200 OK "Key already exist" --     |
|                                   | etcd error message                |
|                                   |                                   |
|                                   | 4xx ERROR HTTP error messages     |
+-----------------------------------+-----------------------------------+

Device driver for Motor1 will use following HTTP connection to watch the
local LDRS:

+-----------------------------------+-----------------------------------+
| Method                            | GET                               |
+===================================+===================================+
| URL                               | <http://localhost:8080/v1.0/Huawe |
|                                   | iProject1/edgecloud/edges/e1/ldrs |
|                                   | /expected/Motor1/?watch=true>     |
+-----------------------------------+-----------------------------------+
| BODY                              | N/A                               |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK                            |
|                                   |                                   |
|                                   | 404 NotFound                      |
|                                   |                                   |
|                                   | 403 Forbidden                     |
|                                   |                                   |
|                                   | BODY:                             |
|                                   |                                   |
|                                   | 1.  Value of Motor1 in JSON       |
|                                   |                                   |
|                                   | 2.  (?recursive=true) values of   |
|                                   |     all keys in Motor1 group in   |
|                                   |     JSON                          |
|                                   |                                   |
|                                   | 3.  "key not exist" if key is not |
|                                   |     defined                       |
|                                   |                                   |
|                                   | e.x:                              |
|                                   |                                   |
|                                   | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "key": "motor1",                  |
|                                   |                                   |
|                                   | "Value": "On",                    |
|                                   |                                   |
|                                   | "Revisions": 43,                  |
|                                   |                                   |
|                                   | "CreateTime": 1522435664          |
|                                   |                                   |
|                                   | "ModifiedTime": 1522435664        |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+

Control devices between edge nodes
----------------------------------

In this scenario, the device driver for Motor2 (DeviceDriverMotor) will
use following HTTP connection to watch its own local LDR:

+-----------------------------------+-----------------------------------+
| Method                            | GET                               |
+===================================+===================================+
| URL                               | <http://localhost:8080/v1.0/Huawe |
|                                   | iProject1/edgecloud/edges/e2/ldrs |
|                                   | /expected/motor2?watch=true>      |
+-----------------------------------+-----------------------------------+
| BODY                              | N/A                               |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK                            |
|                                   |                                   |
|                                   | 404 NotFound                      |
|                                   |                                   |
|                                   | 403 Forbidden                     |
|                                   |                                   |
|                                   | BODY:                             |
|                                   |                                   |
|                                   | 1.  Value of Motor1 in JSON       |
|                                   |                                   |
|                                   | 2.  (?recursive=true) values of   |
|                                   |     all keys in Motor1 group in   |
|                                   |     JSON                          |
|                                   |                                   |
|                                   | 3.  "key not exist" if key is not |
|                                   |     defined                       |
|                                   |                                   |
|                                   | e.x:                              |
|                                   |                                   |
|                                   | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "key": "motor2",                  |
|                                   |                                   |
|                                   | "Value": "On",                    |
|                                   |                                   |
|                                   | "Revisions": 43,                  |
|                                   |                                   |
|                                   | "CreateTime": 1522435664          |
|                                   |                                   |
|                                   | "ModifiedTime": 1522435664        |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+

And the ControlMotorApp application will use following HTTP request to
update the remote LDRS:

+-----------------------------------+-----------------------------------+
| Method                            | PUT                               |
+===================================+===================================+
| URL                               | <http://localhost:8080/v1.0/Huawe |
|                                   | iProject1/edgecloud/edges/e2/ldrs |
|                                   | /expected/motor2/>                |
+-----------------------------------+-----------------------------------+
| BODY                              | \[                                |
|                                   |                                   |
|                                   | {                                 |
|                                   |                                   |
|                                   | "key": "motor2",                  |
|                                   |                                   |
|                                   | "Value": "On",                    |
|                                   |                                   |
|                                   | "Revisions": 43,                  |
|                                   |                                   |
|                                   | "CreateTime": 1522435664          |
|                                   |                                   |
|                                   | "ModifiedTime": 1522435664        |
|                                   |                                   |
|                                   | }                                 |
|                                   |                                   |
|                                   | \]                                |
+-----------------------------------+-----------------------------------+
| RESPONSE                          | 200 OK "Succeed"                  |
|                                   |                                   |
|                                   | 200 OK "Key already exist" --     |
|                                   | etcd error message                |
|                                   |                                   |
|                                   | 4xx ERROR HTTP error messages     |
+-----------------------------------+-----------------------------------+

The suggested resource hierarchy
================================

Central:

-   Root

    -   System

        -   Configuration (Sync from central to edge)

            -   Public

                -   Common

                    -   EdgeList

                        -   ProjectID
                            (*/v1.0/{project\_id}/edgeclouds/Edges/{EdgeName}/{ServiceName}*)

                            -   EdgeName  EdgeDef{Projectid="asdf",
                                EdgeName="edge1", ClusterID=123,
                                EdgeID=456}

                        -   ...

                -   Groups

                    -   

            -   Edges

                -   \<ProjectID\>

                    -   Edge1

                    -   Edge2...

                    -   \<EdgeName\>

                        -   Services

                            -   LDRS

                                -   Schemas

                                -   Device1

                                -   Device2

                                -   ...

                            -   ServiceBus

                            -   MetadataDB

                            -   FunctionEngine

                            -   Rules

                            -   SyncService

                                -   Sync\#1: Source(...) target(...)

                                -   ...

        -   Status (Sync From Edge to Central)

            -   \<ClusterID\>

                -   \<EdgeName\>

                    -   Services

                        -   LDRS

                        -   ServiceBus

                        -   MetadataDB

                        -   ...

        -   Runtime

            -   Edges

                -   Edge1

                -   Edge2

                -   ...

    -   User

        -   ....

Edge: Same as Central

*/v1.0/{project\_id}/edgecloud/Edges/{RemoteEdgeName}/{ServiceName}/*

![](media/image4.png){width="6.5in" height="1.6604166666666667in"}

![](media/image5.png){width="6.5in" height="1.4743055555555555in"}

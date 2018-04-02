package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"time"
	"bufio"
	//"raspberrypi3/devicedrivers/light"
	"wrtnodedriver"
)


type schema struct {
	Key          string `json:"Key,omitempty"`
	Revision     int64  `json:"Revision,omitempty"`
	CreateTime   int64  `json:"CreateTime,omitempty"`
	ModifiedTime int64  `json:"ModifiedTime,omitempty"`
	Value        interface{} `json:"Value,omitempty"`
}

func main() {
	var status uint64
	var lastStatus1 uint64
	var lastStatus2 uint64

	wrtnodedriver.InitDevice()

	// paralelly run target monitor
	go target()

	// logics for source monitor
	status = 0x0100
	lastStatus1 = 0
	lastStatus2 = 0

	target := "http://localhost:8080/v1.0/HuaweiProject1/edgecloud/edges/e3/ldrs/actual/?update=batch"

	for {
		var r int64
		var kvs []schema
		updated := 0
// for sensor 1
		err1 := wrtnodedriver.SetGPIO(0, "00 00 01 00", 1)
		if err1 != nil {
			continue
		}
		ret1 := wrtnodedriver.ReadGPIO(0)

		status = ret1 & 0x0100
		fmt.Printf("Status %08X, %08X\n", ret1, status)

		if status == 0 {
			r = 1
		} else {
			r = 0
		}

		if lastStatus1 != status {
			kv1 := schema{
				Key: "demoboard/coversensor1",
				Value: r,
			}
			kvs = append(kvs, kv1)
			lastStatus1 = status
			updated = 1
		}
// for sensor 2
		err2 := wrtnodedriver.SetGPIO(1, "00 00 01 00", 1)
		if err2 != nil {
			continue
		}
		ret2 := wrtnodedriver.ReadGPIO(1)

		status = ret2 & 0x0100
		fmt.Printf("Status %08X, %08X\n", ret2, status)

		if status == 0 {
			r = 1
		} else {
			r = 0
		}
		if lastStatus2 != status {
			kv2 := schema{
				Key: "demoboard/coversensor2",
				Value: r,
			}
			kvs = append(kvs, kv2)
			lastStatus2 = status
			updated = 1
		}

		if updated == 0 {
			time.Sleep(300 * time.Millisecond)
			continue
		}

// update all to target
		b, err := json.Marshal(kvs)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Body: %+v\n", kvs)
		PostCoverObj(b, target)
		time.Sleep(300 * time.Millisecond)
	}

}

func PostCoverObj(b []byte, url string) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}



type Contents struct {
	Event string
	Key     string
	Revision     int64
	CreateTime   int64
	ModifiedTime int64
	Value        interface{}
}

type Response struct {
	Reversion int64
	Content []Contents
}

func target()  {

		url := "http://localhost:8080/v1.0/HuaweiProject1/edgecloud/edges/e3/ldrs/expected/demoboard/motor?watch=true&recursive=true"
		fmt.Println(url)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Connection", "close")
		resp, _ := http.DefaultClient.Do(req)
		fmt.Println("start test")
		//buf := bytes.NewBuffer(make([]byte, 0, 10000))
		//go io.Copy(buf, resp.Body)
		reader := bufio.NewReader(resp.Body)
		reader.ReadBytes('}')
		for {
			line, _ := reader.ReadBytes('\n')
			w := Response{}
			err :=json.Unmarshal(line, &w)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, c := range w.Content {
				fmt.Printf("In target loop ...\n")
				status := c.Value.(float64)
				fmt.Printf("====> Status: %v\n", status)
				if status > 0 {
					//wrtnodedriver.SetGPIO(0, "00 10 00 00", 1)- alarm
					// motor:
					fmt.Printf("====> Calling SetGPIO on from target\n")
					wrtnodedriver.SetGPIO(0, "00 06 00 00", 1)
				} else {
					//wrtnodedriver.SetGPIO(0, "00 10 00 00", 0)- alarm
					// motor:
					fmt.Printf("====> Calling SetGPIO off from target\n")
					wrtnodedriver.SetGPIO(0, "00 06 00 00", 0)
				}

			}
			fmt.Printf("Got one line: %s\n", string(line))
		}

}

package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"DeviceCertification/raspberrypi3/drivers/gpio/devicedrivers/switchdriver"
	"time"
)

type schema struct {
	Key          string `json:"Key,omitempty"`
	Revision     int64  `json:"Revision,omitempty"`
	CreateTime   int64  `json:"CreateTime,omitempty"`
	ModifiedTime int64  `json:"ModifiedTime,omitempty"`
	Value        interface{} `json:"Value,omitempty"`
}

func main() {
	var status int64

	status = 7
	target := "http://localhost:8080/v1.0/HuaweiProject1/edgecloud/edges/e1/ldrs/actual/?update=batch"
	for {
		// Switch ReadStatus with readStatus for testing
		ret := switchdriver.ReadStatus()
		//ret := readStatus(status)
		if ret == status {
			continue
		}

		fmt.Println("Status ")
		status = ret
		fmt.Println(ret)
		var r int64
		if status == 0 {
			r = 1
		} else {
			r = 0
		}
		kv := schema{
			Key: "switch",
			Value: r,
		}
		var kvs []schema
		kvs = append(kvs, kv)
		b, err := json.Marshal(kvs)
		if err != nil {
			fmt.Println(err)
			return
		}
		PostCoverObj(b, target)
		time.Sleep(2000 * time.Millisecond)
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

// Soft switch supporting tests
func readStatus(status int64) int64 {
	var ret int64
	if status == 1 {
		ret = 0
	} else {
		ret = 1
	}
	return ret
}

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"runtime"
	"DeviceCertification/raspberrypi3/drivers/gpio/devicedrivers/lightdriver"
	"bytes"
)


type Content struct {
	Event string
	Key     string
	Revision     int64
	CreateTime   int64
	ModifiedTime int64
	Value        interface{}
}

type WatchResponse struct {
	Reversion int64
	Content []Content
}




func main()  {
	fmt.Printf("Start light mapper ARCH [%s]\n", runtime.GOARCH)

	url := "http://localhost:8080/v1.0/p1/edgecloud/edges/e1/ldrs/expected/light?watch=true&recursive=true"
	targetUrl:="http://localhost:8081/v1.0/HuaweiProject1/edgecloud/edges/e1/ldrs/actual/light"
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	for {
		// Get json state from switch
		w := WatchResponse{}
		err := json.NewDecoder(resp.Body).Decode(&w)
		if err != nil{
			fmt.Println(err)
			continue
		}
		fmt.Println(len(w.Content))
		// Create json state for light
		for _, c := range w.Content {
			status := c.Value.(float64)
			fmt.Println(status)
			if runtime.GOARCH == "amd64" {
				if status > 0 {
					fmt.Println("Light turned on")
					UpdateActual(c, targetUrl)
				} else {
					fmt.Println("Light turned off")
					UpdateActual(c, targetUrl)
				}
			} else {
				if status > 0 {
					lightdriver.TurnON()
					UpdateActual(c, targetUrl)
				} else {
					lightdriver.TurnOff()
					UpdateActual(c, targetUrl)
				}
			}

		}

	}

}

func UpdateActual(c Content,actualUrl string){
	actualValueUpdation:=Content{
		Key:c.Key,
		Value:c.Value,
	}
	body,_:=json.Marshal(actualValueUpdation)
	fmt.Println("Actual Value updated")
	PostObj(body,actualUrl)
}



func PostObj(b []byte, url string){

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

}
package main

import (
	"fmt"
	"net/http"
	"bufio"
	"encoding/json"
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

func main() {
	url := "http://localhost:8081/v1.0/HuaweiProject1/edgecloud/edges/e1/ldrs/expected/building1/floor1/?watch=true&recursive=true"
	targetUrl := "http://localhost:8081/v1.0/HuaweiProject1/edgecloud/edges/e1/ldrs/actual/"
	//sas++
	//url := "http://localhost:8081/v1.0/HuaweiProject1/edgecloud/edges/e1/ldrs/expected/light"
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "close")
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println("start light mapper")
	reader := bufio.NewReader(resp.Body)
	reader.ReadBytes('}')
	junk, _ := reader.ReadBytes('\n')
	if false {
		fmt.Println(junk)
	}
	for {
		line, _ := reader.ReadBytes('\n')
		w := WatchResponse{}
		err := json.Unmarshal(line, &w)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(len(w.Content))
		for i, c := range w.Content {
			if w.Content[i].Key == "light_1" {
					status := c.Value.(float64)
					fmt.Println(status)
					if status > 0 {
						lightdriver.TurnON(23)
						fmt.Println("Light1 turned on")
						UpdateActual(c, targetUrl)

					} else {
						lightdriver.TurnOff(23)
						fmt.Println("Light1 turned off")
						UpdateActual(c, targetUrl)
					}


			} else if w.Content[i].Key == "light_2" {
				status := c.Value.(float64)
					if status > 0 {
						lightdriver.TurnON(24)
						fmt.Println("Light2 turned on")
						UpdateActual(c, targetUrl)
					} else {
						lightdriver.TurnOff(24)
						fmt.Println("Light2 turned off")
						UpdateActual(c, targetUrl)
					}


			} else if w.Content[i].Key == "light_3" {
					status := c.Value.(float64)
					if status > 0 {
						lightdriver.TurnON(4)
						fmt.Println("Light3 turned on")
						UpdateActual(c, targetUrl)
					} else {
						lightdriver.TurnOff(4)
						fmt.Println("Light3 turned off")
						UpdateActual(c, targetUrl)
					}


			} else if w.Content[i].Key == "light_4" {
					status := c.Value.(float64)
					if status > 0 {
						lightdriver.TurnON(22)
						fmt.Println("Light4 turned on")
						UpdateActual(c, targetUrl)
					} else {
						lightdriver.TurnOff(22)
						fmt.Println("Light4 turned off")
						UpdateActual(c, targetUrl)
					}


			}
		}
	}


}
func UpdateActual(c Content,actualUrl string){
	actualValueUpdation := Content{
		Key: c.Key,
		Value: c.Value,
	}
	body, _ := json.Marshal(actualValueUpdation)
	fmt.Println("Actual Value updated")
	PostObj(body, actualUrl)

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

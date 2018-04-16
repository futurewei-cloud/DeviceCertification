package main

import (
	"fmt"
	"net/http"
	"bufio"
	"encoding/json"
	"DeviceCertification/raspberrypi3/drivers/gpio/devicedrivers/lightdriver"
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

	url := "http://localhost:8080/v1.0/HuaweiProject1/edgecloud/edges/e1/ldrs/expected/light?watch=true&recursive=true"
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "close")
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println("start light mapper")
	reader := bufio.NewReader(resp.Body)
	reader.ReadBytes('}')
	for {
		line, _ := reader.ReadBytes('\n')
		w := WatchResponse{}
		err :=json.Unmarshal(line, &w)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, c := range w.Content {
			status := c.Value.(float64)
			fmt.Println(status)
			if status > 0 {
				lightdriver.TurnON()
				fmt.Println("Light turned on")
			} else {
				lightdriver.TurnOff()
				fmt.Println("Light turned off")
			}

		}
		fmt.Println(string(line))
	}

}

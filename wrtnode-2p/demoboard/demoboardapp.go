package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strings"
	"log"
	//"bufio"
	"io/ioutil"
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

type Schema struct {
	Key          string `json:"Key,omitempty"`
	Revision     int64  `json:"Revision,omitempty"`
	CreateTime   int64  `json:"CreateTime,omitempty"`
	ModifiedTime int64  `json:"ModifiedTime,omitempty"`
	Value        interface{} `json:"Value,omitempty"`
}

func main() {

	source := "http://localhost:8080/v1.0/HuaweiProject1/edgecloud/edges/e3/ldrs/actual/demoboard/coversensor?watch=true&recursive=true"
	target := "http://localhost:8080/v1.0/HuaweiProject1/edgecloud/edges/e3/ldrs/expected/?update=batch"
	//req, _ := http.NewRequest("GET", "http://localhost:8080/Futurewei4/RainerCore/1.0.0/logicaldevices/watch/abc", nil)
	req, _ := http.NewRequest("GET", source, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	fmt.Println("start test")


	for {
		//reader := bufio.NewReader(resp.Body)
		//line, err := reader.ReadBytes('\n')
		//b, err := ioutil.ReadAll(resp.Body)
		//reader.Reset(resp.Body)
		//if err != nil {
		//	fmt.Println(err)
		//	continue
		//}
		w := WatchResponse{}
		err := json.NewDecoder(resp.Body).Decode(&w)
		//err = json.Unmarshal(b, &w)
		if err != nil{
			fmt.Println(err)
			continue
		}
		//log.Println(string(b))
		kvs := []Schema{}
		fmt.Printf(" ---> len(w.Content): %d\n", len(w.Content))
		fmt.Printf(" ===> w.Content: %+v\n", w.Content)
		for _, c := range w.Content {
			fmt.Printf(" ===> c.Value: %v\n", c.Value)
			kv := Schema{
				Key: "demoboard/motor1",
				Value: c.Value,
			}
			fmt.Println("test test")
			fmt.Println(c)
			kvs = append(kvs, kv)
		}
		if len(kvs) == 0 {
			continue
		}
		fmt.Printf("kvs: %+v\n", kvs)
		body, _ := json.Marshal(kvs)
		PostObj(body, target)
	}

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
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

type Message struct {
	Action, Key, Value string
}

func readJson(jsonStream string) {

	dec := json.NewDecoder(strings.NewReader(jsonStream))

	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for dec.More() {
		var m Schema
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v: %v\n", m.Key, m.Value)
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

}

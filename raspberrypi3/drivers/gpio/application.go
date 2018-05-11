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

	source := "http://localhost:8080/v1.0/p1/edgecloud/edges/e1/ldrs/actual/switch?watch=true&recursive=true"
	target := "http://localhost:8080/v1.0/p1/edgecloud/edges/e1/ldrs/expected/?update=batch"
	req, _ := http.NewRequest("GET", source, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	fmt.Println("start application")


	for {
		// Get json state from switch
		w := WatchResponse{}
		err := json.NewDecoder(resp.Body).Decode(&w)
		if err != nil{
			fmt.Println(err)
			continue
		}
		kvs := []Schema{}
		fmt.Println(len(w.Content))
		// Create json state for light
		for _, c := range w.Content {
			kv := Schema{
				Key: "light",
				Value: c.Value,
			}
			fmt.Println(c)
			kvs = append(kvs, kv)
		}
		if len(kvs) == 0 {
			continue
		}
		fmt.Println(kvs)
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

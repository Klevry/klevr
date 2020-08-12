package communicator

import (
	"bytes"
	_ "encoding/json"
	"github.com/NexClipper/logger"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var http_body_buffer string

func Put_http(url, data, api_key_string string) {
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		log.Printf("HTTP PUT Request error: ", err)
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")
	//    client := &http.Client{}
	//	res, err := client.Do(req)
	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
	} else {
		log.Println(err)
	}
}

func Put_Json_http(url string, data []byte, api string) []byte{
	var body []byte
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil{
		log.Printf("HTTP PUT Request error: ", err)
	}

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", api)
	req.Header.Add("X-ZONE-ID", "8282")
	req.Header.Add("X-API-KEY", "wnskslxptmxmrPwjd123")

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		http_body_buffer = string(body)
		//logger.Debugf("%v", http_body_buffer)
	} else{
		log.Println(err)
	}

	return body
}

func Get_http(uri, api_key_string string) string {
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")
	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		http_body_buffer = string(body)
	} else {
		log.Printf("Server connection error: ", err)
	}
	return http_body_buffer
}

func Get_Json_http(url string, api string) []byte{
	var body []byte

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", api)
	req.Header.Add("X-ZONE-ID", "8282")
	req.Header.Add("X-API-KEY", "wnskslxptmxmrPwjd123")

	res, err := http.DefaultClient.Do(req)
	if err == nil{
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		http_body_buffer = string(body)
	} else{
		logger.Error("Server connection error: ", err)
	}
	return body
}

func Delete_http(uri, api_key_string string) {
	req, _ := http.NewRequest("DELETE", uri, nil)
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
}

func Post_http(url, data, api_key_string string) {
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		log.Printf("HTTP POST Request error: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
}

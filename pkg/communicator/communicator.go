package communicator

import (
	"bytes"
	_ "encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NexClipper/logger"
)

var http_body_buffer string

func Put_http(url, data, api_key_string string) {
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
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
		logger.Error(err)
	}
}

func Put_Json_http(url string, data []byte, agent string, api string, zone string) []byte {
	var body []byte
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
	}

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agent)
	req.Header.Add("X-ZONE-ID", zone)
	req.Header.Add("X-API-KEY", api)

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		http_body_buffer = string(body)
	} else {
		logger.Error(err)
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
		logger.Errorf("Server connection error: %v", err)
	}
	return http_body_buffer
}

func Get_Json_http(url string, agent string, api string, zone string) []byte {
	var body []byte

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agent)
	req.Header.Add("X-ZONE-ID", zone)
	req.Header.Add("X-API-KEY", api)

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		http_body_buffer = string(body)
	} else {
		logger.Errorf("Server connection error: ", err)
	}
	return body
}

func Delete_http(uri, api_key_string string) {
	req, _ := http.NewRequest("DELETE", uri, nil)
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	defer res.Body.Close()
}

func Post_http(url, data, api_key_string string) {
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		logger.Errorf("HTTP POST Request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Error(err)
	}
	defer res.Body.Close()
}

func Post_Json_http(url string, data []byte, agent string, api string, zone string) []byte {
	var body []byte
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
	}

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agent)
	req.Header.Add("X-ZONE-ID", zone)
	req.Header.Add("X-API-KEY", api)

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		http_body_buffer = string(body)
	} else {
		logger.Error(err)
	}

	return body
}

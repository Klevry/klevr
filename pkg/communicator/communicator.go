package communicator

import (
	"bytes"
	_ "encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NexClipper/logger"
)

func Put_http(url, data, api_key_string string) error {
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
	} else {
		logger.Errorf("Server connection error: %v", err)
		return err
	}

	return nil
}

func Put_Json_http(url string, data []byte, agent string, api string, zone string) ([]byte, error) {
	var body []byte

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agent)
	req.Header.Add("X-ZONE-ID", zone)
	req.Header.Add("X-API-KEY", api)

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
	} else {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	return body, nil
}

func Get_http(uri, api_key_string string) ([]byte, error) {
	var body []byte

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		logger.Errorf("HTTP GET Request error: %v", err)
		return nil, err
	}

	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
	} else {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	return body, nil
}

func Get_Json_http(url string, agent string, api string, zone string) ([]byte, error) {
	var body []byte

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorf("HTTP GET Request error: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agent)
	req.Header.Add("X-ZONE-ID", zone)
	req.Header.Add("X-API-KEY", api)

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
	} else {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	return body, nil
}

func Delete_http(uri, api_key_string string) error {
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		logger.Errorf("HTTP DELETE Request error: %v", err)
		return err
	}

	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return err
	}
	defer res.Body.Close()

	return nil
}

func Post_http(url, data, api_key_string string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		logger.Errorf("HTTP POST Request error: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("nexcloud-auth-token", api_key_string)
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return err
	}
	defer res.Body.Close()

	return nil
}

func Post_Json_http(url string, data []byte, agent string, api string, zone string) ([]byte, error) {
	var body []byte

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("HTTP POST Request error: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agent)
	req.Header.Add("X-ZONE-ID", zone)
	req.Header.Add("X-API-KEY", api)

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
	} else {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	return body, nil
}

package communicator

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NexClipper/logger"
	"github.com/hashicorp/go-retryablehttp"
)

func Put_http(url, data, apiKey string) error {
	req, err := retryablehttp.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Add("nexcloud-auth-token", apiKey)
	req.Header.Add("cache-control", "no-cache")

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		logger.Errorf("HTTP Error, status code: %d", res.StatusCode)
		return fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
	}

	return nil
}

func Put_Json_http(url string, data []byte, agentKey, apiKey, zoneID string) ([]byte, error) {
	var body []byte

	req, err := retryablehttp.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agentKey)
	req.Header.Add("X-ZONE-ID", zoneID)
	req.Header.Add("X-API-KEY", apiKey)

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		logger.Errorf("HTTP Error, status code: %d", res.StatusCode)
		return nil, fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
	}

	body, _ = ioutil.ReadAll(res.Body)

	return body, nil
}

func Get_http(url, apiKey string) ([]byte, error) {
	var body []byte
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorf("HTTP GET Request error: %v", err)
		return nil, err
	}
	req.Header.Add("nexcloud-auth-token", apiKey)
	req.Header.Add("cache-control", "no-cache")

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		logger.Errorf("HTTP Error, status code: %d", res.StatusCode)
		return nil, fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
	}

	body, _ = ioutil.ReadAll(res.Body)

	return body, nil
}

func Get_Json_http(url, agentKey, apiKey, zoneID string) ([]byte, error) {
	var body []byte

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorf("HTTP GET Request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agentKey)
	req.Header.Add("X-ZONE-ID", zoneID)
	req.Header.Add("X-API-KEY", apiKey)

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		logger.Errorf("HTTP Error, status code: %d", res.StatusCode)
		return nil, fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
	}

	body, _ = ioutil.ReadAll(res.Body)

	return body, nil
}

func Delete_http(url, apiKey string) error {
	req, err := retryablehttp.NewRequest("DELETE", url, nil)
	if err != nil {
		logger.Errorf("HTTP DELETE Request error: %v", err)
		return err
	}
	req.Header.Add("nexcloud-auth-token", apiKey)
	req.Header.Add("cache-control", "no-cache")

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		logger.Errorf("HTTP Error, status code: %d", res.StatusCode)
		return fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
	}

	return nil
}

func Post_http(url, data, apiKey string) error {
	req, err := retryablehttp.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		logger.Errorf("HTTP POST Request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("nexcloud-auth-token", apiKey)
	req.Header.Add("cache-control", "no-cache")

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		logger.Errorf("HTTP Error, status code: %d", res.StatusCode)
		return fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
	}

	return nil
}

func Post_Json_http(url string, data []byte, agentKey, apiKey, zoneID string) ([]byte, error) {
	var body []byte
	req, err := retryablehttp.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("HTTP POST Request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", agentKey)
	req.Header.Add("X-ZONE-ID", zoneID)
	req.Header.Add("X-API-KEY", apiKey)

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Server connection error: %v", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		logger.Errorf("HTTP Error, status code: %d", res.StatusCode)
		return nil, fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
	}

	body, _ = ioutil.ReadAll(res.Body)

	return body, nil
}

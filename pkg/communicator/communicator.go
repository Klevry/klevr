package communicator

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/NexClipper/logger"
	"github.com/hashicorp/go-retryablehttp"
)

type Http struct {
	URL string

	APIKey   string
	AgentKey string
	ZoneID   string

	RetryCount int
	Timeout    int
}

// Timeout이 0이면 기본 타임아웃을 기본값(3초)로 한다.
func (h *Http) request(req *retryablehttp.Request) (*http.Response, error) {
	if h.Timeout == 0 {
		// polling을 5초 주기로 하고 있다. 그 시간안에 응답을 받지 못 하면 타임아웃으로 한다.
		h.Timeout = 3
	}
	client := retryablehttp.NewClient()
	client.RetryMax = h.RetryCount
	if h.Timeout > 0 {
		client.HTTPClient.Timeout = time.Duration(h.Timeout) * time.Second
	}
	client.Logger = nil
	// 디버깅하는 용도로 사용할 수 있도록 주석 처리해서 남겨 놓음
	/*client.RequestLogHook = func(l retryablehttp.Logger, req *http.Request, cnt int) {
		logger.Debugf("%s %s(%d)", req.Method, req.URL.String(), cnt)
	}*/

	return client.Do(req)
}

func (h *Http) newJsonRequest(method string, data []byte) (*retryablehttp.Request, error) {
	m := strings.ToUpper(method)
	req, err := retryablehttp.NewRequest(m, h.URL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "json/application; charset=utf-8")
	req.Header.Add("X-AGENT-KEY", h.AgentKey)
	req.Header.Add("X-ZONE-ID", h.ZoneID)
	req.Header.Add("X-API-KEY", h.APIKey)

	return req, nil
}

func (h *Http) newRequest(method string, data string) (*retryablehttp.Request, error) {
	var rawBody interface{}
	if data == "" {
		rawBody = nil
	} else {
		rawBody = strings.NewReader(string(data))
	}

	req, err := retryablehttp.NewRequest(strings.ToUpper(method), h.URL, rawBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("nexcloud-auth-token", h.APIKey)
	req.Header.Add("cache-control", "no-cache")

	return req, nil
}

func (h *Http) Put(data string) error {
	req, err := h.newRequest("PUT", data)
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
		return err
	}

	res, err := h.request(req)
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

func (h *Http) PutJson(data []byte) ([]byte, error) {
	var body []byte

	req, err := h.newJsonRequest("PUT", data)
	if err != nil {
		logger.Errorf("HTTP PUT Request error: %v", err)
		return nil, err
	}

	res, err := h.request(req)
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

func (h *Http) Get() ([]byte, error) {
	var body []byte

	req, err := h.newRequest("GET", "")
	if err != nil {
		logger.Errorf("HTTP GET Request error: %v", err)
		return nil, err
	}

	res, err := h.request(req)
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

func (h *Http) GetJson() ([]byte, error) {
	var body []byte

	req, err := h.newJsonRequest("GET", nil)
	if err != nil {
		logger.Errorf("HTTP GET Request error: %v", err)
		return nil, err
	}

	res, err := h.request(req)
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

func (h *Http) Delete() error {
	req, err := h.newRequest("DELETE", "")
	if err != nil {
		logger.Errorf("HTTP DELETE Request error: %v", err)
		return err
	}

	res, err := h.request(req)
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

func (h *Http) Post(data string) error {
	req, err := h.newRequest("POST", data)
	if err != nil {
		logger.Errorf("HTTP POST Request error: %v", err)
		return err
	}

	res, err := h.request(req)
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

func (h *Http) PostJson(data []byte) ([]byte, error) {
	var body []byte

	req, err := h.newJsonRequest("POST", data)
	if err != nil {
		logger.Errorf("HTTP POST Request error: %v", err)
		return nil, err
	}

	res, err := h.request(req)
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

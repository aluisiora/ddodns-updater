package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Ipify struct {
	baseURL string
}

func NewIpify() Ipify {
	return Ipify{
		baseURL: "https://api.ipify.org",
	}
}

func (c *Ipify) FetchPublicIPv4() (string, error) {
	resp, err := http.Get(c.baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type ipapiResponse struct {
	IP      string `json:"ip"`
	Version string `json:"version"`
}

type Ipapi struct {
	baseURL string
}

func NewIpapi() Ipapi {
	return Ipapi{
		baseURL: "https://ipapi.co",
	}
}

func (c *Ipapi) FetchPublicIPv4() (string, error) {
	resp, err := http.Get(c.baseURL + "/json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ipResponse ipapiResponse
	err = json.Unmarshal(body, &ipResponse)
	if err != nil {
		return "", err
	}

	if ipResponse.Version != "IPv4" {
		return "", fmt.Errorf("expected IPv4 but got %s with value %s", ipResponse.Version, ipResponse.Version)
	}

	return ipResponse.IP, nil
}

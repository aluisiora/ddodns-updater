package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DomainRecord struct {
	Data string `json:"data"`
	Type string `json:"type"`
	ID   int    `json:"id"`
}

type domainRecordsListResponse struct {
	DomainRecords []DomainRecord `json:"domain_records"`
	Meta          struct {
		Total int `json:"total"`
	} `json:"meta"`
}

type domainRecordResponse struct {
	DomainRecord DomainRecord `json:"domain_record"`
}

type DigitalOceanDNS struct {
	baseURL string
	token   string
}

func NewDigitalOceanDNS(token string) DigitalOceanDNS {
	return DigitalOceanDNS{
		baseURL: "https://api.digitalocean.com",
		token:   token,
	}
}

func (c *DigitalOceanDNS) FindDomainRecord(domain string, record string, recordType string) (DomainRecord, error) {
	path := fmt.Sprintf("/v2/domains/%s/records?type=%s&name=%s", domain, recordType, record)

	resp, err := c.request("GET", path, nil)
	if err != nil {
		return DomainRecord{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DomainRecord{}, err
	}

	var records domainRecordsListResponse

	err = json.Unmarshal(body, &records)
	if err != nil {
		return DomainRecord{}, err
	}

	if records.Meta.Total == 0 {
		return DomainRecord{}, fmt.Errorf("%s record was not found", record)
	}

	return records.DomainRecords[0], nil
}

func (c *DigitalOceanDNS) UpdateDomainRecord(domain string, record DomainRecord, value string) (DomainRecord, error) {
	path := fmt.Sprintf("/v2/domains/%s/records/%d", domain, record.ID)

	payload, err := json.Marshal(map[string]string{
		"type": record.Type,
		"data": value,
	})
	if err != nil {
		return DomainRecord{}, err
	}

	resp, err := c.request("PATCH", path, bytes.NewBuffer(payload))
	if err != nil {
		return DomainRecord{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DomainRecord{}, err
	}

	var jsonBody domainRecordResponse
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return DomainRecord{}, err
	}

	return jsonBody.DomainRecord, nil
}

func (c *DigitalOceanDNS) request(method string, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}

	return client.Do(req)
}

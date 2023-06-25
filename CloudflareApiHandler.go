package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (token *CloudFlareAuth) VerifyToken() (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
	if err != nil {
		return false, err
	}
	token.setHttpHeaders(req)
	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return false, fmt.Errorf("api response is not ok")
	}
	if err != nil {
		return false, err
	}
	body, err := io.ReadAll(res.Body)
	var verifyApiResponse VerifyApiResponse
	err = json.Unmarshal(body, &verifyApiResponse)
	defer res.Body.Close()
	if err != nil {
		return false, err
	}
	return verifyApiResponse.Success, nil
}

func (token *CloudFlareAuth) ListDnsRecords() ([]DnsZoneResult, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.cloudflare.com/client/v4/zones/"+token.ZoneId+"/dns_records",
		nil)
	token.setHttpHeaders(req)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("api response is not ok")
	}
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	var dnsResponse CloudflareListDnsResponse
	err = json.Unmarshal(body, &dnsResponse)
	if err != nil {
		return nil, err
	}

	return dnsResponse.Result, err
}

func (token *CloudFlareAuth) AddOrUpdateARecord(record *Record, createNewRecord bool) error {
	method := determineHttpMethod(createNewRecord)

	url := "https://api.cloudflare.com/client/v4/zones/" + token.ZoneId + "/dns_records"
	byteStream, err := json.Marshal(*record)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(byteStream))
	token.setHttpHeaders(req)
	res, err := httpClient.Do(req)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	var dnsResponse CloudflareCreateADnsRecordApiResponse
	err = json.Unmarshal(body, &dnsResponse)
	if res.StatusCode != 200 {
		return fmt.Errorf("api response is not ok")
	}
	return err
}

func (token *CloudFlareAuth) setHttpHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")
}

func determineHttpMethod(create bool) string {
	var method string
	switch create {
	case false:
		method = "POST"
		break
	case true:
		method = "PUT"
		break
	}
	return method
}

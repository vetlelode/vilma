package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// For some reason there are several APIs that do this for free with the same response
var apiUrls = []string{
	"https://api.seeip.org/jsonip?",
	"https://api.ipify.org?format=json",
}

func getIp() *ipResponse {
	var body []byte
	var err error
	for _, url := range apiUrls {
		err = tryMultipleApiPaths(url, &body)
		if err == nil {
			break
		}
	}
	data := ipResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	return &data
}

func tryMultipleApiPaths(url string, body *[]byte) error {
	fmt.Println("Calling api for public IP")
	resp, err := http.Get(url)
	if err != nil || resp.Status != "200 OK" {
		return fmt.Errorf("api response is not ok")
	}
	*body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error unmarshalling response %v}", err.Error())
	}
	return nil
}

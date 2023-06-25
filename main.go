package main

import (
	"flag"
	"fmt"
	"net/http"
)

// Command line flags
var cloudFlareAuth = flag.String("apikey", "", "Cloudflare API key")
var zoneId = flag.String("zoneid", "", "Cloudflare Zone ID")
var recordName = flag.String("name", "", "Record to set")
var recordContent = flag.String("content", "", "Content to set")
var recordType = flag.String("type", "A", "Type of record to set")

var httpClient = &http.Client{}

func main() {
	dnsRecord, err := parseCommandLineFlags()
	if err != nil {
		panic(err.Error())
	}
	data := getIp()
	dnsRecord.Content = data.IP
	fmt.Println(data.IP)
	cloudFlareSettings := CloudFlareAuth{Token: *cloudFlareAuth, ZoneId: *zoneId}
	success, err := cloudFlareSettings.VerifyToken()
	if err != nil {
		panic(err.Error())
	}
	if !success {
		panic("Token is invalid")
	}
	records, err := cloudFlareSettings.ListDnsRecords()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Existing dns records:")
	var recordExists bool = false
	for _, record := range records {
		fmt.Println(record.ID, record.Type, record.Name, record.Content, record.TTL)
		if record.Name == dnsRecord.Name && record.Type == dnsRecord.Type {
			recordExists = true
			if record.Content == dnsRecord.Content {
				fmt.Println("A record already exists with the same content. Quitting..")
				return
			}
		}
	}
	err = cloudFlareSettings.AddOrUpdateARecord(dnsRecord, recordExists)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func parseCommandLineFlags() (*Record, error) {
	flag.Parse()
	if *recordName == "" {
		return nil, fmt.Errorf("record name is required")
	}
	if *recordContent == "" {
		return nil, fmt.Errorf("record content is required")
	}
	if *cloudFlareAuth == "" {
		return nil, fmt.Errorf("cloudflare api key is required")
	}
	if *zoneId == "" {
		return nil, fmt.Errorf("cloudflare api key is required")
	}
	record := Record{Content: *recordContent, Name: *recordName, Type: *recordType, Proxied: true, TTL: 60}
	return &record, nil
}

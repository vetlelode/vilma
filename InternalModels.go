package main

type Record struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     uint   `json:"ttl"`
}

type ipResponse struct {
	IP string `json:"ip"`
}

type CloudFlareAuth struct {
	Token  string
	ZoneId string
}

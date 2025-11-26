package iplookup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ipData is the feedback interface which is returns
type ipData struct {
	IP            string  `json:"ip"`
	Success       bool    `json:"success"`
	Type          string  `json:"type"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsEu          bool    `json:"is_eu"`
	Postal        string  `json:"postal"`
	CallingCode   string  `json:"calling_code"`
	Capital       string  `json:"capital"`
	Borders       string  `json:"borders"`
	Flag          struct {
		Img          string `json:"img"`
		Emoji        string `json:"emoji"`
		EmojiUnicode string `json:"emoji_unicode"`
	} `json:"flag"`
	Connection struct {
		Asn    int    `json:"asn"`
		Org    string `json:"org"`
		Isp    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`
	Timezone struct {
		ID          string `json:"id"`
		Abbr        string `json:"abbr"`
		IsDst       bool   `json:"is_dst"`
		Offset      int    `json:"offset"`
		Utc         string `json:"utc"`
		CurrentTime string `json:"current_time"`
	} `json:"timezone"`
}

// lookup is used whenever the main lookup process fails
func lookup(ip string) (*Internet, error) {
	data, err := http.Get(fmt.Sprintf("http://ipwho.is/%s", url.QueryEscape(ip)))
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}

	var ipData *ipData = new(ipData)
	if err := json.Unmarshal(bytes, &ipData); err != nil {
		return nil, err
	}

	return &Internet{Query: ip, ASN: uint(ipData.Connection.Asn), Continent: ipData.ContinentCode, City: ipData.City, Country: ipData.Country, Organization: ipData.Connection.Org}, nil
}
package iplookup

import (
	"Nosviak4/source"
	"errors"
	"net"
	"net/url"
	"strconv"
)

// Internet is the returned interface
type Internet struct {
	ASN          uint   `swash:"asn"` //
	Query        string `swash:"query"`
	Continent    string `swash:"continent"` //
	City         string `swash:"city"`
	Country      string `swash:"country"`
	Organization string `swash:"org"`
	Longitude    string `swash:"longitude"`
	Latitude     string `swash:"latitude"`
}

// ErrBadIP is returned when the IP address is not a valid IP address
var ErrBadIP error = errors.New("bad ip address")

// Lookup will index within the databases for the ip provided
func Lookup(ip string) (*Internet, error) {
	address := net.ParseIP(ip)
	if address == nil {
		logic, err := url.ParseRequestURI(ip)
		if err != nil {
			logic = &url.URL{Host: ip}
		}

		resv, err := source.Resolver.LookupHost(logic.Host)
		if err != nil || len(resv) == 0 {
			return nil, ErrBadIP
		}

		address = resv[0]
	}

	mutex.Lock()
	defer mutex.Unlock()
	if lookup, ok := cache[address.String()]; ok && lookup != nil {
		if len(cache) >= 20 {
			cache = make(map[string]*Internet)
		}
		
		return lookup, nil
	}

	lookupASN, err := asn.ASN(address)
	if err != nil {
		return lookup(ip)
	}

	lookupCity, err := city.City(address)
	if err != nil {
		return lookup(ip)
	}

	data := &Internet{
		Query: address.String(), 
		ASN: lookupASN.AutonomousSystemNumber, 
		City: lookupCity.City.Names["en"], 
		Country: lookupCity.Country.Names["en"], 
		Continent: lookupCity.Continent.Code, 
		Organization: lookupASN.AutonomousSystemOrganization, 
		Longitude: strconv.FormatFloat(lookupCity.Location.Longitude, 'g', 5, 64), 
		Latitude: strconv.FormatFloat(lookupCity.Location.Latitude, 'g', 5, 64),
	}

	cache[data.Query] = data
	return data, nil
}

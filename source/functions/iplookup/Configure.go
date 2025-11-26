package iplookup

import (
	"Nosviak4/source"
	"path/filepath"
	"sync"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var (
	// default memory stores
	asn, city *geoip2.Reader

	// cache stores all the past IPLookup records
	cache map[string]*Internet = make(map[string]*Internet)

	// only allows one read or write inside cache.
	mutex sync.Mutex
)

// runs on startup
func init() {
	var err error
	asn, err = geoip2.Open(filepath.Join(source.ASSETS, source.COMMANDS, "bin", "asn.mmdb"))
	if err != nil {
		panic(err)
	}

	city, err = geoip2.Open(filepath.Join(source.ASSETS, source.COMMANDS, "bin", "city.mmdb"))
	if err != nil {
		panic(err)
	}
}
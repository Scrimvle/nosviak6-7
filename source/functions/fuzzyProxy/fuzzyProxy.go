package fuzzyproxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

// Proxy is the FuzzyProxy foundations
type Proxy struct {
	API string
	Key string

	mutex sync.RWMutex
	cache map[int]*IP
	client *http.Client
}

// IP is what exists in the cache
type IP struct{
	ID int

	net.Addr
	Expire time.Time

	Proxy *Proxy
}

// Response is a response from the proxy API
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// New creates a new FuzzyProxy instance
func New(API, Key string) *Proxy {
	return &Proxy{
		API: API,
		Key: Key,
		cache: make(map[int]*IP),
		client: &http.Client{
			Timeout: 1 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}
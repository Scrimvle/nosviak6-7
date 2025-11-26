package fuzzyproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

// fetch will query the Proxy for the real IP of the client
func (IP *IP) fetch() error {
	url := url.URL{
		Scheme: "https",
		Host: IP.Proxy.API,
		Path: fmt.Sprintf("/%d/addr", IP.ID),
	}

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", "Bearer " + IP.Proxy.Key)
	data, err := IP.Proxy.client.Do(request)
	if err != nil {
		return err
	}

	var response *Response
	if err := json.NewDecoder(data.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		return errors.New(response.Message)
	}

	addr, port, err := net.SplitHostPort(response.Message)
	if err != nil {
		return err
	}

	IP.ID, err = strconv.Atoi(port)
	if err != nil {
		return err
	}

	IP.Addr = &net.TCPAddr{
		IP: net.ParseIP(addr),
		Port: IP.ID,
	}
	
	return nil
}
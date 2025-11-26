package client

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//request will perform the request with path
func (C *Client) request(path string, schema string, method string, body io.Reader, headers map[string]string) (*http.Response, error) {
	system := url.URL{Scheme: schema, Host: C.Host, Path: path}
	urlPath := system.String() //gets the url aftermath properly

	//sets http custom params
	clientTransport := http.Client{
		Timeout: 3 * time.Second,
	}

	//prepares to perform the request properly
	request, err := http.NewRequest(method, urlPath, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("bytes", hex.EncodeToString(C.license.Public.Bytes())) //sets the curve bytes
	request.Header.Set("curve", strings.ReplaceAll(C.Key().ToECDSA().Curve.Params().Name, "-", "")) //sets the size

	//ranges through the headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	//performs the request properly
	return clientTransport.Do(request)
}
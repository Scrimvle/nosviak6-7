package client

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type entryPacket struct {
	Success			bool 		`json:"success"`
	Error			string		`json:"error"`
	Commit			string		`json:"commit"`
}

//EnableEntry will create the cookie with main server
func (C *Client) EnableEntry() (*entryPacket, error) {

	//creates our custom URL with out system properly
	target := url.URL{Scheme: C.Schema, Host: C.Host + ":443", Path: "/client/entry"}
	path := target.String() //returns the URL built string properly

	//custom params for http
	client := http.Client{ 
		Timeout: time.Duration(10 * time.Second), //timeout period set
	}

	//NewRequest will create the new request properly
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("fingerprint", hex.EncodeToString(C.Key().Public.Fingerprint()))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	//performs the http request
	result, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	//Reads all the body which was given
	output, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	var packet entryPacket //indexes the output properly
	if err := json.Unmarshal(output, &packet); err != nil {
		return nil, err
	}

	//tries to validate the schema from remote
	if !packet.Success || len(packet.Error) > 0 {
		return nil, errors.New(packet.Error)
	}

	return &packet, nil
}
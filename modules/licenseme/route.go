package client

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
)

//PerformRoute will perform the route and parse into interface
func (C *Client) PerformRoute(path string, i interface{}) error {

	//performs the request properly and safely
	result, err := C.request(path, C.Schema, "GET", nil, make(map[string]string))
	if err != nil {
		return err
	}

	//decides what encryption route we will take
	has := result.Header.Get("encryption")


	//ReadAll will properly read the body
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}


	//encryption enabled
	if has == "true" {
		//decodes the incoming string properly
		pureEncryption, err := hex.DecodeString(string(body))
		if err != nil {
			return err
		}

		//decrypts the body properly and safely
		newest, err := C.license.Decrypt(pureEncryption, C.license.Public.Curve)
		if err != nil {
			return err
		}

		//saves the body properly
		body = newest
	}

	//unmarshals the system properly
	return json.Unmarshal(body, i)
}
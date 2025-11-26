package client

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

type LicenseReturn struct {
	Success				bool				`json:"success"`
	Error				string				`json:"error"`
	Queried 			time.Time			`json:"queried"` 		//when the command was Queried
	Client				*User				`json:"client"`  		//information about the client
	//Created				time.Time			`json:"created"`  //TODO: add
	//Updated				time.Time			`json:"updated"`  //TODO: add
	HardwareBinded		bool				`json:"hardware_binded"`
	LicenseExpiry		time.Time			`json:"expiry"`
	App					string				`json:"app"`
	LicenseSHA			string				`json:"laced_license"`
	Alerts				[]string			`json:"alerts"`
}

type User struct {
	User string `json:"user"`
	Discord string `json:"discord"`
	Email string `json:"email"`
	CID string `json:"cid"`
}


//RunQuery will query the information from the license
func (C *Client) RunQuery(commit string, hardware string, structure string) (*LicenseReturn, error) {

	//headers will be forward
	headers := map[string]string{
		"commit":commit, "fingerprint":hex.EncodeToString(C.Key().Public.Fingerprint()),
		"hardware":hardware, "structure":structure, "license":hex.EncodeToString(C.encodedLicense), "app_name":C.TargetApp,
	}

	//runs the request with the remote system properly
	response, err := C.request("/client/QueryLicense", C.Schema, "GET", nil, headers)
	if err != nil {
		return nil, err
	}

	//ReadAll will read all the bytes properly
	BodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//encryption has enabled at this point properly
	if response.Header.Get("encryption") == "true" {
		//decodes the incoming string properly
		pureEncryption, err := hex.DecodeString(string(BodyBytes))
		if err != nil {
			return nil, err
		}

		//decrypts the body properly and safely
		newest, err := C.license.Decrypt(pureEncryption, C.license.Public.Curve)
		if err != nil {
			return nil, err
		}

		BodyBytes = newest
	}

	var gain LicenseReturn
	//parses the body properly
	//this will ensure its done without issues
	if err := json.Unmarshal(BodyBytes, &gain); err != nil {
		return nil, err
	}

	//success message
	if !gain.Success { //returns error
		return nil, errors.New(gain.Error)
	}

	return &gain, nil

}
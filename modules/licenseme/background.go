package client

import "Nosviak4/modules/licenseme/ecc"

type Client struct {
	Host    			string 		//stores the host string properly
	license 			ecc.Private  	//stores the current instances license key bytes
	Schema			string	//stores what schema the remote server is running with
	encodedLicense		[]byte	//stores what we will consider as the encoded license key properly
	TargetApp			string	//stores what app we are targetting with out license key properly and safely
}


//MakeClient starts the instance
func MakeClient(host string, schema string, targetApp string) *Client {
	return &Client{
		Host: host, Schema: schema, TargetApp: targetApp,
	}
}
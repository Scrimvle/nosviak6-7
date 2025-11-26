package client 

import (
	"Nosviak4/modules/licenseme/env"
)

//Hardware
func (C *Client) Hardware() (string, error) {
	return client_env.GrabHardware(C.TargetApp)
}
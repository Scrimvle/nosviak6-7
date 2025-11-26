package client

import (
	"crypto/sha256"
	"hash"
)

//bytes_string will convert incoming into strings properly
func (C *Client) bytes_string(incoming_bytes []byte, Hash ...hash.Hash) []byte {
	var hasher hash.Hash
	if len(Hash) > 0 {
		hasher = Hash[0]
	} else {
		hasher = sha256.New()
	}

	hasher.Write(incoming_bytes)
	return hasher.Sum(nil)
}

package client

import (
	"Nosviak4/modules/licenseme/ecc"
	"crypto/x509"
	"io/ioutil"
)

//GetFromFile will load the license from the file
func (C *Client) GetFromFile(file string) error {

	//ReadFiles the target file
	target, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	//converts the bytes into a compatible private key
	ecdsa_private, err := x509.ParseECPrivateKey(target)
	if err != nil {
		return err
	}

	//converts and saves into file properly
	C.license = *ecc.PrivateFromECDSA(ecdsa_private)
	C.encodedLicense = C.bytes_string(C.license.D.Bytes())
	//fmt.Println(md5.New(hex.EncodeToString(C.encodedLicense)))
	return nil
}

//SaveLicense will update the store properly
func (C *Client) SaveLicense(private *ecc.Private) {
	C.license = *private
	C.encodedLicense = C.bytes_string(C.license.D.Bytes())
}

//Key allows you to access the private
func (C *Client) Key() *ecc.Private {
	return &C.license
}
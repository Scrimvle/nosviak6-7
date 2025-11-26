package client_env

import machineid "Nosviak4/modules/licenseme/machine"

//GrabHardware will get the hashed machine id
func GrabHardware(App string) (string, error) {
	return machineid.ProtectedID(App)
}
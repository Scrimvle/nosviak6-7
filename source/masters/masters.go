package masters

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"encoding/base64"
	"errors"

	"Nosviak4/source/database"
	"Nosviak4/source/masters/terminal"

	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// Masters controls the main server for users
type Masters struct {
	config *ssh.ServerConfig
	binder *terminal.ServerConfig
	logger *gologr.FileLogger
}

// NewMaster will create the new config for the server.
func NewMaster() (*Masters, error) {
	var NewMasters = new(Masters)
	NewMasters.binder = new(terminal.ServerConfig)
	err := source.OPTIONS.MarshalFromPath(NewMasters.binder, "ssh")
	if err != nil {
		return nil, err
	}

	NewMasters.logger = source.LOGGER.NewFileLogger(filepath.Join(source.ASSETS, "logs", "connections.log"), int64(source.OPTIONS.Ints("branding", "recycle_log")))
	if NewMasters.logger.Err != nil {
		return nil, NewMasters.logger.Err
	}

	NewMasters.config = &ssh.ServerConfig{
		NoClientAuth:  NewMasters.binder.CustomAuth,
		ServerVersion: "SSH-2.0-OpenSSH",
		MaxAuthTries:  NewMasters.binder.MaxAuthAttempts,
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		 	user, err := database.DB.GetUser(conn.User())
		 	if err != nil || !user.IsPassword(password) {
		 		return nil, errors.New("bad password")
		 	}

			return nil, nil
		},

		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			user, err := database.DB.GetUserWithPublicKey([]byte(base64.StdEncoding.EncodeToString(key.Marshal())))
			if err != nil || user == nil {
				return nil, errors.New("bad ssh key")
			}

			return nil, nil
		},
	}

	keyBytes, err := os.ReadFile(NewMasters.binder.ServerKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	NewMasters.config.AddHostKey(privateKey)
	return NewMasters, nil
}

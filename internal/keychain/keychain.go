package keychain

import (
	"errors"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "MyPasswordManager"
	passwordPrefix = "PASSMANAGER:"
)

func SaveToKeychain(username, password string) error {
	return keyring.Set(keyringService, passwordPrefix+username, string(password))
}

func LoadFromKeychain(username string) (string, error) {
	password, err := keyring.Get(keyringService, passwordPrefix+username)

	if errors.Is(err, keyring.ErrNotFound) {
		return password, nil
	} else if err != nil {
		return "", err
	}

	return password, nil
}

package account

import (
	"errors"
	"fmt"

	"github.com/lemmyMwaura/pass/internal/vault"
)

// Create builds a new encrypted vault for username.
func Create(username, password, confirm string) (*vault.Vault, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("master password cannot be empty")
	}
	if password != confirm {
		return nil, errors.New("passwords don't match")
	}

	v, err := vault.Create(username, password)
	if err != nil {
		if errors.Is(err, vault.ErrAlreadyExists) {
			return nil, fmt.Errorf("account %q already exists — login instead", username)
		}
		return nil, err
	}
	return v, nil
}

// Login unlocks an existing vault with the master password.
func Login(username, password string) (*vault.Vault, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("master password cannot be empty")
	}

	v, err := vault.Unlock(username, password)
	if err != nil {
		switch {
		case errors.Is(err, vault.ErrNoVault):
			return nil, fmt.Errorf("no account found for %q — create one first", username)
		case errors.Is(err, vault.ErrWrongPassword):
			return nil, errors.New("wrong password")
		default:
			return nil, err
		}
	}
	return v, nil
}

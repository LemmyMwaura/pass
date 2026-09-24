package account

import (
	"errors"
	"fmt"

	"github.com/lemmyMwaura/pass/internal/reader"
	"github.com/lemmyMwaura/pass/internal/vault"
)

// CreateAccount prompts for credentials and creates an encrypted vault.
func CreateAccount() (*vault.Vault, error) {
	r := reader.NewInputReader()

	username, err := r.ReadUserInput("Enter your username: ")
	if err != nil {
		return nil, err
	}
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	password, err := r.ReadPassword("Enter your master password: ")
	if err != nil {
		return nil, err
	}
	confirm, err := r.ReadPassword("Confirm your master password: ")
	if err != nil {
		return nil, err
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

	fmt.Printf("Account %q created. Vault stored at ~/.pass/\n", username)
	return v, nil
}

// Login unlocks an existing vault with the master password.
func Login() (*vault.Vault, error) {
	r := reader.NewInputReader()

	username, err := r.ReadUserInput("Enter your username: ")
	if err != nil {
		return nil, err
	}
	password, err := r.ReadPassword("Enter your master password: ")
	if err != nil {
		return nil, err
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

	fmt.Printf("Unlocked vault for %q\n", username)
	return v, nil
}

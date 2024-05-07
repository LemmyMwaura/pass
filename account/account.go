package account

import (
	"errors"
	"fmt"

	"github.com/lemmyMwaura/pass/pkg/reader"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/bcrypt"
)

const (
	keyringService = "MyPasswordManager"
	passwordPrefix = "PASSMANAGER:"
)

type Account struct {
	username string
	password []byte
}

func NewAccount(uName string, pass []byte) *Account {
	return &Account{
		username: uName,
		password: pass,
	}
}

// SaveToKeychain saves the account password to the system's keychain.
func (acc *Account) SaveToKeychain() error {
	return keyring.Set(keyringService, passwordPrefix+acc.username, string(acc.password))
}

// LoadFromKeychain loads the account password from the system's keychain.
func (acc *Account) LoadFromKeychain() (string, error) {
	password, err := keyring.Get(keyringService, passwordPrefix+acc.username)

	if err == keyring.ErrNotFound {
		return password, nil
	} else if err != nil {
		return "", err
	}

	return password, nil
}

func CreateAccount() {
	r := reader.NewInputReader()

	username, _ := r.ReadUserInput("Enter your username:")
	mpassword, _ := r.ReadUserInput("Enter your MainPassword:")
	cpassword, _ := r.ReadUserInput("Confirm your MainPassword:")

	if mpassword != cpassword {
		fmt.Println("Passwords don't match.")
		CreateAccount()
		return
	}

	hashedPassword, err := hashPassword(mpassword)

	if err != nil {
		fmt.Printf("something went wrong: %s\n", err)
		return
	}

	account := NewAccount(username, hashedPassword)

	password, err := account.LoadFromKeychain()
	if err != nil {
		fmt.Printf("something went wrong: %s\n", err)
		return
	}

	if password == "" {
		account.SaveToKeychain()
	} else {
		fmt.Println("Account already exists, ...Login in instead")
		Login()
	}
}

func Login() {
	r := reader.NewInputReader()

	username, err1 := r.ReadUserInput("Enter your username:")
	password, err2 := r.ReadUserInput("Enter your password:")

	err := errors.Join(err1, err2)

	if err != nil {
		fmt.Printf("something went wrong: %s\n", err)
		return
	}

	fmt.Println(username, password)
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

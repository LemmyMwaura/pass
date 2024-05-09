package account

import (
	"errors"
	"fmt"

	"github.com/lemmyMwaura/pass/internal/keychain"
	"github.com/lemmyMwaura/pass/internal/reader"
	"golang.org/x/crypto/bcrypt"
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

	password, err := keychain.LoadFromKeychain(account.username)
	if err != nil {
		fmt.Printf("something went wrong: %s\n", err)
		return
	}

	if password == "" {
		keychain.SaveToKeychain(account.username, string(account.password))
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

	hashedPassword, err := keychain.LoadFromKeychain(password)
	if err != nil {
		fmt.Printf("something went wrong: %s\n", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println("Wrong password")
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

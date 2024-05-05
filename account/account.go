package account

import (
	"fmt"

	"github.com/lemmyMwaura/pass/pkg/reader"
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
	password, _ := r.ReadUserInput("Enter your password:")
	cpassword, _ := r.ReadUserInput("Confirm your password:")

	if password != cpassword {
		fmt.Println("Passwords don't match.")
		CreateAccount()
		return
	}

	hashedPassword, err := hashPassword(password)

	if err != nil {
		fmt.Printf("something went wrong: %s\n", err)
		return
	}

	account := NewAccount(username, hashedPassword)
	fmt.Println(account)
}

func Login() {
	r := reader.NewInputReader()

	username, _ := r.ReadUserInput("Enter your username:")
	fmt.Println(username)

	password, _ := r.ReadUserInput("Enter your password:")
	cpassword, _ := r.ReadUserInput("Confirm your password:")

	if password != cpassword {
		fmt.Println("Passwords don't match.")
		CreateAccount()
	}
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

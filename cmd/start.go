package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lemmyMwaura/pass/internal/account"
	"github.com/lemmyMwaura/pass/internal/reader"
	usersession "github.com/lemmyMwaura/pass/internal/user-session"
	"github.com/lemmyMwaura/pass/internal/vault"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:    "start",
	Short:  "Start the interactive password manager",
	Run:    readOption,
	PreRun: preCheckOption,
}

var (
	SelectedOption string
	options        = []string{"create Account", "login", "exit"}
)

func init() {
	startCmd.Flags().StringVarP(&SelectedOption, "option", "o", "", "Select an option: 1=create, 2=login, 3=exit")
	rootCmd.AddCommand(startCmd)
}

func preCheckOption(cmd *cobra.Command, args []string) {
	if validateSelectedOption(SelectedOption) {
		num, err := strconv.Atoi(SelectedOption)
		if err != nil {
			log.Fatal("Invalid input, expected an integer ", err)
		}
		SelectedOption = options[num-1]
		return
	}

	if SelectedOption != "" {
		fmt.Printf("Option %q unavailable, kindly re-select:\n\n", SelectedOption)
	}
	promptUser()
}

func promptUser() {
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	r := reader.NewInputReader()
	text, err := r.ReadUserInput("Enter the number of your choice: ")
	if err != nil {
		log.Fatal(err)
	}

	index, err := strconv.Atoi(text)
	if err != nil {
		log.Fatal("Invalid input, expected an integer: ", err)
	}
	if index < 1 || index > len(options) {
		log.Fatalf("Invalid option %d", index)
	}

	SelectedOption = options[index-1]
}

func readOption(cmd *cobra.Command, args []string) {
	var (
		v   *vault.Vault
		err error
	)

	switch SelectedOption {
	case "login":
		fmt.Println("\nLogin selected")
		v, err = account.Login()
	case "create Account":
		fmt.Println("\nCreate account selected")
		v, err = account.CreateAccount()
	case "exit":
		fmt.Println("\nExiting.")
		os.Exit(0)
	default:
		fmt.Println("Invalid option selected")
		return
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	session := usersession.NewUserSession(v)
	defer session.Close()
	runVaultMenu(session)
}

func validateSelectedOption(selectedOptn string) bool {
	switch selectedOptn {
	case "1", "2", "3":
		return true
	default:
		return false
	}
}

func runVaultMenu(session *usersession.UserSession) {
	r := reader.NewInputReader()
	menu := []string{
		"add",
		"list",
		"get",
		"delete",
		"generate",
		"exit",
	}

	for {
		fmt.Println("\n--- Vault ---")
		for i, item := range menu {
			fmt.Printf("%d. %s\n", i+1, item)
		}

		choice, err := r.ReadUserInput("Choose an option: ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		n, err := strconv.Atoi(choice)
		if err != nil || n < 1 || n > len(menu) {
			fmt.Println("Invalid option")
			continue
		}

		switch menu[n-1] {
		case "add":
			handleAdd(session.Vault, r)
		case "list":
			handleList(session.Vault)
		case "get":
			handleGet(session.Vault, r)
		case "delete":
			handleDelete(session.Vault, r)
		case "generate":
			handleGenerate(r)
		case "exit":
			fmt.Println("Locked vault. Goodbye.")
			return
		}
	}
}

func handleAdd(v *vault.Vault, r *reader.InputReader) {
	service, err := r.ReadUserInput("Service: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	username, err := r.ReadUserInput("Username: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	password, err := r.ReadPassword("Password (leave empty to generate): ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if password == "" {
		password = randomPassword(16, true, true)
		fmt.Printf("Generated password: %s\n", password)
	}
	notes, err := r.ReadUserInput("Notes (optional): ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	entry, err := v.Add(service, username, password, notes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Saved entry for %q\n", entry.Service)
}

func handleList(v *vault.Vault) {
	entries := v.List()
	if len(entries) == 0 {
		fmt.Println("Vault is empty.")
		return
	}
	fmt.Println("Stored services:")
	for _, e := range entries {
		fmt.Printf("  - %s (%s)\n", e.Service, e.Username)
	}
}

func handleGet(v *vault.Vault, r *reader.InputReader) {
	service, err := r.ReadUserInput("Service: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	entry, err := v.Get(service)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Service:  %s\n", entry.Service)
	fmt.Printf("Username: %s\n", entry.Username)
	fmt.Printf("Password: %s\n", entry.Password)
	if entry.Notes != "" {
		fmt.Printf("Notes:    %s\n", entry.Notes)
	}
}

func handleDelete(v *vault.Vault, r *reader.InputReader) {
	service, err := r.ReadUserInput("Service to delete: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	confirm, err := r.ReadUserInput(fmt.Sprintf("Type %q to confirm: ", service))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if confirm != service {
		fmt.Println("Delete cancelled.")
		return
	}
	if err := v.Delete(service); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Deleted %q\n", service)
}

func handleGenerate(r *reader.InputReader) {
	lengthStr, err := r.ReadUserInput("Length (default 16): ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	length := 16
	if lengthStr != "" {
		n, err := strconv.Atoi(lengthStr)
		if err != nil || n < 4 {
			fmt.Println("Invalid length; using 16")
		} else {
			length = n
		}
	}
	fmt.Println(randomPassword(length, true, true))
}

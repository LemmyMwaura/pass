package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lemmyMwaura/cli-test/account"
	"github.com/lemmyMwaura/cli-test/pkg/reader"
	"github.com/lemmyMwaura/cli-test/pkg/validators"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:    "start",
	Short:  "Start the application",
	Run:    readOption,
	PreRun: preCheckOption,
}

var (
	SelectedOption string
	options        = []string{"create Account", "login"}
)

func init() {
	startCmd.Flags().StringVarP(&SelectedOption, "option", "o", "", "Select an option: login or createAccount")
	rootCmd.AddCommand(startCmd)
}

func preCheckOption(cmd *cobra.Command, args []string) {
	isSelected := validators.ValidateSelectedOption(SelectedOption)

	if isSelected {
		option, _ := cmd.Flags().GetString("option")
		num, err := strconv.Atoi(option)

		if err != nil {
			log.Fatal("Invalid input, expected an integer ", err)
		} else if num < 1 || num > len(options) {
			log.Fatalf("Invalid option %d", num)
		}

		SelectedOption = options[num-1]
	} else {
		fmt.Printf("Option %s unavailable, kindly re-select: \n\n", SelectedOption)
		promptUser()
	}
}

func promptUser() {
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	r := reader.NewInputReader()
	text, err := r.ReadUserInput("Enter the number of your choice: ")

	if err != nil {
		return
	}

	index, err := strconv.Atoi(text)

	if err != nil {
		log.Fatal("Invalid input, expected an integer", err)
	} else if index < 1 || index > len(options) {
		log.Fatalf("Invalid option %d", index)
	}

	SelectedOption = options[index-1]
}

func readOption(cmd *cobra.Command, args []string) {
	switch SelectedOption {
	case "login":
		fmt.Println("\nLogin selected")
	case "create Account":
		fmt.Println("\nCreate account selected")
		account.CreateAccount()
	default:
		fmt.Println("Invalid option selected")
	}
}

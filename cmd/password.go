package cmd

import (
	"math/rand"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"genPass"},
	Short:   "generate passwords",
	Long: `Generate random passwords with customizable options:
	For example: 
			
	password gen: -l 12 -d -s
	`,
	Run: generatePassword,
}

func init() {
	generateCmd.Flags().IntP("length", "l", 8, "Length of the generated password (default 8)")
	generateCmd.Flags().BoolP("digits", "d", false, "Include Digits in the generated passwords")
	generateCmd.Flags().BoolP("special-chars", "s", false, "Include Special Chars in the generated passwwords")
	rootCmd.AddCommand(generateCmd)
}

func generatePassword(cmd *cobra.Command, args []string) {
	length, _ := cmd.Flags().GetInt("length")
	isDigits, _ := cmd.Flags().GetBool("digits")
	isSpecialChars, _ := cmd.Flags().GetBool("special-chars")

	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if isDigits {
		charset += "0123456789"
	}

	if isSpecialChars {
		charset += "!@#$%^&*()-_=+{}[]|;:<>,.?/~"
	}

	password := make([]byte, length)

	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	println(string(password))
}

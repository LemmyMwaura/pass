package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"genPass"},
	Short:   "Generate a random password",
	Long: `Generate random passwords with customizable options.

Example:
  pass generate -l 16 -d -s
`,
	Run: generatePassword,
}

func init() {
	generateCmd.Flags().IntP("length", "l", 16, "Length of the generated password")
	generateCmd.Flags().BoolP("digits", "d", true, "Include digits in the generated password")
	generateCmd.Flags().BoolP("special-chars", "s", true, "Include special characters")
	rootCmd.AddCommand(generateCmd)
}

func generatePassword(cmd *cobra.Command, args []string) {
	length, _ := cmd.Flags().GetInt("length")
	isDigits, _ := cmd.Flags().GetBool("digits")
	isSpecialChars, _ := cmd.Flags().GetBool("special-chars")

	fmt.Println(randomPassword(length, isDigits, isSpecialChars))
}

func randomPassword(length int, digits, special bool) string {
	if length < 1 {
		length = 16
	}

	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if digits {
		charset += "0123456789"
	}
	if special {
		charset += "!@#$%^&*()-_=+{}[]|;:<>,.?/~"
	}

	password := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range password {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		password[i] = charset[n.Int64()]
	}
	return string(password)
}

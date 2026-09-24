package cmd

import (
	"fmt"

	pwgen "github.com/lemmyMwaura/pass/internal/password"
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

	fmt.Println(pwgen.Generate(length, isDigits, isSpecialChars))
}

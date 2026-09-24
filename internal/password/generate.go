package password

import (
	"crypto/rand"
	"math/big"
)

// Generate returns a cryptographically random password.
func Generate(length int, digits, special bool) string {
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

	out := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range out {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		out[i] = charset[n.Int64()]
	}
	return string(out)
}

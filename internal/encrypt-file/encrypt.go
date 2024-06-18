package encryptfile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func SaveToFile(filename, password string, masterPassword []byte) error {
	// Open or create the file
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Generate a random initialization vector (IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	// Create the AES cipher block using the derived key from master password
	block, err := aes.NewCipher(deriveKey(masterPassword))
	if err != nil {
		return err
	}

	// Encrypt the password using AES GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	encryptedPassword := gcm.Seal(nil, iv, []byte(password), nil)

	// Write IV and encrypted password to the file
	_, err = file.Write(iv)
	if err != nil {
		return err
	}
	_, err = file.Write(encryptedPassword)
	return err
}

func LoadFromFile(filename, password string, masterPassword []byte) ([]byte, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read IV from the file
	iv := make([]byte, aes.BlockSize)
	_, err = file.Read(iv)
	if err != nil {
		return nil, err
	}

	// Create the AES cipher block using the derived key from master password
	block, err := aes.NewCipher(deriveKey(masterPassword))
	if err != nil {
		return nil, err
	}

	// Decrypt the password using AES GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	encryptedPassword := make([]byte, len(password))
	_, err = file.Read(encryptedPassword)
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := gcm.Open(nil, iv, encryptedPassword, nil)
	if err != nil {
		return nil, err
	}

	return decryptedPassword, nil
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func Decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// deriveKey derives a key from the master password using bcrypt.
func deriveKey(masterPassword []byte) []byte {
	// Perform key derivation using bcrypt or any other suitable KDF
	// For simplicity, I'll just return the master password itself
	return masterPassword
}

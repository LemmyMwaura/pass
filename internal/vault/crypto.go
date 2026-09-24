package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	magic          = "PASS"
	fileVersion    = byte(1)
	saltSize       = 16
	nonceSize      = 12
	argonTime      = 3
	argonMemory    = 64 * 1024 // 64 MiB
	argonThreads   = 4
	argonKeyLength = 32
)

var (
	ErrWrongPassword = errors.New("wrong master password")
	ErrCorruptVault  = errors.New("corrupt vault file")
)

func deriveKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, argonTime, argonMemory, argonThreads, argonKeyLength)
}

func encrypt(plaintext, key []byte) (nonce, ciphertext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return nonce, ciphertext, nil
}

func decrypt(nonce, ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrWrongPassword
	}

	return plaintext, nil
}

// seal wraps plaintext into the on-disk vault format.
// Layout: magic(4) | version(1) | salt(16) | ciphertextLen(4 LE) | nonce(12) | ciphertext
func seal(plaintext []byte, masterPassword string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := deriveKey([]byte(masterPassword), salt)
	nonce, ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 0, 4+1+saltSize+4+nonceSize+len(ciphertext))
	out = append(out, magic...)
	out = append(out, fileVersion)
	out = append(out, salt...)

	lenBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuf, uint32(len(ciphertext)))
	out = append(out, lenBuf...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return out, nil
}

func unseal(data []byte, masterPassword string) ([]byte, error) {
	minSize := 4 + 1 + saltSize + 4 + nonceSize
	if len(data) < minSize {
		return nil, ErrCorruptVault
	}

	if string(data[:4]) != magic {
		return nil, ErrCorruptVault
	}

	if data[4] != fileVersion {
		return nil, fmt.Errorf("%w: unsupported version %d", ErrCorruptVault, data[4])
	}

	offset := 5
	salt := data[offset : offset+saltSize]
	offset += saltSize

	cipherLen := binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	if offset+nonceSize+int(cipherLen) != len(data) {
		return nil, ErrCorruptVault
	}

	nonce := data[offset : offset+nonceSize]
	offset += nonceSize
	ciphertext := data[offset:]

	key := deriveKey([]byte(masterPassword), salt)
	return decrypt(nonce, ciphertext, key)
}

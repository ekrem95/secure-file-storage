package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AesDecrypt returns original data before encryption
func AesDecrypt(ciphertext *[]byte, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(*ciphertext) < aes.BlockSize {
		return errors.New("Text is too short")
	}

	// Get the 16 byte IV
	iv := (*ciphertext)[:aes.BlockSize]

	// Remove the IV from the ciphertext
	*ciphertext = (*ciphertext)[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(*ciphertext, *ciphertext)

	return nil
}

// AesEncrypt returns encrypted string
func AesEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Empty array of 16 + plaintext length
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

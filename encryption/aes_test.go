package encryption

import (
	"fmt"
	"testing"
)

var (
	plaintext  = []byte("plaintext")
	validKey   = []byte("aes aes aes aes ")
	invalidKey = []byte("aes aes aes aes")
)

func TestAesDecrypt(t *testing.T) {
	c, err := AesEncrypt(plaintext, validKey)
	if err != nil {
		t.Error(err)
	}

	if err = AesDecrypt(&c, validKey); err != nil {
		t.Error(err)
	}

	if string(c) != string(plaintext) {
		fmt.Println(string(c))
		fmt.Println(string(plaintext))
		t.Error("decrypted text must be equal to plaintext")
	}
}

func TestAesEncrypt(t *testing.T) {
	if _, err := AesEncrypt(plaintext, validKey); err != nil {
		t.Error(err)
	}
}

func TestInvalidKeySizeForAES(t *testing.T) {
	if _, err := AesEncrypt(plaintext, invalidKey); err.Error() != "crypto/aes: invalid key size 15" {
		t.Error(err)
	}
}

func TestAesEncryptAndDesDecrypt(t *testing.T) {
	c, err := AesEncrypt(plaintext, validKey)
	if err != nil {
		t.Error(err)
	}

	// there must be an error because => len(*src)%bs != 0
	if err = DesDecrypt(&c, validKey2); err.Error() != "crypto/cipher: input not full blocks" {
		t.Error(err)
	}

	if string(c) == string(plaintext) {
		fmt.Println(string(c))
		fmt.Println(string(plaintext))
		t.Error("decrypted text must not be equal to plaintext")
	}
}

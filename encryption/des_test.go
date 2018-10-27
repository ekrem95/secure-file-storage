package encryption

import (
	"fmt"
	"testing"
)

var (
	validKey2   = []byte("des1key6")
	invalidKey2 = []byte("des1key678")
)

func TestDesDecrypt(t *testing.T) {
	c, err := DesEncrypt(plaintext, validKey2)
	if err != nil {
		t.Error(err)
	}

	if err = DesDecrypt(&c, validKey2); err != nil {
		t.Error(err)
	}

	if string(c) != string(plaintext) {
		fmt.Println(string(c))
		fmt.Println(string(plaintext))
		t.Error("decrypted text must be equal to plaintext")
	}
}

func TestDesEncrypt(t *testing.T) {
	if _, err := DesEncrypt(plaintext, validKey2); err != nil {
		t.Error(err)
	}
}

func TestInvalidKey2SizeForDES(t *testing.T) {
	if _, err := DesEncrypt(plaintext, invalidKey2); err.Error() != "crypto/des: invalid key size 10" {
		t.Error(err)
	}
}

func TestDesEncryptAndAesDecrypt(t *testing.T) {
	c, err := DesEncrypt(plaintext, validKey2)
	if err != nil {
		t.Error(err)
	}

	if err = AesDecrypt(&c, validKey); err != nil {
		t.Error(err)
	}

	if string(c) == string(plaintext) {
		fmt.Println(string(c))
		fmt.Println(string(plaintext))
		t.Error("decrypted text must not be equal to plaintext")
	}
}

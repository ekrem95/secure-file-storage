package encryption

import (
	"bytes"
	"crypto/des"
	"errors"
)

// DesDecrypt returns original data before encryption
func DesDecrypt(src *[]byte, key []byte) error {
	block, err := des.NewCipher(key)
	if err != nil {
		return err
	}
	out := make([]byte, len(*src))
	dst := out
	bs := block.BlockSize()
	if len(*src)%bs != 0 {
		return errors.New("crypto/cipher: input not full blocks")
	}
	for len(*src) > 0 {
		block.Decrypt(dst, (*src)[:bs])
		*src = (*src)[bs:]
		dst = dst[bs:]
	}

	*src = zeroUnPadding(out)

	return nil
}

// DesEncrypt returns encrypted string
func DesEncrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	src = zeroPadding(src, bs)

	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

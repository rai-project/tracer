package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	key := []byte("mysecret")
	plaintext := []byte("sometext")

	secrettext, err := Encrypt(key, plaintext)
	assert.NoError(t, err)
	assert.NotEmpty(t, secrettext)

	// pp.Println(base64.StdEncoding.EncodeToString(secrettext))

	pt, err := Decrypt(key, secrettext)
	assert.NoError(t, err)
	assert.NotEmpty(t, pt)

	assert.Equal(t, plaintext, pt)
}

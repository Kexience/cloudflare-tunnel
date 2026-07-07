package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	plaintext := "cf_api_token_1234567890abcdef"
	encrypted, err := Encrypt(plaintext, key)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDeterministic(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	plaintext := "same_input"
	enc1, err := Encrypt(plaintext, key)
	require.NoError(t, err)
	enc2, err := Encrypt(plaintext, key)
	require.NoError(t, err)

	assert.NotEqual(t, enc1, enc2, "AES-GCM should produce different ciphertexts for same input")
}

func TestDecryptWrongKey(t *testing.T) {
	key1, err := GenerateKey()
	require.NoError(t, err)
	key2, err := GenerateKey()
	require.NoError(t, err)

	encrypted, err := Encrypt("secret", key1)
	require.NoError(t, err)

	_, err = Decrypt(encrypted, key2)
	assert.ErrorIs(t, err, ErrDecryptFailed)
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	_, err = Decrypt("not-valid-base64!!!", key)
	assert.Error(t, err)

	_, err = Decrypt("dGVzdA==", key) // "test" in base64, too short for GCM
	assert.ErrorIs(t, err, ErrDecryptFailed)
}

func TestInvalidKeyLength(t *testing.T) {
	shortKey := []byte("short")

	_, err := Encrypt("test", shortKey)
	assert.ErrorIs(t, err, ErrInvalidKeyLength)

	_, err = Decrypt("test", shortKey)
	assert.ErrorIs(t, err, ErrInvalidKeyLength)
}

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)
	assert.Len(t, key, KeyLength)

	keyHex := hex.EncodeToString(key)
	assert.Len(t, keyHex, KeyLength*2)

	key2, err := GenerateKey()
	require.NoError(t, err)
	assert.NotEqual(t, key, key2)
}

func TestEncryptEmptyString(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	encrypted, err := Encrypt("", key)
	require.NoError(t, err)

	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, "", decrypted)
}

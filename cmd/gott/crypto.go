package main

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"math/rand"
	"time"
)

// ----------------------------- CYPHER STUFF ----------------------------------
var encKey *[32]byte

func init() {
	// Seed Random Fung
	rand.Seed(time.Now().UTC().UnixNano())
	// Generate Encryption Key
	// MAYB: Change key every X hours, keep last key till all sessions are timed out
	encKey = newEncryptionKey()
	// Mongo DB setup
}

// Code copied from
// https://github.com/gtank/cryptopasta/blob/master/encrypt.go

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func newEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(cryptorand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func encrypt(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(cryptorand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	// Encode to hex so we avoid err:
	// net/http: invalid byte 'ó' in Cookie.Value; dropping invalid bytes
	ciphertext = gcm.Seal(nonce, nonce, plaintext, nil)
	hexCiphertext := make([]byte, hex.EncodedLen(len(ciphertext)))
	hex.Encode(hexCiphertext, ciphertext)

	return hexCiphertext, nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func decrypt(hexCiphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, hex.DecodedLen(len(hexCiphertext)))
	hex.Decode(ciphertext, hexCiphertext)

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

// ----------------------------- CRYPTO HELPER FUNCS ---------------------------

// generateRandomString generates a pseudo random string from a given set of
// chars s
func generateRandomString(n int, s string) string {
	str := make([]byte, n)
	l := len(s)
	// Create random Participant ID
	for i := range str {
		str[i] = s[rand.Intn(l)]
	}
	return string(str)
}

// generatePassword generates a pseudo random password
func generatePassword() (pw string) {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!?_-&#@&()"
	pw = generateRandomString(10, charSet)
	return
}

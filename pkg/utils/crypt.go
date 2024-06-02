package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/daison12006013/web-golang-101/pkg/env"
	"github.com/forgoer/openssl"
	"golang.org/x/crypto/bcrypt"
)

func HashStr(v string) string {
	hasher := sha256.New()
	hasher.Write([]byte(v))
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type Crypt struct {
	Key string
}

func NewCrypt() *Crypt {
	appKey := env.AppKey()
	if appKey == "" {
		panic("APP_KEY is not set")
	}
	return &Crypt{Key: appKey}
}

func NewCryptWithKey(key string) *Crypt {
	return &Crypt{Key: key}
}

// Encrypt takes a value and a key, both of type string, and returns the encrypted value as a string.
// It uses AES CBC encryption with PKCS7 padding.
func (c *Crypt) Encrypt(value string) (string, error) {
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return "", errors.New("failed to generate IV")
	}

	valueBytes := []byte(value)
	keyBytes := []byte(c.Key)

	res, err := openssl.AesCBCEncrypt(valueBytes, keyBytes, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}

	resVal := base64.StdEncoding.EncodeToString(res)
	resIv := base64.StdEncoding.EncodeToString(iv)

	data := resIv + resVal
	mac := computeHmacSha256(data, keyBytes)

	ticket := make(map[string]interface{})
	ticket["iv"] = resIv
	ticket["mac"] = mac
	ticket["value"] = resVal

	resTicket, err := json.Marshal(ticket)
	if err != nil {
		return "", err
	}

	ticketR := base64.StdEncoding.EncodeToString(resTicket)

	return ticketR, nil
}

// Decrypt takes an encrypted value and a key, both of type string, and returns the decrypted value as a string.
// It uses AES CBC decryption with PKCS7 padding.
func (c *Crypt) Decrypt(value string) (string, error) {
	token, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	tokenJson := make(map[string]string)
	err = json.Unmarshal(token, &tokenJson)
	if err != nil {
		return "", err
	}

	tokenJsonIv, okIv := tokenJson["iv"]
	tokenJsonValue, okValue := tokenJson["value"]
	tokenJsonMac, okMac := tokenJson["mac"]
	if !okIv || !okValue || !okMac {
		return "", errors.New("invalid token: missing iv, value, or mac (1)")
	}

	keyBytes := []byte(c.Key)

	data := tokenJsonIv + tokenJsonValue
	check := checkMAC(data, tokenJsonMac, keyBytes)
	if !check {
		return "", errors.New("invalid token: missing iv, value, or mac (2)")
	}

	tokenIv, err := base64.StdEncoding.DecodeString(tokenJsonIv)
	if err != nil {
		return "", err
	}
	tokenValue, err := base64.StdEncoding.DecodeString(tokenJsonValue)
	if err != nil {
		return "", err
	}

	dst, err := openssl.AesCBCDecrypt(tokenValue, keyBytes, tokenIv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}

	return string(dst), nil
}

// checkMAC compares the expected hash with the actual hash.
func checkMAC(message string, msgMac string, secret []byte) bool {
	expectedMAC := computeHmacSha256(message, secret)
	return hmac.Equal([]byte(expectedMAC), []byte(msgMac))
}

// computeHmacSha256 calculates the HMAC SHA256 value.
func computeHmacSha256(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

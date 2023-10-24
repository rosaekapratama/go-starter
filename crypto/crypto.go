package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"io"
)

const (
	spanGenerateKey = "common.crypto.generateKey"
	spanEncrypt     = "common.crypto.Encrypt"
	spanDecrypt     = "common.crypto.Decrypt"
)

// GenerateKey Return 32 byte encoded key from secret
func generateKey(ctx context.Context, secret string) string {
	_, span := otel.Trace(ctx, spanGenerateKey)
	defer span.End()

	hasher := md5.New()
	hasher.Write([]byte(secret))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(ctx context.Context, secret string, plainText string) (chiperText string, nonce string, err error) {
	ctx, span := otel.Trace(ctx, spanEncrypt)
	defer span.End()

	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	var block cipher.Block
	block, err = aes.NewCipher([]byte(generateKey(ctx, secret)))
	if err != nil {
		log.Error(ctx, err)
		return
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonceB := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonceB); err != nil {
		log.Error(ctx, err, "Generate nonce failed")
		return
	}

	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	ciphertextB := aesgcm.Seal(nil, nonceB, []byte(plainText), nil)
	chiperText = hex.EncodeToString(ciphertextB)
	nonce = hex.EncodeToString(nonceB)
	return
}

func Decrypt(ctx context.Context, secret string, nonce string, cipherText string) (plainText string, err error) {
	ctx, span := otel.Trace(ctx, spanDecrypt)
	defer span.End()

	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	var bCipherText []byte
	bCipherText, err = hex.DecodeString(cipherText)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	var bNonce []byte
	bNonce, err = hex.DecodeString(nonce)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	var block cipher.Block
	block, err = aes.NewCipher([]byte(generateKey(ctx, secret)))
	if err != nil {
		log.Error(ctx, err)
		return
	}

	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	var plainTextB []byte
	plainTextB, err = aesgcm.Open(nil, bNonce, bCipherText, nil)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	plainText = string(plainTextB)
	return
}

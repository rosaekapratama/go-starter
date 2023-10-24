package keycloak

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/utils"
	"golang.org/x/crypto/pbkdf2"
)

func HashPassword(ctx context.Context, plainPassword string) (base64HashedPassword string, base64Salt string, err error) {
	salt, err := utils.GenerateSalt(ctx, integer.I16) // Generate a random salt with 16 bytes
	if err != nil {
		log.Error(ctx, err, "Failed to generate salt")
		return
	}

	base64HashedPassword = hashPassword(plainPassword, salt)
	base64Salt = base64.StdEncoding.EncodeToString(salt)
	return
}

func HashPasswordWithSalt(ctx context.Context, plainPassword string, base64Salt string) (base64HashedPassword string, err error) {
	salt, err := base64.StdEncoding.DecodeString(base64Salt)
	if err != nil {
		log.Error(ctx, err, "Failed to decode base64 salt")
		return
	}

	base64HashedPassword = hashPassword(plainPassword, salt)
	return
}

func ComparePassword(ctx context.Context, plainPassword string, base64HashedPassword string, base64Salt string) (isMatch bool, err error) {
	targetBase64HashedPassword, err := HashPasswordWithSalt(ctx, plainPassword, base64Salt)
	if err != nil {
		log.Error(ctx, err, "Failed to get base64 hashed password")
		return
	}

	isMatch = base64HashedPassword == targetBase64HashedPassword
	return
}

func hashPassword(plainPassword string, salt []byte) (base64HashedPassword string) {
	// Hash the password using PBKDF2 with SHA256 and a key length of 256 bits
	hashedPassword := pbkdf2.Key([]byte(plainPassword), salt, 27500, 256/8, sha256.New)
	base64HashedPassword = base64.StdEncoding.EncodeToString(hashedPassword)
	return
}

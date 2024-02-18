package jwt

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func JwtPublicKeyID(key *rsa.PublicKey) string {
	hash := sha256.Sum256(key.N.Bytes())
	return hex.EncodeToString(hash[:])
}

func JwtPrivateKey(c *config.Config) (*rsa.PrivateKey, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(c.Jwt.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}
	return key, nil
}

func JwtPublicKey(c *config.Config) (*rsa.PublicKey, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(c.Jwt.PublicKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}
	return key, nil
}

// Ref: https://docs.temporal.io/self-hosted-guide/security#authorization
func NewToken(config *config.Config, permissions []string, subject string) (string, error) {
	privateKey, err := JwtPrivateKey(config)
	if err != nil {
		return "", err
	}

	publicKey, err := JwtPublicKey(config)
	if err != nil {
		return "", err
	}

	type WorkerClaims struct {
		jwt.RegisteredClaims
		Permissions []string `json:"permissions"`
	}

	// Create claims with multiple fields populated
	claims := WorkerClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * 30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zdravko",
			Subject:   subject,
		},
		permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = JwtPublicKeyID(publicKey)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

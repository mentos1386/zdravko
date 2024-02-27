package jwt

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func JwtPublicKeyID(key *rsa.PublicKey) string {
	hash := sha256.Sum256(key.N.Bytes())
	return hex.EncodeToString(hash[:])
}

func JwtPrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}
	return key, nil
}

func JwtPublicKey(publicKey string) (*rsa.PublicKey, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}
	return key, nil
}

type Claims struct {
	jwt.RegisteredClaims
	Permissions []string `json:"permissions"`
}

func NewTokenForUser(privateKey string, publicKey string, email string) (string, error) {
	// Create claims with multiple fields populated
	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * 30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zdravko",
			Subject:   "user:" + email,
		},
		// Ref: https://docs.temporal.io/self-hosted-guide/security#authorization
		[]string{"temporal-system:admin", "default:admin"},
	}

	return NewToken(privateKey, publicKey, claims)
}

func NewTokenForServer(privateKey string, publicKey string) (string, error) {
	// Create claims with multiple fields populated
	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * 30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zdravko",
			Subject:   "server",
		},
		// Ref: https://docs.temporal.io/self-hosted-guide/security#authorization
		[]string{"temporal-system:admin", "default:admin"},
	}

	return NewToken(privateKey, publicKey, claims)
}

func NewTokenForWorker(privateKey string, publicKey string, workerGroup *models.WorkerGroup) (string, error) {
	// Create claims with multiple fields populated
	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * 30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zdravko",
			Subject:   "worker-group:" + workerGroup.Slug,
		},
		// Ref: https://docs.temporal.io/self-hosted-guide/security#authorization
		[]string{"default:read", "default:write", "default:worker"},
	}

	return NewToken(privateKey, publicKey, claims)
}

func NewToken(privateKey string, publicKey string, claims Claims) (string, error) {
	privKey, err := JwtPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	pubKey, err := JwtPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = JwtPublicKeyID(pubKey)

	signedToken, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(tokenString string, publicKey string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JwtPublicKey(publicKey)
	})
	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}

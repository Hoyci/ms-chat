package utils

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/types"
)

var (
	PrivateKeyAccess  *rsa.PrivateKey
	PrivateKeyRefresh *rsa.PrivateKey
)

func InitJWT() error {
	pathAccess := filepath.Join(config.Envs.RootPath, "private_key_access.pem")
	pathRefresh := filepath.Join(config.Envs.RootPath, "private_key_refresh.pem")

	keyBytesAccess, err := os.ReadFile(pathAccess)
	if err != nil {
		return err
	}

	keyBytesRefresh, err := os.ReadFile(pathRefresh)
	if err != nil {
		return err
	}

	PrivateKeyAccess, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytesAccess)
	PrivateKeyRefresh, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytesRefresh)
	return err
}

func GenerateTestToken(userID int, username, email string, privateKey *rsa.PrivateKey) string {
	claims := types.CustomClaims{
		ID:       "mocked-id",
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(1 * time.Hour)},
		},
	}
	token, err := CreateJWTFromClaims(claims, privateKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate test token: %v", err))
	}
	return token
}

func CreateJWTFromClaims(claims types.CustomClaims, privateKey *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}

func CreateJWT(userID int, username string, email string, secretKey string, expTimeInSeconds int64, uuidGen types.UUIDGenerator, privateKey *rsa.PrivateKey) (string, error) {
	jti := uuidGen.New()

	claims := types.CustomClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expTimeInSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ms-chat-auth",
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}

func VerifyJWT(tokenString string, publicKey *rsa.PublicKey) (*types.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Confirma que o método de assinatura é RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*types.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

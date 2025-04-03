package utils

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hoyci/ms-chat/core/types"
)

func LoadRSAKey(rootPath, fileName string, isPrivate bool) (any, error) {
	keyPath := filepath.Join(rootPath, fileName)

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	if isPrivate {
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		return privateKey, nil
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return publicKey, nil
}

func GenerateTestToken(userID string, username, email string, privateKey *rsa.PrivateKey) string {
	claims := types.CustomClaims{
		ID:       "mocked-id",
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(1 * time.Hour)},
		},
	}
	token, err := CreateJWTTestTokenFromClaims(claims, privateKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate test token: %v", err))
	}
	return token
}

func CreateJWTTestTokenFromClaims(claims types.CustomClaims, privateKey *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}

func CreateJWT(userID string, username string, email string, secretKey string, expTimeInSeconds int64, uuidGen types.UUIDGenerator, privateKey *rsa.PrivateKey) (string, error) {
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
	token, err := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(token *jwt.Token) (any, error) {
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

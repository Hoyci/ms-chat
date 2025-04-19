package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                          int     `env:"PORT" envDefault:"8080"`
	Environment                   string  `env:"ENVIRONMENT" envDefault:"development"`
	DatabaseURL                   url.URL `env:"DATABASE_URL" envDefault:"postgres://user:password@postgres:5432/postgres?sslmode=disable"`
	AccessJWTSecret               string  `env:"ACCESS_JWT_SECRET" envDefault:"UM_ACCESS_TOKEN_MTO_DIFICIL"`
	AccessJWTExpirationInSeconds  int     `env:"ACCESS_JWT_EXPIRATION" envDefault:"3600"`
	RefreshJWTSecret              string  `env:"REFRESH_JWT_SECRET" envDefault:"UM_REFRESH_TOKEN_MTO_DIFICIL"`
	RefreshJWTExpirationInSeconds int     `env:"REFRESH_JWT_EXPIRATION" envDefault:"604800"`
	PublicKeyAccessPEM            string  `env:"PUBLIC_KEY_ACCESS,required"`
	PrivateKeyAccessPEM           string  `env:"PRIVATE_KEY_ACCESS,required"`
	PublicKeyRefreshPEM           string  `env:"PUBLIC_KEY_REFRESH,required"`
	PrivateKeyRefreshPEM          string  `env:"PRIVATE_KEY_REFRESH,required"`

	PublicKeyAccess   *rsa.PublicKey
	PrivateKeyAccess  *rsa.PrivateKey
	PublicKeyRefresh  *rsa.PublicKey
	PrivateKeyRefresh *rsa.PrivateKey
}

var Envs = initConfig()

func must[T any](val T, err error) T {
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	return val
}

func loadPublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
	pemStr = strings.ReplaceAll(pemStr, "\\n", "\n")
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid PEM block, expected PUBLIC KEY")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("decoded key is not *rsa.PublicKey")
	}

	return rsaPub, nil
}

func loadPrivateKeyFromPEM(pemStr string) (*rsa.PrivateKey, error) {
	pemStr = strings.ReplaceAll(pemStr, "\\n", "\n")
	decoded, err := base64.StdEncoding.DecodeString(pemStr)
	if err != nil {
		decoded = []byte(pemStr)
	}

	block, _ := pem.Decode(decoded)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM: no block found")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed PKCS#8 key is not RSA")
	}
	return rsaKey, nil
}

func findEnv() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		try := filepath.Join(dir, ".env")
		if info, err := os.Stat(try); err == nil && !info.IsDir() {
			return try, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func initConfig() Config {
	if path, err := findEnv(); err == nil {
		_ = godotenv.Load(path)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}

	cfg.PublicKeyAccess = must(loadPublicKeyFromPEM(cfg.PublicKeyAccessPEM))
	cfg.PrivateKeyAccess = must(loadPrivateKeyFromPEM(cfg.PrivateKeyAccessPEM))
	cfg.PublicKeyRefresh = must(loadPublicKeyFromPEM(cfg.PublicKeyRefreshPEM))
	cfg.PrivateKeyRefresh = must(loadPrivateKeyFromPEM(cfg.PrivateKeyRefreshPEM))

	return cfg
}

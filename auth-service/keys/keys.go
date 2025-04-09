package keys

import (
	"crypto/rsa"
	"github.com/hoyci/ms-chat/auth-service/config"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"log"
)

var (
	PrivateKeyAccess  *rsa.PrivateKey
	PublicKeyAccess   *rsa.PublicKey
	PrivateKeyRefresh *rsa.PrivateKey
	PublicKeyRefresh  *rsa.PublicKey
	TestPrivateKey    *rsa.PrivateKey
	TestPublicKey     *rsa.PublicKey
)

func LoadRunKeys() {
	PrivateKeyAccess = loadKey(config.Envs.KeysPath, config.Envs.PrivateKeyAccessFilename, true).(*rsa.PrivateKey)
	PublicKeyAccess = loadKey(config.Envs.KeysPath, config.Envs.PublicKeyAccessFilename, false).(*rsa.PublicKey)

	PrivateKeyRefresh = loadKey(config.Envs.KeysPath, config.Envs.PrivateKeyRefreshFilename, true).(*rsa.PrivateKey)
	PublicKeyRefresh = loadKey(config.Envs.KeysPath, config.Envs.PublicKeyRefreshFilename, false).(*rsa.PublicKey)
}

func LoadTestKeys() {
	TestPrivateKey = loadKey(config.Envs.KeysPath, config.Envs.TestPrivateKeyFilename, true).(*rsa.PrivateKey)
	TestPublicKey = loadKey(config.Envs.KeysPath, config.Envs.TestPublicKeyFilename, false).(*rsa.PublicKey)

	// Set the test keys to the variables used in the application
	PrivateKeyAccess = TestPrivateKey
	PublicKeyAccess = TestPublicKey
	PrivateKeyRefresh = TestPrivateKey
	PublicKeyRefresh = TestPublicKey
}

func loadKey(path, filename string, isPrivate bool) interface{} {
	key, err := coreUtils.LoadRSAKey(path, filename, isPrivate)
	if err != nil {
		log.Fatalf("Failed to load key: %v", err)
	}
	return key
}

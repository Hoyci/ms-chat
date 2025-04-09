package keys

import (
	"crypto/rsa"
	"fmt"
	"github.com/hoyci/ms-chat/contacts-service/config"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"log"
)

var (
	PublicKeyAccess *rsa.PublicKey
	TestPrivateKey  *rsa.PrivateKey
	TestPublicKey   *rsa.PublicKey
)

func LoadRunKeys() {
	PublicKeyAccess = loadKey(config.Envs.KeysPath, config.Envs.PublicKeyFilename, false).(*rsa.PublicKey)
	fmt.Printf("Public key loaded: %v\n", PublicKeyAccess)
}

func LoadTestKeys() {
	TestPrivateKey = loadKey(config.Envs.KeysPath, config.Envs.TestPrivateKeyFilename, true).(*rsa.PrivateKey)
	TestPublicKey = loadKey(config.Envs.KeysPath, config.Envs.TestPublicKeyFilename, false).(*rsa.PublicKey)
}

func loadKey(path, filename string, isPrivate bool) interface{} {
	key, err := coreUtils.LoadRSAKey(path, filename, isPrivate)
	if err != nil {
		log.Fatalf("Failed to load key: %v", err)
	}
	return key
}

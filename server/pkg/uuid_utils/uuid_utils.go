package uuidUtils

import (
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
)

func NewMnemonicFromUuid(uuid uuid.UUID) (string, error) {
	entropy := [16]byte(uuid)
	mnemonic, err := bip39.NewMnemonic(entropy[:])
	return mnemonic, err
}

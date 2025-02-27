package uuidUtils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
)

func NewMnemonicFromUuid(uuid uuid.UUID) (string, error) {
	entropy := [16]byte(uuid)
	mnemonic, err := bip39.NewMnemonic(entropy[:])
	return mnemonic, err
}

func UuidFromMnemonic(mnemonic []string) (uuid.UUID, error) {
	mnemonicString := strings.Join(mnemonic, " ")
	bytes, err := bip39.EntropyFromMnemonic(mnemonicString)
	if err != nil {
		return uuid.UUID([16]byte{}), fmt.Errorf("unable to get id from string %s: %w", mnemonicString, err)
	}
	return uuid.UUID(bytes), nil
}

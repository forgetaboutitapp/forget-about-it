package uuidUtils

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
)

func NewMnemonicFromUuid(uuid uuid.UUID) (string, error) {
	entropy := [16]byte(uuid)
	mnemonic, err := bip39.NewMnemonic(entropy[:])
	return mnemonic, err
}

var ErrUUIDFromMnemonic = errors.New("unable to get id from string")

func UuidFromMnemonic(mnemonic []string) (uuid.UUID, error) {
	mnemonicString := strings.Join(mnemonic, " ")
	bytes, err := bip39.EntropyFromMnemonic(mnemonicString)
	if err != nil {
		slog.Error("cant get Uuid from mnemonic", "mnemonic", mnemonic, "err", err)
		return uuid.UUID([16]byte{}), errors.Join(ErrUUIDFromMnemonic, err)
	}
	return uuid.UUID(bytes), nil
}

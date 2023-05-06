package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func String(input string) (string, error) {
	sha := sha256.New()

	sha.Write([]byte(input))

	hashBytes := sha.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

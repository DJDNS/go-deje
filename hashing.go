package deje

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

func HashObject(object interface{}) (string, error) {
	serialized, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	sum := sha1.Sum(serialized)
	return hex.EncodeToString(sum[:]), nil
}

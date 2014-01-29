package deje

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/hex"
)

func HashObject(object interface{}) (string, error) {
	serialized, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

    sum := sha1.Sum(serialized)
	return hex.EncodeToString(sum[:]), nil
}

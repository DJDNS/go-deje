package util

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

// Serialize an object to compact, deterministic JSON
// (no whitespace, fields in struct order or alphabetical),
// then take the SHA1 of that, and return the hex digest.
//
// This is the algorithm we use for hashing events, topics, etc.
func HashObject(object interface{}) (string, error) {
	serialized, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	hasher := sha1.New()
	_, _ = hasher.Write(serialized)

	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum[:]), nil
}

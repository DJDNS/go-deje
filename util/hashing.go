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
// This is the algorithm we use for hashing events, quorums,
// and IRC locations.
func HashObject(object interface{}) (string, error) {
	serialized, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	sum := sha1.Sum(serialized)
	return hex.EncodeToString(sum[:]), nil
}

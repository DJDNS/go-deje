package deje

import (
    "encoding/json"
    "crypto/sha1"
)

type SHA1Hash [20]byte

func HashObject(object interface{}) (SHA1Hash, error) {
    serialized, err := json.Marshal(object)
    if err != nil {
        return SHA1Hash{}, err
    }

    return sha1.Sum(serialized), nil
}

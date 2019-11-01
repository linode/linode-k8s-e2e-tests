package rand

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// WithRandomSuffix generates a 4-byte random suffix at the end of a given string.
func WithRandomSuffix(prefix string) (string, error) {
	suffixBytes := make([]byte, 4)

	if _, errRead := rand.Read(suffixBytes); errRead != nil {
		return "", errRead
	}

	suffix := hex.EncodeToString(suffixBytes)

	bucketName := fmt.Sprintf("%s%s", prefix, suffix)

	return bucketName, nil
}

package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateVoteHash(
	nim string,
	secret string,
) string {

	h := hmac.New(
		sha256.New,
		[]byte(secret),
	)

	h.Write([]byte(nim))

	return hex.EncodeToString(
		h.Sum(nil),
	)
}

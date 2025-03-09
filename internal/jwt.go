package internal

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
)

type payload struct {
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss"`
}

func newPayload(issuer string) payload {
	now := time.Now()

	return payload{
		// To protect against clock drift, we recommend that you set this 60 seconds in the past [...]
		now.Add(-60 * time.Second).Unix(),
		now.Add(240 * time.Second).Unix(),
		issuer,
	}
}

type payloadSigner struct {
	Signer crypto.Signer
}

func (s payloadSigner) SignPayload(payload any) (string, error) {
	hp, err := s.appendPayload(payload)

	if err != nil {
		return "", err
	}

	return s.appendSignature(hp)
}

func (s *payloadSigner) appendPayload(payload any) (string, error) {
	// `{"typ":"JWT", "alg":"RS256"}` base64 encoded
	const jwtHeader = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9"

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	return jwtHeader + "." + base64.StdEncoding.EncodeToString(jsonPayload), nil
}

func (s *payloadSigner) appendSignature(hp string) (string, error) {
	digest := sha256.Sum256([]byte(hp))

	signature, err := s.Signer.Sign(rand.Reader, digest[:], crypto.SHA256)

	if err != nil {
		return "", err
	}

	return hp + "." + base64.StdEncoding.EncodeToString(signature), nil
}

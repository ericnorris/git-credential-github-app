package internal

import (
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type AppCredentialHelper struct {
	clientID       string
	installationID string

	signer payloadSigner
	client *http.Client
}

func NewAppCredentialHelper(clientID, installationID string, s crypto.Signer) AppCredentialHelper {
	return AppCredentialHelper{
		clientID,
		installationID,
		payloadSigner{s},
		http.DefaultClient,
	}
}

func (h *AppCredentialHelper) Get(in []CredentialAttribute) ([]CredentialAttribute, error) {
	out := make([]CredentialAttribute, 0)

	for _, a := range in {
		if a.Key == "Host" && a.Value != "github.com" {
			return out, nil
		}
	}

	tokenResponse, err := h.fetchAccessToken()

	if err != nil {
		return nil, err
	}

	return []CredentialAttribute{
		{"username", "x-access-token"},
		{"password", tokenResponse.Token},
		{"password_expiry_utc", strconv.FormatInt(tokenResponse.ExpiresAt.Unix(), 10)},
	}, nil
}

type tokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *AppCredentialHelper) fetchAccessToken() (*tokenResponse, error) {
	jwt, err := h.signer.SignPayload(newPayload(h.clientID))

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", h.installationID)

	request, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", "Bearer "+jwt)
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	response, err := h.client.Do(request)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 201 {
		return nil, fmt.Errorf("unexpected status code from GitHub API (%d): %s", response.StatusCode, body)
	}

	var tr tokenResponse

	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("unexpected response body from GitHub API: %s", body)
	}

	return &tr, nil
}

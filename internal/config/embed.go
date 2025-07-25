package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed client_secret.json
var clientSecretJSON []byte

//go:embed maps_config.json
var mapsConfigJSON []byte

// MapsConfig matches your existing JSON structure
type MapsConfig struct {
	MapsAPIKey string `json:"maps_api_key"`
}

// GoogleServiceAccount matches the Google service account JSON structure
type GoogleServiceAccount struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// GetEmbeddedMapsAPIKey returns just the API key string (most common use case)
func GetEmbeddedMapsAPIKey() (string, error) {
	var config MapsConfig
	err := json.Unmarshal(mapsConfigJSON, &config)
	if err != nil {
		return "", fmt.Errorf("failed to parse embedded maps config: %w", err)
	}

	if config.MapsAPIKey == "" {
		return "", fmt.Errorf("maps API key is empty in embedded config")
	}

	return config.MapsAPIKey, nil
}

// GetEmbeddedServiceAccountJSON returns the raw JSON bytes for Google OAuth
func GetEmbeddedServiceAccountJSON() []byte {
	return clientSecretJSON
}

// GetEmbeddedServiceAccount returns the parsed service account config
func GetEmbeddedServiceAccount() (*GoogleServiceAccount, error) {
	var config GoogleServiceAccount
	err := json.Unmarshal(clientSecretJSON, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded service account config: %w", err)
	}
	return &config, nil
}

// GetEmbeddedMapsConfig returns the full maps configuration
func GetEmbeddedMapsConfig() (*MapsConfig, error) {
	var config MapsConfig
	err := json.Unmarshal(mapsConfigJSON, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded maps config: %w", err)
	}
	return &config, nil
}

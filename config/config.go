package config

import (
	"os"
)

type KeycloakConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Realm        string `json:"realm"`
	AuthServer   string `json:"auth_server"`
}

func LoadKeycloakConfig() KeycloakConfig {
	var keycloakConfig KeycloakConfig = KeycloakConfig{
		ClientID:     os.Getenv("KC_CLIENTID"),
		ClientSecret: os.Getenv("KC_CLIENTSECRET"),
		Realm:        os.Getenv("KC_REALM"),
		AuthServer:   os.Getenv("KC_AUTHSERVER"),
	}
	return keycloakConfig
}

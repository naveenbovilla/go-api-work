package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"gitlab.com/kleene/extra-hours/config"
)

type WellKnownConfig struct {
	Issuer                     string `json:"issuer"`
	AuthorizationEndpoint      string `json:"authorization_endpoint"`
	TokenEndpoint              string `json:"token_endpoint"`
	TokenIntrospectionEndpoint string `json:"token_introspection_endpoint"`
	EndSessionEndpoint         string `json:"end_session_endpoint"`
	UserInfoEndpoint           string `json:"userinfo_endpoint"`
	JwksURI                    string `json:"jwks_uri"`
}

func fetchWellKnownConfig(authServer string, realm string) (WellKnownConfig, error) {
	wellKnownURL := fmt.Sprintf("%s/realms/%s/.well-known/openid-configuration", authServer, realm)
	resp, err := http.Get(wellKnownURL)
	if err != nil {
		return WellKnownConfig{}, err
	}
	defer resp.Body.Close()

	var wellKnownConfig WellKnownConfig
	err = json.NewDecoder(resp.Body).Decode(&wellKnownConfig)
	if err != nil {
		return WellKnownConfig{}, err
	}

	return wellKnownConfig, nil
}

func fetchKeySet(jwksURI string) (*jwk.Set, error) {
	keySet, err := jwk.Fetch(context.Background(), jwksURI)
	if err != nil {
		return nil, err
	}
	return &keySet, nil
}

type User struct {
	ID         string              `json:"id"`
	Username   string              `json:"username"`
	FirstName  string              `json:"firstName"`
	LastName   string              `json:"lastName"`
	Email      string              `json:"email"`
	Enabled    bool                `json:"enabled"`
	ImageURL   string              `json:"imageURL"`
	Attributes map[string][]string `json:"attributes"`
	// Altri campi utente, se presenti
}

func GetUserDetails(accessToken string, userIds []string) ([]User, error) {
	users := []User{}
	keycloakConfig := config.LoadKeycloakConfig()

	// Costruisci una stringa di query per gli userIds
	queryParams := ""
	for _, userId := range userIds {
		if queryParams != "" {
			queryParams += "&"
		}
		queryParams += "search=" + userId
	}

	usersURL := fmt.Sprintf("%s/auth/admin/realms/%s/users?%s", keycloakConfig.AuthServer, keycloakConfig.Realm, queryParams)

	req, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keycloakUsers []User
	err = json.NewDecoder(resp.Body).Decode(&keycloakUsers)
	if err != nil {
		return nil, err
	}

	// Costruisci una mappa di utenti per ID per facilitare l'accesso durante il loop successivo
	userMap := make(map[string]User)
	for _, user := range keycloakUsers {
		userMap[user.ID] = user
	}

	// Aggiungi gli utenti nell'ordine corretto
	for _, userId := range userIds {
		user, ok := userMap[userId]
		if !ok {
			// L'utente non è stato trovato, gestisci l'errore di conseguenza
			return nil, fmt.Errorf("Utente non trovato: %s", userId)
		}
		users = append(users, user)
	}

	return users, nil
}

func GetSubClaim(token jwt.Token) (string, error) {
	sub, err := token.Get("sub")
	if !err {
		return "", nil
	}
	return sub.(string), errors.New("Invalid conversion")

}
func KeycloakAuth() gin.HandlerFunc {
	keycloakConfig := config.LoadKeycloakConfig()
	wellKnownConfig, err := fetchWellKnownConfig(keycloakConfig.AuthServer, keycloakConfig.Realm)
	if err != nil {
		panic(err)
	}

	keySet, err := fetchKeySet(wellKnownConfig.JwksURI)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		// Extract the token from the request (e.g., from the Authorization header)
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not found"})
			return
		}
		tokenStr = strings.Split(tokenStr, "Bearer ")[1]
		// Parse the token and validate it
		token, err := jwt.Parse([]byte(tokenStr), jwt.WithKeySet(*keySet))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check if the token is valid
		//if !token.IsValid() {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
		//	return
		//}

		// Add the parsed token to the request context
		c.Set("token", token)

		c.Next()
	}
}

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

// Define flags for remote path, token TTL, and ticker interval.
var (
	remotePath     = flag.String("auth-url", "http://mqtt-panel.test/api/mqtt/auth", "URL of the remote authentication service")
	tokenTTL       = flag.Duration("api-token-ttl", 24*time.Hour, "Time-to-live (TTL) in seconds for each authentication token")
	tickerInterval = flag.Duration("auto-clean-interval", 168*time.Hour, "Interval for the auto-cleaner ticker")
)

// AuthenticatedToken represents a user's authenticated session with an expiration time (TTL).
type AuthenticatedToken struct {
	TeamID       uint64 // Team ID associated with the token
	MqttClientID uint64 // MQTT Client ID associated with the token
	ApiTokenID   uint64 // API Token ID associated with the token
	TTL          uint64 // Time-to-Live (expiration) for the token, in UNIX timestamp format
}

// AuthService manages active authenticated tokens and periodically cleans up expired ones.
type AuthService struct {
	AuthenticatedList map[string]*AuthenticatedToken // Stores tokens by unique authentication key
}

// AuthServiceInstance Global instance of AuthService.
var AuthServiceInstance *AuthService

// AuthServiceInit initializes the AuthService and starts the auto-cleaner.
func AuthServiceInit() *AuthService {
	// Create a new AuthService with an initialized token map.
	AuthServiceInstance = &AuthService{
		AuthenticatedList: make(map[string]*AuthenticatedToken),
	}

	// Start the automatic cleanup of expired tokens.
	setupAutoCleaner(context.Background())
	return AuthServiceInstance
}

// setupAutoCleaner starts a background goroutine that periodically deletes expired tokens.
func setupAutoCleaner(ctx context.Context) {
	go func() {
		// Set a ticker to trigger at the interval specified by the tickerInterval flag.
		ticker := time.NewTicker(*tickerInterval)
		defer ticker.Stop() // Ensure the ticker is stopped when the function exits.

		for {
			select {
			case <-ticker.C:
				// On each tick, get the current time in UNIX format.
				now := uint64(time.Now().Unix())

				// Iterate over each token and delete if it has expired.
				for key, token := range AuthServiceInstance.AuthenticatedList {
					if token.TTL < now {
						delete(AuthServiceInstance.AuthenticatedList, key)
					}
				}
			case <-ctx.Done():
				// If the context is canceled, exit the cleanup loop.
				return
			}
		}
	}()
}

// Authenticate verifies credentials, either by cache lookup or remote authentication.
func (s *AuthService) Authenticate(clientId, username, password string) bool {
	// Generate a unique key for this client using their credentials.
	authKey := clientId + "::" + username

	// Check if the token is already in the cache and hasn't expired.
	if cache := s.AuthenticatedList[authKey]; cache != nil && cache.TTL > uint64(time.Now().Unix()) {
		cache.TTL = newTTL()
		return true // Token is valid in cache, return success and update.
	}

	// ACL authentications use only clientId + username. No cache = unauthenticated
	if password == "" {
		return false
	}

	// Perform remote authentication if token is not in cache or has expired.
	authentication, err := handleRemoteAuthentication(clientId, username, password)
	if err != nil {
		// Log the error and return false if remote authentication fails.
		fmt.Println("Error authenticating:", err)
		return false
	}

	// Cache the new authentication token for future requests.
	s.AuthenticatedList[authKey] = authentication
	return true
}

// handleRemoteAuthentication makes a remote call to validate credentials and returns a token.
func handleRemoteAuthentication(clientId, username, password string) (*AuthenticatedToken, error) {
	// Send an HTTP POST request with the provided credentials.
	response, err := sendRequest(clientId, username, password)
	if err != nil {
		return nil, err // Return an error if the request fails.
	}

	//goland:noinspection GoUnhandledErrorResult
	defer response.Body.Close() // Ensure the response body is closed to prevent memory leaks.

	// Parse the response body.
	var responseContent map[string]uint64
	if err = json.NewDecoder(response.Body).Decode(&responseContent); err != nil {
		return nil, err // Return an error if JSON decoding fails.
	}

	// Create a new AuthenticatedToken using the data from the response.
	return &AuthenticatedToken{
		TeamID:       responseContent["team_id"],
		MqttClientID: responseContent["mqtt_client_id"],
		ApiTokenID:   responseContent["api_token_id"],
		TTL:          newTTL(),
	}, nil
}

// sendRequest sends a POST request with credentials to the authentication endpoint.
func sendRequest(clientId, username, password string) (*http.Response, error) {
	// Create the JSON payload for the request.
	requestData, err := json.Marshal(map[string]string{
		"client_id":  clientId,
		"api_key":    username,
		"api_secret": password,
	})
	if err != nil {
		return nil, err // Return an error if JSON encoding fails.
	}

	// Initialize a new HTTP request with JSON headers.
	request, err := http.NewRequest("POST", *remotePath, bytes.NewReader(requestData))
	if err != nil {
		return nil, err // Return an error if request creation fails.
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Set up an HTTP client with a timeout to prevent indefinite hangs.
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err // Return an error if the request execution fails.
	}

	return response, nil // Return the HTTP response for further processing.
}

// newTTL creates a new timestamp expiry
func newTTL() uint64 {
	// Now + TTL from Flag
	return uint64(time.Now().Unix()) + uint64(tokenTTL.Seconds())
}

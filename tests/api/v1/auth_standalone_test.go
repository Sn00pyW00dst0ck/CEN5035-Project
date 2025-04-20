package v1Test

import (
	"Sector/internal/api"
	v1 "Sector/internal/api/v1"
	"Sector/internal/config"
	"Sector/internal/database"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	"encoding/pem"

	"net/http"
	"net/http/httptest"

	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

// setupAuthSuite sets up a test environment for authentication tests
func setupAuthSuite(t *testing.T) (*httptest.Server, *v1.SectorAPI, v1.Account, *rsa.PrivateKey, func()) {
	config.LoadEnv()

	// Create a test instance of the API
	tmpDir, clean := database.TestingTempDir(t, "sectordb_auth_cache_test")

	router := mux.NewRouter().StrictSlash(true)
	testSectorAPI := v1.NewTestingSector(context.Background(), "log_test.txt", tmpDir, t)
	api.AddV1SectorAPIToRouter(router, testSectorAPI)

	server := httptest.NewServer(router)

	// Generate test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Encode public key to PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	pubKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	)

	// Create test user with the public key
	now := time.Now()
	testUser := v1.Account{
		Id:         uuid.New(),
		CreatedAt:  &now,
		Username:   "testuser",
		ProfilePic: "",
	}

	// Create user data map with pubkey
	userData := v1.StructToMap(testUser)
	userData["pubkey"] = string(pubKeyPEM)

	// Add user to database
	_, err = testSectorAPI.DB.Store.Put(context.Background(), userData)
	require.NoError(t, err)

	return server, testSectorAPI, testUser, privateKey, func() {
		testSectorAPI.DB.Disconnect()
		server.Close()
		clean()
	}
}

func TestAuthenticationStandalone(t *testing.T) {
	// Setup test environment
	server, _, testAccount, privateKey, cleanup := setupAuthSuite(t)
	defer cleanup()

	testClient, err := v1.NewClientWithResponses(server.URL, v1.WithHTTPClient(server.Client()), v1.WithBaseURL(server.URL+"/v1/api"))
	require.NoError(t, err)

	t.Run("Authentication Flow", func(t *testing.T) {
		username := "testuser"

		// Step 1: Get Challenge
		challengeRequestBody := v1.GetChallengeParams{
			Username: username,
		}
		challengeResponse, err := testClient.GetChallengeWithResponse(context.Background(), &challengeRequestBody)
		require.NoError(t, err)
		require.Equal(t, 200, challengeResponse.StatusCode())

		var challengeResponseJson map[string]interface{}
		err = json.Unmarshal(challengeResponse.Body, &challengeResponseJson)
		require.NoError(t, err)

		// Step 2: Sign the challenge
		decoded, err := base64.StdEncoding.DecodeString(challengeResponseJson["challenge"].(string))
		require.NoError(t, err)
		hashed := sha256.Sum256(decoded)
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
		require.NoError(t, err, "Should sign the challenge without error")

		signatureBase64 := base64.StdEncoding.EncodeToString(signature)
		require.NotEmpty(t, signatureBase64, "Signature should not be empty")

		// Step 3: Login with signature
		loginRequestBody := v1.LoginJSONRequestBody{
			Username:  &username,
			Signature: &signatureBase64,
		}
		loginResponse, err := testClient.LoginWithResponse(context.Background(), loginRequestBody)
		require.NoError(t, err)
		require.Equal(t, 200, loginResponse.StatusCode())

		var loginResponseJson map[string]interface{}
		err = json.Unmarshal(loginResponse.Body, &loginResponseJson)
		require.NoError(t, err)

		// Step 4: Test accessing a protected endpoint with JWT Token
		authorizedReq, err := testClient.GetAccountByIDWithResponse(context.Background(), testAccount.Id, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+loginResponseJson["token"].(string))
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 200, authorizedReq.StatusCode())

		// Step 5: Test accessing a protected endpoint without JWT Token
		unauthorizedReq, err := testClient.GetAccountByIDWithResponse(context.Background(), testAccount.Id, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "")
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 401, unauthorizedReq.StatusCode())

		// Step 6: Test accessing a protected endpoint with an invalid JWT Token
		maliciousReq, err := testClient.GetAccountByIDWithResponse(context.Background(), testAccount.Id, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer INVALID_TOKEN")
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 401, maliciousReq.StatusCode())
	})
}

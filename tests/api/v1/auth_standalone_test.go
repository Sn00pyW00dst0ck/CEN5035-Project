package v1Test

import (
	"Sector/internal/api"
	v1 "Sector/internal/api/v1"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

// setupAuthSuite sets up a test environment for authentication tests
func setupAuthSuite(t *testing.T) (*httptest.Server, *v1.SectorAPI, func()) {
	// Create a test instance of the API
	router := mux.NewRouter().StrictSlash(true)
	testSectorAPI := v1.NewTestingSector(context.Background(), "log_test.txt", "test_db_cache", t)
	api.AddV1SectorAPIToRouter(router, testSectorAPI)

	server := httptest.NewServer(router)

	return server, testSectorAPI, func() {
		testSectorAPI.DB.Disconnect()
		server.Close()
	}
}

func TestAuthenticationStandalone(t *testing.T) {
	// Setup test environment
	server, sectorAPI, cleanup := setupAuthSuite(t)
	defer cleanup()

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
	_, err = sectorAPI.DB.Store.Put(context.Background(), userData)
	require.NoError(t, err)

	fmt.Println("Registered routes:")
	router := server.Config.Handler.(*mux.Router)
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		fmt.Printf("%s %s\n", methods, tpl)
		return nil
	})

	t.Run("Authentication Flow", func(t *testing.T) {
		// Step 1: Get Challenge
		req, err := http.NewRequest("GET", server.URL+"/v1/api/challenge?username=testuser", nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK for challenge request")

		var challengeResp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&challengeResp)
		require.NoError(t, err)
		require.Contains(t, challengeResp, "challenge", "Response should contain a challenge")

		challenge := challengeResp["challenge"]
		require.NotEmpty(t, challenge, "Challenge should not be empty")

		// Step 2: Sign the challenge
		hashed := sha256.Sum256([]byte(challenge))
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
		require.NoError(t, err, "Should sign the challenge without error")

		signatureBase64 := base64.StdEncoding.EncodeToString(signature)
		require.NotEmpty(t, signatureBase64, "Signature should not be empty")

		// Step 3: Login with signature
		loginPayload := `{"username":"testuser","signature":"` + signatureBase64 + `"}`
		loginReq, err := http.NewRequest("POST", server.URL+"/v1/api/login", strings.NewReader(loginPayload))
		require.NoError(t, err)
		loginReq.Header.Set("Content-Type", "application/json")

		loginResp, err := http.DefaultClient.Do(loginReq)
		require.NoError(t, err)
		defer loginResp.Body.Close()

		require.Equal(t, http.StatusOK, loginResp.StatusCode, "Should return 200 OK for successful login")

		var tokenResp map[string]string
		err = json.NewDecoder(loginResp.Body).Decode(&tokenResp)
		require.NoError(t, err)
		require.Contains(t, tokenResp, "token", "Response should contain a token")

		token := tokenResp["token"]
		require.NotEmpty(t, token, "Token should not be empty")

		// Step 4: Test accessing a protected endpoint
		protectedReq, err := http.NewRequest("GET", server.URL+"/v1/api/account/"+testUser.Id.String(), nil)
		require.NoError(t, err)
		protectedReq.Header.Set("Authorization", "Bearer "+token)

		protectedResp, err := http.DefaultClient.Do(protectedReq)
		require.NoError(t, err)
		defer protectedResp.Body.Close()

		require.Equal(t, http.StatusOK, protectedResp.StatusCode, "Should allow access to protected endpoint with valid token")

		// Step 5: Test accessing without token
		unauthorizedReq, err := http.NewRequest("GET", server.URL+"/v1/api/account/"+testUser.Id.String(), nil)
		require.NoError(t, err)

		unauthorizedResp, err := http.DefaultClient.Do(unauthorizedReq)
		require.NoError(t, err)
		defer unauthorizedResp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, unauthorizedResp.StatusCode, "Should deny access without token")

		// Step 6: Test with invalid token
		invalidTokenReq, err := http.NewRequest("GET", server.URL+"/v1/api/account/"+testUser.Id.String(), nil)
		require.NoError(t, err)
		invalidTokenReq.Header.Set("Authorization", "Bearer invalidtoken")

		invalidTokenResp, err := http.DefaultClient.Do(invalidTokenReq)
		require.NoError(t, err)
		defer invalidTokenResp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, invalidTokenResp.StatusCode, "Should deny access with invalid token")
	})
}

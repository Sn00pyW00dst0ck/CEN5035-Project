package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Encode private key to PEM format
	privKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	// Encode public key to PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}
	pubKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	)

	// Save private key to file
	err = os.WriteFile("private_key.pem", privKeyPEM, 0600)
	if err != nil {
		log.Fatalf("Failed to write private key: %v", err)
	}

	// Save public key to file
	err = os.WriteFile("public_key.pem", pubKeyPEM, 0644)
	if err != nil {
		log.Fatalf("Failed to write public key: %v", err)
	}

	// Print public key as base64 for easy copying to database
	fmt.Printf("Public Key (base64 encoded):\n%s\n", base64.StdEncoding.EncodeToString(pubKeyPEM))
	fmt.Println("\nKeys saved to 'private_key.pem' and 'public_key.pem'")
}

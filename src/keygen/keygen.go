package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// generateKeys generates an RSA key pair and saves them to files.
func generateKeys() error {
	// Generate a 2048-bit RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Encode and save the private key
	privFile, err := os.Create("../private.pem")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privFile.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	_, err = privFile.Write(privPEM)
	if err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	// Encode and save the public key
	pubFile, err := os.Create("../public.pem")
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer pubFile.Close()

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	_, err = pubFile.Write(pubPEM)
	if err != nil {
		return fmt.Errorf("failed to write public key: %v", err)
	}

	fmt.Println("RSA key pair generated and saved as private.pem and public.pem")
	return nil
}

func main() {
	if err := generateKeys(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

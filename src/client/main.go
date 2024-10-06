package main

import (
	"bufio"
	"chatapp/src/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <server address> <server port>\n", os.Args[0])
		os.Exit(1)
	}

	serverAddr := os.Args[1]
	serverPort := os.Args[2]
	address := fmt.Sprintf("%s:%s", serverAddr, serverPort)

	// Connect to the server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server at %s", address)

	// Step 1: Receive server's public key
	serverPubBytes, err := utils.ReceiveWithLengthPrefix(conn)
	if err != nil {
		log.Fatalf("Failed to receive server's public key: %v", err)
	}

	block, _ := pem.Decode(serverPubBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatalf("Invalid server public key")
	}

	serverPubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse server's public key: %v", err)
	}

	serverPubKey, ok := serverPubInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("Server's public key is not RSA")
	}

	// Step 2: Generate client's RSA key pair
	clientPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate client's private key: %v", err)
	}

	// Encode client's public key in PEM format
	clientPubBytes, err := x509.MarshalPKIXPublicKey(&clientPrivKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to marshal client's public key: %v", err)
	}

	clientPubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: clientPubBytes,
	})

	// Step 3: Send client's public key to the server
	if err := utils.SendWithLengthPrefix(conn, clientPubPEM); err != nil {
		log.Fatalf("Failed to send client's public key: %v", err)
	}

	log.Println("Public key exchange completed.")

	// Step 4: Start a goroutine to listen for incoming messages
	go func() {
		for {
			encryptedMsg, err := utils.ReceiveWithLengthPrefix(conn)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error receiving message: %v", err)
				}
				os.Exit(0)
			}

			// Decrypt the message using client's private key
			decryptedMsg, err := rsa.DecryptPKCS1v15(rand.Reader, clientPrivKey, encryptedMsg)
			if err != nil {
				log.Printf("Error decrypting message: %v", err)
				continue
			}

			fmt.Println(string(decryptedMsg))
		}
	}()

	// Step 5: Read user input and send encrypted messages
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		// Encrypt the message with server's public key
		encryptedMsg, err := rsa.EncryptPKCS1v15(rand.Reader, serverPubKey, []byte(text))
		if err != nil {
			log.Printf("Error encrypting message: %v", err)
			continue
		}

		// Send the encrypted message to the server
		if err := utils.SendWithLengthPrefix(conn, encryptedMsg); err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}
	}
}

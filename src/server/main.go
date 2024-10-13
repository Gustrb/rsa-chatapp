package main

import (
	"chatapp/src/crypt"
	"chatapp/src/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var pemImpl crypt.PEM

// Client represents a connected client
type Client struct {
	Conn      net.Conn
	PublicKey *rsa.PublicKey
	Name      string
}

// Server manages client connections and message broadcasting
type Server struct {
	Clients    map[net.Conn]*Client
	PrivateKey *rsa.PrivateKey
	PubPEM     []byte
	Mutex      sync.Mutex
}

// NewServer initializes a new Server instance
func NewServer(privateKeyPath, publicKeyPath string) (*Server, error) {
	// Load private key
	privBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	privBlock, _ := pemImpl.Decode(privBytes)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key data")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Load public key
	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %v", err)
	}

	pubBlock, _ := pemImpl.Decode(pubBytes)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key data")
	}

	sPubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	_, ok := sPubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}

	return &Server{
		Clients:    make(map[net.Conn]*Client),
		PrivateKey: privateKey,
		PubPEM:     pubBytes, // Send raw public key bytes
	}, nil
}

// Start begins listening on the specified address
func (s *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()
	log.Printf("Server listening on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("New client connected: %s", conn.RemoteAddr())
		go s.handleClient(conn)
	}
}

// handleClient manages communication with a connected client
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	// Step 1: Send server's public key with length prefix
	if err := utils.SendWithLengthPrefix(conn, s.PubPEM); err != nil {
		log.Printf("Error sending public key to %s: %v", conn.RemoteAddr(), err)
		return
	}

	// Step 2: Receive client's public key
	clientPubBytes, err := utils.ReceiveWithLengthPrefix(conn)
	if err != nil {
		log.Printf("Error receiving public key from %s: %v", conn.RemoteAddr(), err)
		return
	}

	block, _ := pemImpl.Decode(clientPubBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Printf("Invalid public key from %s", conn.RemoteAddr())
		return
	}

	cPubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Printf("Error parsing client's public key: %v", err)
		return
	}

	cPubKey, ok := cPubInterface.(*rsa.PublicKey)
	if !ok {
		log.Printf("Client's public key is not RSA")
		return
	}

	// Assign a name to the client (e.g., based on address)
	clientName := conn.RemoteAddr().String()
	client := &Client{
		Conn:      conn,
		PublicKey: cPubKey,
		Name:      clientName,
	}

	// Add client to the server's client list
	s.Mutex.Lock()
	s.Clients[conn] = client
	s.Mutex.Unlock()

	// Notify other clients about the new connection
	s.broadcast([]byte(fmt.Sprintf("[System]: %s has joined the chat.", client.Name)), conn)

	// Listen for incoming messages from this client
	for {
		// Read message with length prefix
		encryptedMsg, err := utils.ReceiveWithLengthPrefix(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from %s: %v", client.Name, err)
			}
			break
		}

		// Decrypt the message using server's private key
		decryptedMsg, err := rsa.DecryptPKCS1v15(rand.Reader, s.PrivateKey, encryptedMsg)
		if err != nil {
			log.Printf("Error decrypting message from %s: %v", client.Name, err)
			continue
		}

		log.Printf("Message from %s: %s", client.Name, string(decryptedMsg))

		// Prepare the message to broadcast
		broadcastMsg := fmt.Sprintf("[%s]: %s", client.Name, string(decryptedMsg))

		// Broadcast the message to all other clients
		s.broadcast([]byte(broadcastMsg), conn)
	}

	// Client disconnected
	s.Mutex.Lock()
	delete(s.Clients, conn)
	s.Mutex.Unlock()
	log.Printf("Client disconnected: %s", client.Name)
	s.broadcast([]byte(fmt.Sprintf("[System]: %s has left the chat.", client.Name)), conn)
}

// broadcast sends a message to all clients except the sender
func (s *Server) broadcast(message []byte, sender net.Conn) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for conn, client := range s.Clients {
		if conn == sender {
			continue
		}

		// Encrypt the message with each client's public key
		encryptedMsg, err := rsa.EncryptPKCS1v15(rand.Reader, client.PublicKey, message)
		if err != nil {
			log.Printf("Error encrypting message for %s: %v", client.Name, err)
			continue
		}

		// Send the encrypted message with length prefix
		if err := utils.SendWithLengthPrefix(conn, encryptedMsg); err != nil {
			log.Printf("Error sending message to %s: %v", client.Name, err)
		}
	}
}

func main() {
	pemImpl = &crypt.PEMStdlibImpl{}

	// Initialize the server with the paths to the RSA keys
	server, err := NewServer("../private.pem", "../public.pem")
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start the server on port 8000
	server.Start(":8000")
}

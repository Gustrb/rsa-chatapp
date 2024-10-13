package utils

import (
	"fmt"
	"io"
	"net"
)

// ReceiveWithLengthPrefix reads data prefixed with its length as 4 bytes
func ReceiveWithLengthPrefix(conn net.Conn) ([]byte, error) {
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBytes)
	if err != nil {
		return nil, err
	}
	length := int(lengthBytes[0])<<24 | int(lengthBytes[1])<<16 | int(lengthBytes[2])<<8 | int(lengthBytes[3])
	if length <= 0 {
		return nil, fmt.Errorf("invalid message length")
	}

	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// sendWithLengthPrefix sends data prefixed with its length as 4 bytes
func SendWithLengthPrefix(conn net.Conn, data []byte) error {
	length := len(data)
	lengthBytes := []byte{
		byte((length >> 24) & 0xFF),
		byte((length >> 16) & 0xFF),
		byte((length >> 8) & 0xFF),
		byte(length & 0xFF),
	}
	_, err := conn.Write(append(lengthBytes, data...))
	return err
}

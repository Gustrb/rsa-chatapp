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

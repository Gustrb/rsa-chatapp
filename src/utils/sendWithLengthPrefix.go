package utils

import "net"

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

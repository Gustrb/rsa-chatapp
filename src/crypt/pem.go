package crypt

import "encoding/pem"

type PEM interface {
	Encode([]byte, string) []byte
	Decode([]byte) (*pem.Block, []byte)
}

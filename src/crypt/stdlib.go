package crypt

import "encoding/pem"

type PEMStdlibImpl struct{}

func (*PEMStdlibImpl) Encode(bytes []byte, t string) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  t,
		Bytes: bytes,
	})
}

func (*PEMStdlibImpl) Decode(bytes []byte) (*pem.Block, []byte) {
	return pem.Decode(bytes)
}

package crypt_test

import (
	"chatapp/src/crypt"
	"testing"
)

func testEncodeSuccess(t *testing.T, plain, expected []string) {
	codec := crypt.Base64Codec{}
	for i := range plain {
		actual := codec.Encode(plain[i])

		if actual != expected[i] {
			t.Errorf("Expected: %s, got: %s", expected[i], actual)
		}
	}
}

func testDecodeSuccess(t *testing.T, plain, expected []string) {
	codec := crypt.Base64Codec{}
	for i := range plain {
		actual := codec.Decode(plain[i])

		if actual != expected[i] {
			t.Errorf("Expected: %s, got: %s", expected[i], actual)
		}
	}
}

func TestShouldEncodeBase64String(t *testing.T) {
	plain := []string{
		"this is a sample text in base 64",
		"my name is gustavo reis bauer",
		"this Is a SUPER cool base 64 ENCODER",
	}
	expected := []string{
		"dGhpcyBpcyBhIHNhbXBsZSB0ZXh0IGluIGJhc2UgNjQ=",
		"bXkgbmFtZSBpcyBndXN0YXZvIHJlaXMgYmF1ZXI=",
		"dGhpcyBJcyBhIFNVUEVSIGNvb2wgYmFzZSA2NCBFTkNPREVS",
	}
	testEncodeSuccess(t, plain, expected)
}

func TestShouldDecodeBase64String(t *testing.T) {
	plain := []string{
		"dGhpcyBpcyBhIHNhbXBsZSB0ZXh0IGluIGJhc2UgNjQ=",
		"bXkgbmFtZSBpcyBndXN0YXZvIHJlaXMgYmF1ZXI=",
		"dGhpcyBJcyBhIFNVUEVSIGNvb2wgYmFzZSA2NCBFTkNPREVS",
	}
	expected := []string{
		"this is a sample text in base 64",
		"my name is gustavo reis bauer",
		"this Is a SUPER cool base 64 ENCODER",
	}
	testDecodeSuccess(t, plain, expected)
}

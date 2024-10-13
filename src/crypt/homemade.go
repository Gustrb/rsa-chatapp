package crypt

/*
	Privacy-Enhanced Mail (PEM) is a de facto file format for storing and sending cryptographic keys, certificates, and other data,
	based on a set of 1993 IETF standards defining "privacy-enhanced mail."

	While the original standards were never broadly adopted and were supplanted by PGP and S/MIME, the textual encoding they defined became very popular.
	The PEM format was eventually formalized by the IETF in RFC 7468.[1]

	The PEM format solves this problem by encoding the binary data using base64.
	PEM also defines a one-line header, consisting of -----BEGIN, a label, and -----, and a one-line footer,
	consisting of -----END, a label, and -----. The label determines the type of message encoded.
	Common labels include CERTIFICATE, CERTIFICATE REQUEST, PRIVATE KEY and X509 CRL.
*/

import "encoding/pem"

type PEMHomeMade struct{}

func (*PEMHomeMade) Encode([]byte, string) []byte {
	return nil
}

func (*PEMHomeMade) Decode([]byte) (*pem.Block, []byte) {
	return nil, nil
}

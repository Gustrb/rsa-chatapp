package crypt

// Implementation was taken from: https://en.wikibooks.org/wiki/Algorithm_Implementation/Miscellaneous/Base64

const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

type Base64EncoderDecoder interface {
	Encode(string) string
	Decode(string) string
}

type Base64Codec struct{}

func getBase64Digits(b string, c int) (int, int, int, int) {
	n := (int(b[c]) << 16) + (int(b[c+1]) << 8) + (int(b[c+2]))
	n1 := (n >> 18) & 63
	n2 := (n >> 12) & 63
	n3 := (n >> 6) & 63
	n4 := n & 63

	return n1, n2, n3, n4
}

func (*Base64Codec) Encode(b string) string {
	result := ""
	padding := ""
	c := len(b) % 3

	if c > 0 {
		for c < 3 {
			padding += "="
			b += "\000"
			c += 1
		}
	}

	for c = 0; c < len(b); c += 3 {
		if c > 0 && (c/3*4)%76 == 0 {
			result += "\r\n"
		}

		n1, n2, n3, n4 := getBase64Digits(b, c)

		result += string(base64Chars[n1]) + string(base64Chars[n2]) + string(base64Chars[n3]) + string(base64Chars[n4])
	}
	return result[0:len(result)-len(padding)] + padding
}

func (*Base64Codec) Decode(s string) string {
	padding := 0
	if len(s) != 0 {
		if string(s[len(s)-1]) == "=" {
			padding++
		}
		if string(s[len(s)-2]) == "=" {
			padding++
		}
	}

	decodedBytes := []byte{}
	temp := 0
	cursor := 0
	for cursor < len(s) {
		for quantumPosition := 0; quantumPosition < 4; quantumPosition++ {
			temp <<= 6
			v := int(s[cursor])
			if v >= 0x41 && v <= 0x5A {
				temp |= v - 0x41
			} else if v >= 0x61 && v <= 0x7A {
				temp |= v - 0x47
			} else if v >= 0x30 && v <= 0x39 {
				temp |= v + 0x04
			} else if v == 0x2B {
				temp |= 0x3E
			} else if v == 0x2F {
				temp |= 0x3F
			} else if string(s[cursor]) == "=" {
				switch len(s) - cursor {
				case 1:
					{
						decodedBytes = append(decodedBytes, byte((temp>>16)&0x000000FF))
						decodedBytes = append(decodedBytes, byte((temp>>8)&0x000000FF))
						return string(decodedBytes)
					}
				case 2:
					{
						decodedBytes = append(decodedBytes, byte((temp>>10)&0x000000FF))
						return string(decodedBytes)
					}
				}
			}
			cursor += 1
		}

		decodedBytes = append(decodedBytes, byte((temp>>16)&0x000000FF))
		decodedBytes = append(decodedBytes, byte((temp>>8)&0x000000FF))
		decodedBytes = append(decodedBytes, byte((temp)&0x000000FF))
	}

	return string(decodedBytes)
}

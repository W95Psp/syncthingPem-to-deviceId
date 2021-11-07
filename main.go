// this file is just a bunch of copy/paste from syncthing sources

package main

import (
	"crypto/tls"
	"strings"
	"crypto/sha256"
	"encoding/base32"
	"log"
	"fmt"
	"os"
)

var luhnBase32 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

func codepoint32(b byte) int {
	switch {
	case 'A' <= b && b <= 'Z':
		return int(b - 'A')
	case '2' <= b && b <= '7':
		return int(b + 26 - '2')
	default:
		return -1
	}
}

func luhn32(s string) (rune, error) {
	factor := 1
	sum := 0
	const n = 32

	for i := range s {
		codepoint := codepoint32(s[i])
		if codepoint == -1 {
			return 0, fmt.Errorf("digit %q not valid in alphabet %q", s[i], luhnBase32)
		}
		addend := factor * codepoint
		if factor == 2 {
			factor = 1
		} else {
			factor = 2
		}
		addend = (addend / n) + (addend % n)
		sum += addend
	}
	remainder := sum % n
	checkCodepoint := (n - remainder) % n
	return rune(luhnBase32[checkCodepoint]), nil
}
func luhnify(s string) (string, error) {
	if len(s) != 52 {
		panic("unsupported string length")
	}

	res := make([]byte, 4*(13+1))
	for i := 0; i < 4; i++ {
		p := s[i*13 : (i+1)*13]
		copy(res[i*(13+1):], p)
		l, err := luhn32(p)
		if err != nil {
			return "", err
		}
		res[(i+1)*(13)+i] = byte(l)
	}
	return string(res), nil
}

func chunkify(s string) string {
	chunks := len(s) / 7
	res := make([]byte, chunks*(7+1)-1)
	for i := 0; i < chunks; i++ {
		if i > 0 {
			res[i*(7+1)-1] = '-'
		}
		copy(res[i*(7+1):], s[i*7:(i+1)*7])
	}
	return string(res)
}

func main() {
     args := os.Args[1:]
     if len(args) != 2 {
	fmt.Fprintf(os.Stderr, "Usage:\n\t%s path/cert.pem path/key.pem\n", os.Args[0])
	os.Exit(1);
     }

     certFile := args[0]
     keyFile  := args[1]
     cert, err := tls.LoadX509KeyPair(certFile, keyFile)
     if err != nil {
     	log.Fatalln("Failed to load X509 key pair:", err)
     }
     n := sha256.Sum256(cert.Certificate[0])
     id0 := base32.StdEncoding.EncodeToString(n[:])
     id0 = strings.Trim(id0, "=")
     id, err := luhnify(id0)
     if err != nil {
     	// Should never happen
     	panic(err)
     }
     fmt.Println(chunkify(id))
}


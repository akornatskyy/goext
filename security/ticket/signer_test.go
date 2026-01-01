package ticket_test

import (
	"crypto/sha1" // #nosec G505
	"encoding/base64"
	"fmt"

	"github.com/akornatskyy/goext/security/ticket"
)

var signerKey = []byte("secret")

func ExampleSigner_Signature() {
	s := ticket.NewSigner(sha1.New, signerKey)

	signature, err := s.Signature([]byte("test"))
	fmt.Println(base64.RawURLEncoding.EncodeToString(signature), err)
	// Output: GqNJWF7X7L07nEhqMAZ-OVyks1Y <nil>
}

func ExampleSigner_EncodeToString() {
	s := ticket.NewSigner(sha1.New, signerKey)

	fmt.Println(s.EncodeToString([]byte("test")))
	// Output: dGVzdBqjSVhe1-y9O5xIajAGfjlcpLNW <nil>
}

func ExampleSigner_Verify() {
	s := ticket.NewSigner(sha1.New, signerKey)

	signature, err := base64.RawURLEncoding.DecodeString(
		"GqNJWF7X7L07nEhqMAZ-OVyks1Y")
	fmt.Println(s.Verify([]byte("test"), signature), err)
	// Output: <nil> <nil>
}

func ExampleSigner_DecodeString() {
	s := ticket.NewSigner(sha1.New, signerKey)

	text, err := s.DecodeString("dGVzdBqjSVhe1-y9O5xIajAGfjlcpLNW")
	fmt.Println(string(text), err)
	// Output: test <nil>
}

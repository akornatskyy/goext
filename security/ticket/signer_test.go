package ticket_test

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/akornatskyy/goext/security/ticket"
)

var key = []byte("secret")

func ExampleSigner_Signature() {
	s := ticket.NewSigner(sha1.New, key)

	signature, err := s.Signature([]byte("test"))
	fmt.Println(base64.RawURLEncoding.EncodeToString(signature), err)
	// Output: GqNJWF7X7L07nEhqMAZ-OVyks1Y <nil>
}

func ExampleSigner_EncodeToString() {
	s := ticket.NewSigner(sha1.New, key)

	fmt.Println(s.EncodeToString([]byte("test")))
	// Output: dGVzdBqjSVhe1-y9O5xIajAGfjlcpLNW <nil>
}

func ExampleSigner_Verify() {
	s := ticket.NewSigner(sha1.New, key)

	signature, err := base64.RawURLEncoding.DecodeString(
		"GqNJWF7X7L07nEhqMAZ-OVyks1Y")
	fmt.Println(s.Verify([]byte("test"), signature), err)
	// Output: <nil> <nil>
}

func ExampleSigner_DecodeString() {
	s := ticket.NewSigner(sha1.New, key)

	text, err := s.DecodeString("dGVzdBqjSVhe1-y9O5xIajAGfjlcpLNW")
	fmt.Println(string(text), err)
	// Output: test <nil>
}

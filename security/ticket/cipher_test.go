package ticket_test

import (
	"fmt"

	"github.com/akornatskyy/goext/security/ticket"
)

var cipherKey = []byte("1234567890123456")

func Example_c_EncodeToString() {
	c := ticket.NewCipher(cipherKey)

	_, err := c.EncodeToString([]byte("test"))

	fmt.Println(err)
	// Output: <nil>
}

func Example_c_DecodeString() {
	c := ticket.NewCipher(cipherKey)

	b, err := c.DecodeString("4ADOyviA4Khjpe3VNCDaD6x5ceLzBNSGsVqBuhBbTpw")

	fmt.Println(string(b), err)
	// Output: test <nil>
}

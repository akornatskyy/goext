package ticket_test

import (
	"fmt"

	"github.com/akornatskyy/goext/security/ticket"
)

var key2 = []byte("1234567890123456")

func Example_c_EncodeToString() {
	c := ticket.NewCipher(key2)

	_, err := c.EncodeToString([]byte("test"))

	fmt.Println(err)
	// Output: <nil>
}

func Example_c_DecodeString() {
	c := ticket.NewCipher(key2)

	b, err := c.DecodeString("4ADOyviA4Khjpe3VNCDaD6x5ceLzBNSGsVqBuhBbTpw")

	fmt.Println(string(b), err)
	// Output: test <nil>
}

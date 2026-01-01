package ticket_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/akornatskyy/goext/security/ticket"
)

func TestTicket_EncodeToString(t *testing.T) {
	signer := ticket.NewSigner(sha256.New, []byte("secret-key"))
	tk := &ticket.Ticket{
		MaxAge: time.Hour,
		Signer: signer,
	}

	// Test successful encoding
	encoded, err := tk.EncodeToString([]byte("test-data"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if encoded == "" {
		t.Fatal("Expected non-empty encoded string")
	}

	// Test encoding with empty data
	encoded2, err := tk.EncodeToString([]byte(""))
	if err != nil {
		t.Fatalf("Expected no error for empty data, got %v", err)
	}
	if encoded2 == "" {
		t.Fatal("Expected non-empty encoded string for empty data")
	}
}

func TestTicket_DecodeString(t *testing.T) {
	signer := ticket.NewSigner(sha256.New, []byte("secret-key"))
	tk := &ticket.Ticket{
		MaxAge: time.Hour,
		Signer: signer,
	}

	// Test successful encode and decode
	original := []byte("test-data")
	encoded, err := tk.EncodeToString(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, remaining, err := tk.DecodeString(encoded)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(decoded) != string(original) {
		t.Fatalf("Expected %q, got %q", string(original), string(decoded))
	}
	if remaining <= 0 || remaining > int(time.Hour.Seconds()) {
		t.Fatalf("Expected remaining time between 0 and %d seconds, got %d", int(time.Hour.Seconds()), remaining)
	}
}

func TestTicket_DecodeString_Expired(t *testing.T) {
	signer := ticket.NewSigner(sha256.New, []byte("secret-key"))
	tk := &ticket.Ticket{
		MaxAge: 1 * time.Millisecond, // Very short expiration
		Signer: signer,
	}

	// Encode a ticket
	encoded, err := tk.EncodeToString([]byte("test-data"))
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Try to decode - should fail with expired error
	decoded, remaining, err := tk.DecodeString(encoded)
	if err != ticket.ErrExpired {
		t.Fatalf("Expected ErrExpired, got %v", err)
	}
	if decoded != nil {
		t.Fatalf("Expected nil decoded data, got %v", decoded)
	}
	if remaining > 0 {
		t.Fatalf("Expected non-positive remaining time, got %d", remaining)
	}
}

func TestTicket_DecodeString_InvalidSignature(t *testing.T) {
	signer1 := ticket.NewSigner(sha256.New, []byte("secret-key-1"))
	signer2 := ticket.NewSigner(sha256.New, []byte("secret-key-2"))

	tk1 := &ticket.Ticket{
		MaxAge: time.Hour,
		Signer: signer1,
	}
	tk2 := &ticket.Ticket{
		MaxAge: time.Hour,
		Signer: signer2,
	}

	// Encode with first ticket
	encoded, err := tk1.EncodeToString([]byte("test-data"))
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Try to decode with second ticket (different key) - should fail
	_, _, err = tk2.DecodeString(encoded)
	if err == nil {
		t.Fatal("Expected error when decoding with different key, got nil")
	}
}

func TestTicket_DecodeString_InvalidInput(t *testing.T) {
	signer := ticket.NewSigner(sha256.New, []byte("secret-key"))
	tk := &ticket.Ticket{
		MaxAge: time.Hour,
		Signer: signer,
	}

	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"invalid base64", "not-valid-base64!!!"},
		{"random string", "xyz123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := tk.DecodeString(tt.input)
			if err == nil {
				t.Fatalf("Expected error for input %q, got nil", tt.input)
			}
		})
	}
}

func TestTicket_RoundTrip(t *testing.T) {
	signer := ticket.NewSigner(sha256.New, []byte("secret-key"))
	tk := &ticket.Ticket{
		MaxAge: 5 * time.Second,
		Signer: signer,
	}

	testCases := [][]byte{
		[]byte("simple"),
		[]byte("with spaces and special chars: !@#$%^&*()"),
		[]byte("unicode: 你好世界 🌍"),
		[]byte(""),
		make([]byte, 1024), // large data
	}

	for _, data := range testCases {
		t.Run("", func(t *testing.T) {
			encoded, err := tk.EncodeToString(data)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			decoded, remaining, err := tk.DecodeString(encoded)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if string(decoded) != string(data) {
				t.Fatalf("Round trip failed. Expected %q, got %q", string(data), string(decoded))
			}

			if remaining <= 0 || remaining > 5 {
				t.Fatalf("Expected remaining time between 0 and 5 seconds, got %d", remaining)
			}
		})
	}
}

func TestTicket_MultipleEncodingsAreDifferent(t *testing.T) {
	signer := ticket.NewSigner(sha256.New, []byte("secret-key"))
	tk := &ticket.Ticket{
		MaxAge: time.Hour,
		Signer: signer,
	}

	data := []byte("test-data")

	// Encode the same data multiple times
	encoded1, err := tk.EncodeToString(data)
	if err != nil {
		t.Fatalf("Encode 1 failed: %v", err)
	}

	time.Sleep(1100 * time.Millisecond) // Ensure different second (unix timestamp granularity)

	encoded2, err := tk.EncodeToString(data)
	if err != nil {
		t.Fatalf("Encode 2 failed: %v", err)
	}

	// Encoded strings should be different due to different timestamps
	if encoded1 == encoded2 {
		t.Fatal("Expected different encoded strings for same data at different times")
	}

	// But both should decode to the same original data
	decoded1, _, err := tk.DecodeString(encoded1)
	if err != nil {
		t.Fatalf("Decode 1 failed: %v", err)
	}

	decoded2, _, err := tk.DecodeString(encoded2)
	if err != nil {
		t.Fatalf("Decode 2 failed: %v", err)
	}

	if string(decoded1) != string(data) || string(decoded2) != string(data) {
		t.Fatal("Decoded data does not match original")
	}
}

package httptoken_test

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akornatskyy/goext/httptoken"
	"github.com/akornatskyy/goext/security/ticket"
)

// mockPrincipal is a test struct that implements BinaryMarshaler/BinaryUnmarshaler
type mockPrincipal struct {
	ID   int
	Name string
}

func (p *mockPrincipal) MarshalBinary() ([]byte, error) {
	buf := make([]byte, binary.MaxVarintLen64+len(p.Name))
	n := binary.PutVarint(buf, int64(p.ID))
	copy(buf[n:], p.Name)
	return buf[:n+len(p.Name)], nil
}

func (p *mockPrincipal) UnmarshalBinary(data []byte) error {
	id, n := binary.Varint(data)
	if n <= 0 {
		return errors.New("invalid data")
	}
	p.ID = int(id)
	p.Name = string(data[n:])
	return nil
}

// failingMarshaler always fails to marshal
type failingMarshaler struct{}

func (f *failingMarshaler) MarshalBinary() ([]byte, error) {
	return nil, errors.New("marshal error")
}

// failingUnmarshaler always fails to unmarshal
type failingUnmarshaler struct{}

func (f *failingUnmarshaler) UnmarshalBinary(data []byte) error {
	return errors.New("unmarshal error")
}

func newCookieToken() *httptoken.CookieToken {
	return &httptoken.CookieToken{
		Name: "auth",
		Ticket: &ticket.Ticket{
			MaxAge: 2 * time.Hour,
			Signer: ticket.NewSigner(sha256.New, []byte("secret-key-32-bytes-long-12345")),
		},
	}
}

func TestCookieToken_Write(t *testing.T) {
	token := newCookieToken()
	principal := &mockPrincipal{ID: 123, Name: "testuser"}

	w := httptest.NewRecorder()
	err := token.Write(w, principal)

	if err != nil {
		t.Errorf("Write() error = %v, want nil", err)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "auth" {
		t.Errorf("cookie.Name = %q, want %q", cookie.Name, "auth")
	}
	if cookie.Path != "/" {
		t.Errorf("cookie.Path = %q, want %q", cookie.Path, "/")
	}
	if !cookie.HttpOnly {
		t.Error("cookie.HttpOnly = false, want true")
	}
	if cookie.Value == "" {
		t.Error("cookie.Value is empty")
	}

	// Check expiration is approximately correct (within 10 seconds)
	expectedExpires := time.Now().Add(2 * time.Hour)
	diff := cookie.Expires.Sub(expectedExpires)
	if diff < -10*time.Second || diff > 10*time.Second {
		t.Errorf("cookie.Expires = %v, want approximately %v", cookie.Expires, expectedExpires)
	}
}

func TestCookieToken_Write_MarshalError(t *testing.T) {
	token := newCookieToken()
	principal := &failingMarshaler{}

	w := httptest.NewRecorder()
	err := token.Write(w, principal)

	if err == nil {
		t.Error("Write() error = nil, want error")
	}
	if err.Error() != "marshal error" {
		t.Errorf("Write() error = %q, want %q", err.Error(), "marshal error")
	}
}

func TestCookieToken_Authorize(t *testing.T) {
	token := newCookieToken()
	principal := &mockPrincipal{ID: 456, Name: "alice"}

	// First, write a cookie
	w := httptest.NewRecorder()
	if err := token.Write(w, principal); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	// Now authorize with that cookie
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookies[0])

	w2 := httptest.NewRecorder()
	result := &mockPrincipal{}
	err := token.Authorize(w2, r, result)

	if err != nil {
		t.Errorf("Authorize() error = %v, want nil", err)
	}
	if result.ID != 456 {
		t.Errorf("result.ID = %d, want 456", result.ID)
	}
	if result.Name != "alice" {
		t.Errorf("result.Name = %q, want %q", result.Name, "alice")
	}
}

func TestCookieToken_Authorize_NilUnmarshaler(t *testing.T) {
	token := newCookieToken()
	principal := &mockPrincipal{ID: 789, Name: "bob"}

	// Write a cookie
	w := httptest.NewRecorder()
	if err := token.Write(w, principal); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	cookies := w.Result().Cookies()

	// Authorize with nil unmarshaler (just validate the token)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookies[0])

	w2 := httptest.NewRecorder()
	err := token.Authorize(w2, r, nil)

	if err != nil {
		t.Errorf("Authorize() error = %v, want nil", err)
	}
}

func TestCookieToken_Authorize_NoCookie(t *testing.T) {
	token := newCookieToken()

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	result := &mockPrincipal{}

	err := token.Authorize(w, r, result)

	if err == nil {
		t.Error("Authorize() error = nil, want error")
	}
	if !errors.Is(err, http.ErrNoCookie) {
		t.Errorf("Authorize() error = %v, want http.ErrNoCookie", err)
	}
}

func TestCookieToken_Authorize_InvalidToken(t *testing.T) {
	token := newCookieToken()

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: "invalid-token-value",
	})

	w := httptest.NewRecorder()
	result := &mockPrincipal{}
	err := token.Authorize(w, r, result)

	if err == nil {
		t.Error("Authorize() error = nil, want error")
	}
}

func TestCookieToken_Authorize_UnmarshalError(t *testing.T) {
	token := newCookieToken()
	principal := &mockPrincipal{ID: 999, Name: "charlie"}

	// Write a valid cookie
	w := httptest.NewRecorder()
	if err := token.Write(w, principal); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	cookies := w.Result().Cookies()

	// Try to authorize with a failing unmarshaler
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookies[0])

	w2 := httptest.NewRecorder()
	result := &failingUnmarshaler{}
	err := token.Authorize(w2, r, result)

	if err == nil {
		t.Error("Authorize() error = nil, want error")
	}
	if err.Error() != "unmarshal error" {
		t.Errorf("Authorize() error = %q, want %q", err.Error(), "unmarshal error")
	}
}

func TestCookieToken_Authorize_RenewsCookieWhenNearExpiry(t *testing.T) {
	// Create a token with short MaxAge to test renewal
	token := &httptoken.CookieToken{
		Name: "auth",
		Ticket: &ticket.Ticket{
			MaxAge: 10 * time.Second, // Short duration to make it near expiry
			Signer: ticket.NewSigner(sha256.New, []byte("secret-key-32-bytes-long-12345")),
		},
	}
	principal := &mockPrincipal{ID: 111, Name: "dave"}

	// Write a cookie
	w := httptest.NewRecorder()
	if err := token.Write(w, principal); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	cookies := w.Result().Cookies()

	// Wait a bit so we're past the halfway point (5 seconds)
	time.Sleep(6 * time.Second)

	// Authorize - should renew the cookie
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookies[0])

	w2 := httptest.NewRecorder()
	result := &mockPrincipal{}
	err := token.Authorize(w2, r, result)

	if err != nil {
		t.Errorf("Authorize() error = %v, want nil", err)
	}

	// Check that a new cookie was set (renewal)
	renewedCookies := w2.Result().Cookies()
	if len(renewedCookies) != 1 {
		t.Errorf("expected 1 renewed cookie, got %d", len(renewedCookies))
	}
}

func TestCookieToken_Authorize_DoesNotRenewWhenNotNearExpiry(t *testing.T) {
	token := newCookieToken() // 2 hour MaxAge
	principal := &mockPrincipal{ID: 222, Name: "eve"}

	// Write a cookie
	w := httptest.NewRecorder()
	if err := token.Write(w, principal); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	cookies := w.Result().Cookies()

	// Authorize immediately - should NOT renew (not past halfway point)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookies[0])

	w2 := httptest.NewRecorder()
	result := &mockPrincipal{}
	err := token.Authorize(w2, r, result)

	if err != nil {
		t.Errorf("Authorize() error = %v, want nil", err)
	}

	// Check that no new cookie was set
	renewedCookies := w2.Result().Cookies()
	if len(renewedCookies) != 0 {
		t.Errorf("expected 0 renewed cookies, got %d", len(renewedCookies))
	}
}

func TestCookieToken_Invalidate(t *testing.T) {
	token := newCookieToken()

	w := httptest.NewRecorder()
	token.Invalidate(w)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "auth" {
		t.Errorf("cookie.Name = %q, want %q", cookie.Name, "auth")
	}
	if cookie.Path != "/" {
		t.Errorf("cookie.Path = %q, want %q", cookie.Path, "/")
	}

	// Check that expiration is set to epoch time (invalidates the cookie)
	expectedExpires := time.Unix(0, 0)
	if !cookie.Expires.Equal(expectedExpires) {
		t.Errorf("cookie.Expires = %v, want %v", cookie.Expires, expectedExpires)
	}
}

func TestCookieToken_Authorize_RenewalBoundary(t *testing.T) {
	// Test the exact boundary condition where remaining == MaxAge/2
	token := &httptoken.CookieToken{
		Name: "auth",
		Ticket: &ticket.Ticket{
			MaxAge: 4 * time.Second,
			Signer: ticket.NewSigner(sha256.New, []byte("secret-key-32-bytes-long-12345")),
		},
	}

	principal := &mockPrincipal{ID: 333, Name: "boundary"}

	// Write a cookie
	w := httptest.NewRecorder()
	if err := token.Write(w, principal); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	cookies := w.Result().Cookies()

	// Wait exactly 2 seconds (halfway point)
	time.Sleep(2 * time.Second)

	// Authorize - might or might not renew depending on exact timing
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookies[0])

	w2 := httptest.NewRecorder()
	result := &mockPrincipal{}
	err := token.Authorize(w2, r, result)

	if err != nil {
		t.Errorf("Authorize() error = %v, want nil", err)
	}

	// Verify the principal was correctly decoded regardless of renewal
	if result.ID != 333 {
		t.Errorf("result.ID = %d, want 333", result.ID)
	}
	if result.Name != "boundary" {
		t.Errorf("result.Name = %q, want %q", result.Name, "boundary")
	}
}

func TestCookieToken_CookieAttributes(t *testing.T) {
	// Test that the cookie helper creates cookies with correct attributes
	token := &httptoken.CookieToken{
		Name: "session",
		Ticket: &ticket.Ticket{
			MaxAge: 1 * time.Hour,
			Signer: ticket.NewSigner(sha256.New, []byte("test-secret-key-for-testing-32")),
		},
	}

	principal := &mockPrincipal{ID: 999, Name: "attrtest"}
	w := httptest.NewRecorder()

	if err := token.Write(w, principal); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]

	// Verify all expected attributes
	if cookie.Name != "session" {
		t.Errorf("cookie.Name = %q, want %q", cookie.Name, "session")
	}
	if cookie.Path != "/" {
		t.Errorf("cookie.Path = %q, want %q", cookie.Path, "/")
	}
	if !cookie.HttpOnly {
		t.Error("cookie.HttpOnly = false, want true")
	}
	if cookie.Secure {
		t.Error("cookie.Secure = true, want false (not set in implementation)")
	}
	if cookie.SameSite != 0 {
		t.Errorf("cookie.SameSite = %v, want 0 (not set in implementation)", cookie.SameSite)
	}
	if cookie.Value == "" {
		t.Error("cookie.Value is empty")
	}

	// Verify expires is in the future
	if cookie.Expires.Before(time.Now()) {
		t.Error("cookie.Expires is in the past")
	}
}

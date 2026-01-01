# goext - Go Extensions and Toolkit

[![tests](https://github.com/akornatskyy/goext/actions/workflows/tests.yaml/badge.svg)](https://github.com/akornatskyy/goext/actions/workflows/tests.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/akornatskyy/goext)](https://goreportcard.com/report/github.com/akornatskyy/goext) [![Go Doc](https://godoc.org/github.com/akornatskyy/goext?status.svg)](https://godoc.org/github.com/akornatskyy/goext)

A collection of practical Go packages that extend the standard library with common functionality for web applications, data validation, security, and more.

## Features

- **Data Binding**: Automatic struct population from HTTP form values
- **Validation**: Fluent API for validating strings and numbers
- **Error State Management**: Structured error handling with detailed error information
- **HTTP JSON Utilities**: Simplified JSON encoding/decoding for HTTP handlers
- **Security**: Token-based authentication with encrypted tickets
- **SQL Utilities**: Helper functions for database operations
- **JSON I/O**: Convenient file I/O operations with JSON encoding

## Installation

```sh
go get github.com/akornatskyy/goext
```

## Packages

### binding

Automatically binds HTTP form values to struct fields using reflection and struct tags.

```go
import "github.com/akornatskyy/goext/binding"

type LoginForm struct {
    Username string `binding:"username"`
    Password string `binding:"password"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    var form LoginForm
    if err := binding.Bind(&form, r.Form); err != nil {
        // Handle binding errors
    }
}
```

Supports: `string`, `int`, `uint`, `bool`, `time.Duration`, `time.Time`, and slices.

### validator

Fluent validation API with chainable rules for strings and numbers.

```go
import "github.com/akornatskyy/goext/validator"

// String validation
emailValidator := validator.String("email").
    Required().
    Email().
    Build()

// Number validation
ageValidator := validator.Number("age").
    Required().
    Min(18).
    Max(100).
    Build()

e := &errorstate.ErrorState{Domain: "registration"}
if !emailValidator.Validate(e, email) {
    // Handle validation errors
}
```

### errorstate

Structured error handling that collects multiple validation or business logic errors.

```go
import "github.com/akornatskyy/goext/errorstate"

e := &errorstate.ErrorState{Domain: "user"}
e.Add(&errorstate.Detail{
    Type:     "field",
    Location: "username",
    Reason:   "required",
    Message:  "Username is required",
})

if err := e.OrNil(); err != nil {
    // Returns nil if no errors, otherwise returns the error state
}
```

### httpjson

Simplified JSON encoding and decoding for HTTP handlers.

```go
import "github.com/akornatskyy/goext/httpjson"

func handler(w http.ResponseWriter, r *http.Request) {
    // Decode request
    var data RequestData
    if err := httpjson.Decode(r, &data); err != nil {
        httpjson.Encode(w, ErrorResponse{Error: err.Error()}, 400)
        return
    }

    // Encode response
    httpjson.Encode(w, ResponseData{Success: true}, 200)
}
```

### httptoken

Token-based authentication interface for HTTP handlers with cookie support.

```go
import "github.com/akornatskyy/goext/httptoken"

type Token interface {
    Write(w http.ResponseWriter, m encoding.BinaryMarshaler) error
    Authorize(w http.ResponseWriter, r *http.Request, p encoding.BinaryUnmarshaler) error
    Invalidate(w http.ResponseWriter)
}
```

### security/ticket

Time-limited encrypted tickets with signing and optional encryption.

```go
import "github.com/akornatskyy/goext/security/ticket"

t := &ticket.Ticket{
    MaxAge: time.Hour,
    Signer: signer,
}

// Encode data into a signed ticket
encoded, err := t.EncodeToString([]byte("sensitive data"))

// Decode and verify ticket
data, remaining, err := t.DecodeString(encoded)
if err == ticket.ErrExpired {
    // Ticket has expired
}
```

### iojson

Convenient JSON file I/O operations.

```go
import "github.com/akornatskyy/goext/iojson"

var config Config
if err := iojson.ReadFile("config.json", &config); err != nil {
    // Handle error
}
```

### sqlx

Database utility functions.

```go
import "github.com/akornatskyy/goext/sqlx"

stmt := sqlx.MustPrepare(db, "SELECT * FROM users WHERE id = ?")
defer stmt.Close()
```

## Requirements

- Go 1.11 or higher

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

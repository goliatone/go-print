# Go Print Utilities

Simple go package to print and save JSON data. It can print regular JSON and secure JSON by masking sensitive values like passwords and tokens.

## Install

```bash
go get github.com/goliatone/go-print
```

## Usage

Basic JSON printing:

```go
data := map[string]string{
    "user": "admin",
    "pass": "secret"
}

// Pretty print JSON
str, err := print.PrettyJSON(data)

// Print JSON, returns error message if fails
str := print.MaybePrettyJSON(data)

// Save to file
err := print.SaveJSONFile("data.json", data)
```

Secure JSON (masks sensitive data):

```go
type User struct {
    Username string `json:"username"`
    Password string `json:"password" mask:"filled4"`  // Will be masked
    APIKey   string `json:"api_key" mask:"filled32"` // Will be masked
}

user := User{
    Username: "admin",
    Password: "secret",
    APIKey:   "1234567890",
}

// Print masked JSON
str, err := print.SecureJSON(user)

// Print masked JSON, returns error message if fails
str := print.MaybeSecureJSON(user)

// Save masked JSON to file
err := print.SaveSecureJSONFile("user.json", user)
```

HTTP Request/Response printing:

```go
// Print HTTP request as JSON
req, _ := http.NewRequest("POST", "https://api.example.com", nil)
str := print.PrintHTTPRequest(req)

// Print HTTP response as JSON
resp, _ := http.Get("https://api.example.com")
str := print.PrintHTTPResponse(resp)
```

## Features

- Pretty prints JSON with proper indentation
- Masks sensitive data (passwords, tokens, keys)
- Saves JSON to files
- Prints HTTP requests and responses as JSON
- Thread safe
- Handles errors gracefully

## Default Masked Fields

These fields are masked by default:
- Password/password
- SigningKey/signing_key
- Authorization/authorization

## License

MIT

Copyright (c) 2024 goliatone

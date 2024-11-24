// Package print provides utilities for printing and saving JSON data with support
// for masking sensitive information.
//
// Basic Usage:
//
//	data := map[string]string{
//	    "user": "admin",
//	    "pass": "secret"
//	}
//
//	// Pretty print JSON
//	str, err := print.PrettyJSON(data)
//
//	// Print JSON, returns error message if fails
//	str := print.MaybePrettyJSON(data)
//
//	// Save to file
//	err := print.SaveJSONFile("data.json", data)
//
// Secure JSON (with masked sensitive data):
//
//	type User struct {
//	    Username string `json:"username"`
//	    Password string `json:"password" mask:"filled4"`  // Will be masked
//	    APIKey   string `json:"api_key" mask:"filled32"` // Will be masked
//	}
//
//	user := User{
//	    Username: "admin",
//	    Password: "secret",
//	    APIKey:   "1234567890",
//	}
//
//	// Print masked JSON
//	str, err := print.SecureJSON(user)
//
//	// Save masked JSON to file
//	err := print.SaveSecureJSONFile("user.json", user)
//
// HTTP Request/Response printing:
//
//	// Print HTTP request as JSON
//	req, _ := http.NewRequest("POST", "https://api.example.com", nil)
//	str := print.PrintHTTPRequest(req)
//
//	// Print HTTP response as JSON
//	resp, _ := http.Get("https://api.example.com")
//	str := print.PrintHTTPResponse(resp)
//
// Default Masked Fields:
//   - Password/password
//   - SigningKey/signing_key
//   - Authorization/authorization
//
// Features:
//   - Pretty prints JSON with proper indentation
//   - Masks sensitive data (passwords, tokens, keys)
//   - Saves JSON to files
//   - Prints HTTP requests and responses as JSON
//   - Thread safe
//   - Handles errors gracefully
package print

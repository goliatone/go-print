package print

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestPrintHTTPRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *http.Request
		wantBody string
		wantErr  bool
	}{
		{
			name:     "nil request",
			request:  nil,
			wantBody: "nil",
		},
		{
			name: "GET request without body",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com/test?key=value", nil)
				req.Header.Add("X-Test", "test-value")
				return req
			}(),
			wantBody: `{
	"method": "GET",
	"url": {
		"Scheme": "https",
		"Opaque": "",
		"User": null,
		"Host": "example.com",
		"Path": "/test",
		"RawPath": "",
		"OmitHost": false,
		"ForceQuery": false,
		"RawQuery": "key=value",
		"Fragment": "",
		"RawFragment": ""
	},
	"header": {
		"X-Test": ["test-value"]
	},
	"content_length": 0
}`,
		},
		{
			name: "POST request with body",
			request: func() *http.Request {
				body := strings.NewReader(`{"key":"value"}`)
				req, _ := http.NewRequest("POST", "https://example.com/api", body)
				req.Header.Add("Content-Type", "application/json")
				return req
			}(),
			wantBody: `{
	"method": "POST",
	"url": {
		"Scheme": "https",
		"Opaque": "",
		"User": null,
		"Host": "example.com",
		"Path": "/api",
		"RawPath": "",
		"OmitHost": false,
		"ForceQuery": false,
		"RawQuery": "",
		"Fragment": "",
		"RawFragment": ""
	},
	"header": {
		"Content-Type": ["application/json"]
	},
	"body": "{\"key\":\"value\"}",
	"content_length": 15
}`,
		},
		{
			name: "request with sensitive headers",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com", nil)
				req.Header.Add("Authorization", "Bearer secret-token")
				req.Header.Add("X-API-Key", "secret-key")
				return req
			}(),
			wantBody: `{
	"method": "GET",
	"url": {
		"Scheme": "https",
		"Opaque": "",
		"User": null,
		"Host": "example.com",
		"Path": "",
		"RawPath": "",
		"OmitHost": false,
		"ForceQuery": false,
		"RawQuery": "",
		"Fragment": "",
		"RawFragment": ""
	},
	"header": {
		"Authorization": ["********************************"],
		"X-Api-Key": ["secret-key"]
	},
	"content_length": 0
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrintHTTPRequest(tt.request)
			if tt.wantBody == "nil" {
				if got != "nil" {
					t.Errorf("PrintHTTPRequest() = %v, want nil", got)
				}
				return
			}

			if !compareJSON(t, got, tt.wantBody) {
				t.Errorf("PrintHTTPRequest() = %v, want %v", got, tt.wantBody)
			}

			// Verify body can still be read if it existed
			if tt.request != nil && tt.request.Body != nil {
				body, err := io.ReadAll(tt.request.Body)
				if err != nil {
					t.Errorf("Failed to read body after printing: %v", err)
				}
				if len(body) == 0 && tt.request.ContentLength > 0 {
					t.Error("Body was not properly reset after reading")
				}
			}
		})
	}
}

func TestPrintHTTPResponse(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		wantBody string
		wantErr  bool
	}{
		{
			name:     "nil response",
			response: nil,
			wantBody: "nil",
		},
		{
			name: "simple response",
			response: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body:          io.NopCloser(strings.NewReader(`{"message":"success"}`)),
				ContentLength: 20,
			},
			wantBody: `{
	"status": "200 OK",
	"status_code": 200,
	"header": {
		"Content-Type": ["application/json"]
	},
	"body": "{\"message\":\"success\"}",
	"content_length": 20,
	"request": null
}`,
		},
		{
			name: "response with request",
			response: func() *http.Response {
				req, _ := http.NewRequest("GET", "https://example.com", nil)
				return &http.Response{
					Status:     "200 OK",
					StatusCode: 200,
					Header: http.Header{
						"Content-Type": []string{"text/plain"},
					},
					Body:          io.NopCloser(strings.NewReader("Hello")),
					ContentLength: 5,
					Request:       req,
				}
			}(),
			wantBody: `{
	"status": "200 OK",
	"status_code": 200,
	"header": {
		"Content-Type": ["text/plain"]
	},
	"body": "Hello",
	"content_length": 5,
	"request": {
		"method": "GET",
		"url": {
			"Scheme": "https",
			"Opaque": "",
			"User": null,
			"Host": "example.com",
			"Path": "",
			"RawPath": "",
			"OmitHost": false,
			"ForceQuery": false,
			"RawQuery": "",
			"Fragment": "",
			"RawFragment": ""
		},
		"header": {},
		"content_length": 0
	}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrintHTTPResponse(tt.response)
			if tt.wantBody == "nil" {
				if got != "nil" {
					t.Errorf("PrintHTTPResponse() = %v, want nil", got)
				}
				return
			}

			if !compareJSON(t, got, tt.wantBody) {
				t.Errorf("PrintHTTPResponse() = %v, want %v", got, tt.wantBody)
			}

			// Verify body can still be read if it existed
			if tt.response != nil && tt.response.Body != nil {
				body, err := io.ReadAll(tt.response.Body)
				if err != nil {
					t.Errorf("Failed to read body after printing: %v", err)
				}
				if len(body) == 0 && tt.response.ContentLength > 0 {
					t.Error("Body was not properly reset after reading")
				}
			}
		})
	}
}

func TestIsInterfaceNil(t *testing.T) {
	tests := []struct {
		name string
		i    interface{}
		want bool
	}{
		{
			name: "nil interface",
			i:    nil,
			want: true,
		},
		{
			name: "nil pointer",
			i:    (*string)(nil),
			want: true,
		},
		{
			name: "non-nil pointer",
			i:    new(string),
			want: false,
		},
		{
			name: "nil slice",
			i:    []string(nil),
			want: true,
		},
		{
			name: "empty slice",
			i:    []string{},
			want: false,
		},
		{
			name: "nil map",
			i:    map[string]string(nil),
			want: true,
		},
		{
			name: "empty map",
			i:    map[string]string{},
			want: false,
		},
		{
			name: "string value",
			i:    "test",
			want: false,
		},
		{
			name: "integer value",
			i:    42,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInterfaceNil(tt.i); got != tt.want {
				t.Errorf("IsInterfaceNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotInterfaceNil(t *testing.T) {
	tests := []struct {
		name string
		i    interface{}
		want bool
	}{
		{
			name: "nil interface",
			i:    nil,
			want: false,
		},
		{
			name: "non-nil value",
			i:    "test",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotInterfaceNil(tt.i); got != tt.want {
				t.Errorf("NotInterfaceNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBodyAsString(t *testing.T) {
	tests := []struct {
		name        string
		body        io.ReadCloser
		length      int64
		wantBody    string
		wantReset   bool
		wantBytes   []byte
		shouldError bool
	}{
		{
			name:      "nil body",
			body:      nil,
			length:    0,
			wantBody:  "",
			wantReset: false,
			wantBytes: []byte{},
		},
		{
			name:      "empty body",
			body:      io.NopCloser(strings.NewReader("")),
			length:    0,
			wantBody:  "",
			wantReset: true,
			wantBytes: []byte{},
		},
		{
			name:      "valid body",
			body:      io.NopCloser(strings.NewReader("test content")),
			length:    12,
			wantBody:  "test content",
			wantReset: true,
			wantBytes: []byte("test content"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, gotReset, gotBytes := bodyAsString(tt.body, tt.length)

			if gotBody != tt.wantBody {
				t.Errorf("bodyAsString() body = %v, want %v", gotBody, tt.wantBody)
			}
			if gotReset != tt.wantReset {
				t.Errorf("bodyAsString() reset = %v, want %v", gotReset, tt.wantReset)
			}
			if !bytes.Equal(gotBytes, tt.wantBytes) {
				t.Errorf("bodyAsString() bytes = %v, want %v", gotBytes, tt.wantBytes)
			}
		})
	}
}

// Helper function to compare JSON strings
func compareJSON(t *testing.T, got, want string) bool {
	var gotJSON, wantJSON interface{}

	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		t.Errorf("Failed to parse got JSON: %v", err)
		return false
	}

	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Errorf("Failed to parse want JSON: %v", err)
		return false
	}

	return reflect.DeepEqual(gotJSON, wantJSON)
}

package print

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password" mask:"filled4"`
	APIKey   string `json:"api_key" mask:"filled32"`
}

type BadJSON struct {
	Ch chan int
}

func TestPrettyJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name: "valid simple object",
			input: map[string]string{
				"hello": "world",
			},
			want: `{
	"hello": "world"
}
`,
			wantErr: false,
		},
		{
			name:  "empty object",
			input: map[string]string{},
			want: `{}
`,
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			want:    "null\n",
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   BadJSON{Ch: make(chan int)},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrettyJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrettyJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrettyJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaybePrettyJSON(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name: "valid json",
			input: map[string]string{
				"hello": "world",
			},
			want: `{
	"hello": "world"
}
`,
		},
		{
			name:  "invalid json",
			input: BadJSON{Ch: make(chan int)},
			want:  "error printing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaybePrettyJSON(tt.input)
			if got != tt.want {
				t.Errorf("MaybePrettyJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecureJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name: "mask sensitive data",
			input: TestUser{
				Username: "john",
				Password: "secret123",
				APIKey:   "abcdef123456",
			},
			want: `{
	"username": "john",
	"password": "****",
	"api_key": "********************************"
}
`,
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   BadJSON{Ch: make(chan int)},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SecureJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecureJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && normalizeJSON(got) != normalizeJSON(tt.want) {
				t.Errorf("SecureJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaybeSecureJSON(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name: "valid sensitive data",
			input: TestUser{
				Username: "john",
				Password: "secret123",
				APIKey:   "abcdef123456",
			},
			want: `{
	"username": "john",
	"password": "****",
	"api_key": "********************************"
}
`,
		},
		{
			name:  "invalid json",
			input: BadJSON{Ch: make(chan int)},
			want:  "error printing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaybeSecureJSON(tt.input)
			if normalizeJSON(got) != normalizeJSON(tt.want) {
				t.Errorf("MaybeSecureJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveJSONFile(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "valid json",
			input: map[string]string{
				"hello": "world",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   BadJSON{Ch: make(chan int)},
			wantErr: false, // NOTE: This doesn't error because MaybeSecureJSON handles the error
		},
	}

	tempDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(tempDir, tt.name+".json")
			err := SaveJSONFile(filename, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify file exists and content
			if !tt.wantErr {
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					t.Errorf("SaveJSONFile() file not created")
				}
			}
		})
	}
}

func TestSaveSecureJSONFile(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "valid sensitive data",
			input: TestUser{
				Username: "john",
				Password: "secret123",
				APIKey:   "abcdef123456",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   BadJSON{Ch: make(chan int)},
			wantErr: false, // NOTE: This doesn't error because MaybeSecureJSON handles the error
		},
	}

	tempDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(tempDir, tt.name+".json")
			err := SaveSecureJSONFile(filename, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveSecureJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify file exists
			if !tt.wantErr {
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					t.Errorf("SaveSecureJSONFile() file not created")
				}
			}
		})
	}
}

// Helper function to normalize JSON strings for comparison
func normalizeJSON(s string) string {
	s = strings.TrimSpace(s)
	var temp interface{}
	if err := json.Unmarshal([]byte(s), &temp); err != nil {
		return s // Return original if not valid JSON
	}
	normalized, _ := json.MarshalIndent(temp, "", "\t")
	return string(normalized)
}

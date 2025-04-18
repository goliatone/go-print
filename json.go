package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
)

const (
	empty              = ""
	tab                = "\t"
	unsupportedMessage = "unsupported field type"
)

// PrettyJSON will pretty print as a JSON string
func PrettyJSON(data any) (string, error) {

	safeData := safeToJSON(data)

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(safeData)
	if err != nil {
		return empty, err
	}
	return buffer.String(), nil
}

// MaybePrettyJSON will return a JSON string, in case
// of a decofing error happening it will return the mesasge:
// error printing
func MaybePrettyJSON(data any) string {
	out, err := PrettyJSON(data)
	if err != nil {
		return fmt.Sprintf("error printing: %s", err)
	}
	return out
}

func MaybeHighlightJSON(data any) string {
	out, err := HighlightJSON(data)
	if err != nil {
		return fmt.Sprintf("error printing: %s", err)
	}
	return out
}

func HighlightJSON(data any) (string, error) {
	out, err := PrettyJSON(data)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = quick.Highlight(&buf, out, "json", "terminal16m", "nord")
	if err != nil {
		return "", fmt.Errorf("error highlighting: %w", err)
	}

	return buf.String(), nil
}

func MaybeSecureHighlightJSON(data any) string {
	out, err := SecureHighlightJSON(data)
	if err != nil {
		return fmt.Sprintf("error printing: %s", err)
	}
	return out
}

func SecureHighlightJSON(data any) (string, error) {
	out, err := SecureJSON(data)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = quick.Highlight(&buf, out, "json", "terminal16m", "nord")
	if err != nil {
		return "", fmt.Errorf("error highlighting: %w", err)
	}

	return buf.String(), nil
}

func SecureJSON(data any) (string, error) {
	maskedData, err := PrintMasker.Mask(data)
	if err != nil {
		return "", fmt.Errorf("error masking data: %w", err)
	}
	out, err := PrettyJSON(maskedData)
	if err != nil {
		return "", fmt.Errorf("error printing data: %w", err)
	}
	return out, nil
}

func MaybeSecureJSON(data any) string {
	out, err := SecureJSON(data)
	if err != nil {
		fmt.Printf("secure JSON error parse: %s", data)
		return "error printing"
	}
	return out
}

// SaveJSONFile will create a new file with content
func SaveJSONFile(name string, data any) error {
	str := MaybePrettyJSON(data)
	b := []byte(str)
	return os.WriteFile(name, b, 0644)
}

// SaveSecureJSONFile will create a new file with content
func SaveSecureJSONFile(name string, data any) error {
	str := MaybeSecureJSON(data)
	b := []byte(str)
	return os.WriteFile(name, b, 0644)
}

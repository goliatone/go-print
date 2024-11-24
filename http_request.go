package print

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type RequestJSON struct {
	Method        string      `json:"method"`
	URL           *url.URL    `json:"url"`
	Header        http.Header `json:"header"`
	Body          string      `json:"body,omitempty"`
	ContentLength int64       `json:"content_length"`
}

// PrintHTTPRequest will print an http.Request instance
// as a JSON struct.
// Note, if we read the body using io.ReadAll, the body
// will be drained. So, we either call PrintHTTPRequest
// before we read the body (since we reset it), or we
// get an error message
func PrintHTTPRequest(req *http.Request) string {
	if req == nil {
		return "nil"
	}

	body := ""
	if NotInterfaceNil(req.Body) {
		if bodyBytes, err := io.ReadAll(req.Body); err != nil {
			if req.ContentLength != 0 {
				body = fmt.Sprintf("Error reading body: %s", err)
			}
		} else {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			body = string(bodyBytes)
		}
	}

	r := RequestJSON{
		Method:        req.Method,
		URL:           req.URL,
		Header:        req.Header,
		Body:          body,
		ContentLength: req.ContentLength,
	}

	return MaybeSecureJSON(r)
}

type ResponseJSON struct {
	Status        string       `json:"status"`
	StatusCode    int          `json:"status_code"`
	Header        http.Header  `json:"header"`
	Body          string       `json:"body,omitempty"`
	ContentLength int64        `json:"content_length"`
	Request       *RequestJSON `json:"request"`
}

func PrintHTTPResponse(resp *http.Response) string {
	if resp == nil {
		return "nil"
	}

	body, reset, bodyBytes := bodyAsString(resp.Body, resp.ContentLength)
	if reset {
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	r := &ResponseJSON{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		Body:          body,
		ContentLength: resp.ContentLength,
	}

	if resp.Request != nil {
		body, reset, bodyBytes = bodyAsString(resp.Request.Body, resp.Request.ContentLength)
		if reset {
			resp.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		r.Request = &RequestJSON{
			Method:        resp.Request.Method,
			URL:           resp.Request.URL,
			Header:        resp.Request.Header,
			Body:          body,
			ContentLength: resp.Request.ContentLength,
		}
	}

	return MaybeSecureJSON(r)
}

func bodyAsString(reader io.ReadCloser, len int64) (string, bool, []byte) {
	body := ""
	reset := false
	bodyBytes := []byte{}
	var err error
	if NotInterfaceNil(reader) {
		if bodyBytes, err = io.ReadAll(reader); err != nil {
			if len != 0 {
				body = fmt.Sprintf("Error reading body: %s", err)
			}
		} else {
			reset = true
			body = string(bodyBytes)
		}
	}

	return body, reset, bodyBytes
}

// IsInterfaceNil will check if an interface is nil
func IsInterfaceNil(i any) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Pointer, reflect.UnsafePointer, reflect.Interface,
		reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	default:
		v := reflect.ValueOf(i)
		k := v.Kind()
		fmt.Printf("IsInterfaceNil default: %s %s\n", k.String(), v.String())
	}

	return false
}

// NotInterfaceNil will negate check
func NotInterfaceNil(i any) bool {
	return !IsInterfaceNil(i)
}

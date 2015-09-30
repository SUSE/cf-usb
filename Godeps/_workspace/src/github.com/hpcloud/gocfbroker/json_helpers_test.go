package gocfbroker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	writeJSON(w, http.StatusTeapot, map[string]int{"test": 27})

	if w.Code != http.StatusTeapot {
		t.Error("It should return http teapot:", w.Code)
	}
	// JSON Marshal adds newline at the end
	if str := strings.TrimSpace(w.Body.String()); str != `{"test":27}` {
		t.Error("Expected a correct JSON representation in the body:", str)
	}
}

func TestWriteJSONEmpty(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	writeJSON(w, http.StatusTeapot, nil)

	if w.Code != http.StatusTeapot {
		t.Error("It should return http teapot:", w.Code)
	}

	if str := w.Body.String(); str != `{}` {
		t.Error("Expected a correct JSON representation in the body:", str)
	}
}

func TestWriteJSONError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	writeJSONError(w, http.StatusTeapot, "reason")

	if w.Code != http.StatusTeapot {
		t.Error("It should return http teapot:", w.Code)
	}
	// JSON Marshal adds newline at the end
	if str := strings.TrimSpace(w.Body.String()); str != `{"description":"reason"}` {
		t.Error("Expected a correct JSON representation in the body:", str)
	}
}

func TestWriteJSONErrorEmpty(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	writeJSONError(w, http.StatusTeapot, "")

	if w.Code != http.StatusTeapot {
		t.Error("It should return http teapot:", w.Code)
	}

	if str := w.Body.String(); str != `{}` {
		t.Error("Expected a correct JSON representation in the body:", str)
	}
}

type testReadCloser struct {
	io.Reader
	closed bool
}

func (t *testReadCloser) Close() error {
	t.closed = true
	return nil
}

func TestReadJSON(t *testing.T) {
	t.Parallel()

	closer := &testReadCloser{
		Reader: strings.NewReader(`{ "test": true }`),
		closed: false,
	}

	r, _ := http.NewRequest("GET", "/", closer)

	var testStruct = struct {
		Test bool `json:"test"`
	}{}

	err := readJSON(r, &testStruct)
	if err != nil {
		t.Error(err)
	}

	if !closer.closed {
		t.Error("Must close the body or memory leaks occur")
	}

	if !testStruct.Test {
		t.Error("The value should be true from having been deserialized")
	}
}

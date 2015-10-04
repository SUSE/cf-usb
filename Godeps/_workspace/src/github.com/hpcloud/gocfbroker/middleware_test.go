package gocfbroker

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/blang/semver"
	"github.com/julienschmidt/httprouter"
)

func init() {
	logger = log.New(ioutil.Discard, "", 0)
}

var testSaveLogger *log.Logger

func mockMWWriter() *bytes.Buffer {
	buf := &bytes.Buffer{}
	testSaveLogger = logger
	logger = log.New(buf, "", 0)
	return buf
}

func restoreMWWriter() {
	logger = testSaveLogger
	testSaveLogger = nil
}

func TestBasicAuthMiddleware(t *testing.T) {
	t.Parallel()

	b := Broker{
		Options: Options{
			AuthUser:     "user",
			AuthPassword: "password",
		},
	}

	called := false
	handler := b.basicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		called = true
	}))

	unauthed, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, unauthed)

	if w.Code != http.StatusUnauthorized {
		t.Error("Expected 401 code:", w.Code)
	}

	if called {
		t.Error("It should have the internal handler from being called")
	}

	authed, _ := http.NewRequest("GET", "/", nil)
	authed.SetBasicAuth(b.Options.AuthUser, b.Options.AuthPassword)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, authed)

	if w.Code != http.StatusOK {
		t.Error("Expected OK code:", w.Code)
	}

	if !called {
		t.Error("It should have called the internal handler")
	}
}

func TestAPIVersionMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		version string
		request string
		should  bool
	}{
		{version: "2.4", request: "2.1", should: true},
		{version: "2.4", request: "2.5", should: true},
		{version: "2.4", request: "3.5", should: false},
		{version: "2.4", request: "fail", should: false},
		{version: "2.4", request: "", should: true},
	}

	for i, test := range tests {
		semVer, err := semver.Make(test.version + ".0")
		if err != nil {
			t.Errorf("could not create semver: %s %v", test.version, err)
		}
		b := Broker{apiVersion: semVer}

		called := false
		handler := b.apiVersionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			called = true
		}))

		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.Header.Set(headerBroker, test.request)

		handler.ServeHTTP(w, r)

		if test.should && (!called || w.Code != http.StatusOK) {
			t.Errorf("%d) Should have called the handler", i)
		} else if !test.should && (called || w.Code != http.StatusPreconditionFailed) {
			t.Errorf("%d) Should not have called the handler", i)
		}
	}
}

func TestErrorMiddleware(t *testing.T) {
	buf := mockMWWriter()
	defer restoreMWWriter()

	fnNoError := errorMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		return nil
	})

	fnError := errorMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		return errors.New("fail")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	fnNoError(w, r, nil)
	if buf.Len() > 0 {
		t.Error("It should log nothing if there is no error:", buf.String())
	}

	r, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	fnError(w, r, nil)

	splits := strings.Split(buf.String(), "\n")
	if splits[1] != "\terror: fail" {
		t.Error("Should have logged the failure:", splits[1])
	}
}

func TestJSONMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		accept string
		should bool
	}{
		{"application/quicktime; q=0.8, text/xml", false},
		{"application/quicktime; q=0.8, application/json", true},
		{"application/quicktime; q=0.8, */*", true},
		{"application/*; q=0.8, */xml", true},
	}

	for i, test := range tests {
		called := false
		handler := jsonMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		}))

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", test.accept)

		handler.ServeHTTP(w, r)
		if test.should {
			if w.Code != http.StatusOK {
				t.Errorf("%d) It should return http OK: %d", i, w.Code)
			}
			if !called {
				t.Errorf("%d) It should have called the handler", i)
			}
			if ct := w.HeaderMap.Get(headerContentType); ct != "application/json" {
				t.Errorf("%d) Wrong header for Content-Type: %v", i, ct)
			}
		} else {
			if w.Code != http.StatusNotAcceptable {
				t.Errorf("%d) It should return http not acceptable: %d", i, w.Code)
			}
			if called {
				t.Errorf("%d) It should not have called the handler", i)
			}
		}
	}
}

func TestAsyncMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		query  string
		should bool
	}{
		{"accepts_incomplete=true", true},
		{"accepts_incomplete=false", false},
		{"accepts_incomplete=1", false},
		{"accepts_incomplete=word", false},
		{"", false},
	}

	for i, test := range tests {
		called := false
		handler := asyncMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			called = true
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		var err error
		if r.URL, err = url.Parse("/?" + test.query); err != nil {
			t.Error(err)
		}

		handler(w, r, nil)
		if test.should {
			if w.Code != http.StatusOK {
				t.Errorf("%d) It should return http OK: %d", i, w.Code)
			}
			if !called {
				t.Errorf("%d) It should have called the handler", i)
			}
		} else {
			if w.Code != statusUnprocessableEntity {
				t.Errorf("%d) It should return http unprocessable entity: %d", i, w.Code)
			}
			if s := w.Body.String(); s != jsonErrAsync {
				t.Errorf("%d) Body was wrong: %s", i, s)
			}
			if called {
				t.Errorf("%d) It should not have called the handler", i)
			}
		}
	}
}

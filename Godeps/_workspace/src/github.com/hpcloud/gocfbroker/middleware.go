package gocfbroker

import (
	"crypto/subtle"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/blang/semver"
	"github.com/julienschmidt/httprouter"
	logging "github.com/unrolled/logger"
	"github.com/unrolled/recovery"
)

const (
	headerBroker           = "x-broker-api-version"
	headerAccept           = "Accept"
	headerContentType      = "Content-Type"
	queryAcceptsIncomplete = "accepts_incomplete"
)

// basicAuthMiddleware checks to ensure that basic auth is used and that the
// username and password is correct.
func (b *Broker) basicAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usr, pwd, ok := r.BasicAuth()
		if !ok || usr != b.Options.AuthUser || subtle.ConstantTimeCompare([]byte(pwd), []byte(b.Options.AuthPassword)) != 1 {
			w.Header().Set("Connection", "close")
			writeJSONError(w, http.StatusUnauthorized, "Not Authorized")
			return
		}

		h.ServeHTTP(w, r)
	})
}

// apiVersionMiddleware checks the "x-broker-api-version" header to ensure
// compatibility with the current service broker.
func (b *Broker) apiVersionMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiVersion := r.Header.Get(headerBroker)

		// Ensure patch is specified on the incoming request as it is a
		// requirment of the semver library
		if strings.Count(apiVersion, ".") == 1 {
			apiVersion += ".0"
		}

		// Ignore header missing.
		if len(apiVersion) > 0 {
			reqSemver, err := semver.Make(apiVersion)
			if err != nil || reqSemver.Major != b.apiVersion.Major {
				writeJSONError(w, http.StatusPreconditionFailed, "Invalid API version")
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

// jsonMiddleware ensures that the request's accept header includes
// application/json in some way. It also sets the response header to use json.
// Panics if it can not write to the response.
func jsonMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := strings.ToLower(r.Header.Get(headerAccept))
		if len(accept) > 0 {
			if !strings.Contains(accept, "*/*") && !strings.Contains(accept, "application/json") && !strings.Contains(accept, "application/*") {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Header().Set(headerContentType, "text/plain")
				_, _ = io.WriteString(w, "clients must be able to accept application/json content-type, check Accept header")
				return
			}
		}

		w.Header().Set(headerContentType, "application/json")
		h.ServeHTTP(w, r)
	})
}

// asyncMiddleware checks that there is an accepts_incomplete=true on the request
// when we get async bind/unbind this can be a regular middleware instead of
// being passed the router
func asyncMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h(w, r, p)
	}
}

// apiHandler is the type of handler used throughout the API
type apiHandler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error

// errorMiddleware is a middleware that handles errors in a generic way
func errorMiddleware(handler apiHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := handler(w, r, p)
		if err == nil {
			return
		}

		switch err.(type) {
		default:
			logger.Printf("%6s \"%s\"\n\terror: %v\n",
				r.Method,
				r.RequestURI,
				err,
			)
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func loggingMW(h http.Handler) http.Handler {
	return logging.New(logging.Options{
		Out:         logOutput,
		OutputFlags: log.Ldate | log.Ltime,
	}).Handler(h)
}

func recoverMW(h http.Handler) http.Handler {
	return recovery.New(recovery.Options{
		Out:              logOutput,
		IncludeFullStack: false,
		StackSize:        8 * 1024,
		OutputFlags:      log.Ldate | log.Ltime,
	}).Handler(h)
}

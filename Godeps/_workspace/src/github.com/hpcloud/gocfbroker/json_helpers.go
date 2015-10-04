package gocfbroker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/oxtoacart/bpool"
)

var (
	emptyJSONObject = []byte("{}")
	bufferPoolJSON  = bpool.NewBufferPool(10)
)

// writeJSON writes a status and a json object to the response
// logging as en error handling mechanism. If obj is nil, it writes "{}" to the
// response.
func writeJSON(w http.ResponseWriter, status int, obj interface{}) {
	if obj == nil {
		w.WriteHeader(status)
		if _, err := w.Write(emptyJSONObject); err != nil {
			logger.Println("failed to write empty json response:", err)
		}
		return
	}

	buf := bufferPoolJSON.Get()
	defer bufferPoolJSON.Put(buf)

	enc := json.NewEncoder(buf)
	err := enc.Encode(obj)
	if err != nil {
		logger.Printf("failed to serialize json body: %#v - %v", obj, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, err = io.Copy(w, buf)
	if err != nil {
		log.Println("failed to write json response:", err)
	}
}

// writeJSONError writes an error to the ResponseWriter in the expected format:
//   { "description": "reason" }
func writeJSONError(w http.ResponseWriter, status int, description string) {
	if len(description) == 0 {
		writeJSON(w, status, nil)
		return
	}

	errResponse := struct {
		Description string `json:"description"`
	}{
		Description: description,
	}

	writeJSON(w, status, errResponse)
}

// readJSON reads json from a request into an object
func readJSON(r *http.Request, dest interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dest)
	// Because dec.Decode will fail if the conn dies mid-way, we don't care
	// if close fails.
	_ = r.Body.Close()
	return err
}

package httpclient

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	assert := assert.New(t)

	client := NewHTTPClient(true)
	reqTest := os.Getenv("REQUEST_TEST")

	if reqTest == "" {
		t.Skip("Skipping test, not all env variables are set:'REQUEST_TEST'")
	}

	req := Request{Verb: "GET", Endpoint: reqTest, APIURL: "", StatusCode: 200}

	response, err := client.Request(req)

	assert.NoError(err)
	assert.NotNil(response)
}

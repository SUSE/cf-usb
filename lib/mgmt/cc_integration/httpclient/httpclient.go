package httpclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//HTTPClient defines a HTTPClient
type HTTPClient interface {
	Request(request Request) ([]byte, error)
}

//BasicAuth holds the basic structure for basic auth
type BasicAuth struct {
	Username string
	Password string
}

//Request defines a request
type Request struct {
	Verb        string
	Endpoint    string
	APIURL      string
	Body        io.ReadSeeker
	Headers     map[string]string
	Credentials *BasicAuth
	StatusCode  int
}

type httpClient struct {
	skipSslValidation bool
}

//NewHTTPClient creates and returns a  HTTPClient
func NewHTTPClient(skipSslValidation bool) HTTPClient {
	return &httpClient{
		skipSslValidation: skipSslValidation,
	}
}

//Request performs a request on a http client
func (client *httpClient) Request(request Request) ([]byte, error) {
	httpResponse, err := client.httpRequest(request)
	if err != nil {
		return nil, err
	}

	return httpResponse, nil
}

func (client *httpClient) httpRequest(req Request) ([]byte, error) {
	request, err := http.NewRequest(req.Verb, req.Endpoint+req.APIURL, req.Body)
	if err != nil {
		return nil, errors.New("Error building request")
	}

	if req.Credentials != nil {
		request.SetBasicAuth(req.Credentials.Username, req.Credentials.Password)
	}

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: client.skipSslValidation},
	}
	httpClient := &http.Client{Transport: tr}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != req.StatusCode {
		return nil, fmt.Errorf("status code: %d, body: %s", response.StatusCode, responseBody)
	}

	return responseBody, nil
}

package httpclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpClient interface {
	Request(request Request) ([]byte, error)
}

type BasicAuth struct {
	Username string
	Password string
}

type Request struct {
	Verb        string
	Endpoint    string
	ApiUrl      string
	Body        io.ReadSeeker
	Headers     map[string]string
	Credentials *BasicAuth
	StatusCode  int
}

type httpClient struct {
	skipTslValidation bool
}

func NewHttpClient(skipTslValidation bool) HttpClient {
	return &httpClient{
		skipTslValidation: skipTslValidation,
	}
}

func (client *httpClient) Request(request Request) ([]byte, error) {
	httpResponse, err := client.httpRequest(request)
	if err != nil {
		return nil, err
	}

	return httpResponse, nil
}

func (client *httpClient) httpRequest(req Request) ([]byte, error) {
	request, err := http.NewRequest(req.Verb, req.Endpoint+req.ApiUrl, req.Body)
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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: client.skipTslValidation},
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
		return nil, errors.New(fmt.Sprintf("status code: %d, body: %s", response.StatusCode, responseBody))
	}

	return responseBody, nil
}

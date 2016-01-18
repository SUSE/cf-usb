// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strings"

	"github.com/go-swagger/go-swagger/client"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/strfmt"
)

// Runtime represents an API client that uses the transport
// to make http requests based on a swagger specification.
type Runtime struct {
	DefaultMediaType      string
	DefaultAuthentication client.AuthInfoWriter
	Consumers             map[string]httpkit.Consumer
	Producers             map[string]httpkit.Producer

	Transport http.RoundTripper
	//Spec      *spec.Document
	Host     string
	BasePath string
	Formats  strfmt.Registry
	Debug    bool

	client *http.Client
}

// New creates a new default runtime for a swagger api client.
func New(host, basePath string, schemes []string) *Runtime {
	var rt Runtime
	rt.DefaultMediaType = httpkit.JSONMime

	// TODO: actually infer this stuff from the spec
	rt.Consumers = map[string]httpkit.Consumer{
		httpkit.JSONMime: httpkit.JSONConsumer(),
	}
	rt.Producers = map[string]httpkit.Producer{
		httpkit.JSONMime: httpkit.JSONProducer(),
	}
	rt.Transport = http.DefaultTransport
	rt.client = http.DefaultClient
	rt.Host = host
	rt.BasePath = basePath
	if !strings.HasPrefix(rt.BasePath, "/") {
		rt.BasePath = "/" + rt.BasePath
	}
	rt.Debug = os.Getenv("DEBUG") == "1"
	if len(schemes) == 0 {
		schemes = append(schemes, "http")
	}
	return &rt
}

// Submit a request and when there is a body on success it will turn that into the result
// all other things are turned into an api error for swagger which retains the status code
func (r *Runtime) Submit(context *client.Operation) (interface{}, error) {
	params, readResponse, auth := context.Params, context.Reader, context.AuthInfo

	request, err := newRequest(context.Method, context.PathPattern, params)
	if err != nil {
		return nil, err
	}

	// TODO: infer most appropriate content type
	request.SetHeaderParam(httpkit.HeaderContentType, r.DefaultMediaType)
	var accept []string
	for k := range r.Consumers {
		accept = append(accept, k)
	}

	request.SetHeaderParam(httpkit.HeaderAccept, accept...)

	if auth == nil && r.DefaultAuthentication != nil {
		auth = r.DefaultAuthentication
	}
	if auth != nil {
		if err := auth.AuthenticateRequest(request, r.Formats); err != nil {
			return nil, err
		}
	}

	req, err := request.BuildHTTP(r.Producers[r.DefaultMediaType], r.Formats)

	// set the scheme
	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}
	schLen := len(context.Schemes)
	if schLen > 0 {
		scheme := context.Schemes[0]
		// prefer https, but skip when not possible
		if scheme != "https" && schLen > 1 {
			for _, sch := range context.Schemes {
				if sch == "https" {
					scheme = sch
					break
				}
			}
		}
		req.URL.Scheme = scheme
	}

	req.URL.Host = r.Host
	req.URL.Path = path.Join(r.BasePath, req.URL.Path)
	if err != nil {
		return nil, err
	}

	r.client.Transport = r.Transport
	if r.Debug {
		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(b))
	}
	res, err := r.client.Do(req) // make requests, by default follows 10 redirects before failing
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if r.Debug {
		b, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(b))
	}
	ct := res.Header.Get(httpkit.HeaderContentType)
	if ct == "" { // this should really really never occur
		ct = r.DefaultMediaType
	}

	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, fmt.Errorf("parse content type: %s", err)
	}
	cons, ok := r.Consumers[mt]
	if !ok {
		// scream about not knowing what to do
		return nil, fmt.Errorf("no consumer: %q", ct)
	}
	return readResponse.ReadResponse(response{res}, cons)
}

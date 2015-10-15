package spec

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-swagger/go-swagger/jsonpointer"
)

// Swagger this is the root document object for the API specification.
// It combines what previously was the Resource Listing and API Declaration (version 1.2 and earlier) together into one document.
//
// For more information: http://goo.gl/8us55a#swagger-object-
type Swagger struct {
	swaggerProps
}

// MarshalJSON marshals this swagger structure to json
func (s Swagger) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.swaggerProps)
}

// UnmarshalJSON unmarshals a swagger spec from json
func (s *Swagger) UnmarshalJSON(data []byte) error {
	var sw Swagger
	if err := json.Unmarshal(data, &sw.swaggerProps); err != nil {
		return err
	}
	*s = sw
	return nil
}

type swaggerProps struct {
	ID                  string                 `json:"id,omitempty"`
	Consumes            []string               `json:"consumes,omitempty"`
	Produces            []string               `json:"produces,omitempty"`
	Schemes             []string               `json:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Swagger             string                 `json:"swagger,omitempty"`
	Info                *Info                  `json:"info,omitempty"`
	Host                string                 `json:"host,omitempty"`
	BasePath            string                 `json:"basePath,omitempty"` // must start with a leading "/"
	Paths               *Paths                 `json:"paths,omitempty"`    // required
	Definitions         Definitions            `json:"definitions,omitempty"`
	Parameters          map[string]Parameter   `json:"parameters,omitempty"`
	Responses           map[string]Response    `json:"responses,omitempty"`
	SecurityDefinitions SecurityDefinitions    `json:"securityDefinitions,omitempty"`
	Security            []map[string][]string  `json:"security,omitempty"`
	Tags                []Tag                  `json:"tags,omitempty"`
	ExternalDocs        *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Dependencies represent a dependencies property
type Dependencies map[string]SchemaOrStringArray

// SchemaOrBool represents a schema or boolean value, is biased towards true for the boolean property
type SchemaOrBool struct {
	Allows bool
	Schema *Schema
}

// JSONLookup implements an interface to customize json pointer lookup
func (s SchemaOrBool) JSONLookup(token string) (interface{}, error) {
	if token == "allows" {
		return s.Allows, nil
	}
	r, _, err := jsonpointer.GetForToken(s.Schema, token)
	return r, err
}

var jsTrue = []byte("true")
var jsFalse = []byte("false")

// MarshalJSON convert this object to JSON
func (s SchemaOrBool) MarshalJSON() ([]byte, error) {
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}

	if s.Schema == nil && !s.Allows {
		return jsFalse, nil
	}
	return jsTrue, nil
}

// UnmarshalJSON converts this bool or schema object from a JSON structure
func (s *SchemaOrBool) UnmarshalJSON(data []byte) error {
	var nw SchemaOrBool
	if len(data) < 4 {
		return nil
	}
	if data[0] == '{' {
		var sch Schema
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		if reflect.DeepEqual(Schema{}, sch) {
			return nil
		}
		nw.Schema = &sch
	}
	nw.Allows = !(data[0] == 'f' && data[1] == 'a' && data[2] == 'l' && data[3] == 's' && data[4] == 'e')
	*s = nw

	return nil
}

// SchemaOrStringArray represents a schema or a string array
type SchemaOrStringArray struct {
	Schema   *Schema
	Property []string
}

// JSONLookup implements an interface to customize json pointer lookup
func (s SchemaOrStringArray) JSONLookup(token string) (interface{}, error) {
	r, _, err := jsonpointer.GetForToken(s.Schema, token)
	return r, err
}

// MarshalJSON converts this schema object or array into JSON structure
func (s SchemaOrStringArray) MarshalJSON() ([]byte, error) {
	if len(s.Property) > 0 {
		return json.Marshal(s.Property)
	}
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}
	return nil, nil
}

// UnmarshalJSON converts this schema object or array from a JSON structure
func (s *SchemaOrStringArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}
	first := data[0]
	var nw SchemaOrStringArray
	if first == '{' {
		var sch Schema
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		if reflect.DeepEqual(Schema{}, sch) {
			return nil
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.Property); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}

// Definitions contains the models explicitly defined in this spec
// An object to hold data types that can be consumed and produced by operations.
// These data types can be primitives, arrays or models.
//
// For more information: http://goo.gl/8us55a#definitionsObject
type Definitions map[string]Schema

// SecurityDefinitions a declaration of the security schemes available to be used in the specification.
// This does not enforce the security schemes on the operations and only serves to provide
// the relevant details for each scheme.
//
// For more information: http://goo.gl/8us55a#securityDefinitionsObject
type SecurityDefinitions map[string]*SecurityScheme

// StringOrArray represents a value that can either be a string
// or an array of strings. Mainly here for serialization purposes
type StringOrArray []string

// Contains returns true when the value is contained in the slice
func (s StringOrArray) Contains(value string) bool {
	for _, str := range s {
		if str == value {
			return true
		}
	}
	return false
}

// JSONLookup implements an interface to customize json pointer lookup
func (s SchemaOrArray) JSONLookup(token string) (interface{}, error) {
	if _, err := strconv.Atoi(token); err == nil {
		r, _, err := jsonpointer.GetForToken(s.Schemas, token)
		return r, err
	}
	r, _, err := jsonpointer.GetForToken(s.Schema, token)
	return r, err
}

// UnmarshalJSON unmarshals this string or array object from a JSON array or JSON string
func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}

	if data[0] == '[' {
		var parsed []string
		if err := json.Unmarshal(data, &parsed); err != nil {
			return err
		}
		*s = StringOrArray(parsed)
		return nil
	}

	var single interface{}
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	if single == nil {
		return nil
	}
	switch single.(type) {
	case string:
		*s = StringOrArray([]string{single.(string)})
		return nil
	default:
		return fmt.Errorf("only string or array is allowed, not %T", single)
	}
}

// MarshalJSON converts this string or array to a JSON array or JSON string
func (s StringOrArray) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal([]string(s)[0])
	}
	return json.Marshal([]string(s))
}

// SchemaOrArray represents a value that can either be a Schema
// or an array of Schema. Mainly here for serialization purposes
type SchemaOrArray struct {
	Schema  *Schema
	Schemas []Schema
}

// Len returns the number of schemas in this property
func (s SchemaOrArray) Len() int {
	if s.Schema != nil {
		return 1
	}
	return len(s.Schemas)
}

// ContainsType returns true when one of the schemas is of the specified type
func (s *SchemaOrArray) ContainsType(name string) bool {
	if s.Schema != nil {
		return s.Schema.Type != nil && s.Schema.Type.Contains(name)
	}
	return false
}

// MarshalJSON converts this schema object or array into JSON structure
func (s SchemaOrArray) MarshalJSON() ([]byte, error) {
	if len(s.Schemas) > 0 {
		return json.Marshal(s.Schemas)
	}
	return json.Marshal(s.Schema)
}

// UnmarshalJSON converts this schema object or array from a JSON structure
func (s *SchemaOrArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}
	first := data[0]
	var nw SchemaOrArray
	if first == '{' {
		var sch Schema
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		if reflect.DeepEqual(Schema{}, sch) {
			return nil
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.Schemas); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}

// vim:set ft=go noet sts=2 sw=2 ts=2:

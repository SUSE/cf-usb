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

package middleware

import (
	"encoding"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/internal/validate"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/strfmt"
	"github.com/go-swagger/go-swagger/swag"
)

const defaultMaxMemory = 32 << 20

var textUnmarshalType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func newUntypedParamBinder(param spec.Parameter, spec *spec.Swagger, formats strfmt.Registry) *untypedParamBinder {
	binder := new(untypedParamBinder)
	binder.Name = param.Name
	binder.parameter = &param
	binder.formats = formats
	if param.In != "body" {
		binder.validator = validate.NewParamValidator(&param, formats)
	} else {
		binder.validator = validate.NewSchemaValidator(param.Schema, spec, param.Name, formats)
	}

	return binder
}

type untypedParamBinder struct {
	parameter *spec.Parameter
	formats   strfmt.Registry
	Name      string
	validator validate.EntityValidator
}

func (p *untypedParamBinder) Type() reflect.Type {
	return p.typeForSchema(p.parameter.Type, p.parameter.Format, p.parameter.Items)
}

func (p *untypedParamBinder) typeForSchema(tpe, format string, items *spec.Items) reflect.Type {
	switch tpe {
	case "boolean":
		return reflect.TypeOf(true)

	case "string":
		if tt, ok := p.formats.GetType(format); ok {
			return tt
		}
		return reflect.TypeOf("")

	case "integer":
		switch format {
		case "int8":
			return reflect.TypeOf(int8(0))
		case "int16":
			return reflect.TypeOf(int16(0))
		case "int32":
			return reflect.TypeOf(int32(0))
		case "int64":
			return reflect.TypeOf(int64(0))
		}

	case "number":
		switch format {
		case "float":
			return reflect.TypeOf(float32(0))
		case "double":
			return reflect.TypeOf(float64(0))
		}

	case "array":
		if items == nil {
			return nil
		}
		itemsType := p.typeForSchema(items.Type, items.Format, items.Items)
		if itemsType == nil {
			return nil
		}
		return reflect.MakeSlice(reflect.SliceOf(itemsType), 0, 0).Type()

	case "file":
		return reflect.TypeOf(&httpkit.File{}).Elem()

	case "object":
		return reflect.TypeOf(map[string]interface{}{})
	}
	return nil
}

func (p *untypedParamBinder) allowsMulti() bool {
	return p.parameter.In == "query" || p.parameter.In == "formData"
}

func (p *untypedParamBinder) readValue(values interface{}, target reflect.Value) ([]string, bool, error) {
	name, in, cf, tpe := p.parameter.Name, p.parameter.In, p.parameter.CollectionFormat, p.parameter.Type
	if tpe == "array" {
		if cf == "multi" {
			if !p.allowsMulti() {
				return nil, false, errors.InvalidCollectionFormat(name, in, cf)
			}
			return values.(url.Values)[name], false, nil
		}

		v := httpkit.ReadSingleValue(values.(httpkit.Gettable), name)
		return p.readFormattedSliceFieldValue(v, target)
	}

	v := httpkit.ReadSingleValue(values.(httpkit.Gettable), name)
	if v == "" {
		return nil, false, nil
	}
	return []string{v}, false, nil
}

func (p *untypedParamBinder) Bind(request *http.Request, routeParams RouteParams, consumer httpkit.Consumer, target reflect.Value) error {
	// fmt.Println("binding", p.name, "as", p.Type())
	switch p.parameter.In {
	case "query":
		data, custom, err := p.readValue(request.URL.Query(), target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}

		return p.bindValue(data, target)

	case "header":
		data, custom, err := p.readValue(request.Header, target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data, target)

	case "path":
		data, custom, err := p.readValue(routeParams, target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data, target)

	case "formData":
		var err error
		var mt string

		mt, _, e := httpkit.ContentType(request.Header)
		if e != nil {
			// because of the interface conversion go thinks the error is not nil
			// so we first check for nil and then set the err var if it's not nil
			err = e
		}

		if err != nil {
			return errors.InvalidContentType("", []string{"multipart/form-data", "application/x-www-form-urlencoded"})
		}

		if mt != "multipart/form-data" && mt != "application/x-www-form-urlencoded" {
			return errors.InvalidContentType(mt, []string{"multipart/form-data", "application/x-www-form-urlencoded"})
		}

		if mt == "multipart/form-data" {
			if err := request.ParseMultipartForm(defaultMaxMemory); err != nil {
				return errors.NewParseError(p.Name, p.parameter.In, "", err)
			}
		}

		if err := request.ParseForm(); err != nil {
			return errors.NewParseError(p.Name, p.parameter.In, "", err)
		}

		if p.parameter.Type == "file" {
			file, header, err := request.FormFile(p.parameter.Name)
			if err != nil {
				return errors.NewParseError(p.Name, p.parameter.In, "", err)
			}
			target.Set(reflect.ValueOf(httpkit.File{Data: file, Header: header}))
			return nil
		}

		if request.MultipartForm != nil {
			data, custom, err := p.readValue(url.Values(request.MultipartForm.Value), target)
			if err != nil {
				return err
			}
			if custom {
				return nil
			}
			return p.bindValue(data, target)
		}
		data, custom, err := p.readValue(url.Values(request.PostForm), target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data, target)

	case "body":
		newValue := reflect.New(target.Type())
		if err := consumer.Consume(request.Body, newValue.Interface()); err != nil {
			if err == io.EOF && p.parameter.Default != nil {
				target.Set(reflect.ValueOf(p.parameter.Default))
				return nil
			}
			tpe := p.parameter.Type
			if p.parameter.Format != "" {
				tpe = p.parameter.Format
			}
			return errors.InvalidType(p.Name, p.parameter.In, tpe, nil)
		}
		target.Set(reflect.Indirect(newValue))
		return nil
	default:
		return errors.New(500, fmt.Sprintf("invalid parameter location %q", p.parameter.In))
	}
}

func (p *untypedParamBinder) bindValue(data []string, target reflect.Value) error {
	if p.parameter.Type == "array" {
		return p.setSliceFieldValue(target, p.parameter.Default, data)
	}
	var d string
	if len(data) > 0 {
		d = data[0]
	}
	return p.setFieldValue(target, p.parameter.Default, d)
}

func (p *untypedParamBinder) setFieldValue(target reflect.Value, defaultValue interface{}, data string) error {
	tpe := p.parameter.Type
	if p.parameter.Format != "" {
		tpe = p.parameter.Format
	}

	if data == "" && p.parameter.Required && p.parameter.Default == nil {
		return errors.Required(p.Name, p.parameter.In)
	}

	ok, err := p.tryUnmarshaler(target, defaultValue, data)
	if err != nil {
		return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
	}
	if ok {
		return nil
	}

	defVal := reflect.Zero(target.Type())
	if defaultValue != nil {
		defVal = reflect.ValueOf(defaultValue)
	}

	if tpe == "byte" {
		if data == "" {
			target.SetBytes(defVal.Bytes())
			return nil
		}

		b, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			b, err = base64.URLEncoding.DecodeString(data)
			if err != nil {
				return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
			}
		}
		target.SetBytes(b)
		return nil
	}

	switch target.Kind() {
	case reflect.Bool:
		if data == "" {
			target.SetBool(defVal.Bool())
			return nil
		}
		b, err := swag.ConvertBool(data)
		if err != nil {
			return err
		}
		target.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if data == "" {
			target.SetInt(defVal.Int())
			return nil
		}
		i, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
		}
		if target.OverflowInt(i) {
			return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
		}

		target.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if data == "" {
			target.SetUint(defVal.Uint())
			return nil
		}
		u, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
		}
		if target.OverflowUint(u) {
			return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
		}
		target.SetUint(u)

	case reflect.Float32, reflect.Float64:
		if data == "" {
			target.SetFloat(defVal.Float())
			return nil
		}
		f, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
		}
		if target.OverflowFloat(f) {
			return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
		}
		target.SetFloat(f)

	case reflect.String:
		value := data
		if value == "" {
			value = defVal.String()
		}
		// validate string
		target.SetString(value)

	case reflect.Ptr:
		if data == "" && defVal.Kind() == reflect.Ptr {
			target.Set(defVal)
			return nil
		}
		newVal := reflect.New(target.Type().Elem())
		if err := p.setFieldValue(reflect.Indirect(newVal), defVal, data); err != nil {
			return err
		}
		target.Set(newVal)

	default:
		return errors.InvalidType(p.Name, p.parameter.In, tpe, data)
	}
	return nil
}

func (p *untypedParamBinder) tryUnmarshaler(target reflect.Value, defaultValue interface{}, data string) (bool, error) {
	// When a type implements encoding.TextUnmarshaler we'll use that instead of reflecting some more
	if reflect.PtrTo(target.Type()).Implements(textUnmarshalType) {
		if defaultValue != nil && len(data) == 0 {
			target.Set(reflect.ValueOf(defaultValue))
			return true, nil
		}
		value := reflect.New(target.Type())
		if err := value.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(data)); err != nil {
			return true, err
		}
		target.Set(reflect.Indirect(value))
		return true, nil
	}
	return false, nil
}

func (p *untypedParamBinder) readFormattedSliceFieldValue(data string, target reflect.Value) ([]string, bool, error) {
	ok, err := p.tryUnmarshaler(target, p.parameter.Default, data)
	if err != nil {
		return nil, true, err
	}
	if ok {
		return nil, true, nil
	}

	return swag.SplitByFormat(data, p.parameter.CollectionFormat), false, nil
}

func (p *untypedParamBinder) setSliceFieldValue(target reflect.Value, defaultValue interface{}, data []string) error {
	if len(data) == 0 && p.parameter.Required && p.parameter.Default == nil {
		return errors.Required(p.Name, p.parameter.In)
	}
	defVal := reflect.Zero(target.Type())
	if defaultValue != nil {
		defVal = reflect.ValueOf(defaultValue)
	}
	if len(data) == 0 {
		target.Set(defVal)
		return nil
	}

	sz := len(data)
	value := reflect.MakeSlice(reflect.SliceOf(target.Type().Elem()), sz, sz)

	for i := 0; i < sz; i++ {
		if err := p.setFieldValue(value.Index(i), nil, data[i]); err != nil {
			return err
		}
	}

	target.Set(value)

	return nil
}

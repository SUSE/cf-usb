// Code generated by go-bindata.
// sources:
// swagger-spec/api.json
// swagger-spec/service_manager_api.json
// DO NOT EDIT!

package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _swaggerSpecApiJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x5b\xcd\x72\xdb\x36\x10\x3e\xcb\x4f\x81\x61\x7b\xd2\x38\xa2\x92\x26\x9d\x69\x4e\x55\xac\xb4\xa3\x49\x26\xc9\xd4\xcd\xc9\xf5\xd8\x10\x09\x49\x48\x49\x82\x01\xc0\x28\xae\x47\xef\x5e\xfc\xf1\x57\x14\x09\x4a\x4a\x4c\xb9\xbd\x54\x35\xb8\x8b\x5d\xec\x7e\xf8\x76\x41\x30\xf7\x67\x03\x87\xad\xe1\x72\x89\xa8\xf3\x12\x38\xcf\x46\x63\xe7\x5c\x0c\xe1\x68\x41\xc4\xdf\xe2\xe9\xc0\xf1\x11\xf3\x28\x8e\x39\x26\x91\x14\xf9\x18\xe1\x2f\x88\x32\x18\x00\x86\xe8\x17\xec\x21\x30\xa7\xe4\x6f\x44\x41\x08\x23\xb8\x44\x21\x8a\x38\x98\x7c\x98\xc9\x69\x06\x8e\x94\x34\x7a\xe3\xd1\x78\xf4\x54\x8f\x72\xcc\x03\xa4\xe6\xba\x7c\x55\x54\x83\x31\x36\x02\x88\x86\xec\xfd\xe2\x52\x1b\x90\x92\x2b\xce\x63\xf6\xd2\x75\x97\x98\xaf\x92\xf9\xc8\x23\xa1\xbb\x8a\xbd\x80\x24\xbe\xeb\x2d\x9e\x24\x6c\xae\x15\x3d\x12\x71\xe8\x71\xe3\xfa\xc0\x41\x21\xc4\x81\xd2\x17\x52\x3e\xfa\xf2\xeb\x2a\x46\x52\xd9\x11\x4f\x37\x4a\x23\x10\x06\x22\x86\x32\x8d\x08\x86\xca\xe0\x24\x86\xde\x0a\x01\x13\x10\xf1\x20\xa1\x41\xea\x88\xf0\x63\xbd\x5e\x8f\xa0\x12\x19\x11\xba\x74\xcd\x2c\xcc\x7d\x3b\xbb\x78\xfd\xee\xf2\xf5\x13\xa1\x37\x5a\xf1\x30\x50\x86\xce\x94\x2d\x87\x21\x2f\xa1\x98\xdf\x4d\xd1\x02\x47\x58\xc6\x93\xa5\x31\x9e\x24\x7c\x45\x28\xfe\x07\x9a\x28\x6b\x5f\xf8\x5d\xac\x7c\x11\x71\x79\x83\xee\x8c\x1f\x58\x45\x73\x85\xa0\x2f\x52\x76\x5e\xf6\xb9\x34\x8b\x7e\x56\xcd\xde\x64\x02\xfe\x14\xe9\x8a\xea\x1c\x13\x02\x57\xb5\xee\x5c\x5d\x0b\xc9\x6b\x29\x2a\x02\xcc\x92\x10\x49\xc7\xaf\xa4\x24\x8c\x63\xb1\x76\x25\xe7\x7e\x62\x44\x4e\xab\xe4\x62\x4a\xfc\xc4\xb3\x90\x83\x7c\x95\x45\xc1\x4d\x62\x1f\x72\x74\x23\xe4\x60\x40\x96\x59\x1c\x62\xc2\xb2\x94\x6e\xaf\x48\xe9\x30\xc0\x57\x19\x12\x8d\xfe\x5f\x26\x06\x03\x87\xc4\x88\x2a\xdb\x33\x5f\xaa\x68\x33\x17\xc6\x8a\x91\xa1\x88\xc5\x44\xe6\x30\xb3\x34\x70\x9e\x8d\xc7\xf9\x5f\x5b\x96\x2f\x13\x4f\xac\x90\x2d\x92\x00\xa4\xca\x8e\x16\xdd\x9c\x9b\x09\x5e\x34\x4e\xf0\x31\x42\x5f\x63\xe4\x71\xe4\x03\x44\x29\x31\xe9\x94\x82\x4c\x00\x2b\x84\x05\xd5\x1c\x0c\x8c\x53\x1c\x2d\x8d\x21\x99\xc2\xfc\x47\xfd\x77\x93\x41\xdb\x2d\xec\xe2\x81\xb3\x44\xbb\x63\xf8\x3b\xe2\x0c\x48\x69\x1a\xaa\x30\x01\x38\x27\x09\x57\x21\x15\x5b\x74\xb4\x33\x90\x62\xce\x99\xb4\x71\xd4\x10\x36\x05\xe1\x47\x8a\x16\x52\xef\x07\xd7\xcf\x37\x91\x5e\x67\x25\x20\xfd\xc8\x80\x4f\x25\x59\xde\xe0\x88\x71\x18\x79\x79\x60\xda\xb3\xa1\x35\x41\xa6\x09\x44\x72\x00\x34\xc3\x4d\x09\x99\x2a\x89\x59\x66\xd1\x08\xc6\x90\x0a\x96\x10\xc4\xaa\xb6\xe4\xf5\xde\x19\x7b\xff\xc6\x2a\x46\x90\x52\x78\x97\x49\x0a\xd6\xe2\x28\x64\x45\xc9\x5d\xb9\xac\x44\x2c\x8d\x72\x1a\xe6\x5e\xe4\xf7\xbc\x9d\x95\x2e\x28\x12\x14\x93\x25\x2c\xcb\xe3\x8e\xbc\x79\x4a\xbc\x9c\xba\xfa\xcc\xa5\x59\x4a\x39\xbf\x1a\xaf\x34\x2e\x15\x7f\x2a\x6e\x00\x4e\xc0\x1c\x01\x6d\xd6\xcf\x94\x74\x6d\x99\x13\x3f\x4b\x9d\x40\xc8\xe7\x04\x53\x24\xbd\xe4\x34\x41\xe9\xf0\x56\xf8\x3a\xa5\xd3\x04\xb3\x19\x85\x4f\x1b\xf2\x9a\xc6\xa8\xba\x82\x0e\xa4\xb1\x03\x68\x55\x7c\x3d\x1f\xff\xd2\xe0\xc7\x05\x89\x16\xa2\xae\xf1\x7e\xf1\x7e\x95\x75\xdc\xfb\xca\xc8\x0d\xf6\x37\x9d\xa9\x48\x94\xfe\x05\x5e\x26\x1a\xb6\xcc\x9e\x82\xf6\xc2\xb1\xf0\xb0\x82\x4a\xd9\x29\xec\x42\xb7\x36\x08\x32\x54\xcc\xa6\x6d\xf8\xad\x8b\x6a\x0b\x1e\x0f\x60\xc5\x03\x11\xf8\xbc\xc1\xf2\x3b\xc2\xc1\x6f\x24\x89\xfc\x3e\x40\xd0\x10\x63\xd2\xd6\xad\x6d\xf3\x62\x4b\xaf\x76\xaa\x88\x3a\x07\x3b\xbc\xd2\x9b\xa9\x89\x78\x2b\x0e\x4d\x7c\x1f\xd4\xeb\x3e\x1c\x43\xb7\x74\x76\xba\xb1\x3b\xa8\xb3\xfb\xb6\x5b\xe4\x74\xb8\x3d\x3d\xcb\x05\x02\xed\x3b\xf7\xd6\x54\x3d\xb6\xee\x39\xf4\x6c\xa7\xba\xb3\x9a\x91\xd9\x04\x07\xe3\x80\x5e\x7e\x0d\x26\x4e\x86\x6d\x3b\x16\x7c\x37\x96\xd3\xd9\x54\xfd\x0f\x42\x50\x1f\xa8\x5b\x8e\x1b\x72\xc6\xc7\x09\xa0\x96\x62\xff\x18\x40\x63\x5e\xdf\x75\x38\x93\x4a\x44\xa0\xaf\x98\x71\x61\x02\xdc\x1a\xfd\xdb\xa6\x4e\xd0\xbc\xc1\x7b\x75\x97\x66\x6b\xe6\x1f\x09\x21\x9f\x13\x44\xf3\x62\x59\x59\x7e\x0b\x1e\xea\x01\x85\x0b\x80\xea\x47\xf5\x33\x11\xee\xe7\xab\x8d\x14\x3e\xee\xbd\xf9\xbf\x6e\x87\x0a\x89\xa5\x0c\x42\x60\x8d\xf9\x4a\x0d\x61\x1f\x0c\x87\x66\x78\x36\x1d\x0e\x5b\xa1\x65\x07\xa7\xdc\xc5\x0e\x44\x33\x9b\x02\xb2\x50\x5e\xb1\x92\xad\xef\xc5\x33\xdf\x05\x44\x27\x44\x5d\xb6\xa7\x8b\x43\xc1\xa5\x8f\x1c\xa7\x82\xaf\xed\x23\x46\x75\xb6\xf6\xc3\x85\x39\x95\x59\xba\x61\x7b\xb2\x28\x23\xaf\x27\x27\x8a\x6f\xb3\x1d\x4e\xe7\x24\x91\x13\x78\x1c\xc0\xa8\x43\xf1\xbf\x55\xf2\xb7\xad\x84\xfc\x41\x4d\xfb\x30\x45\x7e\x01\x03\xd6\x52\xe5\x67\x0f\x58\xe5\x0f\x7e\x39\x2e\x53\xd0\xcb\x37\xe2\x15\x54\xb9\xf7\xf2\x67\x8f\x86\x40\xaa\x15\x08\x7b\x38\x94\x03\x56\x7d\x80\x84\x9d\x1d\xea\x8c\x6b\x7b\x11\x75\x9c\x5b\x39\xad\x2e\xa0\x88\x9c\xff\x4c\x0b\x20\x17\x5d\xa9\xfe\x29\x9e\xd4\x75\x5a\xa1\xf4\xee\xd5\x1b\xf4\x1c\x73\xdb\x9d\x41\x69\x2a\xeb\xb6\xc0\xc6\x01\xdb\x9e\xa0\x80\xc3\x9e\x34\x04\x8f\x64\x67\xe4\x0c\xec\x63\x51\x05\x3b\xd4\x75\x25\xdf\x50\xd7\x27\x41\x30\x55\x53\xf6\xb6\xa6\x3f\xe4\xc9\xfd\xf0\x0b\x6f\x11\xdb\x5e\xd6\x74\xc3\xb5\x76\xb7\xdc\x62\x11\x92\x52\x9b\xaf\xb7\xe5\x4a\xed\x40\x94\x4b\xda\xf0\xd4\x3b\xb4\x06\x25\x95\x43\xaf\x44\xf2\x94\x1c\x7c\x53\x7d\xa4\x7b\x90\x02\x48\xfa\x80\x8d\x0a\xdb\xb8\xf7\xf2\x67\x8f\x7e\x4f\xaa\x95\xfa\x3d\x33\x4f\x53\xc3\xd7\x0d\x46\xfb\x56\x5d\x1b\x3c\xf5\xb2\xd3\x6b\x82\xca\x09\xd5\xb3\x6e\x9d\x9e\xe2\x9f\x72\xa7\xd7\x86\x24\x73\xb1\xdc\x6f\x30\xd5\xdc\x1f\x77\xa3\x46\x1d\x24\xff\x41\xe9\xf1\x7f\xcc\x1f\xfd\x96\x77\x5f\xfa\x34\x77\xbe\xfd\x06\xfd\xfe\x57\xbb\x75\x68\x3a\x5d\x38\x9c\x0d\xb2\x2f\xb4\xfd\xed\x4f\xc6\x4b\x1f\xf7\xa6\x73\x92\xf9\x27\xe1\x86\x36\x5f\x8c\xf9\x95\x8e\xa8\xfe\x3e\xfa\x06\xc6\xf8\x26\xfd\x32\xdf\xc4\x3a\x61\xf3\x6c\x48\x8e\x5c\x1b\xfe\xa5\x12\x40\x1c\x17\x72\x50\x37\x49\x9e\x9f\x7a\x16\xdb\xb6\xd1\xa2\x51\x6e\x32\x58\xf6\xef\x00\x3a\xad\x76\xe7\x51\x44\xc3\xbb\x65\x99\xd8\xb7\x5c\x56\x8d\x99\x9d\x8a\x29\x9c\x42\x1c\xbd\x45\xd1\x52\x6c\xa3\x97\xe0\xa7\x9f\xb3\x51\xf8\xb5\x30\x5a\x32\x32\xc7\x91\x0f\xe7\x01\xaa\x99\x7a\x4e\x48\x80\xd2\x33\xec\xa6\xb8\xc2\x6e\x6e\xd4\x79\xf1\x62\x5c\x5e\x6a\x09\xf6\x36\xc1\xe1\x70\xc9\x6a\x44\x4b\xa7\xa5\xad\xb3\x52\xfd\x16\xd9\x94\x26\x16\x9c\x05\x45\x69\x83\x35\x93\x1b\x58\xd4\x21\x49\x1d\xf6\xbb\xc2\xa8\xc4\x81\x56\xd8\xf9\x36\xd1\xb7\x47\xa4\xf1\xf8\xf8\x30\xec\x0e\x80\x05\x45\xed\xa0\x2d\x27\x49\x15\x8e\xc3\xf7\xfa\x49\xed\xef\xd2\xc7\xb5\x6d\x98\x4e\x4b\x77\x01\xcc\x76\x44\x5a\xfd\x9c\xae\x63\x90\x8f\x4a\x9b\xcd\x7b\xa4\x42\x22\x54\x9c\xfb\x3e\xfe\xf1\xd6\x52\x7e\x9f\x60\x96\xde\x9a\x1d\x95\xa9\x2a\xd5\xcb\x32\x53\x72\x0b\x97\x9b\x87\x1a\xb5\x52\x56\xdb\xa5\x45\x33\x71\xb6\xf9\x37\x00\x00\xff\xff\x56\x07\x4f\x5f\xf2\x37\x00\x00")

func swaggerSpecApiJsonBytes() ([]byte, error) {
	return bindataRead(
		_swaggerSpecApiJson,
		"swagger-spec/api.json",
	)
}

func swaggerSpecApiJson() (*asset, error) {
	bytes, err := swaggerSpecApiJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "swagger-spec/api.json", size: 14322, mode: os.FileMode(420), modTime: time.Unix(1461671462, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _swaggerSpecService_manager_apiJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x59\xdf\x6f\xdb\x36\x10\x7e\x96\xff\x0a\x42\x1b\xb0\x0d\x08\xe2\xb4\x18\xf6\xb0\xb7\x22\x29\x86\x0c\x58\x17\x34\x05\xf6\xd0\x04\x01\x2d\x9d\x6c\xd6\x32\xa9\x91\x94\xbd\x20\xf0\xff\xbe\x23\x45\xea\xa7\x2d\xc9\xcb\x80\xda\x4d\x1e\xea\x5a\xd4\xf1\x78\xf7\x7d\xdf\x9d\xce\xca\xd3\x24\x08\xd5\x86\xce\xe7\x20\xc3\x5f\x49\xf8\xf6\xfc\x22\x3c\xc3\x25\xc6\x13\x81\xd7\x78\x37\x08\xd7\x20\x15\x13\xdc\xdc\x5e\xbf\x31\x77\x83\x50\x33\x9d\x82\x59\xb8\xa4\x9a\xa6\x62\x4e\x6e\x41\xae\x59\x04\xe4\x0f\xca\x29\xba\x22\xef\x6e\xae\x0b\xcb\x18\x54\x24\x59\xa6\x9d\x83\x4f\x0b\x20\x09\x93\x4a\x13\xe7\x95\x88\x84\x68\x5c\xec\x71\x74\x4e\xee\xf8\xa7\x05\x53\xe6\x3b\xd9\xb0\x34\x25\x33\x20\x74\x4d\x59\x4a\x67\x29\x10\xf4\xd1\xe7\x20\x12\x5c\x53\xc6\xf1\xdb\x66\xc1\xa2\xc5\x1d\x97\x39\x57\x04\x4d\xf9\x9c\x28\x16\x03\x79\x14\xb9\x24\xca\xed\xa2\x3c\xb6\xdf\x41\x11\x25\x56\xe0\xa3\xf3\xb7\x57\xd6\xe9\x0a\xb8\xc6\x98\x22\x9a\xd1\x19\x4b\x99\x66\xa0\xce\xef\x78\x38\x09\xb6\x06\x3a\x3c\x50\xe5\x2b\x50\x98\xee\x67\x83\x00\xcd\xb2\x94\x45\xd4\x20\x30\xfd\xa2\x84\xb1\xbb\x37\x76\x99\x14\x71\x1e\x0d\xdb\xa9\x68\x01\x35\x77\x0b\xad\xb3\xd2\x05\xd5\x0b\xe5\x69\x9a\xd2\x18\x31\xd5\x4c\x81\x5b\xc1\xfb\x42\x69\x7f\x11\x84\x22\x03\x69\xdd\x5f\xc7\x86\x8a\xd2\xdc\x21\x67\xf9\xea\x32\xf6\xce\x9b\x95\x20\x88\x24\x01\xc9\x78\xb9\x21\xa3\x92\xae\x40\x23\x9f\x26\xc6\x7b\xb7\x2a\x41\x65\x88\x04\xa8\x32\x80\x20\x7c\x7b\x71\x51\x5d\xf5\x9c\x14\x39\x2e\xc5\x32\x2c\x6c\xb7\x67\x13\xbf\x25\xa1\x79\xaa\x7b\xbc\xcc\x01\xb9\x66\x11\x01\x29\x85\x24\x3e\x0a\x17\x6b\xe0\xe0\xa4\x35\x07\x41\xf8\xbd\x84\xc4\x6c\xfd\x6e\x8a\xee\x19\x67\xc6\x95\x9a\xbe\x37\x0e\xdc\xf9\xc1\x76\x52\xfb\xcf\x7e\x9a\x0f\x1b\x56\x38\xdd\x08\xb9\x54\x19\x8d\xaa\x5c\x5b\xc8\x6b\x3a\xf7\xfc\x99\xcb\xd2\xbe\xf0\xee\x11\x6b\xf1\x13\x49\xa0\x1a\xfe\x2a\x6d\x77\xd3\x73\x69\xad\x08\x87\x0d\xd9\xb4\x4d\x9b\xc4\x78\x12\x38\xae\x55\xfe\x1f\xca\x5d\x0f\x12\xfe\xce\x01\xc3\xf6\x58\x33\x7b\xc0\x4c\xc4\x8f\x61\x05\x7f\xa7\x9a\xbd\x2a\x7e\xbf\xfd\xf3\x83\xa9\x25\xb2\xa1\x58\x1d\x5a\x10\x0b\x81\xdf\xd8\x41\x7d\x0f\xe8\xae\x7a\x5d\xf1\x96\xc9\x17\x59\x7e\x74\x01\xb6\x34\x61\xe2\x66\x12\x0c\x68\x5a\xe6\x50\x50\xd4\x2f\xc3\x37\x3d\x02\x2a\x70\xe9\xa0\x79\x88\x74\xf6\x64\xf1\xd1\x8b\xb1\x25\xaa\x13\x11\xf7\xf4\xa9\xd2\x0a\x8b\xb7\xa5\xd6\xe7\xf0\x5c\xa9\xa3\x87\x21\x9d\xff\x06\xda\x76\xe2\x18\xb0\x9b\xa7\x8a\x24\x98\xbf\xed\xcc\x19\x44\x2c\x61\xc8\xfe\x28\xd1\xd7\x33\x68\x09\xdd\x34\xd3\x7d\x42\x2f\xa3\x23\xd7\x57\xe1\x1e\xdd\xf9\x65\xfd\x98\xd9\xa3\x94\xb6\x4d\x72\x8c\x1c\xfb\xba\xa2\x4f\x18\x9f\x44\x65\xae\xa4\x8d\xd6\xcb\x12\xa7\xf5\x88\xa7\xa7\x48\xf3\x33\x95\x57\x38\x19\x12\xdf\x95\xb5\xaa\xe1\x7f\x60\xab\x3d\x45\xd5\xd9\x8c\x5b\x20\x1e\xa5\x1a\x46\xb4\xaa\x29\xce\x64\x1c\x22\xeb\x6f\xdc\x23\xba\xda\x30\xe2\x19\x7d\x59\x19\x0f\x3e\xa4\xa3\x8e\xed\x71\x4a\xe7\x8c\x74\xe6\x85\x32\xf2\x07\x37\x3a\x1c\xdf\xc0\x50\x31\xf1\x15\x26\x86\x0e\xb5\xff\xbd\x2b\x57\x79\x9c\x50\x5b\x3e\xb0\x10\xa7\x4f\x35\x45\x8d\x9d\x27\x46\xd6\x25\xba\x18\x2c\xca\x13\x9f\x28\x7a\x0b\xf4\xa0\x50\xaa\x7d\x47\x32\xdd\xbc\xf4\x42\x1a\x35\xdf\x8c\xac\x84\xc2\xcb\x60\x31\x74\x26\x9c\x6f\xe2\x39\x75\x42\x65\x60\xe1\x6f\x73\x7a\x94\xea\x9c\xd8\x7f\xf6\x55\x57\x6d\x9b\x73\x1c\x0e\xfd\xba\xf0\x4d\xde\x03\x28\x66\x5f\x30\xe3\x22\xcc\x3a\xd8\x9f\x9d\xe8\xa4\xc0\xa7\x88\x42\x8c\x1f\xec\x0e\x87\xb0\xd2\x54\xe7\xca\x06\x59\x80\x6e\x0c\x33\xf3\xf6\xa8\x06\x7a\x67\x73\xc5\x46\x8b\x3e\x8f\x32\xf0\x7c\x55\x95\x58\x10\x7e\x10\xbc\x86\xe0\x95\xe3\xa0\x5c\x78\xff\x8f\x06\xae\x2a\xb6\xee\xf7\x28\xeb\xa6\x8c\x83\xd8\x83\x27\x35\x6a\x7d\x2a\x07\xc7\xc6\x1b\xb1\xe5\x7c\xc9\xc5\xa6\xde\x2f\xf3\xc8\x9c\x99\xe4\x69\xb5\x96\x60\xab\x45\x70\x5d\xb0\x8d\x28\x5c\x1f\xde\x11\x46\x9d\xa0\x6e\x6a\x4b\x78\x24\x6b\x9a\xe6\xe6\xd5\x68\x46\x36\x4c\x2f\x6a\x22\x2e\x1f\xad\x3f\xde\x36\xde\x9f\x4a\xfc\x5d\xc1\x7f\xd0\x84\x71\xec\x26\x99\x6c\x3d\x87\xdd\x2b\x57\xe3\x8f\xe9\xe2\x7d\xaf\x02\x1e\xe3\x22\x53\x64\x46\xa3\xa5\x99\x15\x8d\x85\x1b\x3f\xd1\x1d\x55\x84\xa9\x9f\xca\x18\x69\x1c\x5b\x51\xd2\xf4\xa6\xab\x0a\xd3\x0b\x34\xac\x54\xa3\x12\x76\x95\x73\xff\x8c\x33\xea\x45\xd5\x81\x62\x6f\xf4\xd0\x01\x6d\x37\x6c\x07\xc5\xb3\xb7\xeb\xba\x83\x5e\x75\xf0\xff\xe8\x60\xc7\xb8\xf1\xda\xf0\x5e\x1b\xde\xb7\x2c\xf4\x67\x75\xbc\xe6\xac\x36\xa0\xee\xa6\xf1\xa1\x3d\xaf\x8a\xf8\xb5\xe9\x3d\x57\x0b\xc5\x88\x78\x20\xd7\x2b\xac\x4e\xcc\x7f\x98\xe5\x78\x57\xe3\x32\x70\x99\xbf\x51\xfb\x34\xf1\xd7\xfa\x8a\x6a\x77\xe7\x97\x9f\x9b\x84\xfa\xa3\xf6\x6a\xa4\x95\x14\x8e\xb3\x93\xed\xbf\x01\x00\x00\xff\xff\x5c\xdd\x9a\x6b\x0a\x1f\x00\x00")

func swaggerSpecService_manager_apiJsonBytes() ([]byte, error) {
	return bindataRead(
		_swaggerSpecService_manager_apiJson,
		"swagger-spec/service_manager_api.json",
	)
}

func swaggerSpecService_manager_apiJson() (*asset, error) {
	bytes, err := swaggerSpecService_manager_apiJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "swagger-spec/service_manager_api.json", size: 7946, mode: os.FileMode(420), modTime: time.Unix(1461151802, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"swagger-spec/api.json": swaggerSpecApiJson,
	"swagger-spec/service_manager_api.json": swaggerSpecService_manager_apiJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"swagger-spec": &bintree{nil, map[string]*bintree{
		"api.json": &bintree{swaggerSpecApiJson, map[string]*bintree{}},
		"service_manager_api.json": &bintree{swaggerSpecService_manager_apiJson, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


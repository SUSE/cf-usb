// Code generated by go-bindata.
// sources:
// swagger-spec/api.json
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

var _swaggerSpecApiJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x5d\x4d\x8f\xdb\xbc\x11\xbe\xe7\x57\x10\x6e\x4f\x46\x62\x6d\xde\xe6\x0d\xd0\x9c\xea\xc4\x69\x61\x24\x48\x02\x6c\x73\x4a\x17\xbb\xb4\x44\xdb\x4c\x65\x49\x91\xa8\xfd\xe8\x62\xff\x7b\x49\x8a\xd2\xea\x83\xb4\x49\x59\xb2\xe5\x98\x3e\x24\x0b\x79\x38\x9a\x19\xce\x3c\x1c\x8e\x86\xd6\xe3\x0b\x40\x3f\xa3\xe4\x0e\xae\x56\x28\x1e\xbd\x03\xa3\x3f\x26\x17\xa3\x97\xd9\x55\x1c\x2c\x43\x7a\x29\xa3\xe1\x57\x3c\x94\xb8\x31\x8e\x08\x0e\x03\x46\xfb\x3d\xc0\xb7\x28\x4e\xa0\x0f\x12\x14\xdf\x62\x17\x81\x45\x1c\xfe\x17\xc5\x60\x03\x03\xb8\x42\x1b\x14\x10\x30\xfd\x36\x17\xfc\x38\x07\x46\x2f\x46\x5f\x4c\x2e\x26\xaf\xcb\xdf\x11\x4c\x7c\xc4\xf9\x5e\xbe\x2f\xb3\x80\x11\xae\x90\xa1\x78\x93\x7c\x5d\x5e\x66\xb7\x64\xf4\x6b\x42\xa2\xe4\x9d\xe3\xac\x30\x59\xa7\x8b\x89\x1b\x6e\x9c\x75\xe4\xfa\x61\xea\x39\xee\xf2\x55\x9a\x2c\xca\xc3\xdd\x30\x20\xd0\x25\x15\xc5\xf8\x17\x68\x03\xb1\xcf\xd9\xd1\x41\x1e\xba\xfd\xc7\x3a\x42\x8c\xd7\xa8\x20\x7b\x2a\xb1\xf1\xe9\xbd\x83\x04\x35\xd9\x04\x70\xc3\x85\x9a\x46\xd0\x5d\x23\xf0\x6c\xcf\x82\x22\x8d\xfd\x5c\x6a\x2a\xf4\xdd\xdd\xdd\x04\x72\xda\x49\x18\xaf\x1c\xc1\x37\x71\x3e\xcf\x3f\x7c\xfc\x72\xf9\xf1\x15\x65\x30\x59\x93\x8d\x5f\x12\xe3\x45\x49\x98\x51\x82\xdc\x34\xc6\xe4\x61\x86\x96\x38\xc0\x6c\x6a\x92\xea\xa4\x4d\x53\xb2\x0e\x63\xfc\x3f\x28\xa6\xad\x26\x2f\x79\x88\xb8\xbc\xd4\xca\x9f\xd0\x43\x5d\x56\xcc\xa7\x6a\x8d\xa0\x47\xdd\xe3\xa5\x42\xd3\xca\x0d\x6a\x44\x75\x97\x99\x4e\xc1\xbf\xa9\x8f\x04\xbb\xd4\xa1\xc4\x3f\x0a\x92\x9a\xc8\x75\x8d\x7e\x5c\xd5\x98\x5d\x09\x66\x74\xae\x93\x74\x83\x92\x0a\x33\xaa\x69\x44\xad\xcc\x07\x3b\x3f\x93\x50\x88\x92\x8f\x89\xe2\xd0\x4b\x5d\xc3\x31\x90\xac\x6b\x56\x77\xd2\xc8\x83\x04\x5d\xd3\x31\xd0\x0f\x57\x4d\xbb\x47\x61\xd2\x74\x42\xa9\xc9\x38\xa3\x04\x90\x75\x11\x5f\x82\xe9\x7f\xea\xd6\xe6\xc3\xc3\x08\xc5\x5c\xd0\xb9\xc7\x86\x67\x72\x7c\x10\x62\x48\xe8\x63\x94\x44\x21\x73\x39\xa9\x34\x9c\xe4\x8f\x8b\x0b\xe5\x97\x32\x91\x2f\x53\x97\x5a\x30\x59\xa6\x3e\xc8\xb9\x8f\xa4\x83\x9f\x9a\xf2\x70\x7e\x7f\x1a\xde\xf0\x7b\x80\xee\x23\xe4\x12\xe4\x01\x14\xc7\x61\xdd\x55\x2b\x43\x13\x1a\x6a\x1b\xb8\x95\x3d\xa7\xcb\x03\x23\x21\x31\x0e\x56\x72\xf9\xb9\x0e\x72\xcd\x1a\x57\xab\x57\x9e\xa4\xa8\xe2\x34\x30\x97\x5f\x5e\x21\x3d\x57\xf9\x17\x22\x09\x60\x2c\xe2\x0d\xf7\x00\x00\x17\x61\x4a\xb8\xe7\x50\x4c\x9d\x68\xf9\x0b\xbd\xd7\x9c\x09\x71\x3c\x4f\xe9\x60\xee\xfe\x1a\xa3\x25\x63\xff\x17\xc7\x7b\x86\xc5\xcc\xb8\xa6\xf3\x68\x3d\xb4\xea\xa1\x5e\xcc\x17\xfd\x5e\x9c\x14\xde\xd2\x25\x18\x2e\x7c\x04\x6e\xc4\x6d\x6e\x74\x7d\x76\x26\xc4\x3a\x69\xb7\x2d\xd6\xe2\x38\x86\xf5\xa5\xb8\x41\x8c\x09\xda\xa8\x55\xaa\x90\xca\xc3\x21\x33\xb1\xda\x6d\xd8\x47\xee\x3a\xea\x6f\xce\x38\x5c\x5e\xb6\x5c\xdf\x3f\xc4\x88\x2e\xd0\x00\x82\x00\xdd\x81\x6c\x4e\xb4\x9c\xde\xe5\xe3\x32\xbf\x97\x91\x47\x30\xa6\x09\x1a\xc9\x62\xf5\x87\x54\xa7\x2d\xf3\x91\x27\x77\x9e\x8a\x7f\x41\x99\xa5\x88\x8b\xd0\xdb\xe6\xb1\x75\xa5\x33\xb1\x01\x09\xc1\x02\x81\x4c\x13\x6f\xdb\xf0\x18\xfd\x4a\x71\x8c\x98\xe2\x24\x4e\x51\x6f\x4b\xc4\xae\x98\xd0\xf5\x8d\xab\xb6\x38\xf4\xda\x28\x4c\x84\x19\x35\x0c\x78\x14\xb3\x28\xc0\xe0\xcd\xc5\xdf\x8d\xb4\x9c\x8a\xb0\x00\x77\x74\x83\xc7\x17\x8a\x84\x7a\x27\x60\xf1\x0c\xa0\x4f\x75\xf7\x1e\x00\xba\xc7\x09\x49\x6c\x96\x69\xb2\x86\x3b\x8f\xd9\x1f\xd7\xd8\x7b\xda\x6b\x3d\x07\x4a\x94\x50\xae\xd4\x7d\x23\x16\xd5\x69\x37\x68\xb1\x7d\x9b\x39\x68\xcd\x67\xdd\x20\x95\xce\xd4\x77\x88\x2c\x66\xbe\xfe\xf5\xd3\xe9\xa0\xc9\x1b\x23\xcd\xbe\x84\x04\xfc\x33\x4c\x03\xcf\x82\xc5\x8e\x0c\x26\x35\x29\x50\xe8\x63\x40\x56\x91\xb0\x30\x20\x28\xb5\x60\x40\x3e\x7e\x20\xf9\x5b\x36\xa1\x36\x7f\x33\x47\x59\x61\x46\x0d\x03\x5a\xc4\xdd\x79\x9b\xd3\x47\x5c\x0f\xf9\x14\xf7\x7a\x01\xdd\x8c\xb5\x05\x5d\x41\x79\xe0\xdc\xcb\x2c\x5e\x84\xae\xd9\x94\x19\x06\x8d\x8d\xcd\xa6\x0e\x3d\x6f\x9d\x9c\x05\x26\x92\x7a\xa8\x7e\xfa\xe4\x87\xd0\xcb\xb7\xb8\x9c\x97\x4e\x0e\xc5\x06\xd9\x70\x16\x94\xfd\xe6\x50\xc9\x1a\xee\x36\x00\xab\x67\xcf\x20\xd9\x4a\x59\x33\xc2\x12\xfb\x08\x50\xe6\xaf\xc1\x02\x26\xe8\xed\x1b\x80\x02\x37\xf4\xba\x4a\xa4\xfa\xb5\x09\x93\xbd\x17\xa3\x08\xcf\x40\xf7\xc8\x4d\x09\x7b\x16\xd0\xad\x35\xb8\xdc\x7b\xc0\xbd\xf4\x91\x72\x85\x62\x93\xfa\x04\xd3\xf8\x23\x0e\xd3\xfe\x95\xc7\xd4\x3f\x7a\x7e\xc9\xd0\xc2\x2e\x25\x4a\xba\xc1\x2c\x25\xd4\xbd\x96\x78\x75\xad\x90\xbb\x45\x4d\x0e\x64\x1c\x81\xe0\x68\x50\xa1\xbb\x54\x8e\xb0\x8b\x4b\xe3\xf3\x1b\xd7\xe9\x72\x67\xb4\x9b\xc7\x7c\xe8\xc9\xa1\x8a\x87\xa1\xdf\x0d\xa6\x60\xd6\xec\x67\x02\x25\x74\x80\x05\x12\x41\x79\xbe\x40\x52\xf2\x3f\x0b\x23\xf9\xd0\x13\x80\x91\x6b\x1c\x24\x04\x06\xae\xc4\xd1\xcc\xfa\x7d\x44\x32\x52\xb0\x03\x34\x37\x06\xd0\xa4\xdd\xa1\xc8\x4b\xe6\x85\x48\xc7\x45\x94\x36\x58\x91\xa1\xd0\xaf\x14\xc5\x5b\x8b\xed\x16\x52\x4a\x74\xc7\x68\x8a\x2a\xfc\xde\x76\x47\x69\x6a\xb6\x0b\x65\xf6\xee\x8e\xaa\x21\x88\x61\x83\xd4\x7c\xcb\xb0\x6e\x31\x43\x7d\x23\x95\x8a\x35\xcd\xf4\x9b\xa2\xf4\x9e\xdc\x1d\xf6\xd1\x9b\x46\xe4\x0c\xac\x87\x2a\x77\x8d\x43\x75\x51\xb5\xb7\xd0\xf9\xa2\x87\x5e\x8e\xf2\xbc\xe7\x11\x57\xf6\x6e\x6c\x4a\xaa\x55\x94\x34\x83\x97\xa4\x5d\xc2\x72\x30\xec\xe9\x73\x2b\x94\x07\x8b\xdd\x13\x0d\x13\x0c\xec\xbe\xa8\xa9\x43\x37\x19\x8b\x59\x37\x54\x23\x61\x31\x38\xac\x65\x61\x63\x10\x4d\x53\xd7\x19\xe8\x77\xde\x3b\x35\xf5\xf2\xa7\xbd\x1a\x77\x38\xbf\xec\xcd\xf8\x24\x4e\x76\x10\xe7\x10\x07\xc8\x2c\x66\xab\x86\x0e\x13\xb3\x0d\xfa\xa9\x66\x9c\xb4\xd5\x3e\xb3\xdc\x5a\x65\x41\x7b\x38\xb9\x9e\x6d\xba\x3a\xbd\x00\xde\x63\xa3\xe7\x44\x4c\xa0\xd6\xbb\xbd\x6f\x74\x74\x76\xc8\xde\xa0\x28\xcd\x6e\x69\x03\x7f\x68\x81\x6f\xbc\xc9\xb3\xc1\x2e\xa7\x3b\x4a\xb0\x8b\xdf\x92\xd9\xf7\x89\x13\x8b\x64\x7e\x28\x90\x4a\x0e\x6e\x04\x53\xed\xe3\xe4\xe2\xd7\x65\xde\x3f\xe4\xb1\x30\x97\x45\xdb\x11\x23\x7b\xe7\x33\xa4\xda\xdc\x75\x12\xc7\x72\xb8\xc0\x25\xb8\x18\x56\x68\x1f\x74\x73\x20\x5c\xcc\x56\x75\xd5\x57\x8d\xe2\xdf\x79\x14\x7f\x75\x50\xc5\x65\x60\x50\x60\xc0\xf3\x31\x62\xec\x81\xf1\x58\x5c\x9e\xcf\xc6\x63\x23\x6c\xe8\x0f\x0f\x9e\xf5\xee\x7c\x85\x9f\xcf\x40\xb8\xcc\x8e\x50\x2b\xb5\x28\xc6\xfe\x26\x0b\xfc\x29\xa3\x80\x4d\x36\x9a\x3a\x74\x52\x1a\x30\x2b\xe7\x76\x89\x20\x59\x8d\xd7\x82\x48\x83\xb2\xe7\x43\x1c\xbb\x45\x6d\x55\xcd\x15\xf5\xfe\x4e\x2d\x71\x1c\x08\xb2\xf8\xcc\x3e\x16\x9f\x15\x74\x47\x49\x06\x23\x1f\x06\xfb\xee\x04\x6f\x38\x93\x1b\xa3\xe4\xee\x1b\xbf\xef\x19\xec\xf8\x96\xd0\x4f\xcc\xb7\x7c\xf3\x73\xdb\xf2\x1d\xb6\x09\x91\x39\xac\xed\x3c\xd4\xd4\xcc\x34\xf1\x33\xee\x3c\xe4\xb3\xa1\xdb\x6e\x58\xc2\x8f\xfe\xe0\x43\x21\x51\x41\xd7\xfa\xb1\x74\xa6\xac\x68\xce\x1e\x8f\x59\xdf\xbe\x22\xab\x2d\xb8\x1c\x2a\xab\xd9\x1e\x13\x03\xeb\x29\x3c\x68\x4a\xd3\xca\x32\xe7\x8b\x07\x5b\x12\x0d\xe7\x91\xfd\xd7\x55\xbd\x89\xf1\x2a\x6d\x15\xc7\x63\x76\xc1\xb8\xcc\xd4\x3f\x92\xf4\xbb\x45\xdc\x85\x55\xb6\xc8\xb4\xe5\x73\x80\x88\xb7\x3b\x98\xa6\x0e\xdd\x24\x1a\x2d\x2a\x4c\x7c\xf5\xad\x16\x97\x72\xd0\xe0\x8b\x72\xa9\xe6\xb2\x77\xe9\xc9\x02\x4b\x4e\xd9\x6f\xe1\xa9\x97\x64\x4d\x54\x9d\xba\xb3\xc1\x49\x27\x67\x16\xaa\xa5\x8a\x5a\xa8\x56\x5c\xd9\xb7\x4f\x70\x57\x82\x57\x60\xb5\x29\x40\x67\x82\x58\x80\xae\x51\x0e\xba\x71\x50\x06\x27\x36\xac\xe5\x74\xc7\xe9\x1e\xc4\xd0\xdf\xbb\x86\xcc\x99\x68\xd6\x90\xa7\xbe\x3f\xe3\xf7\xb4\xf5\x63\xdb\x32\x94\xd1\x1d\xf8\x10\x3b\xf5\x3e\x5b\x3f\xd6\xd4\xcc\x74\x5b\x67\x7e\x72\x9d\xfd\x68\xce\x52\xaa\xab\xe2\xc8\x3a\x9b\xbe\xfe\xb0\x43\xce\xbd\xa0\x6b\xb5\x1f\xf9\xc2\xde\x5d\xb2\x83\xf1\xc1\x8e\x32\x6d\xf5\xfe\x73\xae\x14\xb7\xb2\xcc\xf9\x46\xfe\x96\x74\xc2\x79\xe4\xbf\x69\xd4\x51\xa5\x98\xf1\xaa\x6c\x24\x04\x73\xdd\x52\x71\xff\x80\xd1\xef\x4e\xa1\x3b\xe4\xb0\x35\xe2\x3e\x62\xdd\xee\x50\x9a\x3a\x74\x93\x4c\xb4\xa8\x11\xf3\x74\xa2\x5a\x23\x36\x81\x0b\x71\xc4\xdc\x22\x86\xa0\xec\xf9\x24\x79\x1f\xc9\x56\xe6\x0d\x9e\x4d\xb8\x80\x05\x61\xb5\xa2\x16\x84\x15\x57\x3a\xaa\xfe\xee\x95\xb4\x89\x33\xe3\x16\x85\x05\xa5\xad\xf0\xda\xd0\x95\x5e\x55\x6e\xc9\xf8\x5f\xf9\x5b\xe3\x3d\xd5\xcb\xef\xe5\xef\xd3\xce\x85\x0c\x17\x3f\xa9\xfa\x35\xb5\xcb\x5e\xdb\x8c\xbd\x11\xfb\x71\x65\x66\xbf\xca\x37\x57\xf5\xbc\x2e\x66\xf1\x4e\xb0\xc2\x0b\x0b\x26\x4a\x17\xdd\x65\x46\x9d\x9d\x6a\xde\xea\xde\xa9\xfa\x5a\x95\xe7\x0c\x82\xf6\xb2\x10\xf6\xf6\x30\x8e\x44\x22\x89\xdc\xba\xfc\x15\xd1\xb6\xc1\xc1\x67\x14\xac\x28\x74\xbe\x03\x7f\x7b\xab\x22\x82\xf7\x25\x22\x2d\x49\x17\x38\xf0\xf8\x7b\x2f\x76\xca\xb7\x08\x43\x1f\xc9\x1e\x2f\xcb\xd8\x8a\x55\xa1\x3b\x95\x35\x34\xfe\xf3\x42\x6f\x6e\x2a\xb0\xd4\xe1\xa4\x13\xb8\xda\xb2\x0c\xe8\xd4\xe4\x35\x6a\xf1\x2d\xd7\x2e\x99\xbc\x74\xcd\x87\xfc\x7d\x21\x3b\x65\x16\xa1\xdb\x0a\x17\xf8\x73\xd9\x6e\x41\x41\x99\x69\x74\x80\x04\x03\x76\xdc\xce\x41\x4a\xd8\x71\xa0\xc8\xd4\x5b\x9c\x2e\x63\xb4\x17\xdc\x69\x38\x3d\x4f\x55\xfb\x5e\x09\xed\x8a\xd7\xca\xaf\x2a\x3f\xae\xba\x07\xf6\xc9\x76\x52\x32\xac\x53\xa9\xdd\xce\xb1\x6a\xbf\x09\xd7\x8b\x8f\x9d\x58\x92\x35\x58\x4f\x33\x5a\x4b\x8e\xed\xbb\xf2\xbe\x95\x06\xcb\xe1\xa4\x2f\xaa\x0d\x87\xf6\xad\xf4\xe3\xad\xdb\x30\xe3\x6e\xa1\xf6\x67\xce\x7a\x50\x81\x36\xe0\x9c\xa8\x6c\xb3\x41\x4b\xa8\x7e\xb5\x49\x43\xcc\x23\x85\x98\x66\x5a\xa3\x7c\xb3\x93\xea\xb6\x92\xe5\xcb\x98\xc5\x8b\xec\xdf\xa7\xff\x07\x00\x00\xff\xff\xcf\x80\x8e\x41\x42\x92\x00\x00")

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

	info := bindataFileInfo{name: "swagger-spec/api.json", size: 37442, mode: os.FileMode(436), modTime: time.Unix(1453118722, 0)}
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


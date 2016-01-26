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

var _swaggerSpecApiJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x5d\x5f\x93\xdb\xb6\x11\x7f\xf7\xa7\xc0\xa8\x7d\xd2\xd8\xe2\x39\x75\x32\x53\x3f\x55\xb1\xd2\x8e\x26\x19\xc7\x33\xd7\x3c\xb9\x9e\x3b\x88\x84\x24\xa4\x14\xc9\x90\xa0\xcf\xd7\x9b\xfb\xee\x05\x40\x90\xc7\x3f\x80\x04\x50\xa0\x44\x45\xd0\x83\x7d\x43\x2d\x96\xbb\x8b\xdd\x1f\x16\xcb\x05\xf5\xf4\x0a\xd0\xcf\x24\x7b\x80\x9b\x0d\x4a\x27\xef\xc1\xe4\xbb\xd9\xcd\xe4\x75\x71\x15\x47\xeb\x98\x5e\x2a\x68\xf8\x95\x00\x65\x7e\x8a\x13\x82\xe3\x88\xd1\xfe\x16\xe1\xaf\x28\xcd\x60\x08\x32\x94\x7e\xc5\x3e\x02\xab\x34\xfe\x2f\x4a\xc1\x0e\x46\x70\x83\x76\x28\x22\x60\xfe\x69\x29\xf8\x71\x0e\x8c\x5e\x8c\xbe\x99\xdd\xcc\xde\xd6\xbf\x23\x98\x84\x88\xf3\xbd\xfd\xb1\xce\x02\x26\xb8\x41\x86\xd2\x5d\xf6\xeb\xfa\xb6\xb8\x25\xa3\xdf\x12\x92\x64\xef\x3d\x6f\x83\xc9\x36\x5f\xcd\xfc\x78\xe7\x6d\x13\x3f\x8c\xf3\xc0\xf3\xd7\x6f\xf2\x6c\x55\x1f\xee\xc7\x11\x81\x3e\x69\x28\xc6\xbf\x40\x3b\x88\x43\xce\x8e\x0e\x0a\xd0\xd7\x7f\x6c\x13\xc4\x78\x4d\x2a\xb2\xe7\x1a\x9b\x90\xde\x3b\xca\x50\x97\x4d\x04\x77\x5c\xa8\x79\x02\xfd\x2d\x02\x2f\xf6\xac\x28\xf2\x34\x2c\xa5\xa6\x42\x3f\x3c\x3c\xcc\x20\xa7\x9d\xc5\xe9\xc6\x13\x7c\x33\xef\x97\xe5\x87\x9f\x3e\xde\xfe\xf4\x86\x32\x98\x6d\xc9\x2e\xac\x89\xf1\xaa\x26\xcc\x24\x43\x7e\x9e\x62\xf2\xb8\x40\x6b\x1c\x61\x36\x35\x59\x73\xd2\xe6\x39\xd9\xc6\x29\xfe\x1f\x14\xd3\xd6\x92\x97\x3c\x26\x5c\x5e\x6a\xe5\x9f\xd1\x63\x5b\x56\xcc\xa7\x6a\x8b\x60\x40\xdd\xe3\xb5\x42\xd3\xc6\x0d\x5a\x44\x6d\x97\x99\xcf\xc1\xbf\xa9\x8f\x44\x87\xd4\xa1\xc4\x9f\x2b\x92\x96\xc8\x6d\x8d\x3e\x7f\x69\x31\xfb\x22\x98\xd1\xb9\xce\xf2\x1d\xca\x1a\xcc\xa8\xa6\x09\xb5\x32\x1f\xec\xfd\x9e\xc5\x42\x94\x72\x4c\x92\xc6\x41\xee\x1b\x8e\x81\x64\xdb\xb2\xba\x97\x27\x01\x24\xe8\x8e\x8e\x81\x61\xbc\xe9\xda\x3d\x89\xb3\xae\x13\x4a\x4d\xc6\x19\x65\x80\x6c\xab\xf8\x12\x4c\xff\xd3\xb6\x36\x1f\x1e\x27\x28\xe5\x82\x2e\x03\x36\xbc\x90\xe3\x83\x10\x43\x42\x9f\xa2\x2c\x89\x99\xcb\x49\xa5\xe1\x24\xdf\xdd\xdc\x28\xbf\x94\x89\x7c\x9b\xfb\xd4\x82\xd9\x3a\x0f\x41\xc9\x7d\x22\x1d\xfc\xdc\x95\x87\xf3\xfb\xde\xf0\x86\xbf\x45\xe8\x5b\x82\x7c\x82\x02\x80\xd2\x34\x6e\xbb\x6a\x63\x68\x46\x43\x6d\x07\xf7\xb2\xe7\x74\x65\x60\x64\x24\xc5\xd1\x46\x2e\x3f\xd7\x41\xae\x59\xe7\x6a\xf3\xca\xb3\x14\x55\xbc\x0e\xe6\xf2\xcb\x1b\xa4\xe7\x2a\xff\x42\x24\x03\x8c\x45\xba\xe3\x1e\x00\xe0\x2a\xce\x09\xf7\x1c\x8a\xa9\x33\x2d\x7f\xa1\xf7\x5a\x32\x21\xce\xe7\x29\x16\xe6\xee\xaf\x29\x5a\x33\xf6\x7f\xf1\x82\x17\x58\x2c\x8c\x6b\x3a\x8f\xce\x43\x9b\x1e\x1a\xa4\x7c\xd1\x1f\xc4\x49\xe1\x57\xba\x04\xc3\x55\x88\xc0\xbd\xb8\xcd\xbd\xae\xcf\x2e\x84\x58\x17\xed\xb6\xd5\x5a\x9c\xa6\xb0\xbd\x14\x77\x88\x31\x41\x3b\xb5\x4a\x0d\x52\x79\x38\x14\x26\x56\xbb\x0d\xfb\xc8\x5d\x47\xfd\xcd\x15\x87\xcb\xeb\x9e\xeb\xfb\x87\x14\xd1\x05\x1a\x40\x10\xa1\x07\x50\xcc\x89\x96\xd3\xfb\x7c\x5c\xe1\xf7\x32\xf2\x04\xa6\x34\x41\x23\x45\xac\x7e\x96\xea\xb4\x67\x3e\xca\xe4\x2e\x50\xf1\xaf\x28\x8b\x14\x71\x15\x07\xfb\x3c\xb6\xad\x74\x21\x36\x20\x31\x58\x21\x50\x68\x12\xec\x1b\x9e\xa2\x3f\x72\x9c\x22\xa6\x38\x49\x73\x34\xd8\x12\x71\x28\x26\x74\x7d\xe3\x4b\x5f\x1c\x7a\x6b\x14\x26\xc2\x8c\x1a\x06\x3c\x8b\x59\x14\x60\xf0\xee\xe6\xef\x46\x5a\xce\x45\x58\x80\x07\xba\xc1\xe3\x0b\x45\x46\xbd\x13\xb0\x78\x06\x30\xa4\xba\x07\x8f\x00\x7d\xc3\x19\xc9\x5c\x96\x69\xb2\x86\x7b\x4f\xc5\x1f\x77\x38\x78\x3e\x6a\x3d\x07\x4a\x94\x50\xae\xd4\x43\x23\x16\xd5\xe9\x30\x68\xb1\x7d\x9b\x39\x68\x2d\x17\x76\x90\x4a\x67\xea\x2d\x22\x8b\x99\xaf\xff\xfa\xf3\xe5\xa0\xc9\x3b\x23\xcd\x3e\xc6\x04\xfc\x33\xce\xa3\xc0\x81\xc5\x81\x0c\x26\x37\x29\x50\xe8\x63\x40\x51\x91\x70\x30\x20\x28\xb5\x60\x40\x3e\x7e\x24\xf9\x5b\x31\xa1\x2e\x7f\x33\x47\x59\x61\x46\x0d\x03\x5e\x21\xe2\xba\x34\xf1\x85\xee\xb4\xc8\x1f\xa0\x90\xe2\xef\x20\xe0\x5f\xb0\x76\xe0\x2f\x28\x4f\x9c\x03\x9a\xc5\xad\xd0\xb5\x98\x32\xe3\xe0\x75\x59\x59\x5b\x87\x81\xb7\x70\xde\x0a\x13\x49\x5d\x56\x3f\x8d\x0b\x63\x18\x94\x18\xca\x79\xe9\xe4\x72\x6c\x90\x0b\x67\x41\x39\x6c\x2e\x97\x6d\xe1\x61\x03\xb0\xba\xfa\x02\x92\xbd\x94\x2d\x23\xac\x71\x48\xd7\xca\x2d\x7c\x0b\x56\x30\x43\x3f\xbc\x03\x28\xf2\xe3\xc0\x56\x42\x37\xac\x4d\x98\xec\x83\x18\x45\x78\x06\xfa\x86\xfc\x9c\xb0\x67\x12\x76\xad\xc1\xe5\x3e\x02\xee\xa5\x8f\xb6\x1b\x14\xbb\x3c\x24\x98\xc6\x1f\xf1\x98\xf6\x6f\x02\xa6\xfe\xd9\xf3\x5c\x86\x16\x6e\x29\x51\xd2\x8d\x66\x29\xa1\xee\xb5\xc6\x9b\x3b\x85\xdc\x3d\x6a\x83\xa0\xe0\x08\x04\x47\x83\x4a\xe1\xad\x72\x84\x5b\x5c\x3a\x9f\x3f\x71\xbd\xb0\x74\xc6\xcb\xda\xc4\x3a\x54\xa9\xa1\x4a\x80\x61\x68\x07\x53\x30\x6b\x3a\x34\x81\x12\x3a\xc0\x01\x89\xa0\xbc\x5e\x20\xa9\xf9\x9f\x83\x91\x72\xe8\x05\xc0\xc8\x1d\x8e\x32\x02\x23\x5f\xe2\x68\x66\x7d\x47\x22\x19\xa9\xd8\x01\x9a\x1b\x03\x68\xd2\x76\x51\xe5\x25\xcb\x4a\xa4\xf3\x22\x4a\x1f\xac\x28\x50\xe8\x8f\x1c\xa5\x7b\x8b\xfe\x0e\x52\x6a\x74\xe7\x68\xce\xaa\xfc\xde\x75\x69\x69\x6a\x76\x08\x65\x8e\xee\xd2\x6a\x21\x88\x61\xa3\xd6\x72\xcf\x30\xbb\x98\xa1\xbe\x91\x4a\xc5\x96\x66\xfa\xcd\x59\x7a\x4f\x10\x4f\xfb\x08\x50\x23\x72\x46\xd6\xcb\x55\xba\xc6\xa9\xba\xb9\xfa\x5b\xc8\xd2\xf3\xba\x0f\x71\xb4\x0e\xb1\x4f\x5c\x26\xd4\x2b\x13\x7a\xd9\x59\x89\x2b\x47\xb7\x71\x65\xcd\x5a\x4d\x5e\x80\x58\xd6\x2f\x2d\x3a\x19\xc2\x0d\xb9\xe1\x2a\x43\xd2\xed\xbc\xc6\x0a\x39\x6e\xf7\xd5\xd6\xc1\x4e\x5e\x64\xd6\xfb\xd5\x49\x8b\x0c\x8e\xa6\x39\xd8\x18\x45\x8b\xd8\x5d\x01\xfa\xd6\x3b\xc5\xe6\x41\xf9\x4c\x59\xe3\x0e\xd7\x97\x23\x1a\x9f\x3b\x2a\x8e\x1d\x9d\xe2\xb8\xdc\x65\x63\xb6\xcb\x46\x3b\x3a\x58\x59\x19\x0c\x7a\xc3\x16\x9c\xb4\xd7\x9e\xb9\xde\x26\xe6\x96\x86\xf1\x64\x94\xae\x81\xec\xf2\x02\xf8\x88\xed\xa4\x97\x30\x81\x7a\xef\x29\x3f\xd1\xd1\xc5\x8b\x0b\x0c\x0a\xec\xec\x96\x2e\xf0\xc7\x16\xf8\xc6\x5b\x49\x17\xec\x72\xba\xb3\x04\xbb\x78\x3f\xcf\xb1\x4f\xcf\x58\x24\xf3\x0e\x7a\x2a\x39\xb8\x17\x4c\xb5\x8f\xe8\x8b\x37\xf6\xfc\xf8\x58\xc6\xc2\x52\x16\x6d\x67\x8c\xec\x83\xcf\xc3\x5a\x73\x67\x25\x8e\xe5\x70\x81\x6b\x70\x31\xae\xd0\x3e\xe9\x16\x44\xb8\x98\x7b\x69\x87\xfa\xaa\x51\xfc\x7b\x4f\xe2\x2f\x0b\xb5\x62\x06\x06\x15\x06\xbc\x9c\xb9\xc1\x01\x98\x4e\xc5\xe5\xe5\x62\x3a\x35\xc2\x86\xe1\xf0\xe0\x45\x6f\xeb\x2b\xfc\x72\x01\xe2\x75\x71\xde\x48\xa9\x45\x35\xf6\x4f\xb2\xc0\x5f\x32\x0a\xb8\x64\xa3\xab\x83\x95\xd2\x80\x59\xd1\xd8\x26\x82\x14\x95\x64\x07\x22\x1d\xca\x81\x0f\xa4\x1c\x16\xb5\x57\xcd\x58\x3c\x55\xb0\x6a\x89\xf3\x40\x90\xc3\x67\xf6\x19\x25\x3e\xbb\x02\x71\x47\x07\x9b\x29\x67\x12\xc2\xe8\xd8\xfd\xe6\x3d\x67\x72\x6f\x94\x42\x7e\xe2\xf7\xbd\x82\x7d\xe5\x1a\x86\x99\xf9\xc6\x72\x79\x6d\x1b\xcb\xd3\xb6\x6d\x32\x87\x75\xbd\x9a\x9a\x9a\x99\xa6\x97\xc6\xbd\x9a\x7c\x36\x74\x1b\x34\x6b\xf8\x31\x1c\x7c\x28\x24\xaa\xe8\x7a\x3f\x62\x2f\x94\x15\xed\xec\xd3\x29\x3b\xe9\xa0\xc8\x9d\x2b\x2e\xa7\xca\x9d\xf6\xc7\xc4\xc8\xba\x30\x4f\x9a\x38\xf5\xb2\xcc\xf5\xe2\xc1\x9e\x44\xc3\x7b\x62\xff\xd9\xaa\x6a\x31\x5e\xb5\x0d\xe9\x74\xca\x2e\x18\x17\xb3\x86\x47\x92\x61\x37\xa2\x87\xb0\xca\x95\xb2\xf6\x7c\x4e\x10\xf1\xae\x8e\xd5\xd5\xc1\x4e\xa2\xd1\xa3\x8e\xc5\x57\xdf\x66\x09\xab\x04\x0d\xbe\x28\xd7\x2a\x3b\x47\x17\xb8\x1c\xb0\x94\x94\xc3\x96\xb7\x06\x49\xd6\x44\x6d\xcb\x9e\x0d\x2e\x3a\x39\x73\x50\x2d\x55\xd4\x41\xb5\xe2\xca\xb1\xdd\x88\x87\x12\xbc\x0a\xab\x4d\x01\xba\x10\xc4\x01\x74\x8b\x72\xd4\xed\x89\x32\x38\x71\x61\x2d\xa7\x3b\x4f\x8f\x22\x86\xe1\xd1\x35\x64\xce\x44\xb3\x86\x3c\x0f\xc3\x05\xbf\xa7\xab\x1f\xbb\xc6\xa4\x82\xee\xc4\xc7\xfe\xa9\xf7\xb9\xfa\xb1\xa6\x66\xa6\xdb\x3a\xf3\xb3\xfe\xec\x35\x43\x6b\xa9\xae\x8a\x43\xfe\x6c\xfa\x86\xc3\x0e\x39\xf7\x8a\xae\xd7\x7e\xe4\x23\xfb\xd5\x99\x03\x8c\x4f\x76\x2c\x6b\xaf\xf7\x5f\x73\xa5\xb8\x97\x65\xae\x37\xf2\xf7\xa4\x13\xde\x13\x7f\x0b\x94\xa5\x4a\x31\xe3\xd5\xd8\x48\x08\xe6\xba\xa5\xe2\xe1\x01\x63\xd8\x9d\x82\x3d\xe4\x70\x35\xe2\x21\x62\xdd\xed\x50\xba\x3a\xd8\x49\x26\x7a\xd4\x88\x79\x3a\xd1\xac\x11\x9b\xc0\x85\x38\x2e\xef\x10\x43\x50\x0e\x7c\x2a\x7e\x88\x64\xab\xf0\x86\xc0\x25\x5c\xc0\x81\xb0\x5a\x51\x07\xc2\x8a\x2b\x96\xaa\xbf\x47\x25\x6d\xe2\x64\xba\x43\x61\x41\xe9\x2a\xbc\x2e\x74\xa5\x57\x95\x5b\x32\xfe\x57\xf9\x7b\xff\x35\xcc\x6d\x88\x38\x91\xff\x12\x7a\x29\x64\xbc\xfa\x9d\xaa\xdf\x52\xbb\xee\xb5\xdd\xd8\x9b\xb0\xd7\x51\x33\xfb\x35\xbe\xf9\xd2\xce\xeb\x52\x16\xef\x04\x2b\xbc\xb0\x62\xa2\x74\xd1\x43\x66\xd4\xd9\xa9\x96\x0d\xf5\x56\xd5\xd7\xaa\x3c\x17\x10\x74\x94\x85\x70\x70\x84\x71\x24\x12\x49\xe4\xd6\xe5\xaf\x88\xb6\x1d\x8e\x7e\x41\xd1\x86\x42\xe7\x7b\xf0\xb7\x1f\x54\x44\xf0\x5b\x8d\x48\x4b\xd2\x15\x8e\x02\xfe\x4b\x21\x07\xe5\x5b\xc5\x71\x88\x64\x8f\x97\x65\x6c\xc5\xaa\x60\x4f\x65\x0d\x8d\xbf\xbf\xd1\x9b\x9b\x06\x2c\x59\x9c\x74\x02\x37\x7b\x96\x01\x9d\x9a\xbc\x46\x2d\xbe\xe7\xda\x25\x93\x97\xae\xf9\x90\xff\xc2\xca\x41\x99\x45\xe8\xf6\xc2\x05\xfe\x5c\xd6\x2e\x28\x28\x33\x0d\x0b\x48\x30\x62\xc7\xb5\x0e\x52\xc2\x8e\x23\x45\xa6\xc1\xe2\x74\x9d\xa2\xa3\xe0\x4e\xc3\xe9\x79\xaa\x3a\xf4\x4a\xe8\x56\xbc\x5e\x7e\xd5\x78\x51\xec\x11\xd8\x27\xdb\x49\xc9\xb0\x4e\xa5\x76\x3f\xc7\x6a\xbd\xdf\x6e\x10\x1f\xbb\xb0\x24\x6b\xb4\x9e\x66\xb4\x96\x9c\xdb\x77\xe5\x7d\x2b\x1d\x96\xe3\x49\x5f\x54\x1b\x0e\xed\x5b\xe9\xc7\x9b\xdd\x30\xe3\x6e\xa1\xf6\x67\xce\x7a\x54\x81\x36\xe2\x9c\xa8\x6e\xb3\x51\x4b\xa8\xfe\x31\x98\x8e\x98\x67\x0a\x31\xcd\xb4\x46\xf9\x5b\x58\xaa\xdb\x4a\x96\x2f\x63\x16\xaf\x8a\x7f\x9f\xff\x1f\x00\x00\xff\xff\x9d\xdb\xac\x0e\xfc\x93\x00\x00")

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

	info := bindataFileInfo{name: "swagger-spec/api.json", size: 37884, mode: os.FileMode(436), modTime: time.Unix(1453804552, 0)}
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


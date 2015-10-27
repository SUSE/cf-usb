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
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
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
	name string
	size int64
	mode os.FileMode
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

var _swaggerSpecApiJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x5d\x5f\x6f\xe3\xb8\x11\x7f\xdf\x4f\x21\xb8\x7d\x0a\xba\x71\xee\xd0\xa7\x7b\x6a\x6e\x73\x28\x8c\x16\x77\x0b\x04\xfb\xd4\x2e\x12\xda\xa2\x6d\x5e\x65\x49\x47\x51\xc9\x6e\x83\xfd\xee\xc7\x3f\x92\xa3\x3f\xa4\x44\x4a\x94\x57\x8e\x27\x2f\xb1\xe5\xe1\x88\x33\x9c\xf9\xcd\x70\x86\x96\x5f\xde\x05\xfc\x6f\x91\x3d\xa3\xdd\x0e\xd3\xc5\x4f\xc1\xe2\xc7\xeb\x9b\xc5\xdf\xd4\x55\x12\x6f\x13\x7e\x49\xd1\xc8\x2b\x21\xce\x36\x94\xa4\x8c\x24\xb1\xa0\xfd\x14\x93\x27\x4c\x33\x14\x05\x19\xa6\x4f\x64\x83\x83\x35\x4d\xfe\x87\x69\x70\x40\x31\xda\xe1\x03\x8e\x59\x70\xfb\x71\x55\xf0\x93\x1c\x04\x7d\x31\xfa\xe6\xfa\xe6\xfa\x87\xea\x67\x8c\xb0\x08\x4b\xbe\xf7\x3f\x57\x59\xa0\x94\xd4\xc8\x30\x3d\x64\xbf\x6d\xef\xd5\x2d\x05\xfd\x9e\xb1\x34\xfb\x69\xb9\xdc\x11\xb6\xcf\xd7\xd7\x9b\xe4\xb0\xdc\xa7\x9b\x28\xc9\xc3\xe5\x66\xfb\x3e\xcf\xd6\xd5\xe1\x9b\x24\x66\x68\xc3\x6a\x82\xc9\x0f\xf0\x01\x91\x48\xb2\xe3\x83\x42\xfc\xf4\x8f\x7d\x8a\x05\xaf\xc5\x91\xec\x5b\x85\x4d\xc4\xef\x1d\x67\xb8\xcd\x26\x46\x07\x39\xa9\xdb\x14\x6d\xf6\x38\x78\xd5\xe7\x91\x22\xa7\x51\x39\x6b\x3e\xe9\xe7\xe7\xe7\x6b\x24\x69\xaf\x13\xba\x5b\x16\x7c\xb3\xe5\xbf\x57\x1f\x7e\xf9\xf5\xfe\x97\xf7\x9c\xc1\xf5\x9e\x1d\xa2\xca\x34\xde\x55\x26\x23\xe4\xc9\xf2\x03\xce\x38\xcb\xff\xbc\x4e\x0f\xa5\x29\xe7\x84\xc4\x4a\x2d\x7f\xcf\xb8\xc2\xe5\x47\x9f\x8b\x31\x29\x4d\xc2\x7c\xe3\x38\x06\xb1\x7d\x56\x37\x87\xa5\x5a\xee\xb6\x0e\x76\xb8\xad\x5f\x9d\x01\xfd\x13\xb3\x2c\x60\xfb\xa3\xdd\xf0\x95\x0e\xb8\x3c\x31\xde\x08\x8a\x40\x18\xe0\x7f\xe3\x86\xfa\x24\x9f\x24\xc5\x54\x4e\x74\x15\x0a\x3e\xfc\x7e\x3f\x4b\x0e\x2b\x61\xb2\x1a\xfa\x14\x51\xbe\x2c\xdc\x72\xea\x22\x57\xff\xda\xb3\x3d\x8e\x2e\x97\x14\xe5\x6c\x9f\x50\xf2\x7f\x79\x67\xcd\x6d\x4c\x62\xde\x56\xc7\x05\x8c\x4f\xb4\x73\x34\x91\x83\xf6\x18\x85\x5c\xb7\x1d\x74\x14\xff\x91\x13\x8a\x85\x06\x18\xcd\x71\x07\x25\xfb\x9a\x4a\x01\x32\x46\x49\xbc\x5b\x68\x09\xbf\xb5\xae\x7e\xd6\x28\x92\xe2\x2c\x4d\x84\x7d\x6a\xd7\x57\x92\xfc\x78\x73\x63\xfc\x50\x12\x34\xb4\x73\x9f\x6f\xb8\x29\x66\xdb\x3c\x0a\x4a\xee\x5d\x52\x67\xdc\x53\x0e\xa8\xf3\x0e\x92\xee\xaf\x14\x6f\x05\xfb\xbf\x2c\x43\xbc\x25\x31\x11\xb7\xcb\x4a\x8b\x35\x0e\x6d\x6b\x41\x5e\xd5\xcf\x87\x4b\xb2\x45\x79\xa4\xb7\x75\x93\xb8\x9f\x62\xfc\x25\xe5\x06\x8e\xc3\x00\x53\x9a\x74\xae\xf0\x38\x59\x15\x7b\x57\x51\x5b\x57\xeb\x57\x1a\xaa\x58\xa4\xb9\x9d\xab\xdf\x7b\x72\xf5\x3c\x0d\x11\xc3\xe0\xed\x63\xbd\x5d\x3f\xde\x42\x2d\x85\x07\xf5\x4a\xb4\x4e\xc2\xaf\x0e\x5a\xfb\x24\xd7\x35\x0c\xfa\xd9\xdb\x2b\xe2\xbb\x40\x05\xc0\x68\x00\x30\xda\x0f\xa3\xda\xa4\xb2\xd0\xeb\x52\x81\xdc\x03\xcf\xc5\x50\x94\xec\x86\xe7\x57\xca\xa9\x6a\xb8\x5b\xf0\x74\xc0\xda\x0f\xc5\x2c\x00\x68\x9b\x94\x67\x97\x56\x81\x57\x76\x5f\x31\x78\x65\x6b\x23\x2c\x2f\x3b\x6d\x72\x04\x0b\x7a\x50\x96\x8a\xd6\x49\xce\xa4\x4f\xf2\x8d\xee\xb5\xed\x06\x07\x92\x1d\x2d\xe5\xd9\xf9\xe0\x64\x1e\x21\xad\x14\x22\xb2\x67\xdf\x0f\xa9\xac\x71\x4d\xe2\xfe\xe8\x09\x91\x08\xad\x23\x1c\x3c\x16\xb7\x79\xb4\x45\x83\xbb\x62\x5a\x00\x08\x4d\xca\x4b\x04\x84\x52\x66\x44\x29\xea\xda\x74\x49\x62\xc2\xf0\xc1\x2c\x52\x8d\x54\xef\x7e\xca\x56\xcd\xfe\x27\xfe\xf4\x3e\x68\xfe\x04\x80\xe8\xf5\x5d\xb3\xc2\x92\x64\x76\x48\xf3\x81\x62\x9e\xae\x07\x28\x88\xf1\x73\xa0\x16\xc9\x0a\x4e\x36\x72\x9c\x42\x14\x00\x94\x16\xe5\xb4\xe5\x94\xd0\xa4\xf6\x86\x44\x8e\xe5\x14\xb5\x9a\x5c\x5d\xc1\x1a\x07\x6a\x81\xc3\x39\x14\x55\xfa\xb0\x63\xf2\xa2\xca\x0f\x4e\x50\x52\xa8\xd1\x42\x81\xdf\x45\x2d\x00\x9a\xed\x77\xba\xec\x6d\xf9\xa2\x5e\x3c\x90\xf0\xdb\xa8\x4c\x2e\x30\xba\xab\x31\x47\x03\x44\x6d\x51\x9e\x02\x51\xf9\x52\xf7\x0b\x25\x7a\x9a\xee\xa0\xba\xba\xbb\x80\x4c\xf5\xb7\x7f\x01\xda\xf9\x90\x75\x56\x4d\x38\x55\x0f\xb6\x07\x31\x55\x00\x06\x1c\x33\x50\x02\x8e\x49\xca\xf3\x4f\x90\x95\x9d\x43\x82\xec\x1e\x26\x0a\x35\x5a\x28\x10\x42\x86\xad\xac\x27\x09\x19\x21\x8e\x38\x70\x4f\x12\x35\x14\x6b\x88\x1a\x06\x4a\x88\x1a\x92\xf2\xc4\xd9\xef\xdf\x87\xc0\x9a\xb2\xe4\x10\xfa\xa6\xde\x76\xdf\xcb\x35\x61\x9a\x66\x8a\x7d\x02\x1b\x25\x28\x2c\xa0\x28\x90\xbc\x6c\xb2\x58\x31\x08\xf0\xc8\x40\x09\x78\x24\x29\xa7\x55\x43\xb6\x47\xfd\x0a\x10\x3d\xc2\x3b\xc4\x3a\x29\x1b\x4a\xd8\x92\x08\x07\x3d\xcc\xe7\xa2\x03\x31\xd7\x49\x94\x50\x58\x02\xfe\x82\x37\x39\x13\xfd\x54\xbf\xda\x90\xf3\x1e\x11\x9f\xb4\x27\xf4\x6b\x14\x07\x1e\x2f\x08\x87\x21\xb6\x14\xd2\xbf\x0f\x85\xf8\xdf\x3d\xa3\x17\xa0\x09\xb1\xcf\x67\xec\x33\xcc\x7e\x40\x01\x5a\x1c\x9d\xde\x92\x5d\x50\x70\x74\x28\x47\xdf\x1b\x47\x40\x18\x84\x30\x78\xf1\x45\xe9\x87\x82\x09\x14\x1a\x26\xc1\xc3\x07\x12\x67\x0c\xc5\x1b\x8d\x51\xb8\x1d\xa8\x42\x51\x54\x87\xc2\x5c\xa1\x5d\xe6\x76\x80\x6a\x75\x9c\x0e\x00\x62\x93\x72\x06\x80\x38\x04\xea\x94\x22\xfe\xc8\x31\xed\x2c\x07\x57\xf4\xb0\x45\x51\x76\xe1\x90\xf8\x3d\x0e\x90\x1d\xb1\xc0\xf3\x49\xb2\x3e\x7c\x1a\x7d\xbc\xaa\xc0\x9d\xe3\xf4\xdd\x4e\x58\xad\x3a\x86\x01\xe2\x9c\x00\x71\xcc\xfa\x6f\x88\xd6\xd7\x51\x3a\x6d\x4b\xc8\xc2\x5d\x66\x76\x78\xaa\xb4\xf4\x53\x1d\x9f\x1a\xae\x21\xc8\xde\xda\xef\x3a\xb3\xb7\xd7\x6d\x6d\x71\x65\xf4\xc9\xaa\xcc\x67\x3a\x07\xd8\xda\xa2\x3c\x29\xb6\x4e\xb9\xcf\x2d\x41\x05\x36\xbc\x00\x9a\xee\x42\xcf\xf1\x38\x56\x2b\xa3\x74\xf8\x72\x2e\xe0\x1e\xe0\x5e\x75\xf0\xcc\x14\xa3\x82\xb9\xf7\xc3\x5b\xb7\x61\xd9\x03\xb6\xb8\xc3\xe5\xa5\xe9\xce\xdf\x49\x53\x5f\x49\x3b\xc5\x97\x54\x21\xe8\x58\x8b\xea\x18\x74\x1c\x0e\x74\xdd\x49\xd2\x41\x95\x8c\xea\xd9\x2e\x88\x3a\x46\x4a\x88\x3a\x06\x4a\x38\xf5\xf5\x86\x10\x68\x44\xad\x60\x99\x8a\xe5\x1f\x5c\x30\xf8\xc8\x47\xab\xc7\xda\x38\x7c\xcf\x55\xdc\x12\x90\x0b\x90\xab\x3a\x78\xb6\xc8\xe5\x5c\x27\x00\xb4\xea\xbe\x62\x40\xab\xe2\xc9\xb5\x63\xfb\xd1\x02\x8a\xf0\x17\x92\x31\x6e\x1b\xc1\x63\xc1\xd4\xfa\x69\x1e\xf7\xe5\x24\x00\x8f\x9a\x94\xf3\xc3\xa3\xde\xd6\x72\x63\xca\x7e\x9a\xd0\x7a\x94\x23\x15\x94\x9b\x17\x22\x4d\xb4\xab\x3c\x6d\xa3\xba\xf0\x63\x78\xd4\xc9\xac\x1f\x75\x72\x84\xdb\xe0\x75\xc2\xd6\x4d\xf9\xf2\x31\xe2\x80\xbc\x4d\xca\x89\x0f\xc4\x1b\xf5\xde\x10\x69\x40\x69\xb0\x7c\x1a\x7d\xa7\x39\x1c\x19\x9c\xaa\x3e\xd8\x0b\x26\x33\x6b\xdf\x9f\xb4\x2e\x38\x54\x39\x80\xa5\xed\x77\xda\xfc\x76\xf9\x52\xbc\xf2\xd0\xa9\x17\xc9\xee\x2b\xe8\x3e\x13\xb6\x97\x97\x48\x18\x5c\x5d\x15\x97\x57\x77\x57\x57\x4e\xb9\x2f\x00\x70\x8b\xf2\x24\x00\x3c\xc5\x16\x7c\x75\x17\x24\x5b\x69\x12\x16\x28\xff\x46\x76\xe0\x80\x96\xe7\x85\x96\xa3\x5a\xf6\x3e\x21\x50\xf5\xf1\x01\x05\x4d\x94\x80\x82\x75\xca\xf3\x4c\xcb\x8b\xa3\x2e\x5e\x35\x71\xee\xe9\x38\x04\x98\x2e\x79\xdf\x58\x80\x71\x6e\xcf\x67\x01\x72\x2c\x24\xab\x5b\x40\x24\x31\x51\x42\x24\xa9\x53\xce\xba\x17\xaf\x83\x3b\xc0\x9c\xee\x2b\x86\x12\x40\x1a\xa1\x78\x6c\x7f\xeb\x51\x32\x79\x74\xda\xd2\x7f\x94\xf7\x05\x1c\x6a\x52\x42\x4b\x4b\x91\xf6\x37\xee\x2f\x22\x03\x3b\x6d\x4b\x4b\xf8\x31\xf4\xb3\x66\xdb\xcf\x92\xcb\xe3\xd8\xbb\xfa\x68\x18\x03\x30\x3b\x21\xcc\x1a\x16\xaa\x21\xcf\x80\xe6\x95\xb2\x81\x60\x9b\x50\xfe\xf2\xea\x2a\x24\x28\x32\x54\x92\x8e\x5c\x4e\xb5\x65\xee\xc6\x8e\x4b\x6e\x5f\x0d\xd2\x0c\xe0\x66\xfb\x5d\x3b\x71\x5d\xbe\x88\x7f\xbe\xba\x56\x82\x57\xa5\x5e\x7b\x75\x25\x2e\x08\xff\x92\x0e\x27\x2f\xbd\x96\x6f\x9d\xb3\x5d\x40\xe1\x16\xe5\xf4\x28\x3c\xed\x8e\xbb\x0f\xe7\x67\xbb\xdd\x9e\x45\x6e\x0b\x68\x39\xd3\xde\x95\xcc\x31\xea\x6d\xab\x16\x12\x96\xc7\x68\xc6\x36\xb5\x00\x19\xb5\x94\x80\x8c\x25\xe5\x19\x66\xea\x45\x3f\xcb\x9f\x0e\xce\x3a\x33\x87\x58\x63\x14\xf6\x8d\xc5\x1a\xd7\x36\x96\x7d\xda\xed\x18\x61\x6a\xcd\x2e\x88\x30\x5a\x4a\x88\x30\x25\x25\xb4\xba\xde\x1a\x2e\xe9\x2b\x06\xa2\x5a\x36\xba\xd5\x25\x99\x58\xb6\xba\x6e\xa3\xe8\x4e\xde\x13\xd0\xa7\x49\x09\x6d\x2e\x45\x0a\xdf\xdc\x92\x74\x27\x7e\xc4\x28\x77\x4a\x68\x73\xcd\xb6\xcd\x25\x96\x47\xa4\x7e\xd6\xad\x2e\x01\xb2\x80\xb1\x2d\xca\x89\x31\x56\xaf\xf4\x86\x3c\x8e\x3b\xe7\x5f\xc5\x8f\x93\xf7\x30\x3e\xd9\x23\x9b\x3a\x51\xe2\x92\x1b\x5a\x83\x34\x03\x08\xd9\x7e\xd7\x4e\x4f\x97\x2f\xe2\x9f\xaf\x86\x96\xe0\x55\xdb\x59\x17\xcc\x6d\xbf\x7e\x05\xc0\xaa\xa5\x9c\x1e\x58\xa7\xdd\x3a\xfb\x43\x58\x68\x5b\x01\x26\xda\x89\xea\x9a\x35\x0e\x68\x5b\xc9\xbc\xb1\xde\xb6\x72\xc1\xbb\xe2\x99\xa9\x00\x79\x3a\x4a\x80\xbc\x92\xf2\x0c\xb3\x6a\xe5\x24\x21\x64\xd6\x01\x44\x91\x6e\x61\xdf\x58\x14\x19\xd8\x90\x1a\x95\x36\x17\x0f\x41\x85\x30\xa2\xa3\x84\x30\x52\x52\x42\xd3\xe9\xad\x61\x8f\xd5\xae\xde\xcf\xcf\x7c\x8a\x44\xd7\xe9\xd7\x3d\xf9\x00\xf8\x6d\x4f\x03\x25\x60\x52\x49\x09\xbb\xf9\x6a\x1e\x06\xbf\xf6\xd9\x2b\xaa\x3d\x24\xca\x57\x85\x6a\x16\x95\x1b\xd7\xe6\xbb\x58\x53\xee\xfe\xb4\x0d\x8e\xa5\x61\x26\xeb\xdf\xb9\x36\x1a\x5a\xa8\x1a\x78\x1b\xbd\x16\x11\xc9\x98\x16\x52\x44\xf3\x28\xc4\x31\x93\x1d\xfa\xda\xa7\x9f\x9b\xc5\x08\x2a\x20\x95\x11\x83\x39\x97\xb7\x30\x9a\x7a\x9f\x5f\x69\x4c\x66\x81\xbf\x70\x40\x8e\xb9\x15\xe6\x34\xf2\xcb\x39\xc4\x4f\x0f\x87\x24\xd4\xa7\xc5\x35\xae\xeb\x24\x89\xb0\xee\x84\x9c\x8e\x6d\x55\x9b\x46\xce\x9d\x4b\x75\xa4\xca\x33\x21\xfa\xa1\xd3\xbf\x53\x94\x65\xcf\x09\x35\x3c\x88\x5d\x83\x49\x35\xc9\xb4\x86\xf4\xca\xbb\x7b\xbd\xdb\xf3\xb4\xee\x6f\x77\x21\xab\xf8\x33\x80\x87\x9a\x55\x29\xb1\xbf\xbb\x79\x4c\x74\x48\xbc\x4d\xfc\x3a\xae\xf8\x85\x74\x81\x83\xa3\x5c\xb3\x64\x32\xdc\x83\x2c\x64\x2f\x1f\xbc\xe0\x55\x7c\xab\x33\x2c\x2a\xcd\x18\xa5\x21\x62\x36\xa8\x61\xf0\xd2\x9e\xb7\x57\xfe\x6b\x12\x87\x68\x1d\xf9\x86\xaf\x4e\x3f\x1e\x88\xb3\xd5\x58\xee\x91\x31\x43\xbb\x0e\x88\xb5\x39\x4a\x63\x71\x84\x66\x60\x32\xa8\x9b\x2f\xdf\x5b\xa0\x10\x31\x73\x6e\xd2\xf4\x93\x41\x4e\x28\xcf\x83\xfa\xf5\x40\x63\xea\xee\xc1\xed\xfc\xdb\x9b\x77\x47\x2e\xc4\x3f\x0f\xa7\xd8\x52\x3c\x0a\x12\x2c\x2c\x4c\x6e\xb4\xa6\xc6\xf8\x8b\xc2\xf2\xda\xcf\xb9\x8e\x80\x07\x5d\x51\x43\x07\x07\xf6\x93\xb5\x31\x87\xc6\xaf\x95\x4d\x62\x19\x67\x16\xf4\x3d\xb3\xf5\x0f\x92\x13\x5a\x9c\xfe\x78\x79\x8b\xe5\x7c\xe2\xb2\x29\x6d\xb5\xbe\x95\xbd\x97\xf8\x75\x8e\xf9\xd9\xff\x04\xd9\xa3\xf2\xa8\x62\xe8\x74\x98\x3e\x5f\x8b\xb5\x8c\xc8\x0f\xa6\x8a\xba\xc9\x87\x35\x18\x3e\x86\x85\xaa\x9c\x39\xda\x77\x9f\x7d\x1e\x70\x96\xa1\x9d\x77\xec\xb3\x29\xfa\x90\x98\xe1\x1d\xd6\xd4\x02\x7b\x6a\x7b\xef\xbe\xfd\x19\x00\x00\xff\xff\xa8\x3b\x78\x54\x30\xc0\x00\x00")

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

	info := bindataFileInfo{name: "swagger-spec/api.json", size: 49200, mode: os.FileMode(436), modTime: time.Unix(1446112610, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	if (err != nil) {
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
	Func func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"swagger-spec": &bintree{nil, map[string]*bintree{
		"api.json": &bintree{swaggerSpecApiJson, map[string]*bintree{
		}},
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


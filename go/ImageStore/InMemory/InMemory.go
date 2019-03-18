/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Package InMemory is a wrapper around Read, Remove and Store APIs of ImageStore.
package InMemory

import (
	"IEdgeInsights/ImageStore/go/ImageStore/InMemory/redis"
	"errors"
	"io"

	"github.com/golang/glog"
)

// InMemory is a struct used to comprise all Inmemory methods in it's scope
type InMemory struct {
	memType      string
	redisConnect *(redis.RedisConnect) //TODO: This should actually be an interface referring to respective concrete classes
}

// memoryDB is used for module level check with memory type
const memoryDB string = "redis"

// NewInmemory is used to initialize the connection based on DataAgent settings.
//
// Parameters:
// 1. config : map[string]string
//    Refers to the ImageStore config
//
// Returns:
// 1. *InMemory
//    Returns the InMemory instance
// 2. error
//    Returns an error message if initialization fails.
func NewInmemory(config map[string]string) (*InMemory, error) {
	inMemoryType := config["InMemory"]
	if memoryDB == inMemoryType {
		redisConnect, err := redis.NewRedisConnect(config)

		if err != nil {
			glog.Errorf("Redis connect failed, %v", err)
			return nil, err
		}
		return &InMemory{memType: inMemoryType, redisConnect: redisConnect}, nil
	} else {
		msg := "Currently the memory type: " + inMemoryType + " is not supported"
		glog.Errorf(msg)
		err := errors.New(msg)
		return nil, err
	}

}

// Read is used to read the data from InMemory.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be read.
//
// Returns:
// 1. *io.Reader
//    Returns an instance of io.Reader object of the consolidated image handle.
// 2. error
//    Returns an error message if read fails.
func (pInMemory *InMemory) Read(keyname string) (*io.Reader, error) {
	return pInMemory.redisConnect.Read(keyname)
}

// Remove is used to remove the stored data from InMemory.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be removed.
//
// Returns:
// 1. error
//    Returns an error message if remove fails.
func (pInMemory *InMemory) Remove(keyname string) error {
	return pInMemory.redisConnect.Remove(keyname)
}

// Store  is used to store the data in InMemory.
//
// Parameters:
// 1. value : []byte
//    Refers to the image buffer to be stored in ImageStore.
//
// Returns:
// 1. string
//    Returns the image handle of the image stored.
// 2. error
//    Returns an error message if store fails.
func (pInMemory *InMemory) Store(value []byte) (string, error) {
	return pInMemory.redisConnect.Store(value)
}

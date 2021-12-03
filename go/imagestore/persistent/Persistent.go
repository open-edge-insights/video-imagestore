/*
Copyright (c) 2021 Intel Corporation

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Package persistent provides concrete implementation of persistent storage types like minio etc., for ImageStore
package persistent

import (
	"imagestore/go/imagestore/persistent/minio"
	"errors"
	"io"
	"strings"
	"github.com/golang/glog"
)

// Storage is a struct used to comprise all Persistent package methods in it's scope
type Storage interface {
	// Read the given key from stroage
	Read(keyname string) (io.ReadCloser, error)

	// Remove the given key from the storage
	Remove(keyname string) error

	// Store the given byte array to the storage and return the key under which
	// it is not being stored.
	Store(data []byte, key string) (string, error)
}

// Persistent storage structure
type Persistent struct {
	storage Storage
}

// MINIO is used for module level check with memory type
const MINIO string = "minio"

// NewPersistent is used to initialize the connection based on DataAgent settings
//
// Parameters:
// 1. storageType : string
//    Returns the ImageStore storage type.
// 2. config : map[string]string
//    Refers to the persistent config.
//
// Returns:
// 1. *Persistent
//    Returns the Persistent instance
// 2. error
//    Returns an error object if initialization fails.
func NewPersistent(storageType string, config map[string]string) (*Persistent, error) {
	if storageType == MINIO {
		storage, err := minio.NewMinioStorage(config)
		if err != nil {
			glog.Errorf("Error initializing Minio storage: %v", err)
			return nil, err
		}

		return &Persistent{storage: storage}, nil
	}
	msg := "Persistent storage type not supported: " + storageType
	glog.Errorf(msg)
	err := errors.New(msg)
	return nil, err
}

// GetConfgKey is used to get the key to retrieve the configuration from gRPC.
// Parameters:
// 1. storageType : string
//    Returns the Imagestore storage type.
//
// Returns:
// 1. string
//    Returns consolidated config key based on storage type.
// 2. error
//    Returns an error object if fetching config fails.
func GetConfgKey(storageType string) (string, error) {
	storageType = strings.ToLower(storageType)

	if storageType == "minio" {
		return "MinioCfg", nil
	}

	return "", errors.New("Unknown persistent storage type")
}

// Read is used to read the data from Persistent memory.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be read.
//
// Returns:
// 1. *io.Reader
//    Returns an instance of io.Reader object of the consolidated image handle.
// 2. error
//    Returns an error object if read fails.
func (pStorage *Persistent) Read(keyname string) (io.ReadCloser, error) {
	return pStorage.storage.Read(keyname)
}

// Remove is used to remove the stored data from Persistent memory.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be removed.
//
// Returns:
// 1. error
//    Returns an error object if remove fails.
func (pStorage *Persistent) Remove(keyname string) error {
	return pStorage.storage.Remove(keyname)
}

// Store  is used to store the data in Persistent memory.
//
// Parameters:
// 1. data : []byte
//    Refers to the image buffer to be stored in ImageStore.
//
// Returns:
// 1. string
//    Returns the image handle of the image stored.
// 2. error
//    Returns an error object if store fails.
func (pStorage *Persistent) Store(data []byte, key string) (string, error) {
	return pStorage.storage.Store(data, key)
}

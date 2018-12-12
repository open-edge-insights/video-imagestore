/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

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
package Persistent

import (
	"ElephantTrunkArch/ImageStore/go/ImageStore/Persistent/minio"
	"errors"
	"strings"

	"github.com/golang/glog"
)

// Generic interface for underlying storage technologies to implement.
type Storage interface {
	// Read the given key from stroage
	Read(keyname string) (string, error)

	// Remove the given key from the storage
	Remove(keyname string) error

	// Store the given byte array to the storage and return the key under which
	// it is not being stored.
	Store(data []byte) (string, error)
}

// Persistent storage structure
type Persistent struct {
	storage Storage
}

// Consts for string names of underlying storage types
const MINIO string = "minio"

func NewPersistent(storageType string, config map[string]string) (*Persistent, error) {
	if storageType == MINIO {
		storage, err := minio.NewMinioStorage(config)
		if err != nil {
			glog.Errorf("Error initializing Minio storage: %v", err)
			return nil, err
		}

		return &Persistent{storage: storage}, nil
	} else {
		msg := "Persistent storage type not supported: " + storageType
		glog.Errorf(msg)
		err := errors.New(msg)
		return nil, err
	}
}

// NewPersistentMinimal : Initialize a minimal version of the persistent storage.
// This is designed to allow disabling longer initialization tasks that are not
// required by a certain client, such as enforcing a retention policy.
func NewPersistentMinimal(storageType string, config map[string]string) (*Persistent, error) {
	if storageType == MINIO {
		storage, err := minio.NewMinioStorageMinimal(config)
		if err != nil {
			glog.Errorf("Error initializing Minio storage: %v", err)
			return nil, err
		}

		return &Persistent{storage: storage}, nil
	} else {
		msg := "Persistent storage type not supported: " + storageType
		glog.Errorf(msg)
		err := errors.New(msg)
		return nil, err
	}
}

// Get the key to retrieve the configuration from gRPC
func GetConfgKey(storageType string) (string, error) {
	storageType = strings.ToLower(storageType)

	if storageType == "minio" {
		return "MinioCfg", nil
	}

	return "", errors.New("Unknown persistent storage type")
}

// Retrieve object from the image store with the given object name
// (i.e. keyname)
func (pStorage *Persistent) Read(keyname string) (string, error) {
	return pStorage.storage.Read(keyname)
}

// Remove an object from the persistent image store
func (pStorage *Persistent) Remove(keyname string) error {
	return pStorage.storage.Remove(keyname)
}

// Store a new object in the persistent image store
func (pStorage *Persistent) Store(data []byte) (string, error) {
	return pStorage.storage.Store(data)
}

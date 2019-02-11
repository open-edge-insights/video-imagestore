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

// Package Persistent is a wrapper around Read, Remove and Store APIs of ImageStore.
package Persistent

import (
	"IEdgeInsights/ImageStore/go/ImageStore/Persistent/minio"
	"errors"
	"io"
	"strings"

	"github.com/golang/glog"
)

// Storage is a struct used to comprise all Persistent package methods in it's scope
type Storage interface {
	// Read the given key from stroage
	Read(keyname string) (*io.Reader, error)

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

// MINIO is used for module level check with memory type
const MINIO string = "minio"

// NewPersistent is used to initialize the connection based on DataAgent settings
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

// GetConfgKey is used to get the key to retrieve the configuration from gRPC
func GetConfgKey(storageType string) (string, error) {
	storageType = strings.ToLower(storageType)

	if storageType == "minio" {
		return "MinioCfg", nil
	}

	return "", errors.New("Unknown persistent storage type")
}

// Read is used to read the data from Persistent memory.
//
// It accepts keyname as input.
//
// It returns the image of the consolidated keyname.
func (pStorage *Persistent) Read(keyname string) (*io.Reader, error) {
	return pStorage.storage.Read(keyname)
}

// Remove is used to remove the data from Persistent memory.
//
// It accepts keyname as input.
//
// It returns an error if removing the consolidated image fails.
func (pStorage *Persistent) Remove(keyname string) error {
	return pStorage.storage.Remove(keyname)
}

// Store is used to store the data in Persistent memory.
//
// It accepts value of image to be stored as input.
//
// It returns image handle of image stored.
func (pStorage *Persistent) Store(data []byte) (string, error) {
	return pStorage.storage.Store(data)
}

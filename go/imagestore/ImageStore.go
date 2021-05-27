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

// Package imagestore exports Read, Remove and Store APIs.
package imagestore

import (
	persistent "IEdgeInsights/ImageStore/go/imagestore/persistent"
	"io"
	"github.com/golang/glog"
	common "IEdgeInsights/ImageStore/common"
)

// ImageStore :  ImageStore is a struct, used for store & retrieve operations.
type ImageStore struct {
	storageType       string
	persistentStorage *(persistent.Persistent)
}

// NewImageStore : This is the Constructor type method which initialises the Object for ImageStore Operations
//
// Returns:
// 1. ImageStore object
//	  Returns an ImageStore object with config.
// 2. error
//	  Returns an error object if initialization fails.
func NewImageStore(securityDisable bool) (*ImageStore, error) {

	//TODO: This call is failing when trying to connect to gRPC server running in the same container.
	var err error
	storageConfig := make(map[string]string)
	storageConfig["Host"] = common.MinioHost
	storageConfig["Port"] = common.MinioPort
	persistentStorage, err := persistent.NewPersistent("MINIO", storageConfig)
	if err != nil {
		glog.Errorf("Error initializing persistent memory storage: %v", err)
		return nil, err
	}

	return &ImageStore{storageType: "", persistentStorage: persistentStorage}, nil
}

// GetImageStoreInstance is the constructor type method which takes the image store config
// and initialises the Object for ImageStore Operations.
//
// Parameters:
// 1. cfg : map[string]string
// 2. persistCfg : map[string]string
//    Refers to the minio(persistent) config
//
// Returns:
// 1. *ImageStore
//    Returns the ImageStore instance
// 2. error
//    Returns an error object if initialization fails.
func GetImageStoreInstance(persistCfg map[string]string) (*ImageStore, error) {
	persistentStorage, err := persistent.NewPersistent("minio", persistCfg)
	if err != nil {
		return nil, err
	}

	return &ImageStore{storageType: "", persistentStorage: persistentStorage}, nil
}

// Read is used to read the stored data from memory.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be read.
//
// Returns:
// 1. *io.Reader
//    Returns the image of the consolidated image handle.
// 2. error
//    Returns an error object if read fails.
func (pImageStore *ImageStore) Read(keyname string) (io.ReadCloser, error) {
	return pImageStore.persistentStorage.Read(keyname)
}

// Remove is used to remove the stored data from memory.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be removed.
//
// Returns:
// 1. error
//    Returns an error object if remove fails.
func (pImageStore *ImageStore) Remove(keyname string) error {
	return pImageStore.persistentStorage.Remove(keyname)
}

// Store  is used to store the data in selected memory based on SetStorageType API.
//
// Parameters:
// 1. value : []byte
//    Refers to the image buffer to be stored in ImageStore.
//
// Returns:
// 1. string
//    Returns the image handle of the image stored.
// 2. error
//    Returns an error object if store fails.
func (pImageStore *ImageStore) Store(value []byte, keyname string) (string, error) {
	return pImageStore.persistentStorage.Store(value, keyname)
}

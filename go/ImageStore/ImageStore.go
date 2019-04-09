/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Package ImageStore exports Read, Remove, Store and SetStorageType APIs.
package ImageStore

import (
	client "IEdgeInsights/DataAgent/da_grpc/client/go/client_internal"
	inmemory "IEdgeInsights/ImageStore/go/ImageStore/InMemory"
	persistent "IEdgeInsights/ImageStore/go/ImageStore/Persistent"
	"errors"
	"io"
	"strings"

	"github.com/golang/glog"
)

// inMemKeyPattern is the key pattern used for inmemory key reference
const inMemKeyPattern string = "inmem"

// persistKeyPattern is the key pattern used for persistent key references.
const persistKeyPattern string = "persist"

// ImageStore :  ImageStore is a struct, used for store & retrieve operations.
type ImageStore struct {
	storageType       string
	inMemory          *(inmemory.InMemory) //TODO: This should actually be an interface referring to respective concrete classes
	persistentStorage *(persistent.Persistent)
}

// grpc client certificates
const (
	RootCA     = "/etc/ssl/grpc_int_ssl_secrets/ca_certificate.pem"
	ClientCert = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_certificate.pem"
	ClientKey  = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_key.pem"
)

// NewImageStore is the constructor type method which initialises the object for ImageStore Operations.
//
// Returns:
// 1. ImageStore object
//    Returns an ImageStore object with config.
// 2. error
//    Returns an error object if initialization fails.
func NewImageStore() (*ImageStore, error) {

	//TODO: This call is failing when trying to connect to gRPC server running in the same container.
	grpcClient, err := client.NewGrpcInternalClient(ClientCert, ClientKey, RootCA, "localhost", "50052")
	if err != nil {
		glog.Errorf("Error connecting to gRPC: %v", err)
		return nil, err
	}

	config, err := grpcClient.GetConfigInt("RedisCfg")
	if err != nil {
		glog.Errorf("GetConfigInt() Error:%v", err)
		return nil, err
	}

	persistConfig, err := grpcClient.GetConfigInt("PersistentImageStore")
	if err != nil {
		glog.Errorf("Error retrieving persistent storage config: %v", err)
		return nil, err
	}

	config["InMemory"] = "redis"
	inMemory, err := inmemory.NewInmemory(config)
	if err != nil {
		return nil, err
	}

	configKey, err := persistent.GetConfgKey(persistConfig["Type"])
	if err != nil {
		return nil, err
	}

	storageConfig, err := grpcClient.GetConfigInt(configKey)
	if err != nil {
		return nil, err
	}

	persistentStorage, err := persistent.NewPersistent(persistConfig["Type"], storageConfig)
	if err != nil {
		glog.Errorf("Error initializing persistent memory storage: %v", err)
		return nil, err
	}

	return &ImageStore{storageType: "", inMemory: inMemory, persistentStorage: persistentStorage}, nil
}

// GetImageStoreInstance is the constructor type method which takes the image store config
// and initialises the Object for ImageStore Operations.
//
// Parameters:
// 1. cfg : map[string]string
//    Refers to the redis(inmemory) config
// 2. persistCfg : map[string]string
//    Refers to the minio(persistent) config
//
// Returns:
// 1. *ImageStore
//    Returns the ImageStore instance
// 2. error
//    Returns an error object if initialization fails.
func GetImageStoreInstance(cfg map[string]string, persistCfg map[string]string) (*ImageStore, error) {
	cfg["InMemory"] = "redis"
	inMemory, err := inmemory.NewInmemory(cfg)
	if err != nil {
		return nil, err
	}

	persistentStorage, err := persistent.NewPersistent("minio", persistCfg)
	if err != nil {
		return nil, err
	}

	return &ImageStore{storageType: "", inMemory: inMemory, persistentStorage: persistentStorage}, nil
}

// SetStorageType sets the storageType for Write Operation to store on selected memory.
//
// Parameters:
// 1. memoryType : string
//    Refers to the storage type.
//    It can either be inmemory or persistent to store the buffer
//    in Redis or Minio respectively.
//
// Returns:
// 1. error
//    Returns an error object if initialization fails.
func (pImageStore *ImageStore) SetStorageType(memoryType string) error {
	memoryType = strings.ToLower(memoryType)

	if memoryType == "inmemory" {
		pImageStore.storageType = memoryType
		return nil
	} else if memoryType == "persistent" {
		pImageStore.storageType = memoryType
		return nil
	} else if memoryType == "inmemory_persistent" {
		pImageStore.storageType = memoryType
		return nil
	}
	return errors.New("MemoryType: " + memoryType + " not supported")
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
func (pImageStore *ImageStore) Read(keyname string) (*io.Reader, error) {
	if strings.Contains(keyname, inMemKeyPattern) {
		return pImageStore.inMemory.Read(keyname)
	} else if strings.Contains(keyname, persistKeyPattern) {
		return pImageStore.persistentStorage.Read(keyname)
	}
	return nil, errors.New("Unknown key pattern for key: " + keyname)
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

	if strings.Contains(keyname, inMemKeyPattern) {
		return pImageStore.inMemory.Remove(keyname)
	} else if strings.Contains(keyname, persistKeyPattern) {
		return pImageStore.persistentStorage.Remove(keyname)
	}
	return errors.New("Unknown key pattern for key: " + keyname)
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
func (pImageStore *ImageStore) Store(value []byte) (string, error) {
	if pImageStore.storageType == "inmemory" {
		return pImageStore.inMemory.Store(value)
	} else if pImageStore.storageType == "persistent" {
		return pImageStore.persistentStorage.Store(value)
	} else if pImageStore.storageType == "inmemory_persistent" {
		inmemHandle, err := pImageStore.inMemory.Store(value)
		if err != nil {
			return "", err
		}
		persistHandle, err := pImageStore.persistentStorage.Store(value)
		if err != nil {
			return "", err
		}
		inmemHandle += "|" + persistHandle
		return inmemHandle, nil
	}
	return "", errors.New("Memory type: " + pImageStore.storageType + " is not supported. Please set it before using ImageStore.SetStorageType API")
}

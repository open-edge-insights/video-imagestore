/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package ImageStore

import (
	client "ElephantTrunkArch/DataAgent/da_grpc/client/go"
	inmemory "ElephantTrunkArch/ImageStore/go/ImageStore/InMemory"
	"errors"
	"strings"

	"github.com/golang/glog"
)

// This const pattern used for inmemory key reference
const keyPattern string = "inmem"

// ImageStore :  ImageStore is a struct, used for store & retrieve operations
type ImageStore struct {
	storageType string
	inMemory    *(inmemory.InMemory) //TODO: This should actually be an interface referring to respective concrete classes
}

// NewImageStore : This is the Constructor type method which initialises the Object for ImageStore Operations
func NewImageStore() (*ImageStore, error) {

	//TODO: This call is failing when trying to connect to gRPC server running in the same container.
	grpcClient, err := client.NewGrpcClient("localhost", "50051")
	config, err := grpcClient.GetConfigInt("RedisCfg")

	if err != nil {
		glog.Errorf("GetConfigInt() Error:%v", err)
		return nil, err
	}

	config["InMemory"] = "redis"
	inMemory, err := inmemory.NewInmemory(config)
	if err != nil {
		return nil, err
	}
	return &ImageStore{storageType: "", inMemory: inMemory}, nil
}

// GetImageStoreInstance : This is the Constructor type method which takes the image store config
// and initialises the Object for ImageStore Operations
func GetImageStoreInstance(cfg map[string]string) (*ImageStore, error) {
	cfg["InMemory"] = "redis"
	inMemory, err := inmemory.NewInmemory(cfg)
	if err != nil {
		return nil, err
	}
	return &ImageStore{storageType: "", inMemory: inMemory}, nil
}

// SetStorageType : This helps to set the storageType for Write Operation to store on Selected Memory.
func (pImageStore *ImageStore) SetStorageType(memoryType string) error {
	memoryType = strings.ToLower(memoryType)

	if memoryType == "inmemory" {
		pImageStore.storageType = memoryType
		return nil
	}
	return errors.New("MemoryType: " + memoryType + " not supported")
}

// Read : This helps to read the stored data from memory. It accepts keyname as input.
func (pImageStore *ImageStore) Read(keyname string) (string, error) {
	if strings.Contains(keyname, keyPattern) {
		return pImageStore.inMemory.Read(keyname)
	}
	return "", errors.New("keyname " + keyname + " doesn't have " + "keypattern " + keyPattern)
}

// Remove : This helps to remove the stored data from memory. It accepts keyname as input
func (pImageStore *ImageStore) Remove(keyname string) error {

	if strings.Contains(keyname, keyPattern) {
		return pImageStore.inMemory.Remove(keyname)
	}
	return errors.New("keyname " + keyname + " doesn't have " + "keypattern " + keyPattern)
}

// Store : This helps to persist the data in selected memory based SetStorageType API. This Accepts value as input
func (pImageStore *ImageStore) Store(value []byte) (string, error) {
	if pImageStore.storageType == "inmemory" {
		return pImageStore.inMemory.Store(value)
	}
	return "", errors.New("Memory type: " + pImageStore.storageType + " is not supported. Please set it before using ImageStore.SetStorageType API")
}

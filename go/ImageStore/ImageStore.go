package ImageStore

import (
	client "iapoc_elephanttrunkarch/DataAgent/da_grpc/client"
	inmemory "iapoc_elephanttrunkarch/ImageStore/go/ImageStore/InMemory"
	"strings"

	"github.com/golang/glog"
)

// This const pattern used for inmemory key reference
const keyPattern string = "inmem"

var inMemory *(inmemory.InMemory)

// ImageStore :  ImageStore is a struct, used for store & retrieve operations
type ImageStore struct {
	storageType string
}

// NewImageStore : This is the Constructor type method which initialises the Object for ImageStore Operations
func NewImageStore() (*ImageStore, error) {

	config, err := client.GetConfigInt("RedisCfg")

	if err != nil {
		return nil, err
	}

	config["InMemory"] = "redis"

	inMemory, err = inmemory.NewInmemory(config)

	if err != nil {
		return nil, err
	}
	return &ImageStore{storageType: "none"}, nil
}

// SetStorageType : This helps to set the storageType for Write Operation to store on Selected Memory.
func (pImageStore *ImageStore) SetStorageType(memoryType string) (bool, string) {
	status, message := false, "FAILED"
	memoryType = strings.ToLower(memoryType)
	if memoryType != "" {
		if memoryType == "inmemory" {
			pImageStore.storageType = memoryType
			status, message = true, "Selected Storage Type : "+memoryType
		} else {
			glog.Info("Not Suppported ", memoryType)
			status, message = false, "Failed: Currently it not supports : "+memoryType
		}
	} else {
		glog.Info("Not Suppported ", memoryType)
		status, message = false, "Failed: Currently it not supports : "+memoryType
	}

	return status, message
}

// Read : This helps to read the stored data from memory. It accepts keyname as input.
func (pImageStore *ImageStore) Read(keyname string) (bool, string) {
	status, message := false, "FAILED"
	if strings.Contains(keyname, keyPattern) {
		status, message = inMemory.GetDataFromInmemory(keyname)
	} else {
		status, message = false, "FAILED : Invalid Key"
	}

	return status, message
}

// Remove : This helps to remove the stored data from memory. It accepts keyname as input
func (pImageStore *ImageStore) Remove(keyname string) (bool, string) {
	status, message := false, "FAILED"

	if strings.Contains(keyname, keyPattern) {
		status, message = inMemory.RemoveDataFromInmemory(keyname)
	} else {
		status, message = false, "FAILED"
	}

	return status, message
}

// Store : This helps to persist the data in selected memory based SetStorageType API. This Accepts value as input
func (pImageStore *ImageStore) Store(value []byte) (bool, string) {
	status, message := false, "FAILED"
	if pImageStore.storageType == "inmemory" {
		status, message = inMemory.StoreDatainInmemory(value)
	} else {
		glog.Info("Not Suppported : ", " Please set the StorageType using SetStorageType API")
	}

	return status, message
}

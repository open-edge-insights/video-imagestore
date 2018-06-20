package InMemory

import (
	"errors"
	"iapoc_elephanttrunkarch/ImageStore/go/ImageStore/InMemory/redis"
)

// InMemory : This struct is used to comprise all Inmemory methods in it's scope
type InMemory struct {
	InMemoryType string
}

// memoryDB : This should be used for Module Level Check with Memory Type
const memoryDB string = "redis"

// Global Declaration
var redisConnect *(redis.RedisConnect)

// NewInmemory : This method used to initialize the connection based on dataAgent settings
func NewInmemory(config map[string]string) (*InMemory, error) {

	// status, message := false, "FAILED"
	var err error
	err = nil
	inMemoryType := config["InMemory"]
	if memoryDB == inMemoryType {
		redisConnect, err = redis.NewRedisConnect(config)

		if err != nil {
			return nil, err
		}
		return &InMemory{InMemoryType: inMemoryType}, nil
	} else {
		err = errors.New("Currently System Not Supports ")
		return nil, err
	}

}

// GetDataFromInmemory : This helps to read the data from InMemory, It Accepts keyname as input
func (pInMemory *InMemory) GetDataFromInmemory(keyname string) (bool, string) {
	status, message := false, "FAILED"
	if pInMemory.InMemoryType == memoryDB {
		status, message = redisConnect.GetDataFromRedis(keyname)
	} else {
		status, message = false, "FAILED"
	}

	return status, message
}

// RemoveDataFromInmemory : This helps to Remove the data from InMemory, It Accepts keyname as input
func (pInMemory *InMemory) RemoveDataFromInmemory(keyname string) (bool, string) {
	status, message := false, "FAILED"
	if pInMemory.InMemoryType == memoryDB {
		status, message = redisConnect.RemoveFromRedis(keyname)
	} else {
		status, message = false, "FAILED  fdfd"
	}

	return status, message
}

// StoreDatainInmemory : This helps to store the data in InMemory, It Accepts value as input
func (pInMemory *InMemory) StoreDatainInmemory(value []byte) (bool, string) {
	status, message := false, "FAILED"

	if pInMemory.InMemoryType == memoryDB {
		status, message = redisConnect.StoreDatainRedis(value)
	} else {
		status, message = false, "FAILED"
	}

	return status, message
}

package InMemory

import (
	"errors"
	"iapoc_elephanttrunkarch/ImageStore/go/ImageStore/InMemory/redis"

	"github.com/golang/glog"
)

// InMemory : This struct is used to comprise all Inmemory methods in it's scope
type InMemory struct {
	memType      string
	redisConnect *(redis.RedisConnect) //TODO: This should actually be an interface referring to respective concrete classes
}

// memoryDB : This should be used for Module Level Check with Memory Type
const memoryDB string = "redis"

// NewInmemory : This method used to initialize the connection based on dataAgent settings
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

// Read : This helps to read the data from InMemory, It Accepts keyname as input
func (pInMemory *InMemory) Read(keyname string) (string, error) {
	return pInMemory.redisConnect.Read(keyname)
}

// Remove : This helps to Remove the data from InMemory, It Accepts keyname as input
func (pInMemory *InMemory) Remove(keyname string) error {
	return pInMemory.redisConnect.Remove(keyname)
}

// Store : This helps to store the data in InMemory, It Accepts value as input
func (pInMemory *InMemory) Store(value []byte) (string, error) {
	return pInMemory.redisConnect.Store(value)
}

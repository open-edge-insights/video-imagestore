package redis

import (
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/google/uuid"
)

var client *redis.Client

// RedisConnect : This Struct used to have default variables used for redis. Also to comprise methods of redis to it's scope
type RedisConnect struct {
	Retention string
}

// NewRedisConnect : This is a constructor function to connect redis database
func NewRedisConnect(config map[string]string) (*RedisConnect, error) {

	client = redis.NewClient(&redis.Options{
		Addr:     config["Host"] + ":" + config["Port"],
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	return &RedisConnect{Retention: config["Retention"]}, err
}

// GetDataFromRedis : This helps to read the data from Redis, It Accepts keyname as input
func (pRedisConnect *RedisConnect) GetDataFromRedis(keyname string) (bool, string) {

	status, message := false, "FAILED"

	binarydata, err := client.Get(keyname).Result()
	if err == redis.Nil {
		status, message = false, "Key Not Found"
	} else if err != nil {
		status, message = false, err.Error()
	} else {
		status, message = true, binarydata
	}

	return status, message
}

// RemoveFromRedis : This helps to remove the data from Redis, It Accepts keyname as input
func (pRedisConnect *RedisConnect) RemoveFromRedis(keyname string) (bool, string) {

	status, message := false, "FAILED"

	removecode, err := client.Del(keyname).Result()
	if removecode == 1 {
		status, message = true, "Return Code : "+string(removecode)
	} else if removecode == 0 {
		status, message = false, "Not Removed , Key Not Found"
	} else if err != nil {
		status, message = false, err.Error()
	}

	return status, message
}

// generateKeyName : This used to generate the keyname
func (pRedisConnect *RedisConnect) generateKeyName() string {
	keyname := "inmem_" + uuid.New().String()[:8]
	return keyname
}

// StoreDatainRedis : This helps to store the data in redis, It Accepts value as input
func (pRedisConnect *RedisConnect) StoreDatainRedis(value []byte) (bool, string) {
	var err error
	status, message := false, "FAILED"

	ttl, err := time.ParseDuration(pRedisConnect.Retention)
	if err != nil {
		status, message = false, err.Error()

	} else {
		keyname := pRedisConnect.generateKeyName()

		if ttl >= 1 {
			err = client.Set(keyname, value, ttl).Err()
		} else {
			err = client.Set(keyname, value, 0).Err()
		}

		if err != nil {
			status, message = false, err.Error()
		} else {
			status, message = true, keyname
		}
	}

	return status, message
}

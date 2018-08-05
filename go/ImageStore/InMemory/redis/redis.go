package redis

import (
	"time"
	"errors"
	"github.com/go-redis/redis"
	uuid "github.com/google/uuid"
)

var client *redis.Client

// RedisConnect : This Struct used to have default variables used for redis. Also to comprise methods of redis to it's scope
type RedisConnect struct {
	retention string
}

// NewRedisConnect : This is a constructor function to connect redis database
func NewRedisConnect(config map[string]string) (*RedisConnect, error) {

	client = redis.NewClient(&redis.Options{
		Addr:     config["Host"] + ":" + config["Port"],
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	return &RedisConnect{retention: config["Retention"]}, err
}

// Read : This helps to read the data from Redis, It Accepts keyname as input
func (pRedisConnect *RedisConnect) Read(keyname string) (string, error) {
	binarydata, err := client.Get(keyname).Result()
	if err == redis.Nil {
		return "", errors.New("Key Not Found")
	} else if err != nil {
		return "", err
	}
	return binarydata, err
}

// Remove : This helps to remove the data from Redis, It Accepts keyname as input
func (pRedisConnect *RedisConnect) Remove(keyname string) error {
	_, err := client.Del(keyname).Result()
	if err != nil {
		return err
	} 
	return nil
}

// generateKeyName : This used to generate the keyname
func (pRedisConnect *RedisConnect) generateKeyName() string {
	keyname := "inmem_" + uuid.New().String()[:8]
	return keyname
}

// Store : This helps to store the data in redis, It Accepts value as input
func (pRedisConnect *RedisConnect) Store(value []byte) (string, error) {
	
	ttl, err := time.ParseDuration(pRedisConnect.retention)
	if err != nil {
		return "", err

	} else {
		keyname := pRedisConnect.generateKeyName()

		if ttl >= 1 {
			err = client.Set(keyname, value, ttl).Err()
		} else {
			err = client.Set(keyname, value, 0).Err()
		}

		if err != nil {
			return "", err
		} else {
			return keyname, nil
		}
	}
}

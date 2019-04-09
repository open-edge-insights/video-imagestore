/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Package redis exports the Read, Remove & Store of APIs of Redis DB.
package redis

import (
	"errors"
	"io"
	"strings"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/google/uuid"
)

var client *redis.Client

// RedisConnect is a struct used to have default variables used for redis and to comprise methods of redis to it's scope
type RedisConnect struct {
	retention string
}

// NewRedisConnect is a constructor function to connect to redis database.
//
// Parameters:
// 1. config : map[string]string
//    Refers to the redis config.
//
// Returns:
// 1. *RedisConnect
//    Returns the RedisConnect instance
// 2. error
//    Returns an error object if initialization fails.
func NewRedisConnect(config map[string]string) (*RedisConnect, error) {

	client = redis.NewClient(&redis.Options{
		Addr:     config["Host"] + ":" + config["Port"],
		Password: config["Password"],
		DB:       0,
	})
	_, err := client.Ping().Result()
	return &RedisConnect{retention: config["Retention"]}, err
}

// Read is used to read the data from Redis.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be read.
//
// Returns:
// 1. io.Reader
//    Returns an instance of io.Reader object of the consolidated image handle.
// 2. error
//    Returns an error object if read fails.
func (pRedisConnect *RedisConnect) Read(keyname string) (*io.Reader, error) {
	binarydata, err := client.Get(keyname).Result()
	if err == redis.Nil {
		return nil, errors.New("Key Not Found")
	} else if err != nil {
		return nil, err
	}
	outputStr := strings.NewReader(binarydata)
	data := io.Reader(outputStr)
	return &data, nil
}

// Remove is used to remove the stored data from Redis.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be removed.
//
// Returns:
// 1. error
//    Returns an object object if remove fails.
func (pRedisConnect *RedisConnect) Remove(keyname string) error {
	_, err := client.Del(keyname).Result()
	if err != nil {
		return err
	}
	return nil
}

// generateKeyName is used to generate the keyname
//
// Returns:
// 1. string
//    Returns an unique uuid.
func (pRedisConnect *RedisConnect) generateKeyName() string {
	keyname := "inmem_" + uuid.New().String()[:8]
	return keyname
}

// Store  is used to store the data in Redis.
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
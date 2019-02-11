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

// Package minio exports the Read, Remove & Store of APIs of Minio DB.
package minio

import (
	"bytes"
	"errors"
	"io"

	"github.com/golang/glog"
	uuid "github.com/google/uuid"
	minio "github.com/minio/minio-go"
)

// Constant bucket name in Minio, may want to make this variable in the future
const bucketName string = "image-store-bucket"

// Constant for the region in Minio
const region string = "gateway"

// MinioStorage is a struct used to have default variables used for minio and to comprise methods of minio to it's scope
type MinioStorage struct {
	client *(minio.Client)
}

// missingKeyError is helper method for reporting a missing key in the Minio configuration
func missingKeyError(key string) error {
	msg := "Minio config missing key: " + key
	glog.Errorf(msg)
	return errors.New(msg)
}

func initClient(config map[string]string) (*minio.Client, error) {
	glog.Infof("Pulling out config values")
	host, ok := config["Host"]
	if !ok {
		return nil, missingKeyError("Host")
	}

	port, ok := config["Port"]
	if !ok {
		return nil, missingKeyError("Port")
	}

	accessKey, ok := config["AccessKey"]
	if !ok {
		return nil, missingKeyError("AccessKey")
	}

	secretKey, ok := config["SecretKey"]
	if !ok {
		return nil, missingKeyError("SecretKey")
	}

	sslStr, ok := config["Ssl"]
	if !ok {
		return nil, missingKeyError("Ssl")
	}

	ssl := true
	if sslStr == "true" {
		ssl = true
	} else if sslStr == "false" {
		ssl = false
	} else {
		msg := "Ssl key in Minio config must be true or false, not :" + sslStr
		glog.Errorf(msg)
		return nil, errors.New(msg)
	}

	glog.Infof("Config: Host=%s, Port=%s, ssl=%v", host, port, ssl)

	glog.Infof("Initializing Minio client")
	client, err := minio.NewWithRegion(
		host+":"+port, accessKey, secretKey, ssl, region)
	if err != nil {
		glog.Errorf("Failed to connect to Minio server: %v", err)
		return nil, err
	}

	// Check if the bucket exists
	glog.Infof("Checking if Minio bucket already exists")
	found, err := client.BucketExists(bucketName)
	if err != nil {
		glog.Errorf("Failed to verify existence of bucket: %v", err)
		return nil, err
	}

	if !found {
		// Create the bucket if it does not exist
		glog.Infof("Creating bucket")
		client.MakeBucket(bucketName, region)
	}

	return client, nil
}

// NewMinioStorage is used to create a new instance of the MinioStorage
func NewMinioStorage(config map[string]string) (*MinioStorage, error) {
	client, err := initClient(config)
	if err != nil {
		// Error has already been logged
		return nil, err
	}

	minioStorage := &MinioStorage{client: client}

	glog.Infof("Initialization finished")
	return minioStorage, nil
}

// Read is used to read the stored data from Minio.
//
// It accepts keyname as input.
//
// It returns the image of consolidated keyname.
func (pMinioStorage *MinioStorage) Read(keyname string) (*io.Reader, error) {
	// Get the object from the store
	obj, err := pMinioStorage.client.GetObject(
		bucketName, keyname, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	data := io.Reader(obj)
	return &data, nil
}

// Remove is used to remove the data from Minio.
//
// It accepts keyname as input.
//
// It returns an error if removing the consolidated image fails.
func (pMinioStorage *MinioStorage) Remove(keyname string) error {
	return pMinioStorage.client.RemoveObject(bucketName, keyname)
}

// Store  is used to store the data in Minio.
//
// It accepts value to be stored as parameter.
//
// It returns image handle of respective image stored.
func (pMinioStorage *MinioStorage) Store(data []byte) (string, error) {
	key := generateKeyName()
	buffer := bytes.NewBuffer(data)
	buffLen := int64(buffer.Len())
	n, err := pMinioStorage.client.PutObject(bucketName, key, buffer,
		buffLen, minio.PutObjectOptions{})
	if err != nil {
		glog.Errorf("Failed to put object into Minio: %v", err)
		return "", err
	}
	if n < buffLen {
		msg := "Failed to push all of the bytes to Minio"
		glog.Errorf(msg)
		return "", errors.New(msg)
	}

	return key, nil
}

// generateKeyName is used to generate the keyname
func generateKeyName() string {
	keyname := "persist_" + uuid.New().String()[:8]
	return keyname
}

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

// Max number of buffers in the channel and workers for consuming it
const (
	maxBuffers int = 100
	maxWorkers int = 100
)

// Struct for holding the buffers for the store workers
type DataBuffer struct {
	buffer []byte
	key    string
}

// MinioStorage is a struct used to have default variables used for minio and to comprise methods of minio to it's scope
type MinioStorage struct {
	client   *(minio.Client)
	dataChan chan DataBuffer
}

// missingKeyError is helper method for reporting a missing key in the Minio configuration
//
// Parameters:
// 1. key : string
//    Refers to the minio image handle.
//
// Returns:
// 1. error
//    Returns an error object if initialization fails.
func missingKeyError(key string) error {
	msg := "Minio config missing key: " + key
	glog.Errorf(msg)
	return errors.New(msg)
}

// missingKeyError is helper method for reporting a missing key in the Minio configuration
//
// Parameters:
// 1. config : map[string]string
//    Refers to the minio config.
//
// Returns:
// 1. *minio.Client
//    Returns a minio client instance.
// 1. error
//    Returns an error object if initialization fails.
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
// Parameters:
// 1. config : map[string]string
//    Refers to the minio config.
//
// Returns:
// 1. *MinioStorage
//    Returns the MinioStorage instance
// 2. error
//    Returns an error object if initialization fails.
func NewMinioStorage(config map[string]string) (*MinioStorage, error) {
	client, err := initClient(config)
	if err != nil {
		// Error has already been logged
		return nil, err
	}

	// Creating data channel for store workers
	dataChan := make(chan DataBuffer, maxBuffers)

	minioStorage := &MinioStorage{client: client, dataChan: dataChan}

	// Start store workers
	for i := 0; i < maxWorkers; i++ {
		go storeWorker(minioStorage)
	}

	glog.Infof("Initialization finished")
	return minioStorage, nil
}

// Read is used to read the stored data from Minio.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be read.
//
// Returns:
// 1. *io.Reader
//    Returns the ip.Reader instance.
// 2. error
//    Returns an error object if read fails.
func (pMinioStorage *MinioStorage) Read(keyname string) (io.ReadCloser, error) {
	// Get the object from the store
	obj, err := pMinioStorage.client.GetObject(
		bucketName, keyname, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	data := io.ReadCloser(obj)
	return data, nil
}

// Remove is used to remove the data from Minio.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be removed.
//
// Returns:
// 1. error
//    Returns an object object if remove fails.
func (pMinioStorage *MinioStorage) Remove(keyname string) error {
	return pMinioStorage.client.RemoveObject(bucketName, keyname)
}

// Store  is used to store the data in Minio.
//
// Parameters:
// 1. data : []byte
//    Refers to the image buffer to be stored in ImageStore.
//
// Returns:
// 1. string
//    Returns the image handle of the image stored.
// 2. error
//    Returns an error object if store fails.
func (pMinioStorage *MinioStorage) Store(data []byte) (string, error) {
	key := generateKeyName()

	pMinioStorage.dataChan <- DataBuffer{data, key}

	return key, nil
}

// generateKeyName is used to generate the keyname
//
// Returns:
// 1. string
//    Returns an unique uuid.
func generateKeyName() string {
	keyname := "persist_" + uuid.New().String()[:8]
	return keyname
}

// storeWorker is the worker function storing data into the Minio DB
// We start maxWorkers number of workers to ingest data to the DB.
//
// Parameters:
// 1. pMinioStorage : MinioStorage
//    Context of the Minio Image Store
func storeWorker(pMinioStorage *MinioStorage) {

	for {
		buf := <-pMinioStorage.dataChan

		buffer := bytes.NewReader(buf.buffer)
		bufLen := int64(buffer.Len())

		n, err := pMinioStorage.client.PutObject(bucketName, buf.key, buffer,
			bufLen, minio.PutObjectOptions{})
		if err != nil {
			glog.Errorf("Failed to put object into Minio for %s: %v", buf.key, err)
		}
		if n < bufLen {
			glog.Errorf("Failed to push all of the bytes to Minio for key %s", buf.key)
		}
		buffer = nil
		buf.buffer = nil
	}
}

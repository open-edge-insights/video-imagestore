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
package minio

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/golang/glog"
	uuid "github.com/google/uuid"
	minio "github.com/minio/minio-go"
)

// Constant bucket name in Minio, may want to make this variable in the future
const bucketName string = "image-store-bucket"

// Constant for the region in Minio
const region string = "gateway"

// Minio storage abstraction
type MinioStorage struct {
	client        *(minio.Client)
	retentionTime time.Duration
}

// Helper method for reporting a missing key in the Minio configuration
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

	glog.Infof("Config: Host=%s, Port=%d, ssl=%s", host, port, ssl)

	glog.Infof("Initializing Minio client")
	client, err := minio.NewWithRegion(
		host+":"+port, accessKey, secretKey, ssl, region)
	if err != nil {
		glog.Errorf("Failed to connect to Minio server: ", err)
		return nil, err
	}

	// Check if the bucket exists
	glog.Infof("Checking if Minio bucket already exists")
	found, err := client.BucketExists(bucketName)
	if err != nil {
		glog.Errorf("Failed to verify existence of bucket: ", err)
		return nil, err
	}

	if !found {
		// Create the bucket if it does not exist
		glog.Infof("Creating bucket")
		client.MakeBucket(bucketName, region)
	}

	return client, nil
}

// Create a new instance of the MinioStorage
func NewMinioStorage(config map[string]string) (*MinioStorage, error) {
	retentionTimeStr, ok := config["RetentionTime"]
	if !ok {
		return nil, missingKeyError("RetentionTime")
	}

	retentionTime, err := strconv.ParseInt(retentionTimeStr, 10, 64)
	if err != nil {
		msg := "Retention time must be an integer in Minio config"
		glog.Errorf(msg)
		return nil, errors.New(msg)
	}

	client, err := initClient(config)
	if err != nil {
		// Error has already been logged
		return nil, err
	}

	minioStorage := &MinioStorage{
		client:        client,
		retentionTime: time.Duration(retentionTime) * time.Second}

	// Clear out Minio of old data on start, this also starts the timer for
	// the clean up procedure to run again based on the retention time
	glog.Infof("Cleaning store")
	err = minioStorage.cleanStore()
	if err != nil {
		glog.Error("Failed to clear object store: ", err)
		return nil, err
	}

	glog.Infof("Initialization finished")
	return minioStorage, nil
}

// Initialize minio storage which does not attempt to clean up the image store.
// This assumes that the user already initializes the Minio ImageStore in
// other points in the systsm (i.e. video ingestion)
func NewMinioStorageMinimal(config map[string]string) (*MinioStorage, error) {
	retentionTimeStr, ok := config["RetentionTime"]
	if !ok {
		return nil, missingKeyError("RetentionTime")
	}

	retentionTime, err := strconv.ParseInt(retentionTimeStr, 10, 64)
	if err != nil {
		msg := "Retention time must be an integer in Minio config"
		glog.Errorf(msg)
		return nil, errors.New(msg)
	}

	client, err := initClient(config)
	if err != nil {
		// Error has already been logged
		return nil, err
	}

	minioStorage := &MinioStorage{
		client:        client,
		retentionTime: time.Duration(retentionTime) * time.Second}

	return minioStorage, nil
}

// Retrieve the object with the given name from Minio
func (pMinioStorage *MinioStorage) Read(keyname string) (string, error) {
	// Get the object from the store
	obj, err := pMinioStorage.client.GetObject(
		bucketName, keyname, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}

	data := bytes.NewBuffer(nil)

	if _, err = io.Copy(data, obj); err != nil {
		glog.Errorf("Failed to retrieve data from Minio: ", err)
		return "", err
	}

	return data.String(), nil
}

// Remove the objec with the given name from Minio
func (pMinioStorage *MinioStorage) Remove(keyname string) error {
	return pMinioStorage.client.RemoveObject(bucketName, keyname)
}

// Store the given object in Minio
func (pMinioStorage *MinioStorage) Store(data []byte) (string, error) {
	key := generateKeyName()
	buffer := bytes.NewBuffer(data)
	buffLen := int64(buffer.Len())
	n, err := pMinioStorage.client.PutObject(bucketName, key, buffer,
		buffLen, minio.PutObjectOptions{})
	if err != nil {
		glog.Errorf("Failed to put object into Minio: ", err)
		return "", err
	}
	if n < buffLen {
		msg := "Failed to push all of the bytes to Minio"
		glog.Errorf(msg)
		return "", errors.New(msg)
	}

	return key, nil
}

// Clean up the image store
func (pMinioStorage *MinioStorage) cleanStore() error {
	// Channel for objects to be removed from Minio
	objectsCh := make(chan string)
	objectsErrCh := make(chan error, 1)
	defer close(objectsErrCh)

	// Routine to find objects to remove and send them over the `objectsCh`
	go func() {
		// Defer channel close to when the function exits
		defer close(objectsCh)
		now := time.Now()

		for obj := range pMinioStorage.client.ListObjects(bucketName, "", false, nil) {
			if obj.Err != nil {
				glog.Errorf("Failed retrieving objects from Minio: ", obj.Err)
				objectsErrCh <- obj.Err
				return
			}
			elapsed := now.Sub(obj.LastModified)
			if elapsed > pMinioStorage.retentionTime {
				objectsCh <- obj.Key
			}
		}

		objectsErrCh <- nil
	}()

	for rErr := range pMinioStorage.client.RemoveObjects(bucketName, objectsCh) {
		glog.Errorf("Error removing objects from Minio: ", rErr)
		return errors.New("Failed removing objects from Minio")
	}

	if err := <-objectsErrCh; err != nil {
		return err
	}

	// Start timer for next clean up
	time.AfterFunc(
		time.Duration(pMinioStorage.retentionTime), func() {
			err := pMinioStorage.cleanStore()
			if err != nil {
				glog.Errorf("Failed to clear Minio object store: ", err)
			}
		})

	return nil
}

// generateKeyName : This used to generate the keyname
func generateKeyName() string {
	keyname := "persist_" + uuid.New().String()[:8]
	return keyname
}

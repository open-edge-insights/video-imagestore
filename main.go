/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

Explicit permissions are required to publish, distribute, sublicense, and/or sell copies of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	client "ElephantTrunkArch/DataAgent/da_grpc/client/go/client_internal"
	server "ElephantTrunkArch/ImageStore/server"
	util "ElephantTrunkArch/Util"
	"os"
	"time"

	"os/exec"

	"github.com/golang/glog"
	minio "github.com/minio/minio-go"
)

// grpc client certificates
const (
	RootCA        = "/etc/ssl/grpc_int_ssl_secrets/ca_certificate.pem"
	ClientCert    = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_certificate.pem"
	ClientKey     = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_key.pem"
	daServiceName = "ia_data_agent"
	daPort        = "50052"
)

func main() {
	// Wait for DA to be up
	ret := util.CheckPortAvailability(daServiceName, daPort)
	if !ret {
		glog.Error("DataAgent is not up, so exiting...")
		os.Exit(-1)
	}

	grpcClient, errr := client.NewGrpcInternalClient(ClientCert, ClientKey, RootCA, daServiceName, daPort)
	if errr != nil {
		glog.Errorf("Error while obtaining GrpcClient object...")
		os.Exit(-1)
	}
	configRedis := "RedisCfg"
	respMapRedis, err := grpcClient.GetConfigInt(configRedis)
	if err != nil {
		glog.Errorf("GetConfigInt failed...")
		os.Exit(-1)
	}
	configMinio := "MinioCfg"
	respMapMinio, err := grpcClient.GetConfigInt(configMinio)
	if err != nil {
		glog.Errorf("GetConfigInt failed...")
		os.Exit(-1)
	}

	glog.Infof("**************STARTING IMAGESTORE GRPC SERVER**************")
	done := make(chan bool)
	go StartRedis(respMapRedis)
	go StartMinio(respMapMinio)
	go StartMinioRetentionPolicy(respMapMinio)
	go server.StartGrpcServer(respMapRedis, respMapMinio)
	<-done
	glog.Infof("**************Exiting**************")
}

// StartRedis starts redis server
func StartRedis(redisConfigMap map[string]string) {
	redisPort := os.Getenv("REDIS_PORT")
	cmd := exec.Command("redis-server", "--port", redisPort, "--requirepass", redisConfigMap["Password"])
	err := cmd.Run()
	if err != nil {
		glog.Errorf("Not able to start redis server: %v", err)
		os.Exit(-1)
	}
}

// StartMinio starts minio server
func StartMinio(minioConfigMap map[string]string) {
	os.Setenv("MINIO_ACCESS_KEY", minioConfigMap["AccessKey"])
	os.Setenv("MINIO_SECRET_KEY", minioConfigMap["SecretKey"])
	os.Setenv("MINIO_REGION", "gateway")
	minioPort := os.Getenv("MINIO_PORT")
	glog.Infof("Minio port: %v", minioPort)
	// TODO: Need to see a way to pass port while bring
	// as --address switch didn't work as expected
	cmd := exec.Command("./minio", "server", "/data")
	err := cmd.Run()
	if err != nil {
		glog.Errorf("Not able to start minio server: %v", err)
		os.Exit(-1)
	}
}

// Helper method for reporting a missing key in the Minio configuration
func missingKeyError(key string) {
	msg := "Minio config missing key: " + key
	glog.Errorf(msg)
	os.Exit(-1)
}

// Clean up the image store
func StartMinioRetentionPolicy(config map[string]string) {
	glog.Infof("Running minio retention policy")

	minioPort := os.Getenv("MINIO_PORT")
	portUp := util.CheckPortAvailability("", minioPort)
	if !portUp {
		glog.Errorf("Minio port: %s not up, so exiting...", minioPort)
		os.Exit(-1)
	}

	glog.Infof("Initializing Minio client")

	region := "gateway"
	bucketName := "image-store-bucket"
	host := "localhost"

	retentionTimeStr, ok := config["RetentionTime"]
	if !ok {
		missingKeyError("RetentionTime")
	}

	retentionTime, err := time.ParseDuration(retentionTimeStr)
	if err != nil {
		glog.Errorf("Failed to parse retention time duration: %v", err)
		os.Exit(-1)
	}

	pollIntervalStr, ok := config["RetentionPollInterval"]
	if !ok {
		missingKeyError("RetentionPollInterval")
	}

	pollInterval, err := time.ParseDuration(pollIntervalStr)
	if err != nil {
		glog.Errorf("Failed to parse retention poll interval duration: %v", err)
		os.Exit(-1)
	}

	port, ok := config["Port"]
	if !ok {
		missingKeyError("Port")
	}

	accessKey, ok := config["AccessKey"]
	if !ok {
		missingKeyError("AccessKey")
	}

	secretKey, ok := config["SecretKey"]
	if !ok {
		missingKeyError("SecretKey")
	}

	sslStr, ok := config["Ssl"]
	if !ok {
		missingKeyError("Ssl")
	}

	ssl := true
	if sslStr == "true" {
		ssl = true
	} else if sslStr == "false" {
		ssl = false
	} else {
		msg := "Ssl key in Minio config must be true or false, not :" + sslStr
		glog.Errorf(msg)
		os.Exit(-1)
	}

	glog.Infof("Config: Host=%s, Port=%s, ssl=%v", host, port, ssl)

	glog.Infof("Initializing Minio client")
	client, err := minio.NewWithRegion(
		host+":"+port, accessKey, secretKey, ssl, region)
	if err != nil {
		glog.Errorf("Failed to connect to Minio server: %v", err)
		os.Exit(-1)
	}

	// Check if the bucket exists
	glog.Infof("Checking if Minio bucket already exists")
	found, err := client.BucketExists(bucketName)
	if err != nil {
		glog.Errorf("Failed to verify existence of bucket: %v", err)
		os.Exit(-1)
	}

	if !found {
		// Create the bucket if it does not exist
		glog.Infof("Creating bucket")
		client.MakeBucket(bucketName, region)
	}

	// Channel for objects to be removed from Minio
	objectsCh := make(chan string)
	objectsErrCh := make(chan error, 1)
	defer close(objectsErrCh)

	// Routine to find objects to remove and send them over the `objectsCh`
	go func() {
		glog.Infof("Finding objects in Minio to delete")

		// Defer channel close to when the function exits
		defer close(objectsCh)

		for obj := range client.ListObjects(bucketName, "", false, nil) {
			if obj.Err != nil {
				glog.Errorf("Failed retrieving objects from Minio: %v", obj.Err)
				objectsErrCh <- obj.Err
				return
			}

			now := time.Now()
			elapsed := now.Sub(obj.LastModified)

			if elapsed > retentionTime {
				glog.Infof("Deleting key: %s", obj.Key)
				objectsCh <- obj.Key
			}
		}

		objectsErrCh <- nil
	}()

	for rErr := range client.RemoveObjects(bucketName, objectsCh) {
		glog.Errorf("Error removing objects from Minio: %v", rErr)
		os.Exit(-1)
	}

	if err := <-objectsErrCh; err != nil {
		os.Exit(-1)
	}

	// Start timer for next clean up
	time.AfterFunc(
		time.Duration(pollInterval), func() {
			glog.Infof("Executing retention policy")
			StartMinioRetentionPolicy(config)
		})
}

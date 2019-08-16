/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

Explicit permissions are required to publish, distribute, sublicense, and/or sell copies of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	eismsgbus "EISMessageBus/eismsgbus"
	common "IEdgeInsights/ImageStore/common"
	imagestore "IEdgeInsights/ImageStore/go/ImageStore"
	isConfigMgr "IEdgeInsights/ImageStore/isConfigMgr"
	subManager "IEdgeInsights/ImageStore/subManager"
	util "IEdgeInsights/Util"
	cpuidutil "IEdgeInsights/Util/cpuid"
	configmgr "IEdgeInsights/libs/ConfigManager"
	"flag"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	minio "github.com/minio/minio-go"
)

const (
	chunkSize    = 1024 * 1024 * 8  // 8 MB
	maxFrameSize = 1024 * 1024 * 64 // 64MB
)

// IsServer is a struct used to implement ImageStore.IsServer
type IsServer struct {
	is *imagestore.ImageStore
}

func main() {

	flag.Parse()
	devMode, _ := strconv.ParseBool(os.Getenv("DEV_MODE"))
	// Initializing Etcd to set env variables
	cfgMgrConfig := map[string]string{
		"certFile":  "",
		"keyFile":   "",
		"trustFile": "",
	}
	if devMode != true {
		cfgMgrConfig = map[string]string{
			"certFile":  "/run/secrets/etcd_ImageStore_cert",
			"keyFile":   "/run/secrets/etcd_ImageStore_key",
			"trustFile": "/run/secrets/ca_etcd",
		}
	}
	_ = configmgr.Init("etcd", cfgMgrConfig)

	flag.Lookup("alsologtostderr").Value.Set("true")
	flag.Set("stderrthreshold", os.Getenv("GO_LOG_LEVEL"))
	flag.Set("v", os.Getenv("GO_VERBOSE"))

	glog.Infof("=============== STARTING imagestore ===============")

	vendorName := cpuidutil.Cpuid()
	if vendorName != "GenuineIntel" {
		glog.Errorf("*****Software runs only on Intel's hardware*****")
		os.Exit(-1)
	}

	common.DevMode, _ = strconv.ParseBool(os.Getenv("DEV_MODE"))
	minIoConfig, err := isConfigMgr.ReadMinIoConfig()
	if err != nil {
		glog.Errorf("Error while reading config :" + err.Error())
		os.Exit(-1)
	}

	defer glog.Flush()
	respMapMinio := make(map[string]string)
	// Converting struct to MapchunkSize
	respMapMinio["AccessKey"] = minIoConfig.AccessKey
	respMapMinio["SecretKey"] = minIoConfig.SecretKey
	respMapMinio["RetentionTime"] = minIoConfig.RetentionTime
	respMapMinio["RetentionPollInterval"] = minIoConfig.RetentionPollInterval
	respMapMinio["Ssl"] = minIoConfig.Ssl
	respMapMinio["Port"] = common.MinioPort
	respMapMinio["Host"] = common.MinioHost
	respMapMinio["ReplyEndpoint"] = minIoConfig.ReplyEndpoint

	done := make(chan bool)

	serviceConfig, err := isConfigMgr.ReadServiceConfig()
	if err != nil {
		glog.Errorf("Error in processing the serviceConfig")
		os.Exit(-1)
	}

	go StartMinio(respMapMinio)
	go StartMinioRetentionPolicy(respMapMinio)
	go startReqReply(respMapMinio, serviceConfig)
	go startSubScriber(respMapMinio)
	<-done
	glog.Infof("**************Exiting**************")
}

func startSubScriber(minioConfigMap map[string]string) {

	glog.Infof("**************In startSubScriber**************")

	topics := os.Getenv("SubTopics")
	topicArray := strings.Split(topics, ",")

	if len(topicArray) <= 0 {
		glog.Errorf("suscriber list empty")
		os.Exit(-1)
	}

	subConfig, err := isConfigMgr.ReadSubConfig(topicArray)
	if err != nil {
		glog.Errorf("Error in processing the config")
		os.Exit(-1)
	}

	subMgr := subManager.NewSubManager()
	subMgr.RegSubscriberList(subConfig)
	subMgr.StartAllSubscribers(topicArray)

	for _, topic := range topicArray {
		is, err := imagestore.GetImageStoreInstance(minioConfigMap)
		if err != nil {
			glog.Errorf("%v", err)
		}
		subMgr.RegWriterInterface(topic, is)
	}
	subMgr.ReceiveFromAll()
}

func startReqReply(minioConfigMap map[string]string, serviceConfig map[string]interface{}) {

	var ser IsServer
	is, err := imagestore.GetImageStoreInstance(minioConfigMap)
	ser.is = is
	if err != nil {
		glog.Errorf("Error while GetImageStoreInstance %v", err)
		os.Exit(-1)
	}

	client, err := eismsgbus.NewMsgbusClient(serviceConfig)
	if err != nil {
		glog.Errorf("-- Error initializing message bus context: %v\n", err)
		os.Exit(-1)
	}
	defer client.Close()

	serviceName := os.Getenv("AppName")
	glog.Infof("-- Initializing service %s\n", serviceName)
	service, err := client.NewService(serviceName)
	if err != nil {
		glog.Errorf("-- Error initializing service: %v\n", err)
		os.Exit(-1)
	}
	defer service.Close()

	glog.Infof("-- Running service %s\n", serviceName)

	for {
		var errMessage string
		msg, err := service.ReceiveRequest(-1)

		if err != nil {
			errMessage = "-- Error receiving request: " + err.Error()
			glog.Errorf(errMessage)
			continue
		}
		command, ok := msg.Data[common.Command].(string)
		if ok == false {
			errMessage = "Missing " + common.Command
			handleError(service, errMessage)
			continue
		}

		imgHandle, ok := msg.Data[common.ImageHandle].(string)
		if ok == false {
			errMessage += "Missing " + common.ImageHandle
			handleError(service, errMessage)
			continue
		}

		if len(errMessage) > 0 {
			handleError(service, errMessage)
		} else if command == common.ReadCode {
			handleReadCommand(imgHandle, service, ser)
		} else if command == common.StoreCode {
			if msg.Blob != nil {
				handleStoreCommand(imgHandle, service, ser, msg.Blob)
			} else {
				errMessage = "Can not store empty image for handle " + imgHandle
				handleError(service, errMessage)
			}
		} else {
			errMessage = "Invalid Command " + command
			handleError(service, errMessage)
		}
	}
}

func handleError(service *eismsgbus.Service, errMessage string) {
	glog.Errorf(errMessage)
	service.Response(map[string]interface{}{common.Error: errMessage})
}

func handleReadCommand(imgHandle string, service *eismsgbus.Service, ser IsServer) {

	frame, err := ser.Read(imgHandle)

	if err != nil {
		error := "Reading image failed for handle " + imgHandle + " Error :" + err.Error()
		glog.Errorf(error)
		service.Response(map[string]interface{}{common.Error: error})
	} else {
		response := make([]interface{}, 2)
		response[0] = map[string]interface{}{common.ImageHandle: imgHandle}
		response[1] = frame
		service.Response(response)
		message := "Successfully read frame with handle:" + imgHandle
		glog.Infof(message)
	}
}

func handleStoreCommand(imgHandle string, service *eismsgbus.Service, ser IsServer, imgFrame []byte) {
	key, err := ser.StoreData(imgFrame, imgHandle)
	if err != nil {
		error := "Store image failed for handle " + imgHandle + " Error :" + err.Error()
		glog.Errorf(error)
		service.Response(map[string]interface{}{common.Error: error})
	} else {
		service.Response(map[string]interface{}{common.ImageHandle: key})
		message := "Successfully stored frame with handle:" + imgHandle
		glog.Infof(message)
	}
}

// StoreData is used to store image buffer in minio.
//
// 1. keyname : []byte
//    Refers to the image frame to be stored.
// 2. keyname : string
//    Refers to the image handle of the image to be stored.
//
// Returns:
// 1. error
//    Returns an error object if store fails.
func (s *IsServer) StoreData(blob []byte, keyname string) (string, error) {
	key, err := s.is.Store(blob, keyname)
	if err != nil {
		glog.Errorf("Store failed")
		return "", err
	}
	return key, nil
}

// Read is used to read image buffer from minio.
//
// Parameters:
// 1. keyname : string
//    Refers to the image handle of the image to be read.
//
// Returns:
// 1. []byte
//    Returns the byte array of image buffer.
// 2. error
//    Returns an error object if read fails.
func (s *IsServer) Read(key string) ([]byte, error) {
	output, err := s.is.Read(key)
	if err != nil {
		glog.Errorf("Read failed: %v", err)
		return nil, err
	}

	bufLen := 0
	buf := make([]byte, chunkSize, maxFrameSize)
	outputByteArr := make([]byte, chunkSize)
	for {
		// TODO : if len(output handle data) > outputByteArr,
		// Read API crashes. This need to be fixed.
		// Currently the 8MB is max size of image
		n, err := (output).Read(outputByteArr)
		if err != nil {
			if err == io.EOF {
				// This is to send the last remaining chunk
				copy(buf[bufLen:n], outputByteArr[:n])
				bufLen += n
				break
			}
			glog.Errorf("Error for ioReader.Read(): %v for key: %v \n", err, key)
			break
		}
		break
		copy(buf[bufLen:n], outputByteArr[0:n])
		bufLen += n
	}

	output.Close()
	output = nil

	var outputBuff []byte
	if bufLen > 0 {
		outputBuff = make([]byte, bufLen)
		copy(outputBuff, buf[0:bufLen])
	}

	return outputBuff, nil
}

// StartMinio starts the minio server.
//
// Parameters:
// 1. minioConfigMap : map[string]string
//    Refers to the minio config.
func StartMinio(minioConfigMap map[string]string) {
	os.Setenv("MINIO_ACCESS_KEY", minioConfigMap["AccessKey"])
	os.Setenv("MINIO_SECRET_KEY", minioConfigMap["SecretKey"])
	os.Setenv("MINIO_REGION", "gateway")
	glog.Infof("Minio port: %v\n", common.MinioPort)

	// TODO: Need to see a way to pass port while bring
	// as --address switch didn't work as expected
	cmd := exec.Command("./minio", "server", "--address", common.MinioHost+":"+common.MinioPort, "/data")
	err := cmd.Run()
	if err != nil {
		glog.Errorf("Not able to start minio server: %v", err)
		os.Exit(-1)
	}
}

// missingKeyError is a helper method to report a missing key in Minio config
//
// Parameters:
// 1. key : string
//    Refers to Image handle.
func missingKeyError(key string) {
	msg := "Minio config missing key: " + key
	glog.Errorf(msg)
	return
}

// StartMinioRetentionPolicy cleans up the ImageStore
//
// Parameters:
// 1. config : map[string]string
//    Refers to the minio config
func StartMinioRetentionPolicy(config map[string]string) {
	defer glog.Flush()
	glog.Infof("Running minio retention policy")
	minioPort := common.MinioPort
	portUp := util.CheckPortAvailability("", minioPort)
	if !portUp {
		glog.Errorf("Minio port: %s not up, so exiting...", minioPort)
		os.Exit(-1)
	}

	region := "gateway"
	bucketName := "image-store-bucket"
	port := common.MinioPort
	host := common.MinioHost

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

	glog.V(1).Infof("Config: Host=%s, Port=%s, ssl=%v", host, port, ssl)

	client, err := minio.NewWithRegion(
		host+":"+port, accessKey, secretKey, ssl, region)
	if err != nil {
		glog.Errorf("Failed to connect to Minio server: %v", err)
		os.Exit(-1)
	}

	// Check if the bucket exists
	glog.V(1).Infof("Checking if Minio bucket already exists")
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
	removeObjects := func() {
		objectsCh := make(chan string)
		objectsErrCh := make(chan error, 1)
		defer close(objectsErrCh)

		// Routine to find objects to remove and send them over the `objectsCh`
		go func() {
			glog.V(1).Infof("Finding objects in Minio to delete")

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
					glog.V(1).Infof("Deleting key: %s", obj.Key)
					objectsCh <- obj.Key
				} else {
					glog.V(2).Infof("Not deleting key: %s", obj.Key)
				}
			}

			objectsErrCh <- nil
		}()

		for rErr := range client.RemoveObjects(bucketName, objectsCh) {
			glog.Errorf("Error removing objects from Minio: %v", rErr)
			return
		}

		if err := <-objectsErrCh; err != nil {
			return
		}
	}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	removeObjects()

	for range ticker.C {
		removeObjects()
	}
	glog.Infof("Exiting StartMinioRetentionPolicy()...")
}

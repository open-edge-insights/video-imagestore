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
	"os"

	"os/exec"

	"github.com/golang/glog"
)

func main() {
	grpcClient, errr := client.NewGrpcClient("ia_data_agent", "50052")
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
	go server.StartGrpcServer(respMapRedis, respMapMinio)
	<-done
	glog.Infof("**************Exiting**************")
}

// StartRedis starts redis server
func StartRedis(redisConfigMap map[string]string) {
	cmd := exec.Command("redis-server", "--requirepass", redisConfigMap["Password"])
	err := cmd.Run()
	if err != nil {
		glog.Errorf("Command failed: %s", err)
		os.Exit(-1)
	}
}

// StartMinio starts minio server
func StartMinio(minioConfigMap map[string]string) {
	os.Setenv("MINIO_ACCESS_KEY", minioConfigMap["AccessKey"])
	os.Setenv("MINIO_SECRET_KEY", minioConfigMap["SecretKey"])
	os.Setenv("MINIO_REGION", "gateway")
	cmd := exec.Command("./minio", "server", "/data")
	err := cmd.Run()
	if err != nil {
		glog.Errorf("Command failed: %s", err)
		os.Exit(-1)
	}
}

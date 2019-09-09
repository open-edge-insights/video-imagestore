/*
Copyright (c) 2019 Intel Corporation.

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

package isConfigMgr

import (
	configmgr "IEdgeInsights/libs/ConfigManager"
	"encoding/json"

	//"fmt"
	//"strconv"
	common "IEdgeInsights/ImageStore/common"
	msgbusutil "IEdgeInsights/util/msgbusutil"
	"os"

	"github.com/golang/glog"
)

type Configuration struct {
	Minio struct {
		AccessKey             string `json:"accessKey"`
		SecretKey             string `json:"secretKey"`
		RetentionTime         string `json:"retentionTime"`
		RetentionPollInterval string `json:"retentionPollInterval"`
		Ssl                   string `json:"ssl"`
		ReplyEndpoint         string `json:"replyEndpoint"`
		Host                  string `json:"host"`
	} `json:"minio"`
}

type Minio struct {
	AccessKey             string
	SecretKey             string
	RetentionTime         string
	RetentionPollInterval string
	Ssl                   string
	ReplyEndpoint         string
	Host                  string
}

func ReadMinIoConfig() (Minio, error) {

	var minIoConfig Minio
	var tempConfig Configuration

	config := common.GetConfigInfoMap()
	confHandler := configmgr.Init("etcd", config)
	appName := os.Getenv("AppName")
	minioConfigPath := "/" + appName + "/config"
	value, err := confHandler.GetConfig(minioConfigPath)
	if err != nil {
		glog.Infof("Error while getting value of %s, err %s\n", minioConfigPath, err.Error())
		return minIoConfig, err
	}
	err = json.Unmarshal([]byte(value), &tempConfig)
	if err != nil {
		glog.Infof("Error while json.Unmarshal")
		return minIoConfig, err
	}

	minIoConfig.AccessKey = tempConfig.Minio.AccessKey
	minIoConfig.SecretKey = tempConfig.Minio.SecretKey
	minIoConfig.RetentionTime = tempConfig.Minio.RetentionTime
	minIoConfig.RetentionPollInterval = tempConfig.Minio.RetentionPollInterval
	minIoConfig.Ssl = tempConfig.Minio.Ssl
	return minIoConfig, nil
}

func ReadSubConfig(topicArray []string) (map[string]interface{}, error) {

	subsInfoMap := make(map[string]interface{})
	cfgMgrConfig := common.GetConfigInfoMap()
	glog.Info("config for etcd client : %v", cfgMgrConfig)
	for _, topic := range topicArray {
		subsInfoMap[topic] = msgbusutil.GetMessageBusConfig(topic, "sub", common.DevMode, cfgMgrConfig)
	}

	return subsInfoMap, nil
}

func ReadServiceConfig() (map[string]interface{}, error) {
	serviceName := os.Getenv("AppName")

	cfgMgrConfig := common.GetConfigInfoMap()
	glog.Info("config for etcd client : %v", cfgMgrConfig)
	return msgbusutil.GetMessageBusConfig(serviceName, "server", common.DevMode, cfgMgrConfig), nil
}

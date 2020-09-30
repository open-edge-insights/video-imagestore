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

package isconfigmgr

import (
	util "IEdgeInsights/common/util"
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

// Configuration type struct
type Configuration struct {
	Minio struct {
		AccessKey             string `json:"accessKey"`
		SecretKey             string `json:"secretKey"`
		RetentionTime         string `json:"retentionTime"`
		RetentionPollInterval string `json:"retentionPollInterval,omitempty"`
		Ssl                   string `json:"ssl"`
		ReplyEndpoint         string `json:"replyEndpoint"`
		Host                  string `json:"host"`
	} `json:"minio"`
}

// Minio type struct
type Minio struct {
	AccessKey             string
	SecretKey             string
	RetentionTime         string
	RetentionPollInterval string
	Ssl                   string
	ReplyEndpoint         string
	Host                  string
}

// ReadMinIoConfig - function to read Minio configuration
func ReadMinIoConfig(conf map[string]interface{}) (Minio, error) {

	var minIoConfig Minio
	var tempConfig Configuration
	value, err := json.Marshal(conf)
	if err != nil {
		glog.Errorf("Error:Conversion from json to string")
		return minIoConfig, err
	}

	// Reading schema json
	schema, err := ioutil.ReadFile("./schema.json")
	if err != nil {
		glog.Errorf("Schema file not found")
		return minIoConfig, err
	}

	// Validating config json
	if util.ValidateJSON(string(schema), string(value)) != true {
		return minIoConfig, err
	}

	err = json.Unmarshal([]byte(string(value)), &tempConfig)
	if err != nil {
		glog.Errorf("Error while json.Unmarshal")
		return minIoConfig, err
	}

	minIoConfig.AccessKey = tempConfig.Minio.AccessKey
	minIoConfig.SecretKey = tempConfig.Minio.SecretKey
	minIoConfig.RetentionTime = tempConfig.Minio.RetentionTime
	minIoConfig.RetentionPollInterval = tempConfig.Minio.RetentionPollInterval
	minIoConfig.Ssl = tempConfig.Minio.Ssl
	return minIoConfig, nil
}

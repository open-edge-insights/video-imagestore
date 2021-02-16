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

package main

import (
	eiimsgbus "EIIMessageBus/eiimsgbus"
	common "IEdgeInsights/ImageStore/common"
	envconfig "EnvConfig"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
)

func main() {
	/*
		Test Description: Test steps are
		1.This test program construct the raw frame in memory and publishes to
		the topic to which ImageStore is subscribed for.

		2.In the next step, test program sends the read command with the known
		image handle.

		3.It then compares the frame stored and frame read and prints the result
	*/

	if publishFrame() {
		fmt.Println("\nFrame published successfully\n")
		readAndCompareFrame()
	} else {
		fmt.Println("Publishing frame failed.Test failed.")
	}
}

func publishFrame() bool {

	retVal := false

	devModeStr := os.Getenv("DEV_MODE")
	devMode, err := strconv.ParseBool(devModeStr)
	common.DevMode, _ = strconv.ParseBool(os.Getenv("DEV_MODE"))
	if err != nil {
		glog.Errorf("string to bool conversion error")
	}

	cfgMgrConfig := common.GetConfigInfoMap()
	os.Setenv("camera1_stream_results_cfg", "zmq_tcp,127.0.0.1:65013")

	if !devMode {
		os.Setenv("AppName", "VideoAnalytics")
		os.Setenv("Clients", "ImageStore")
	}
	topic := "camera1_stream_results"
	config := envconfig.GetMessageBusConfig(topic, "pub", devMode, cfgMgrConfig)

	fmt.Println("-- Initializing message bus context %v\n", config)
	client, err := eiimsgbus.NewMsgbusClient(config)
	if err != nil {
		fmt.Printf("-- Error initializing message bus context: %v\n", err)
		return retVal
	}
	defer client.Close()

	fmt.Printf("-- Creating publisher for topic %s\n", topic)
	publisher, err := client.NewPublisher(topic)
	if err != nil {
		fmt.Printf("-- Error creating publisher: %v\n", err)
		return retVal
	}
	defer publisher.Close()

	fmt.Println("-- Running...\n")
	sz := 1024 * 1024 * 4
	frame := make([]byte, sz)
	frame[0] = 0
	frame[1] = '|'
	frame[sz-2] = '|'
	frame[sz-1] = 0
	msg := make([]interface{}, 2)
	msg[0] = map[string]interface{}{"img_handle": "pubTest1"}
	msg[1] = frame

	err = publisher.Publish(msg)
	if err != nil {
		fmt.Printf("-- Failed to publish message: %v\n", err)
		return retVal
	}
	time.Sleep(1 * time.Second)
	return true
}

func readAndCompareFrame() {

	fmt.Printf("\n -- Going to read frame -- \n")

	devModeStr := os.Getenv("DEV_MODE")
	devMode, err := strconv.ParseBool(devModeStr)
	if err != nil {
		glog.Errorf("string to bool conversion error")
	}

	cfgMgrConfig := map[string]string{
		"certFile":  "",
		"keyFile":   "",
		"trustFile": "",
	}
	if !devMode {
		os.Setenv("AppName", "ImageStore")
	}
	serviceName := "ImageStore"
	config := envconfig.GetMessageBusConfig(serviceName, "client", devMode, cfgMgrConfig)

	fmt.Println("-- Initializing message bus context %v\n", config)
	client, err := eiimsgbus.NewMsgbusClient(config)
	if err != nil {
		fmt.Printf("-- Error initializing message bus context: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Printf("-- Initializing service requester %s\n", serviceName)
	service, err := client.GetService(serviceName)
	if err != nil {
		fmt.Printf("-- Error initializing service requester: %v\n", err)
		return
	}
	defer service.Close()

	// construct frame data
	sz := 1024 * 1024 * 4
	frame := make([]byte, sz)
	frame[0] = 0
	frame[1] = '|'
	frame[sz-2] = '|'
	frame[sz-1] = 0

	// Send Read command & get the frame data
	response := map[string]interface{}{"command": "read", "img_handle": "pubTest1"}
	err1 := service.Request(response)
	if err1 != nil {
		fmt.Printf("-- Error sending request: %v\n", err1)
		fmt.Printf("--Test Failed--\n")
		return
	}

	fmt.Printf("-- Waiting for read command response")

	resp, err := service.ReceiveResponse(-1)
	if err != nil {
		fmt.Printf("-- Error receiving response: %v\n", err)
		fmt.Printf("--Test Failed--\n")
		return
	}
	fmt.Printf("\nFrame size is : %d \n", len(frame))
	fmt.Printf("\nFrame read and frame size is : %d \n", len(resp.Blob))

	// Compare frame data and declare result
	if bytes.Compare(frame, resp.Blob) == 0 {
		fmt.Printf("\nBinary data stored and read matches. Test passed\n")
	} else {
		fmt.Printf("\nBinary data stored and read does not matches. Test failed\n")
	}
}

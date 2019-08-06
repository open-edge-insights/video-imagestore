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
	eismsgbus "EISMessageBus/eismsgbus"
	util "IEdgeInsights/libs/common/go"
	"github.com/golang/glog"
	"strconv"
	"os"
	"flag"
	"fmt"
	"bytes"
)

func main() {
	/*
	Test Description: Test steps are  
	1.This test program construct the raw frame in memory and sends the store 
	command with known image handle.

	2.In the next step, test program sends the read command with the known 
	image handle.

	3.It then compares the frame stored and frame read and prints the result
	*/

	/*configFile := flag.String("configFile", "", "JSON configuration file")
	serviceName := flag.String("serviceName", "", "Subscription topic")
	flag.Parse()

	if *configFile == "" {
		fmt.Println("-- Config file must be specified")
		return
	}

	if *serviceName == "" {
		fmt.Println("-- Service name must be specified")
		return
	}

	fmt.Printf("-- Loading configuration file %s\n", *configFile)
	config, err := eismsgbus.ReadJsonConfig(*configFile)
	if err != nil {
		fmt.Printf("-- Failed to parse config: %v\n", err)
		return
	}*/

	topic := flag.String("topic", "", "topic")
	devModeStr := flag.String("devmode", "", "devmode")
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")
	defer glog.Flush()

	var config map[string]interface{}
	devMode, err := strconv.ParseBool(*devModeStr)
	if err != nil {
		glog.Errorf("string to bool conversion error")
	}

	os.Setenv("AppName", "ImageStore")
	os.Setenv(*topic+"_cfg", "zmq_tcp,127.0.0.1:5669")
	os.Setenv("Subscribers", "ImageStore")
	pubTopics := util.GetTopics("PUB")

	cfgMgrConfig := map[string]string{
        "certFile":  "",
        "keyFile":   "",
        "trustFile": "",
	}
	config = util.GetMessageBusConfig(pubTopics[0], "sub", devMode, cfgMgrConfig)

	fmt.Println("-- Initializing message bus context %v\n", config)
	client, err := eismsgbus.NewMsgbusClient(config)
	if err != nil {
		fmt.Printf("-- Error initializing message bus context: %v\n", err)
		return
	}
	defer client.Close()
    serviceName := os.Getenv("AppName")
	fmt.Printf("-- Initializing service requester %s\n", serviceName)
	service, err := client.GetService(serviceName)
	if err != nil {
		fmt.Printf("-- Error initializing service requester: %v\n", err)
		return
	}
	defer service.Close()

	// construct frame data
	sz := 1024*1024*4
	frame := make([]byte, sz)
	frame[0] = 0
	frame[1] = '|'
	frame[sz-2] = '|'
	frame[sz-1] = 0

	// send store command to store frame of size sz
	fmt.Printf("-- Sending request")
	response := make([]interface{}, 2)
	response[0] = map[string]interface{}{"command": "store","img_handle":"testHandle"}
	response[1] = frame
	err = service.Request(response)
	if err != nil {
		fmt.Printf("-- Error sending request: %v\n", err)
		fmt.Printf("--Test Failed--\n")
		return
	}

	fmt.Printf("-- Waiting for store command response")
	resp, err := service.ReceiveResponse(-1)
	if err != nil {
		fmt.Printf("-- Error receiving response: %v\n", err)
		fmt.Printf("--Test Failed--\n")
		return
	}

	fmt.Printf("--  Received response for store : %v\n", resp)

    // Send Read command & get the frame data
	response[0] = map[string]interface{}{"command": "read","img_handle":"testHandle"}
	err = service.Request(response)
	if err != nil {
		fmt.Printf("-- Error sending request: %v\n", err)
		fmt.Printf("--Test Failed--\n")
		return
	}

	fmt.Printf("-- Waiting for read command response")
	
	resp, err = service.ReceiveResponse(-1)
	if err != nil {
		fmt.Printf("-- Error receiving response: %v\n", err)
		fmt.Printf("--Test Failed--\n")
		return
	}
	
	fmt.Printf("frame size after read : %d \n", len(resp.Blob))
	
	// Compare frame data and declare result
	if bytes.Compare(frame, resp.Blob) == 0 {
		fmt.Printf("Binary data stored and read matches. Test passed\n")
	}else{
		fmt.Printf("Binary data stored and read does not matches. Test failed\n")
	}
}

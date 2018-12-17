/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	client "ElephantTrunkArch/ImageStore/client/go"
	"flag"
	"os"

	"github.com/golang/glog"
)

// Client Certificates
const (
	RootCA     = "/etc/ssl/grpc_internal/ca_certificate.pem"
	ClientCert = "/etc/ssl/imagestore/imagestore_client_certificate.pem"
	ClientKey  = "/etc/ssl/imagestore/imagestore_client_key.pem"
)

func main() {
	var outputFile string
	flag.StringVar(&outputFile, "output_file", "", "output file that gets"+
		"created from ImageStore read")

	flag.Parse()

	if len(os.Args) < 1 {
		glog.Errorf("provide output file path via output_file option")
	}
	grpcClient, err := client.NewImageStoreClient(RootCA, ClientCert, ClientKey, "localhost", "50055")
	if err != nil {
		glog.Errorf("Error while obtaining GrpcClient object...")
		os.Exit(-1)
	}

	glog.Infof("******Go client gRPC testing******")
	somevar := "inmem_720a633f"
	respMap, err := grpcClient.Read(somevar)
	//Writing files
	writeFile(outputFile, respMap)
	somebytes := []byte(somevar)
	respMapp, err := grpcClient.Store(somebytes, "inmemory")
	glog.Infof(respMapp)	
}

func writeFile(filename string, message []byte) {
	f, err := os.Create(filename)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	defer f.Close()
	n3, _ := f.Write(message)
	glog.Infof("wrote %d bytes\n", n3)
	f.Sync()
}

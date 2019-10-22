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
package main

import (
	persistent "IEdgeInsights/ImageStore/go/ImageStore/Persistent"
	"fmt"
	"flag"
	"github.com/golang/glog"
)



func TestPersistentMinio() {
	config := map[string]string{
		"Host":          "localhost",
		"Port":          "9000",
		"AccessKey":     "adminss",
		"SecretKey":     "passwordss",
		"RetentionTime": "10",
		"Ssl":           "false",
	}

	pims, err := persistent.NewPersistent("minio", config)
	if err != nil {
		glog.Errorf("Initializing persistent storage failed: %v", err)
		return
	}

	data := make([]byte, 512) // Make dummy 512 bytes of data

	// Store the data in the image store
	fmt.Println("-- Storing data")
	key, err := pims.Store(data, "random_key_one")
	if err != nil {
		glog.Errorf("Failed to store data in persistent storage")
		return
	}

	// Read the data from the image store
	fmt.Println("-- Reading data")
	rdata, err := pims.Read(key)
	if err != nil {
		glog.Errorf("Failed to read data from persistent storage")
		return
	}

	fmt.Println("-- Verifying data")
	// Convert the string data that was read into bytes
	outputByteArr := make([]byte, 512)
	result, err := (rdata).Read(outputByteArr)
	brdata := []byte(outputByteArr)

	// Verify that the read data is the same as the stored data
	if len(brdata) != len(data) {
		glog.Errorf("Lengths of retrieved vs original data do not match: %d != %d",
		result, len(data))
		return
	}

	for i := range data {
		if brdata[i] != data[i] {
			glog.Errorf("Retrieved data is different that original data")
			return
		}
	}

	// Remove the value from the image store
	fmt.Println("-- Removing data")
	err = pims.Remove(key)
	if err != nil {
		glog.Errorf("Failed to remove key from image store")
		return
	}

	// Verify that it no longer exists in the image store
	fmt.Println("-- Verifying removed data (an error should printed)")
	_, err = pims.Read(key)
	if err == nil {
		glog.Errorf("The object still exists in image store after removal")
		return
	}

	// Passed!
}

func main() {
	flag.Parse()
	TestPersistentMinio()
}
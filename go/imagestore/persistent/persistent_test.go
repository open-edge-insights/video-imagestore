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
package persistent

import (
	"fmt"
	"testing"
	"time"
)

func TestPersistentMinio(t *testing.T) {
	config := map[string]string{
		"Host":          "localhost",
		"Port":          "9000",
		"AccessKey":     "admin",
		"SecretKey":     "password",
		"RetentionTime": "10",
		"Ssl":           "false",
	}

	pims, err := NewPersistent("minio", config)
	if err != nil {
		t.Errorf("Initializing persistent storage failed: %v", err)
		return
	}

	data := make([]byte, 512) // Make dummy 512 bytes of data

	// Store the data in the image store
	fmt.Println("-- Storing data")
	key, err := pims.Store(data)
	if err != nil {
		t.Errorf("Failed to store data in persistent storage")
		return
	}

	// Read the data from the image store
	fmt.Println("-- Reading data")
	rdata, err := pims.Read(key)
	if err != nil {
		t.Errorf("Failed to read data from persistent storage")
		return
	}

	fmt.Println("-- Verifying data")
	// Convert the string data that was read into bytes
	brdata := []byte(rdata)

	// Verify that the read data is the same as the stored data
	if len(brdata) != len(data) {
		t.Errorf("Lengths of retrieved vs original data do not match: %d != %d",
			len(brdata), len(data))
		return
	}

	for i := range data {
		if brdata[i] != data[i] {
			t.Errorf("Retrieved data is different that original data")
			return
		}
	}

	// Remove the value from the image store
	fmt.Println("-- Removing data")
	err = pims.Remove(key)
	if err != nil {
		t.Errorf("Failed to remove key from image store")
		return
	}

	// Verify that it no longer exists in the image store
	fmt.Println("-- Verifying removed data (an error should printed)")
	_, err = pims.Read(key)
	if err == nil {
		t.Errorf("The object still exists in image store after removal")
		return
	}

	// Passed!
}

func TestPersistentMinioRetention(t *testing.T) {
	config := map[string]string{
		"Host":          "localhost",
		"Port":          "9000",
		"AccessKey":     "admin",
		"SecretKey":     "password",
		"RetentionTime": "5",
		"Ssl":           "false",
	}

	pims, err := NewPersistent("minio", config)
	if err != nil {
		t.Errorf("Initializing persistent storage failed: %v", err)
		return
	}

	data := make([]byte, 512) // Make dummy 512 bytes of data

	// Store the data in the image store
	fmt.Println("-- Storing data")
	key, err := pims.Store(data)
	if err != nil {
		t.Errorf("Failed to store data in persistent storage")
		return
	}

	fmt.Println("-- Sleepint 15 seconds for retention procedure to run")
	time.Sleep(time.Duration(15) * time.Second)

	// Verify that it no longer exists in the image store
	fmt.Println("-- Verifying removed data (an error should printed)")
	_, err = pims.Read(key)
	if err == nil {
		t.Errorf("The object still exists in image store after removal")
		return
	}

	// Passed!
}

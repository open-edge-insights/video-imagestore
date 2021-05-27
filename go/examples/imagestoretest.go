/*
Copyright (c) 2021 Intel Corporation

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

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
	imagestore "IEdgeInsights/ImageStore/go/ImageStore"
	"os"
	"flag"
	"github.com/golang/glog"
)

func checkErr(resp string, err error) {
	if err != nil {
		glog.Errorf("Error: %v", err)
	} else {
		if resp != "" {
			glog.Infof("Response: %s", resp)
		}
	}
}

func readFile(filename string) []byte {

	file, err := os.Open(filename)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		glog.Errorf("Error: %v", err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	_, err = file.Read(buffer)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	return buffer
}

func writeFile(filename string, message string) {
	f, err := os.Create(filename)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	defer f.Close()
	n3, _ := f.WriteString(message)
	glog.Infof("wrote %d bytes\n", n3)
	f.Sync()
}

func main() {

	var inputFile string
	var outputFile string
	flag.StringVar(&inputFile, "input_file", "", "input file path to write to ImageStore")
	flag.StringVar(&outputFile, "output_file", "", "output file that gets" +
				   "created from ImageStore read")

	flag.Parse()

	if len(os.Args) < 2 {
		glog.Errorf("Usage: go run DataAgent/da_grpc/test/clientTest.go " +
			"-input_file=<input_file_path> [-output_file=<output_file_path>]")
		os.Exit(-1)
	}

	flag.Set("logtostderr", "true")
	defer glog.Flush()


	imagestore, err := imagestore.NewImageStore()
	if err != nil {
		glog.Errorf("Failed to instantiate ImageStore. Error: %s", err)
	} else {
		var err error
		var data string
		var keyname string

		data, err = imagestore.Read("inmem")
		checkErr(data, err)

		imagestore.SetStorageType("inmemory")
		keyname, err = imagestore.Store([]byte("vivek"))
		checkErr(keyname, err)

		data, err = imagestore.Read(keyname)
		checkErr("Read success", err)

		err = imagestore.Remove(keyname)
		checkErr("Keyname " + keyname + " removed successfully", err)

		//Reading Files
		inputData := readFile(inputFile)
		keyname, err = imagestore.Store(inputData)
		checkErr(keyname, err)

		data, err = imagestore.Read(keyname)
		checkErr("Read success", err)

		//Writing files
		writeFile(outputFile, data)
	}
}

package main

import (
	"fmt"
	imagestore "iapoc_elephanttrunkarch/ImageStore/go/ImageStore"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(filename string) []byte {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err, bytesread)
	}

	return buffer

}

func writeFile(filename string, message string) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()
	n3, _ := f.WriteString(message)
	fmt.Printf("wrote %d bytes\n", n3)
	f.Sync()
}

func main() {

	imagestore, err := imagestore.NewImageStore()
	if err != nil {
		fmt.Println("Some Issue in Connection", err)
	} else {

		status, message := imagestore.Read("inmem")
		fmt.Println(status, message)

		imagestore.SetStorageType("inmemory")
		status, keyname := imagestore.Store([]byte("vivek"))
		fmt.Println(status, keyname)

		status, message = imagestore.Read(keyname)
		fmt.Println(status, message)

		status, message = imagestore.Remove(keyname)
		fmt.Println(status, message)

		//Reading Files
		readFile := readFile("test.jpg")
		status, keyname = imagestore.Store(readFile)

		status, message = imagestore.Read(keyname)
		//Writing files
		writeFile("fromredis.jpg", message)
	}

}

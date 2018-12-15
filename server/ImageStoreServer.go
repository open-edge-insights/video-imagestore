/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

Explicit permissions are required to publish, distribute, sublicense, and/or sell copies of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	config "ElephantTrunkArch/DataAgent/config"
	pb "ElephantTrunkArch/ImageStore/protobuff/go"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"

	"context"
	"net"
	"os"
	"time"

	imagestore "ElephantTrunkArch/ImageStore/go/ImageStore"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var gRPCImageStoreHost = "localhost"

const (
	gRPCImageStorePort = "50055"
	chunkSize          = 4095 * 1024 // 4 MB
)

// Server Certificates
const (
	RootCA     = "Certificates/ca/ca_certificate.pem"
	ServerCert = "Certificates/imagestore/imagestore_server_certificate.pem"
	ServerKey  = "Certificates/imagestore/imagestore_server_key.pem"
)

// ImgDaCfg stores parsed DataAgent config
var ImgDaCfg config.DAConfig

// IsServer is used to implement ImageStore.IsServer
type IsServer struct{
	is *imagestore.ImageStore
}

// StartGrpcServer starts the ImageStore grpc server
func StartGrpcServer(IsCfg config.DAConfig) {

	ImgDaCfg = IsCfg
	ipAddr, err := net.LookupIP("ia_imagestore")
	if err != nil {
		glog.Errorf("Failed to fetch the IP address for host: %v, error:%v", ipAddr, err)
	} else {
		gRPCImageStoreHost = ipAddr[0].String()
	}

	flag.Parse()

	flag.Lookup("alsologtostderr").Value.Set("true")

	defer glog.Flush()
	if len(os.Args) < 1 {
		glog.Infof("No args passed.")
	}
	addr := gRPCImageStoreHost + ":" + gRPCImageStorePort


	jsonStr, err := getConfig("RedisCfg")
	if err != nil {
		glog.Errorf("getConfig(\"RedisCfg\") method failed. Error: %v", err)
		os.Exit(-1)
	}

	var configMap map[string]string
	configBytes := []byte(jsonStr)
	json.Unmarshal(configBytes, &configMap)

	minioJsonStr, err := getConfig("MinioCfg")
	if err != nil {
		glog.Errorf("Failed to retrieve minio config: %v", err)
		os.Exit(-1)
	}

	var minioConfigMap map[string]string
	minioConfigBytes := []byte(minioJsonStr)
	json.Unmarshal(minioConfigBytes, &minioConfigMap)

	glog.Infof("Minio cfg in read: %+v", minioConfigMap)

	// Manually set the host to localhost since we are inside the docker network
	minioConfigMap["Host"] = "localhost"

	// Read certificate binary
	certPEMBlock, err := ioutil.ReadFile(ServerCert)
	if err != nil {
		glog.Errorf("Failed to Read Server Certificate : %s", err)
		os.Exit(-1)
	}

	keyPEMBlock, err := ioutil.ReadFile(ServerKey)
	if err != nil {
		glog.Errorf("Failed to Read Server Key : %s", err)
		os.Exit(-1)
	}

	// Load the certificates from binary
	certificate, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		glog.Errorf("Failed to Load ServerKey Pair : %s", err)
		os.Exit(-1)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(RootCA)
	if err != nil {
		glog.Errorf("Failed to Read CA Certificate : %s", err)
		os.Exit(-1)
	}

	// Append the certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		glog.Errorf("Failed to Append CA Certificate")
		os.Exit(-1)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Errorf("failed to listen: %v", err)
		os.Exit(-1)
	}

	// Create the TLS configuration to pass to the GRPC server
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	//Create the gRPC server
	s := grpc.NewServer(grpc.Creds(creds))

	time.Sleep(1 * time.Second)

	imgStore, err := imagestore.GetImageStoreInstance(configMap, minioConfigMap)
	if err != nil {
		glog.Errorf("Failed to instantiate GetImageStoreInstance(). Error: %v", err)
		os.Exit(-1)
	}

	//Register the handle object
	pb.RegisterIsServer(s, &IsServer{is: imgStore})

	glog.Infof("Secure gRPC server Started & Listening at: %s", addr)

	//Serve and listen
	if err := s.Serve(lis); err != nil {
		glog.Errorf("grpc serve error: %s", err)
		os.Exit(-1)
	}
}

// Read implementation
func (s *IsServer) Read(in *pb.ReadReq, srv pb.Is_ReadServer) error {
	output, err := s.is.Read(in.ReadKeyname)
	if err != nil {
		glog.Infof("Read failed")
		return err
	}
	outputtwo := []byte(output)
	chnk := &pb.ReadResp{}
	//Iterating through the ByteArray for every 64 KB of chunks
	for currentByte := 0; currentByte < len(outputtwo); currentByte += chunkSize {
		if currentByte+chunkSize > len(outputtwo) {
			chnk.Chunk = outputtwo[currentByte:len(outputtwo)]
		} else {
			chnk.Chunk = outputtwo[currentByte : currentByte+chunkSize]
		}
		if err := srv.Send(chnk); err != nil {
			outputtwo = nil
			return err
		}
	}
	return nil
}

// Store implementation
func (s *IsServer) Store(rcv pb.Is_StoreServer) error {
	blob := []byte{}
	memType := ""
	for {
		point, err := rcv.Recv()
		if err != nil {
			if err == io.EOF {
				glog.Infof("Transfer of %d bytes successful", len(blob))
				break
			}
			glog.Errorf("Error while receiving: %v", err)
			return err
		}
		blob = append(blob, point.Chunk...)
		memType = point.MemoryType
	}
	s.is.SetStorageType(memType)
	output, err := s.is.Store(blob)
	if err != nil {
		glog.Infof("Store failed")
		glog.Info("imgHandle of data stored is: ", output)
		return err
	}
	return rcv.SendAndClose(&pb.StoreResp{
		StoreKeyname: output,
	})
}

// Remove implementation
func (s *IsServer) Remove(ctx context.Context, in *pb.RemoveReq) (*pb.RemoveResp, error) {
	errr := s.is.Remove(in.RemKeyname)
	if errr != nil {
		glog.Infof("gRPC Remove failed")
	}
	return &pb.RemoveResp{}, errr
}

// CloseGrpcServer closes gRPC server
func CloseGrpcServer(done chan bool) {
	done <- true
}

func getConfig(cfgType string) (string, error) {

	var buf []byte
	var err error
	err = nil

	switch cfgType {
	case "InfluxDBCfg":
		buf, err = json.Marshal(ImgDaCfg.InfluxDB)
	case "RedisCfg":
		buf, err = json.Marshal(ImgDaCfg.Redis)
	case "MinioCfg":
		buf, err = json.Marshal(ImgDaCfg.Minio)
	default:
		return "", errors.New("Not a valid config type")
	}

	return string(buf), err
}

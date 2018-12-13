/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

Explicit permissions are required to publish, distribute, sublicense, and/or sell copies of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	pb "ElephantTrunkArch/ImageStore/protobuff/go"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"flag"
	"io"
	"path/filepath"

	"context"
	"net"
	"os"

	client "ElephantTrunkArch/DataAgent/da_grpc/client/go/client_internal"
	imagestore "ElephantTrunkArch/ImageStore/go/ImageStore"
	util "ElephantTrunkArch/Util"

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
	RootCA     = "/etc/ssl/grpc_int_ssl_secrets/ca_certificate.pem"
	ServerCert = "imagestore_server_certificate.pem"
	ServerKey  = "imagestore_server_key.pem"

	ClientCert = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_certificate.pem"
	ClientKey  = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_key.pem"
)

// IsServer is used to implement ImageStore.IsServer
type IsServer struct {
	is *imagestore.ImageStore
}

// StartGrpcServer starts the ImageStore grpc server
func StartGrpcServer(redisConfigMap map[string]string, minioConfigMap map[string]string) {

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

	// Manually set the host to localhost since we are inside the docker network
	minioConfigMap["Host"] = "localhost"
	redisConfigMap["Host"] = "localhost"

	grpcClient, err := client.NewGrpcInternalClient(ClientCert, ClientKey, RootCA, "ia_data_agent", "50052")
	if err != nil {
		glog.Errorf("Error while obtaining GrpcClient object...")
		os.Exit(-1)
	}

	data, err := grpcClient.GetConfigInt("ImgStoreServerCert")
	if err != nil {
		glog.Errorf("Unable to read SERVER certificate from DataAgent %s", err)
	}

	// Read certificate binary
	certPEMBlock, certValid := data[ServerCert]
	if !certValid {
		glog.Errorf("Failed to Read Server Certificate")
		os.Exit(-1)
	}

	keyPEMBlock, keyValid := data[ServerKey]
	if !keyValid {
		glog.Errorf("Failed to Read Server Key : %s", err)
		os.Exit(-1)
	}

	// Load the certificates from binary
	keyPEMBlockByte, _ := b64.StdEncoding.DecodeString(keyPEMBlock)
	certPEMBlockByte, _ := b64.StdEncoding.DecodeString(certPEMBlock)
	certificate, err := tls.X509KeyPair(certPEMBlockByte, keyPEMBlockByte)
	if err != nil {
		glog.Errorf("Failed to Load ServerKey Pair : %s", err)
		os.Exit(-1)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, caValid := data[filepath.Base(RootCA)]
	if !caValid {
		glog.Errorf("Failed to Read CA Certificate")
		os.Exit(-1)
	}

	// Append the certificates from the CA
	caByte, _ := b64.StdEncoding.DecodeString(ca)
	if ok := certPool.AppendCertsFromPEM(caByte); !ok {
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

	glog.Infof("Waiting for redis port to be up...")
	// Wait until Redis port is up
	redisPort := os.Getenv("REDIS_PORT")
	portUp := util.CheckPortAvailability("", redisPort)
	if !portUp {
		glog.Errorf("Redis port: %s not up, so exiting...", redisPort)
		os.Exit(-1)
	}

	glog.Infof("Waiting for minio port to be up...")
	// Wait until Minio port is up
	minioPort := os.Getenv("MINIO_PORT")
	portUp = util.CheckPortAvailability("", minioPort)
	if !portUp {
		glog.Errorf("Minio port: %s not up, so exiting...", minioPort)
		os.Exit(-1)
	}

	imgStore, err := imagestore.GetImageStoreInstance(redisConfigMap, minioConfigMap)
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

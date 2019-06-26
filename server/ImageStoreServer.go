/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

Explicit permissions are required to publish, distribute, sublicense, and/or sell copies of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Package server is the server library of ImageStore APIs
package server

import (
	client "IEdgeInsights/DataAgent/da_grpc/client/go/client_internal"
	util "IEdgeInsights/Util"
	"context"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	imagestore "IEdgeInsights/ImageStore/go/ImageStore"
	pb "IEdgeInsights/ImageStore/protobuff/go"
)

var gRPCImageStoreHost = os.Getenv("IMAGESTORE_GRPC_SERVER")
var gRPCImageStorePort = os.Getenv("IMAGESTORE_PORT")

const (
	chunkSize = 4095 * 1024 // 4 MB
)

// Server Certificates
const (
	RootCA     = "/etc/ssl/grpc_int_ssl_secrets/ca_certificate.pem"
	ServerCert = "imagestore_server_certificate.pem"
	ServerKey  = "imagestore_server_key.pem"

	ClientCert = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_certificate.pem"
	ClientKey  = "/etc/ssl/grpc_int_ssl_secrets/grpc_internal_client_key.pem"
)

// IsServer is a struct used to implement ImageStore.IsServer
type IsServer struct {
	is *imagestore.ImageStore
}

// StartGrpcServer is used to start the ImageStore grpc server
// Parameters:
// 1. redisConfigMap : map[string]string
//    Refers to the ImageStore redis configurations
// 2. minioConfigMap : map[string]string
//    Refers to the ImageStore minio configurations
func StartGrpcServer(redisConfigMap map[string]string, minioConfigMap map[string]string) {

	var s *grpc.Server

	defer glog.Flush()
	if len(os.Args) < 1 {
		glog.Infof("No args passed.")
	}
	addr := ":" + gRPCImageStorePort

	minioConfigMap["Host"] = os.Getenv("IMAGESTORE_GRPC_SERVER")
	redisConfigMap["Host"] = minioConfigMap["Host"]

	// Currently setting the ports here to support dev mode
	redisConfigMap["Port"] = os.Getenv("REDIS_PORT")
	minioConfigMap["Port"] = os.Getenv("MINIO_PORT")

	devMode := os.Getenv("DEV_MODE")
	securityDisable, err := strconv.ParseBool(devMode)
	if err != nil {
		glog.Errorf("Fail to read Development Mode environment variable(DEV_MODE): %s", err)
		os.Exit(-1)
	}

	if !securityDisable {
		grpcClient, err := client.NewGrpcInternalClient(ClientCert, ClientKey, RootCA, os.Getenv("DATA_AGENT_GRPC_SERVER"), os.Getenv("GRPC_INTERNAL_PORT"))
		if err != nil {
			glog.Errorf("Error while obtaining GrpcClient object : %s", err)
			os.Exit(-1)
		}

		data, err := grpcClient.GetConfigInt("ImgStoreServerCert")
		if err != nil {
			glog.Errorf("Unable to read SERVER certificate from DataAgent %s", err)
			os.Exit(-1)
		}

		// Read certificate binary
		certPEMBlock, certValid := data[ServerCert]
		if !certValid {
			glog.Error("Failed to Read Server Certificate")
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

		// Create the TLS configuration to pass to the GRPC server
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

		//Create the gRPC server
		s = grpc.NewServer(grpc.Creds(creds))
	} else {
		s = grpc.NewServer()
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Errorf("failed to listen: %v", err)
		os.Exit(-2)
	}
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

// Read is a wrapper around gRPC go server implementation for
// Read interface.
//
// Parameters:
// 1. in : Protobuf ReadReq struct
//    Refers to the protobuf struct comprising the ImageStore server Read APIs
// 2. srv : Protobuf Is_ReadServer interface
//    Refers to the protobuf interface comprising the ImageStore server Read APIs
//
// Returns:
// 1. error
//    An error object if read fails.
func (s *IsServer) Read(in *pb.ReadReq, srv pb.Is_ReadServer) error {
	key := in.ReadKeyname
	glog.V(1).Infof("Read request for key: %v", key)
	output, err := s.is.Read(key)
	if err != nil {
		glog.Errorf("Read failed: %v for key: %v", err, in.ReadKeyname)
		return err
	}

	chnk := &pb.ReadResp{}
	outputByteArr := make([]byte, chunkSize)
	for {
		n, err := (*output).Read(outputByteArr)
		if err != nil {
			if err == io.EOF {
				// This is to send the last remaining chunk
				chnk.Chunk = outputByteArr[:n]
				if err := srv.Send(chnk); err != nil {
					return err
				}
				break
			}
			glog.Errorf("Error for ioReader.Read(): %v for key: %v", err, key)
			break
		}
		chnk.Chunk = outputByteArr[:n]
		if err := srv.Send(chnk); err != nil {
			return err
		}
	}
	return nil
}

// Store is a wrapper around gRPC go server implementation for
// Store interface.
//
// Parameters:
// 1. rcv : Protobuf Is_StoreServer interface
//    Refers to the protobuf interface comprising the ImageStore server Store APIs
//
// Returns:
// 1. error
//    An error object if store fails.
func (s *IsServer) Store(rcv pb.Is_StoreServer) error {
	blob := []byte{}
	memType := ""
	for {
		point, err := rcv.Recv()
		if err != nil {
			if err == io.EOF {
				glog.V(1).Infof("[Store]Transfer of %d bytes successful", len(blob))
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
		glog.Errorf("Store failed for key: %v", output)
		return err
	}
	glog.V(1).Infof("ImgHandle of data stored is: %v", output)
	return rcv.SendAndClose(&pb.StoreResp{
		StoreKeyname: output,
	})
}

// Remove is a wrapper around gRPC go server implementation for
// Remove interface.
//
// Parameters:
// 1. ctx : context
//    Refers to the server context
// 2. in : Protobuf RemoveReq struct
//    Refers to the protobuf struct comprising the ImageStore server Remove APIs
//
// Returns:
// 1. *pb.RemoveResp
//    Refers to the response struct instance containing the response message
// 2. error
//    An error object if remove fails.
func (s *IsServer) Remove(ctx context.Context, in *pb.RemoveReq) (*pb.RemoveResp, error) {
	key := in.RemKeyname
	err := s.is.Remove(key)
	if err != nil {
		glog.Infof("Remove failed: %v for key: %v", err, key)
	}
	return &pb.RemoveResp{}, err
}

// CloseGrpcServer closes gRPC server
//
// Parameters:
// 1. done : chan bool
//    Refers to the channel used to close the server
func CloseGrpcServer(done chan bool) {
	done <- true
}

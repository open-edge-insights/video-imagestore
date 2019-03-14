/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
#include <chrono>
#include <iostream>
#include <memory>
#include <random>
#include <string>
#include <thread>
#include <stdio.h>
#include <grpcpp/grpcpp.h>
#include <grpc/support/log.h>
#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>
#include <sstream>
#include <fstream>
#include "../../client/cpp/ImageStoreClient.cc"

using namespace std;
using grpc::Channel;
using grpc::ClientAsyncResponseReader;
using grpc::ClientContext;
using grpc::CompletionQueue;
using grpc::Status;
using grpc::Channel;
using grpc::ClientContext;
using grpc::ClientReader;
using grpc::ClientReaderWriter;
using grpc::ClientWriter;
using grpc::Status;
using ImageStore::ReadReq;
using ImageStore::ReadResp;
using ImageStore::is;

void read(const std::string& filename, std::string& data)
{
  std::ifstream file(filename.c_str(),std::ios::in);
	if (file.is_open())
	{
		std::stringstream ss;
		ss << file.rdbuf ();
		file.close ();
		data = ss.str ();
	}
	return;
}

void test_case(int argc, char** argv, const std::string& imgHandle)
{
  ReadReq request;
  ReadResp reply;
  std::string root;
  std::string key;
  std::string cert;
  int exitCondition = 1;

  if(argc < 7)
  {
    cout << "Usage: ./clientTest <imgstore_host> <imgstore_port> <img_client_cert> <img_client_key> <ca_cert> <input_file> <output_file>" << endl;
    exit(exitCondition);
  }

  read(argv[3], cert);
  read(argv[4], key);
  read(argv[5], root);

  grpc::SslCredentialsOptions opts =
		{
			root,
			key,
			cert
		};

  std::string endpoint;
  endpoint = argv[1];
  endpoint += ":";
  endpoint += argv[2];
  std::cout << "Endpoint: " << endpoint << std::endl;
  ImageStoreClient gclient(grpc::CreateChannel(endpoint,
                        grpc::SslCredentials(opts)));

  std::string storedata;
  std::ifstream in;
  in.read(argv[6], std::ios::binary);
  in >> storedata;

  std::cout << "-------------- Calling Store --------------" << std::endl;
  std::string keyname = gclient.Store(imgHandle, storedata);
  cout << "Image handle:" << keyname << endl ;

  int iterations = 20;
  float totalTimeTaken = 0.0;

  std::cout << "-------------- Calling Read --------------" << std::endl;
  for(int i = 0; i < iterations; i++) {

    const clock_t begin_time = clock();
    std::string response = gclient.Read(keyname);
    const clock_t end_time = clock();
    float timeTaken = float(end_time - begin_time) / CLOCKS_PER_SEC;
    std::cout << "Time taken for one read call:" << timeTaken << std::endl;
    totalTimeTaken += timeTaken;
    std::ofstream out;
    out.open(argv[7], std::ios::binary);
    out << response;
    out.close();
  }
  std::cout << "Total time taken for "+std::to_string(iterations)+" read calls:" << totalTimeTaken / iterations << std::endl;
  
  std::cout << "-------------- Calling Remove --------------" << std::endl;
  bool remove_response = gclient.Remove(keyname);
  cout << "Remove status :" << remove_response << endl;
  return;
}

int main(int argc, char** argv) {
  int returnCondition = 0;

  // Testing redis gRPC calls
  test_case(argc, argv, "inmemory");

  // Testing minio gRPC calls
  test_case(argc, argv, "persistent");

  return returnCondition;
}

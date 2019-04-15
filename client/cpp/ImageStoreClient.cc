/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

#include <iostream>
#include <memory>
#include <string>
#include <grpcpp/grpcpp.h>
#include <grpc/support/log.h>
#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>
#include <cstdlib>
#include "../../protobuff/cpp/ImageStore.grpc.pb.h"

using namespace std;
using grpc::ClientAsyncResponseReader;
using grpc::ClientContext;
using grpc::CompletionQueue;
using grpc::Status;
using grpc::Channel;
using grpc::ClientContext;
using grpc::ClientReader;
using grpc::ClientReaderWriter;
using grpc::ClientWriter;
using ImageStore::ReadReq;
using ImageStore::ReadResp;
using ImageStore::RemoveReq;
using ImageStore::RemoveResp;
using ImageStore::StoreReq;
using ImageStore::StoreResp;
using ImageStore::is;

class ImageStoreClient{
  public:
  ImageStoreClient(std::shared_ptr<Channel> channel)
        : _stub(is::NewStub(channel)) {}

  /*
      Read is a wrapper around gRPC C++ client implementation
      for Read gRPC interface.
      Arguments:
      imgHandle(string): key for ImageStore
      Returns:
      The consolidated string(value from ImageStore) associated with
      that imgHandle
  */
  std::string Read(const std::string& imgHandle)
  {
      ReadReq request;
      request.set_readkeyname(imgHandle);
      ReadResp reply;
      ClientContext context;
      std::cout << __FILE__ << " : " << __LINE__ << " : " << "ImageHandle: " << imgHandle << std::endl;
      std::unique_ptr<grpc::ClientReader<ReadResp> > reader(_stub->Read(&context, request));
      std::string response = "";
      while (reader->Read(&reply)) {
        response = response + reply.chunk();
      }
      Status status = reader->Finish();
      bool statusCheck = status.ok();
      if (statusCheck != 1) {
        std::cout << "Transfer failed." << std::endl;
        return "";
      }
      return response;
  }

  /*
      Remove is a wrapper around gRPC C++ client implementation
      for Remove gRPC interface.
      Arguments:
      imgHandle(string): key for ImageStore
      Returns:
      The consolidated boolean if whether the consolidated
      value was removed.
  */
  bool Remove(const std::string& imgHandle)
  {
      RemoveReq request;
      request.set_remkeyname(imgHandle);
      RemoveResp reply;
      ClientContext context;
      std::cout << __FILE__ << " : " << __LINE__ << " : " << "ImageHandle: " << imgHandle << std::endl;
      Status status = _stub->Remove(&context, request, &reply);
      bool statusCheck = status.ok();
      if (statusCheck != 1) {
        std::cout << "Remove failed." << std::endl;
        return false;
      }
      return true;
  }

  /*
      Store is a wrapper around gRPC C++ client implementation
      for Store gRPC interface.
      Arguments:
      memoryType(string): memoryType for ImageStore
      imgFrame(string): image to be stored in ImageStore
      Returns:
      The consolidated imgHandle string of imgFrame stored.
  */
  std::string Store(const std::string& memoryType, const std::string& imgFrame)
  {
      StoreReq request;
      StoreResp reply;
      ClientContext context;
      request.set_memorytype(memoryType);
      request.set_chunk(imgFrame);
      std::unique_ptr<ClientWriter<StoreReq>> writer(_stub->Store(&context, &reply));
      if (!writer->Write(request)) {
        std::cout << "Store failed." << endl;
        exit(1);
      }
      writer->WritesDone();
      Status status = writer->Finish();
      if (status.ok()) {
          std::cout << "Store succeeded." << std::endl;
          return reply.storekeyname();
      }
  }

  private:
    std::unique_ptr<is::Stub> _stub;
};

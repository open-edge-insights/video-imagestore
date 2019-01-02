"""
Copyright (c) 2018 Intel Corporation.

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
"""

# Python grpc client implementation

import grpc
import json
import logging as log
import sys
import os
from ImageStore.protobuff.py import ImageStore_pb2 as is_pb2
from ImageStore.protobuff.py import ImageStore_pb2_grpc as is_pb2_grpc
from Util.crypto.encrypt_decrypt import SymmetricEncryption

chunk_size = 4095*1024


class GrpcImageStoreClient(object):

    def __init__(self, clientCert, clientKey,
                 caCert,  hostname="localhost", port="50055"):
        """
            GrpcImageStoreClient constructor.
            Keyword Arguments:
            clientCert - refers to the imagestore client certificate
            clientKey - refers to the imagestore client key
            caCert - refers to the ca certificate
            hostname - refers to hostname/ip address of the m/c
                       where DataAgent module of ETA is running
                       (default: localhost)
            port     - refers to gRPC port (default: 50055)
        """
        self.hostname = hostname
        self.port = port
        if 'IMAGESTORE_SERVER' in os.environ:
            self.hostname = os.environ['IMAGESTORE_SERVER']
        addr = "{0}:{1}".format(self.hostname, self.port)
        log.debug("Establishing secure grpc channel to %s", addr)

        if 'grpc_int_ssl_secrets' in caCert:
            key = os.environ["SHARED_KEY"]
            nonce = os.environ["SHARED_NONCE"]
            symEncrypt = SymmetricEncryption(key)
            ca_certs = symEncrypt.DecryptFile(caCert, nonce)
        else:
            with open(caCert, 'rb') as f:
                ca_certs = f.read()

        with open(clientKey, 'rb') as f:
            client_key = f.read()

        with open(clientCert, 'rb') as f:
            client_certs = f.read()

        try:
            credentials = grpc.ssl_channel_credentials(
                root_certificates=ca_certs, private_key=client_key,
                certificate_chain=client_certs)

        except Exception as e:
            log.error("Exception Occured : ", e.msg)
            raise Exception

        channel = grpc.secure_channel(addr, credentials)
        self.stub = is_pb2_grpc.isStub(channel)

    def Read(self, imgHandle):
        """
            Read is a wrapper around gRPC python client implementation
            for Read gRPC interface.
            Arguments:
            config(string): imgHandle
            Returns:
            The byte stream corresponding to the config value
        """
        log.debug("Inside Read() client wrapper...")
        response = self.stub.Read(is_pb2.ReadReq(readKeyname=imgHandle),
                                  timeout=1000)
        outputBytes = b''
        for resp in response:
            outputBytes += resp.chunk
        log.debug("Sending the response to the caller...")
        return outputBytes

    def Store(self, byteStream, memType):
        """
            Store is a wrapper around gRPC python client implementation
            for Store gRPC interface.
            Arguments:
            config(string): byte stream to be stored
            Returns:
            The imgHandle corresponding to the stored config value
        """
        log.debug("Inside Store() client wrapper...")
        data = self._chunkfunction(byteStream, memType)
        response = self.stub.Store(data, timeout=1000)
        log.debug("Sending the response to the caller...")
        return response.storeKeyname

    def Remove(self, imgHandle):
        """
            Remove is a wrapper around gRPC python client implementation
            for Remove gRPC interface.
            Arguments:
            config(string): imgHandle
            Returns:
            None
        """
        log.debug("Inside Remove() client wrapper...")
        response = self.stub.Remove(is_pb2.RemoveReq(remKeyname=imgHandle),
                                    timeout=1000)
        log.debug("Sending the response to the caller...")
        return None

    def _chunkfunction(self, byteStream, memType):
        """
            ChunkFunction is used to return the generator object which
            is required by the gRPC store server interface.
        """
        for i in range(0, len(byteStream), chunk_size):
            yield is_pb2.StoreReq(chunk=byteStream[i:i + chunk_size],
                                  memoryType=memType)

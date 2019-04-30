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

import grpc
import json
import logging as log
import sys
import os
# from ImageStore.protobuff.py import ImageStore_pb2 as is_pb2
# from ImageStore.protobuff.py import ImageStore_pb2_grpc as is_pb2_grpc
# from Util.crypto.encrypt_decrypt import SymmetricEncryption
from ImageStore.client.py.client import GrpcImageStoreClient
from DataAgent.da_grpc.client.py.client_internal.client \
    import GrpcInternalClient
from Util.util import write_certs

chunk_size = 4095*1024

ROOTCA_CERT = '/etc/ssl/ca/ca_certificate.pem'
IM_CLIENT_CERT = '/etc/ssl/imagestore/imagestore_client_certificate.pem'
IM_CLIENT_KEY = '/etc/ssl/imagestore/imagestore_client_key.pem'

GRPC_CERTS_PATH = "/etc/ssl/grpc_int_ssl_secrets"
CLIENT_CERT = GRPC_CERTS_PATH + "/grpc_internal_client_certificate.pem"
CLIENT_KEY = GRPC_CERTS_PATH + "/grpc_internal_client_key.pem"
CA_CERT = GRPC_CERTS_PATH + "/ca_certificate.pem"


class GrpcImageStoreInternalClient(GrpcImageStoreClient):
    """
    This class represents grpc ImageStore client
    """

    def __init__(self, hostname="localhost", port="50055", dev_mode=False):
        try:
            if dev_mode:
                GrpcImageStoreClient.__init__(self, hostname=hostname,
                                              port=port)
            else:
                client = GrpcInternalClient(CLIENT_CERT, CLIENT_KEY, CA_CERT)
                self.resp = client.GetConfigInt("ImgStoreClientCert")

                # Write File
                file_list = [IM_CLIENT_CERT, IM_CLIENT_KEY]
                write_certs(file_list, self.resp)
                GrpcImageStoreClient.__init__(self, IM_CLIENT_CERT,
                                              IM_CLIENT_KEY,
                                              ROOTCA_CERT, hostname=hostname,
                                              port=port)
        except Exception as e:
            log.error("Exception Occured : " + str(e))
            raise Exception

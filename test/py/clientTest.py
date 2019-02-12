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

# === Python grpc client test ===

import logging
import argparse
import hashlib
import datetime
import time
import sys
import os
import time
from ImageStore.client.py.client import \
    GrpcImageStoreClient

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s : %(levelname)s : \
                    %(name)s : [%(filename)s] :' +
                    '%(funcName)s : in line : [%(lineno)d] : %(message)s')
log = logging.getLogger("GRPC_TEST")

# === gRPC test library parse_args method ===
def parse_args():
    """
    parse_args is used to parse the command line arguments
        
    Parameters
    ----------
    1. None

    Returns
    -------
    1. Set of cli arguments

    """
    parser = argparse.ArgumentParser()

    parser.add_argument('--hostname', dest='hostname',
                        help='ip address of the node running IEI')

    parser.add_argument('--port', dest='port',
                        help='CA_Cert')

    parser.add_argument('--ca-cert', dest='ca_cert',
                        help='CA_Cert')

    parser.add_argument('--client-key', dest='client_key',
                        help='Client Key')

    parser.add_argument('--client-cert', dest='client_cert',
                        help='Client_Cert')

    parser.add_argument('--input_file', help='input image file')

    parser.add_argument('--output_file', help='output image file')

    return parser.parse_args()

def test_case(imgHandle):

    args = parse_args()
    client = GrpcImageStoreClient(args.client_cert, args.client_key,
                                  args.ca_cert, hostname=args.hostname)
    
    inputFile = args.input_file
    outputFile = args.output_file

    inputBytes = None
    with open(inputFile, "rb") as f:
            inputBytes = f.read()

    # Testing Store("value") gRPC call
    keyname = client.Store(inputBytes, imgHandle)
    totalTime = 0.0

    # Testing Read("imgHandle") gRPC call
    iter1 = 20
    for i in range(iter1):
        start = time.time()
        outputBytes = client.Read(keyname)
        end = time.time()
        timeTaken = end - start
        log.info("Time taken for one read call: %f secs", timeTaken)
        totalTime += timeTaken

    log.info("Average time taken for Read() %d calls: %f secs",
             iter1, totalTime / iter1)

    log.info("Writing the binary data received into a file: %s",
                 outputFile)
    with open(outputFile, "wb") as outfile:
        outfile.write(outputBytes)

    digests = []
    for filename in [inputFile, outputFile]:
        hasher = hashlib.md5()
        with open(filename, 'rb') as f:
            buf = f.read()
            hasher.update(buf)
            a = hasher.hexdigest()
            digests.append(a)
            log.info("Hash for filename: %s is %s", filename, a)

    if digests[0] == digests[1]:
        log.info("md5sum for the files match")
    else:
        log.info("md5sum for the files doesn't match")

    # Testing Remove("imgHandle") gRPC call
    client.Remove(keyname)

if __name__ == '__main__':

    # Testing Redis gRPC calls
    test_case('inmemory')

    # Testing Minio gRPC calls
    test_case('persistent')

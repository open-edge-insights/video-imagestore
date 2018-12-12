#!/usr/bin/python3
# Copyright (c) 2018 Intel Corporation.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
"""Entrypoint for starting minio in Docker.
"""
import os
import sys
import subprocess as sub
from da_grpc.client.py.client_internal.client import GrpcInternalClient


def main():
    """Main method
    """
    try:
        client = GrpcInternalClient()
        config = client.GetConfigInt('MinioCfg')
        os.environ['MINIO_ACCESS_KEY'] = config['AccessKey']
        os.environ['MINIO_SECRET_KEY'] = config['SecretKey']
        os.environ['MINIO_REGION'] = 'gateway'
        sub.check_call(['./minio', 'server', '/data'])
    except Exception as e:
        print('Error with Minio server:', e)
        sys.exit(-1)


if __name__ == '__main__':
    main()

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

from ImageStore.py.imagestore import ImageStore
import os
import argparse


def parse_args():
    '''Parse cmd line arguments
    '''
    parser = argparse.ArgumentParser()
    parser.add_argument('--input_file', help='input image file')
    parser.add_argument('--output_file', help='output image file')

    return parser.parse_args()

try:
    args = parse_args()
    inputFile = args.input_file
    outputFile = args.output_file
    storeInMem = ImageStore()
    storeInMem.setStorageType('inmemory')
    keyname = storeInMem.store(bytes(0x00))
    print("Return of Store Operation Keyname : ", keyname)
    print("Reading the Stored Data", storeInMem.read(keyname))  # binaryvalue
    print("Removing the Stored data", storeInMem.remove(keyname))  # status
    with open(inputFile, 'rb') as f:
        binary = f.read()
        try:
            storeBinaryKeyName = storeInMem.store(binary)
        except Exception as e:
            print("Exception :", e)
        f.close()
    print("Binary File Stored KeyName : ", storeBinaryKeyName)
    print("Reading & Storing Back to New File")
    with open(outputFile, 'wb') as f:
        try:
            binaryvalue = storeInMem.read(storeBinaryKeyName)
        except Exception as e:
            print("Exception :", e)
        # Store Binary Contains 1 Value, Keyname
        if binaryvalue is not None:
            f.write(binaryvalue)
            print("Retrieved binary Stored in this path ", os.path.realpath(
                                                                f.name))
        f.close()

except Exception as e:
    print("Exception is here", e)

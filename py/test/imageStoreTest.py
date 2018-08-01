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
    keyname = storeInMem.store("Value")
    print("Return of Store Operation Keyname : ", keyname)
    print("Reading the Stored Data", storeInMem.read(keyname))  #binaryvalue
    print("Removing the Stored data", storeInMem.remove(keyname)) #status
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
        if binaryvalue != None:
            f.write(binaryvalue)
            print("Retrieved binary Stored in this path ", os.path.realpath(f.name))
        f.close()

except Exception as e:
    print("Exception is here", e)

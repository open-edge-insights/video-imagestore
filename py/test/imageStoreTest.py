from ImageStore.py.imagestore import ImageStore
import os

try:
    storeInMem = ImageStore()
    storeInMem.setStorageType('inmemory')
    status, keyname = storeInMem.store("Value")
    print("Return of Store Operation Status : ", status, "Keyname : ", keyname)
    print("Reading the Stored Data", storeInMem.read(keyname))  #Status, #binaryvalue
    print("Removing the Stored data", storeInMem.remove(keyname)) #statys, ,#message
    with open("./ImageStore/py/test/test.jpg", 'r+b') as f:
        binary = f.read()
        storeBinaryStatus, storeBinaryKeyName = storeInMem.store(binary)
        f.close()
    print("Binary File Stored Status : ", storeBinaryStatus, "KeyName : ", storeBinaryKeyName)
    print("Reading & Storing Back to New File")
    with open("./ImageStore/py/test/fromredis.jpg", 'wb') as f:
        status, binaryvalue = storeInMem.read(storeBinaryKeyName)
        # Store Binary Contains 2 Values, 1st is Status and Second is the Keyname
        if status:
            f.write(binaryvalue)
            print("Retrieved binary Stored in this path ", os.path.realpath(f.name))
        f.close()

except Exception as e:
    print("Exception is here", e)

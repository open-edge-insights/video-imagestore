from ImageStore.py.inmemory.inmemorystore import InMemory
from ImageStore.py import output as output
from DataAgent.da_grpc.client.client import GrpcClient
from Util.exception import DAException


class ImageStore():
    """
        This is Derived from inMemory  base Implementations
        This Class gives the abstracted implementation of various Storage
        api implementations of inMemory and also filesystem.

    """
    def __init__(self):
        """
            Instantiate the objects based on the memoryType and the config of
            particular memoryType's storage system. Based on this other Storage
            operations get's handled

        """
        try:
            self.config = GrpcClient.GetConfigInt("RedisCfg")
            self.config["InMemory"] = "redis"
            self._initializeinMemory()
        except Exception as e:
            raise DAException("Seems to be some issue with gRPC Server. Exception: {0}".format(e))
        # TODO: plan a better approach to set this config later, not to be in
        # DataAgent.conf file as it's not intended to
        # be known to user

    def _initializeinMemory(self):
        """
            This API is to initialize the memory Object and it will Instantiate
            with config parameters

        """
        # config = config.value
        # policy = self.config.retentionpolicy
        try:
            self.inmemoryredis = InMemory(self.config)
        except Exception as e:
            raise e

    def setStorageType(self, memoryType):
        """
            This api to set the StorageType for doing store operations.
            It Accepts storagememory type as argument. Either 'inmemory' or
            'fs'. Currently system Supports only inMemory.
        """
        if memoryType is not None:
            memoryType = memoryType.lower()
            if memoryType == 'inmemory':
                self.memoryType = memoryType
            else:
                raise Exception(output.handleOut('NotSupported',
                                                  self.memoryType))
        else:
            raise Exception(output.handleOut('NotSupported',
                                             memoryType))

    def getKeyList(self):
        """

            This is to get the Key's list from inmemory or File Storage
            It Accepts no attribute. Based on the Instantiated object's
            corresponding storage. It retrieves the entire key list.

            It return data in tuple with 2 values. 1st is the Status
            2nd is the value or description

        """
        returndata = None
        try:
            if self.memoryType == 'inmemory':
                returndata = self.inmemoryredis.getKeyList()
            else:
                returndata = output.handleOut('NotSupported', self.memoryType)
        except Exception as e:
            raise e
        return returndata

    def read(self, keyname):
        """
            This is to read the binary data from inmemory or File Storage
            based on keyname. keyname is the attribute to this method.

            Based on the Instantiated object's corresponding storage.
            It retrieves the data from the storage.

            It return data in tuple with 2 values. 1st is the Status
            2nd is the returned value or description.

        """
        returndata = None
        try:
            if 'inmem' in keyname:
                returndata = self.inmemoryredis.read(keyname)
            else:
                returndata = output.handleOut('NotSupported', 'keyname is not having any\
                                                inMemory key pattern')
        except Exception as e:
            raise e
        return returndata

    def store(self, binarydata):
        """
            This is to persist the data in inmemory or File Storage. It
            Accepts binarydata as its argument. Based on the Instantiated
            Storage, It stores the binarydata and returns the keyname.

            It return data in tuple with 2 values. 1st is the Status
            2nd is the returned value or description.

        """
        returndata = None
        try:
            if self.memoryType is not None:
                if self.memoryType == 'inmemory':
                    returndata = self.inmemoryredis.store(binarydata)
                else:
                    returndata = output.handleOut('NotSupported', self.memoryType)
            else:
                returndata = output.handleOut('error', 'Please use \
                            setStorageType() api before using store operations')
        except Exception as e:
            raise DAException("Seems to be some issue with Redis. Exception: {0}".format(e))
        return returndata

    def remove(self, keyname):
        """
            This is to remove the stored file from inmemory or File Storage.
            It Accepts keyname as attribute. Based on the Instantiated Storage,
            It removes the data based on the keyname.

        """
        returndata = None
        try:
            if 'inmem' in keyname:
                returndata = self.inmemoryredis.remove(keyname)
            else:
                returndata = output.handleOut('NotSupported', self.memoryType)
        except Exception as e:
            raise e
        return returndata

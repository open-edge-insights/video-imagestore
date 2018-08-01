from ImageStore.py.inmemory.redisStore.redis import RedisConnect
from DataAgent.da_grpc.client.client import GrpcClient
from ImageStore.py import output as output


class InMemory():

    """
        This is Derived from various redisStore base Implementations
        This Class gives the abstracted implementation of various inMemory
        storage. Class get's Instantiated based on the memoryType request
        received from api.

    """

    def __init__(self, config):
        """
            Instantiate the Objects based on inMemoryType and set defaults
            This Accepts inMemoryType and Policy as it's arguments. Based on
            policy & memory type corresponding inmemory storage operations
            will be handled

        """
        try:
            self.inMemoryType = config["InMemory"].lower()
            if self.inMemoryType == "redis":
                self.redisStore = RedisConnect(config)
            else:
                raise Exception(output.handleOut('NotSupported',
                                                  self.inMemoryType))
        except Exception as e:
            raise e

    def getKeyList(self):
        """
            Get's the keys list from inmemory storage. This has no
            attributes based on the Instantiated storage it will get from
            corresponding storage stored keys.

        """
        returndata = None
        try:
            returndata = self.redisStore.getKeyList()
        except Exception as e:
            raise e
        return returndata

    def read(self, keyname):
        """
            retrieve's the stored data from inmemory based on the key passed
            as attribute to this method.

        """
        returndata = None
        try:
            returndata = self.redisStore.read(keyname)
        except Exception as e:
            raise e
        return returndata

    def store(self, binarydata):
        """
            Stores data inMemory. This Accepts binarydata as it's argument.
            Based on the Instantiated storage it's store the binarydata and
            return the key

        """
        returndata = None
        try:
            returndata = self.redisStore.store(binarydata)
        except Exception as e:
            raise e
        return returndata

    def remove(self, keyname):
        """
            Removes data from inmemory based on keyname. It Accepts
            keyname as argument. Based on the Instantiated storage. It
            removes key from corresponding storage.

        """
        returndata = None
        try:
            returndata = self.redisStore.remove(keyname)
        except Exception as e:
            raise e
        return returndata
from ImageStore.py.inmemory.redisStore.redis import redisConnect
from ImageStore.py import settings as config
from ImageStore.py import output as output


class InMemory():

    """
        This is Derived from various redisStore base Implementations
        This Class gives the abstracted implementation of various inMemory
        storage. Class get's Instantiated based on the memoryType request
        received from api.

    """

    def __init__(self, inMemoryType, policy):
        """
            Instantiate the Objects based on inMemoryType and set defaults
            This Accepts inMemoryType and Policy as it's arguments. Based on
            policy & memory type corresponding inmemory storage operations
            will be handled

        """
        self.inMemoryType = inMemoryType.lower()
        self.config = config.value
        if self.inMemoryType == self.config.inmemory:
            self.config = config.value
            self.policy = policy
            self.redisStore = redisConnect(host=self.config.host, \
            port=self.config.port, policy=self.policy)
        else:
            print(output.handleOut('NotSupported', self.inMemoryType))

    def getKeyListInMemory(self):
        """
            Get's the keys list from inmemory storage. This has no
            attributes based on the Instantiated storage it will get from
            corresponding storage stored keys.

        """
        returndata = ()
        if self.inMemoryType == self.config.inmemory:
            returndata = self.redisStore.getKeyListfromRedis()
        else:
            returndata = output.handleOut('NotSupported', self.inMemoryType)

        return returndata

    def getDataFromMemory(self, keyname):
        """
            retrieve's the stored data from inmemory based on the key passed
            as attribute to this method.

        """
        returndata = ()
        if self.inMemoryType == self.config.inmemory:
            returndata = self.redisStore.getDataFromRedis(keyname)
        else:
            returndata = output.handleOut('NotSupported', self.inMemoryType)

        return returndata

    def storeDatainMemory(self, binarydata):
        """
            Stores data inMemory. This Accepts binarydata as it's argument.
            Based on the Instantiated storage it's store the binarydata and
            return the key

        """
        returndata = ()
        if self.inMemoryType == self.config.inmemory:
            returndata = self.redisStore.storeDatainRedis(binarydata)
        else:
            returndata = output.handleOut('NotSupported', self.inMemoryType)

        return returndata

    def removeFromMemory(self, keyname):
        """
            Removes data from inmemory based on keyname. It Accepts
            keyname as argument. Based on the Instantiated storage. It
            removes key from corresponding storage.

        """
        returndata = ()
        if self.inMemoryType == self.config.inmemory:
            returndata = self.redisStore.removeFromRedis(keyname)
        else:
            returndata = output.handleOut('NotSupported', self.inMemoryType)

        return returndata

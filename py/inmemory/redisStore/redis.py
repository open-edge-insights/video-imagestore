import redis
import uuid
from pytimeparse.timeparse import timeparse


class RedisConnect:

    def __init__(self, config):
        try:
            self.redis_db = redis.StrictRedis(host=config["Host"],
                                              port=config["Port"],
                                              db=0)
        except Exception as e:
            raise DAException("Seems to be some issue with Redis. Exception: {0}".format(e))
        self.retention = timeparse(config["Retention"])

    def getDataFromRedis(self, keyname):
        """
            Base implementation to get the data from redis
        """
        returndata = ()
        try:
            if self.isKeyExistsinRedis(keyname)[1]:
                returndata = True, self.redis_db.get(keyname)
            else:
                returndata = False, 'This is key is not in inmemory (redis)'
        except Exception as e:
            raise e
        return returndata

    def getKeyListfromRedis(self):
        """
            Base implementation to read the data from redis
        """
        returndata = ()
        try:
            fileslist = True, self.redis_db.keys()
        except Exception as e:
            raise e
        return returndata

    def isKeyExistsinRedis(self, keyname):
        """
            Base implementation to check key exists or not
        """
        status = ()
        try:
            status = True, self.redis_db.exists(keyname)
        except Exception as e:
            raise e
        return status

    def removeFromRedis(self, keyname):
        """
            Base implementation to remove data from redis

        """
        returndata = ()
        try:
            if self.isKeyExistsinRedis(keyname)[1]:
                if self.redis_db.delete(keyname) == 1:
                    returndata = True, 'Removed Succesfully'
                else:
                    returndata = False, 'Not Removed'
            else:
                returndata = False, 'This is key is not in inmemory (redis)'
        except Exception as e:
            raise e
        return returndata

    def generateRedisKey(self):
        """
            This generates key to store data in redis
        """
        keyname = 'inmem_' + str(uuid.uuid1())[:8]

        return keyname

    def storeDatainRedis(self, binarydata):
        """
            Base Implementation to store data in redis
        """
        keyname = self.generateRedisKey()
        returndata = ()
        try:
            if(not self.isKeyExistsinRedis(keyname)[1]):
                if self.retention:
                    store = self.redis_db.set(keyname,
                                            binarydata,
                                            self.retention)
                else:
                    store = self.redis_db.set(keyname, binarydata)

                returndata = True, keyname
            else:
                returndata = False, 'Already Key exists'
        except Exception as e:
            raise e

        return returndata

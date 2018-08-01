import redis
import uuid
from pytimeparse.timeparse import timeparse
from Util.exception import DAException


class RedisConnect:

    def __init__(self, config):
        try:
            self.redis_db = redis.StrictRedis(host=config["Host"],
                                              port=config["Port"],
                                              db=0)
        except Exception as e:
            raise DAException("Seems to be some issue with Redis. Exception: {0}".format(e))
        self.retention = timeparse(config["Retention"])

    def read(self, keyname):
        """
            Base implementation to get the data from redis
        """
        returndata = None
        try:
            if self.isKeyExists(keyname):
                returndata = self.redis_db.get(keyname)
        except Exception as e:
            raise e
        return returndata

    def getKeyList(self):
        """
            Base implementation to read the data from redis
        """
        try:
            fileslist = self.redis_db.keys()
        except Exception as e:
            raise e
        return fileslist

    def isKeyExists(self, keyname):
        """
            Base implementation to check key exists or not
        """
        status = False
        try:
            status = self.redis_db.exists(keyname)
        except Exception as e:
            raise e
        return status

    def remove(self, keyname):
        """
            Base implementation to remove data from redis

        """
        returndata = False
        try:
            if self.isKeyExists(keyname):
                if self.redis_db.delete(keyname) == 1:
                    returndata = True
        except Exception as e:
            raise e
        return returndata

    def generateKey(self):
        """
            This generates key to store data in redis
        """
        try:
            keyname = 'inmem_' + str(uuid.uuid1())[:8]
        except Exception as e:
            raise e
        return keyname

    def store(self, binarydata):
        """
            Base Implementation to store data in redis
        """
        keyname = self.generateKey()
        returndata = None
        try:
            if(self.isKeyExists(keyname) == False):
                if self.retention:
                    store = self.redis_db.set(keyname,
                                            binarydata,
                                            self.retention)
                else:
                    store = self.redis_db.set(keyname, binarydata)

                returndata = keyname
        except Exception as e:
            raise e
        return returndata

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
            raise DAException("Seems to be some issue with Redis." +
                              " Exception: {0}".format(e))
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
            if(self.isKeyExists(keyname) is False):
                if self.retention:
                    store = self.redis_db.set(
                                            keyname,
                                            binarydata,
                                            self.retention)
                else:
                    store = self.redis_db.set(keyname, binarydata)

                returndata = keyname
        except Exception as e:
            raise e
        return returndata

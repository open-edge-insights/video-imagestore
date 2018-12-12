# Copyright (c) 2018 Intel Corporation.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
"""Persistent image store module.
"""
from ImageStore.py.persistent.minio import MinioStorage
from Util.exception import DAException


def get_config_key(storage_type):
    """Get the configuration to retrieve over gRPC to configure the persistent
    storage.

    :param storage_type: Storage type for persistent storage
    :type: str
    :return: String value to request configuration from gRPC
    :rtype: str
    """
    storage_type.lower()
    if storage_type == 'minio':
        return 'MinioCfg'
    else:
        raise DAException('Unknown storage type: {}'.format(storage_type))


class PersistentImageStore:
    """PersistentImageStore object which provides an abstraction layer over
    the underlying storage technology in use.
    """
    def __init__(self, storage_type, config):
        """Constructor.

        :param config: Image storage configuration object
        :type: dict
        """
        try:
            storage_type = storage_type.lower()

            if storage_type == 'minio':
                self.storage = MinioStorage(config)
            else:
                raise DAException(
                        'Unknown persistent storage type: {}'.format(
                            storage_type))
        except KeyError as e:
            raise DAException('Config missing key: {}'.format(e))

    def getKeyList(self):
        """Get the list of keys in the persistent storage.
        """
        return self.storage.getKeyList()

    def read(self, key):
        """Retrieve the data stored for the given key. This method will raise
        an exception if the key does not exist in the storage.

        :param key: Image key
        :type: str
        :return: numpy.ndarray
        """
        return self.storage.read(key)

    def store(self, data):
        """Store the given binary data in the persistent image store.

        :param data: Binary blob of data
        :type: bytearray
        :return: Key for the stored value
        :rtype: str
        """
        return self.storage.store(data)

    def remove(self, key):
        """Remove the binary blob with the given key from the image store.

        :param key: Key of the binary blob to remove
        :type: str
        """
        return self.storage.remove(key)

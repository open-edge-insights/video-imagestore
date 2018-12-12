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
"""Implementation of persistent storage using Minio
"""
import uuid
import urllib
import datetime
import io
import logging
from threading import Timer
from minio import Minio
from minio.error import ResponseError, NoSuchKey, BucketAlreadyExists, \
        BucketAlreadyOwnedByYou

from Util.exception import DAException


class MinioStorage:
    """Minio persistent storage interface.
    """
    def __init__(self, config):
        """Constructor.

        :param config: Minio connection configuration
        :type: dict
        """
        try:
            region = 'gateway'
            self.log = logging.getLogger(__name__)
            self.timer = None
            # May want to put in the config
            self.bucket_name = 'image-store-bucket'
            self.retention_time = datetime.timedelta(
                    seconds=int(config['RetentionTime']))

            # Creating connection to the Minio bucket
            ssl_str = config['Ssl'].lower()
            if ssl_str == 'false':
                ssl = False
            elif ssl == 'true':
                ssl = True
            else:
                raise DAException('Ssl key must be true or false')

            self.client = Minio(
                    '{0}:{1}'.format(config['Host'], config['Port']),
                    access_key=config['AccessKey'],
                    secret_key=config['SecretKey'],
                    secure=ssl,
                    region=region)

            try:
                # Make the bucket in Minio if it does not exist
                self.client.make_bucket(self.bucket_name, location=region)
            except (BucketAlreadyOwnedByYou, BucketAlreadyExists):
                pass
            except ResponseError:
                raise

            # Remove any objects in Minio that should be expired based on the
            # retention policies
            self._clean_store()
        except KeyError as e:
            raise DAException('Config missing key: {}'.format(e))

    def read(self, key):
        """Retrieve a specific binary blob from Minio.

        :param key: Blob handle
        :type: str
        :return: Binary blob in the database, or None if the object does not
            exist
        :rtype: bytearray or None
        """
        try:
            self.log.info('Getting frame %s', key)
            data = self.client.get_object(self.bucket_name, key)
            return data.read()
        except NoSuchKey:
            return None

    def remove(self, key):
        """Remove a given key from the database.

        :param key: Key of the object to remove
        :type: str
        """
        self.client.remove_object(self.bucket_name, key)

    def store(self, data):
        """Store the given data blob in the database.

        :param data: Binary blob to store in Minio
        :type: bytearray
        :return: Key assigned to the blob in storage
        :rtype: str
        """
        key = self._gen_key()
        data_len = len(data)
        data = io.BytesIO(data)  # Minio requires an io.RawIOBase object
        self.client.put_object(self.bucket_name, key, data, data_len)
        return key

    def getKeyList(self):
        """Get list of blob keys in the database.

        :return: List of keys
        :type: list
        """
        return map(lambda o: o.object_name,
                   self.client.list_objects(self.bucket_name))

    def isKeyExists(self, key):
        """Check if a given key exists in the database.

        :return: True/False
        :type: bool
        """
        return True if key in self.getKeyList() else False

    def _gen_key(self):
        """Helper method to generate a key for a new blob, and verify that it
        is unique in the database.

        :return: New key
        :rtype: str
        """
        # Generate initial key
        key = 'persist_' + str(uuid.uuid1())[:8]

        # If the key exists in the database, keep generating keys until it
        # does not exist

        # TODO: Should probably limit this to a certain number of tries to
        # guarentee that it can never loop forever
        while self.isKeyExists(key):
            key = 'persist_' + str(uuid.uuid1())[:8]

        return key

    def _clean_store(self):
        """Private helper method to execute the retention policies on the
        objects in Minio.
        """
        self.log.debug('Cleaning up Minio object store')

        objects = self.client.list_objects(self.bucket_name)
        now = datetime.datetime.now(datetime.timezone.utc)
        keys_to_remove = []

        for obj in objects:
            if (now - obj.last_modified) > self.retention_time:
                keys_to_remove.append(obj.object_name)

        self.log.debug('Expired objects: %s', keys_to_remove)
        errs = self.client.remove_objects(self.bucket_name, keys_to_remove)

        for err in errs:
            self.log.error('Minio failed to remove object "%s"', err)

        # Start the tim er for the clean up proceadure to run in the given
        # amount of seconds
        self.timer = Timer(
                self.retention_time.total_seconds(), self._clean_store)
        self.timer.daemon = True
        self.timer.start()

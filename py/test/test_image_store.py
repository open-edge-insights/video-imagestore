# Copyright (c) 2018 Intel Corporation.
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
"""Unit tests for testing the ImageStore functionality.
"""
import time
import unittest
from persistent import PersistentImageStore 


class TestImageStore(unittest.TestCase):
    """ImageStore unit tests
    """
    def setUp(self):
        """Unit test setup method (automatically called).
        """
        # This configuration assumes that Minio is running on localhost with
        # port 9000, and can be accessed using 'admin' and 'password' for
        # authentication
        self.persistent_minio_config = {
            'Host': 'localhost',
            'Port': '9000',
            'AccessKey': 'admin',
            'SecretKey': 'password',
            'RetentionTime': 10,  # 10 Seconds 
            'Ssl': 'false' 
        }

    def test_persistent_minio(self):
        """Unit test to test the persistent image store using Minio.
        """
        # Initialize persistent image store
        pims = PersistentImageStore('minio', self.persistent_minio_config)

        # Create random 512 byte blob
        data = b'\x00' * 512

        # Store the data blob
        key = pims.store(data)

        # Read the value from the store
        rdata = pims.read(key)

        # Verify that the data is not None, and that the data is the same as
        # was inserted into the image store
        self.assertIsNotNone(rdata)
        self.assertEqual(data, rdata)
        
        # Remove the value from the store
        pims.remove(key)
        
        # Try getting the key again (should be None)
        rdata = pims.read(key)
        self.assertIsNone(rdata)

    def test_persistent_minio_retention(self):
        """Unit test to test that the retention policy implementation is
        working correctly.
        """
        config = dict(self.persistent_minio_config)
        config['RetentionTime'] = 5  # Change to 5 seconds

        # Initialize persistent image store
        pims = PersistentImageStore('minio', config)

        # Create random 512 byte blob
        data = b'\x00' * 512

        # Store the data blob
        key = pims.store(data)

        # Sleep to let the image store to remove object from the store
        # Sleeping a little extra time (10 seconds) to allow time for the 
        # removal to happen in the image store
        time.sleep(15)

        # Try getting the key again (should be None)
        rdata = pims.read(key)
        self.assertIsNone(rdata)


if __name__ == '__main__':
    unittest.main(verbosity=2)


# ImageStore
ImageStore Module provides APIs for image read, store and remove in both in-memory(redis) and persistence storage(minio).

## How to Test from present working directory:
A test script is available under ImageStore/test/py/

```
python3.6 clientTest.py --client-cert cert-tool/Certificates/imagestore/imagestore_client_certificate.pem --client-key cert-tool/Certificates/imagestore/imagestore_client_key.pem --ca-cert cert-tool/Certificates/ca/ca_certificate.pem
```
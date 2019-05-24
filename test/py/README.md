# ImageStore
ImageStore Module provides APIs for image read, store and remove in both in-memory(redis) and persistence storage(minio).

## Pre-requisites (present working directory - ImageStore/test/py/)

* `ia_imagestore` container should be running
* Install imagestore_requirements.txt (for python) by running cmd: `sudo -H pip3.6 install -r imagestore_requirements.txt`
* Set `PYTHONPATH` env variable
    ```sh
        export PYTHONPATH=../../../:../../../ImageStore/protobuff/py:../../../ImageStore/client/py
    ```
    > **Note**: `../../../` refers to parent directory of ImageStore. This should be adjusted accordingly based on where the `clientTest.py` resides w.r.t `ImageStore` folder
* Set `no_proxy` env variable
    ```sh
    export no_proxy=$no_proxy,<IEI_node_ip_address>
    ```
* Copying certs and keys:
    * Copy ImageStore client cert and key to /etc/ssl/imagestore
    * Copy CA cert to /etc/ssl/ca

    > **Note**: If one wish to provide a diff cert/key path, they can do so by providing the right cert/key path while running `clientTest.py` script below


## How to Test from present working directory (ImageStore/test/py/)

```sh
    python3.6 clientTest.py --hostname <IEI_node_ip_address> \
                            --port 50055 \
                            --client-cert /etc/ssl/imagestore/imagestore_client_certificate.pem \
                            --client-key /etc/ssl/imagestore/imagestore_client_key.pem \
                            --ca-cert /etc/ssl/ca/ca_certificate.pem \
                            --input_file <Path to input file> \
                            --output_file <Path to output file>

```

> **Note**: Right now, the IEI ImageStore gRPC server runs at port `50055`. If one wish to change this, ensure that IEI's ImageStore gRPC server
> also is listening at that port.

## To view the documentation navigate to the below link

```sh
/opt/intel/iei/dist_libs/ImageStore/client/py/doc/html/files.html
```

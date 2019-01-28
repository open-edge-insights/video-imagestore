# ImageStore
ImageStore Module provides C++ APIs for image read and remove in both in-memory(redis) and persistence storage(minio). We are facing some issues in getting this working

## Pre-requisites

* `ia_imagestore` container should be running
* Set `no_proxy` env variable
    ```sh
    export no_proxy=$no_proxy,<ETA_node_ip_address>
    ```
* Copying certs and keys:
    * Copy ImageStore client cert and key to /etc/ssl/imagestore
    * Copy CA cert to /etc/ssl/ca

    > **Note**: If one wish to provide a diff cert/key path, they can do so by providing the right cert/key path in the
    > `ImageStore/test/cpp/Makefile`. Also, one can configure the host, port, img handle and output file in the same Makefile
* Change the `HOST` value from `localhost` to ip address of the node running ETA/imagestore in the `ImageStore/test/cpp/Makefile`.


## How to Test from present working directory (ImageStore/test/cpp/)

```sh
    sudo make
```

> **Note**:
> 1. Right now, the ETA ImageStore gRPC server runs at port `50055`. If one wish to change this, ensure that ETA's ImageStore gRPC server
>    also is listening at that port.
> 2. Run command `docker exec -it ia_imagestore bash` to get inside ia_imagestore container ans fetch the redis `inmem` image handle stored in
> redis by following below steps:
>
>    ```sh
>    I have no name!@ia_imagestore:/ETA/go/src/IEdgeInsights$ ./redis-5.0.2/src/redis-cli
>    127.0.0.1:6379> AUTH redis123
>    OK
>    127.0.0.1:6379> KEYS *
>   ```
> 3. Go to `/opt/intel/eta/data/image-store-bucket` path to fetch the blob for `perist` image handle stored in minio.
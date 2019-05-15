# ImageStore
ImageStore Module provides C++ APIs for image read and remove in both in-memory(redis) and persistence storage(minio).

## Pre-requisites

* `ia_imagestore` container should be running
* **Setting up C++ dev env**
    * Run the [setup_ubuntu_dev_env_cpp.sh](setup_ubuntu_dev_env_cpp.sh) script file after copying it to the directory where you
      need grpc to be installed. Give necessary permissions required.
    * In case of any issues running the above script file, use the following guide
        (https://github.com/grpc/grpc/blob/master/BUILDING.md)
    * To verify successfull installation, try running gRPC C++ HelloWorld example:
        * cd grpc/examples/cpp/helloworld
        * make
        * ./greeter_server
        * ./greeter_client (In a separate terminal)
          Terminal should display `Greeter received: Hello world` on correct installation.

* Set `no_proxy` env variable
    ```sh
    export no_proxy=$no_proxy,<IEI_node_ip_address>
    ```

* Set `INPUT_FILE` & `OUTPUT_FILE` variables in Makefile with absolute paths of input & output images.

* Copying certs and keys:
    * Copy ImageStore client cert and key to /etc/ssl/imagestore
    * Copy CA cert to /etc/ssl/ca

    > **Note**: If one wish to provide a diff cert/key path, they can do so by providing the right cert/key path in the
    > `ImageStore/test/cpp/Makefile`. Also, one can configure the host, port, img handle and output file in the same Makefile

* Change the `HOST` value from `localhost` to ip address of the node running IEI/imagestore in the `ImageStore/test/cpp/Makefile`.


## How to Test from present working directory (ImageStore/test/cpp/)

```sh
    sudo make
```

> **Note**:
> 1. Right now, the IEI ImageStore gRPC server runs at port `50055`. If one wish to change this, ensure that IEI's ImageStore gRPC server
>    also is listening at that port.
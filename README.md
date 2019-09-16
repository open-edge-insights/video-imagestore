# `ImageStore Module`

1. ImageStore will subscribe to the Video Analytics result. So ImageStore
   receives the frame data and stores it in minio.
   The payload format expected by image store is as follows
   map ("img_handle":"$handle_name"),[]byte($binaryImage)

2. ImageStore starts the server which provides the read and store interfaces.
   The payload format is as follows

   Request : map ("command": "store" , "img_handle":"$handle_name"),[]byte($binaryImage)
   Response : map ("img_handle":"$handle_name", "error":"$error_msg")
   ("error" is optional and available only in case of error in execution.)

   Request : map ("command": "read" , "img_handle":"$handle_name")
   Response : map ("img_handle":"$handle_name", "error":"$error_msg"),[]byte($binaryImage) 
   ("error" is optional and available only in case of error in execution.
   And $binaryImage is available only in case of successfull read)
   
## `Configuration`

All the ImageStore module configuration are added into etcd (distributed
key-value data store) under `AppName` as mentioned in the
environment section of this app's service definition in docker-compose.

If `AppName` is `ImageStore`, then the app's config would look like as below
 for `/ImageStore/config` key in Etcd:
 ```
    "/ImageStore/config": {
        "minio":{  
           "accessKey":"admin",
           "secretKey":"password",
           "retentionTime":"1h",
           "retentionPollInterval":"60s",
           "ssl":"false"
        }
    }
 ```
For more details on Etcd and MessageBus endpoint configuration, visit [Etcd_and_MsgBus_Endpoint_Configuration](../Etcd_and_MsgBus_Endpoint_Configuration.md).

## `Installation`

* Follow [provision/README.md](../README#provision-eis.md) for EIS provisioning
  if not done already as part of EIS stack setup

* Run ImageStore

  Present working directory to try out below commands is: `[repo]/ImageStore`

    1. Build and Run VideoAnalytics as container
        ```
        $ cd [repo]/docker_setup
        $ docker-compose up --build ImageStore
        ```

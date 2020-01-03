# `ImageStore Module`

The Image Store component of EIS comes as a separate container which primarily
subscribes to the stream that comes out of the VideoAnalytics app via EIS
MessageBus and stores the frame into minio for historical analysis.

The high level logical flow of ImageStore is as below:
1. The messagebus subscriber in ImageStore will subscribe to the VideoAnalytics
   published classified result (metadata, frame) on the messagebus.
   The img_handle is extracted out of the metadata and is used as the key and
   the frame is stored as a value for that key in minio persistent storage.
2. For historical analysis of the stored classified images, ImageStore starts
   the messagebus server which provides the read and store interfaces.
   The payload format is as follows for:
   * Store interface:
     ```
        Request: map ("command": "store","img_handle":"$handle_name"),[]byte($binaryImage)
        Response : map ("img_handle":"$handle_name", "error":"$error_msg") ("error" is optional and available only in case of error in execution.)
     ```
   * Read interface:
     ```
        Request : map ("command": "read", "img_handle":"$handle_name")
        Response : map ("img_handle":"$handle_name", "error":"$error_msg"),[]byte($binaryImage) ("error" is optional and available only in case of error in execution. And $binaryImage is available only in case of successful read)
     ```

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

### `Detailed description on each of the keys used`
|  Key	        | Description 	                                                                                           | Possible Values  	                      |Required/Optional |
|---	        |---	                                                                                                   |---	                                      |---	             |
|  accessKey 	|   Username required to access Minio DB	                                                               | Any suitable value                       | Required	     |
|  secretKey 	|   Password required to access Minio DB	                                                               | Any suitable value             	      | Required         |
|  retentionTime|   The retention parameter specifies the retention policy to apply for the images stored in Minio DB. | Suitable duration string value as mentioned at https://golang.org/pkg/time/#ParseDuration |   Required        |
|  retentionPollInterval | Used to set the time interval for checking images for expiration. Expired images will become candidates for deletion and no longer retained |	Suitable duration string value as mentioned at https://golang.org/pkg/time/#ParseDuration  |   Required        |
|  ssl          |  If "true", establishes a secure connection with Minio DB else a non-secure connection                   | "true" or "false"                        |   Required        |

For more details on Etcd and MessageBus endpoint configuration, visit [Etcd_Secrets_and_MsgBus_Endpoint_Configuration](../Etcd_Secrets_and_MsgBus_Endpoint_Configuration.md).

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

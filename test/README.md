# Testing ImageStore service
go run ./ImageStoreTestService.go --configFile ./service.json --serviceName ImageStore

# Testing ImageStore subscriber
go run ./ImageStoreSubTest.go -pubConfigFile ./publisher.json -pubTopic camera1_stream_results -servConfigFile ./service.json -serviceName ImageStore

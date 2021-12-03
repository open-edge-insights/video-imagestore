module imagestore

go 1.15

replace imagestore => ./

replace github.com/open-edge-insights/eii-configmgr-go => ../../ConfigMgr/

replace github.com/open-edge-insights/eii-messagebus-go => ../../EIIMessageBus/

require (
        github.com/go-ini/ini 6ed8d5f64cd79a498d1f3fab5880cc376ce41bbe
        github.com/minio/minio-go v6.0.10
        github.com/mitchellh/go-homedir ae18d6b8b3205b561c79e8e5f69bff09736185f4
        github.com/xeipuuv/gojsonschema v1.2.0
        golang.org/x/crypto ff983b9c42bc9fbf91556e191cc8efb585c16908
        golang.org/x/net 26e67e76b6c3f6ce91f7c52def5af501b4e0f3a2
        github.com/golang/text v0.3.0
        github.com/golang/sys d0be0721c37eeb5299f245a996a483160fc36940
)


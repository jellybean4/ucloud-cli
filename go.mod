module github.com/ucloud/ucloud-cli

go 1.12

require (
	github.com/jellybean4/ucloud-sdk-go v0.12.7
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/ucloud/ucloud-sdk-go v0.11.1
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace (
	github.com/spf13/cobra v0.0.3 => github.com/lixiaojun629/cobra v0.0.9
	github.com/spf13/pflag v1.0.3 => github.com/lixiaojun629/pflag v1.0.5
)

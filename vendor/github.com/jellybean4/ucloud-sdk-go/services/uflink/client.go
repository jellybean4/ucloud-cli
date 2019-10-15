package uflink

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

// UDPNClient is the client of UDPN
type UFlinkClient struct {
	*ucloud.Client
}

// NewClient will return a instance of UDPNClient
func NewClient(config *ucloud.Config, credential *auth.Credential) *UFlinkClient {
	meta := ucloud.ClientMeta{Product: "UFlink"}
	client := ucloud.NewClientWithMeta(config, credential, meta)
	return &UFlinkClient{
		client,
	}
}

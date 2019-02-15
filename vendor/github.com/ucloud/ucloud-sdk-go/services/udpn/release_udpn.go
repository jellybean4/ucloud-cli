//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDPN ReleaseUDPN

package udpn

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// ReleaseUDPNRequest is request schema for ReleaseUDPN action
type ReleaseUDPNRequest struct {
	request.CommonBase

	// [公共参数] 地域。 参见 [地域和可用区列表](../summary/regionlist.html)
	// Region *string `required:"false"`

	// [公共参数] 项目ID。不填写为默认项目，子帐号必须填写。 请参考[GetProjectList接口](../summary/get_project_list.html)
	// ProjectId *string `required:"false"`

	// UDPN 资源 Id
	UDPNId *string `required:"true"`
}

// ReleaseUDPNResponse is response schema for ReleaseUDPN action
type ReleaseUDPNResponse struct {
	response.CommonBase
}

// NewReleaseUDPNRequest will create request of ReleaseUDPN action.
func (c *UDPNClient) NewReleaseUDPNRequest() *ReleaseUDPNRequest {
	req := &ReleaseUDPNRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// ReleaseUDPN - 释放 UDPN
func (c *UDPNClient) ReleaseUDPN(req *ReleaseUDPNRequest) (*ReleaseUDPNResponse, error) {
	var err error
	var res ReleaseUDPNResponse

	err = c.client.InvokeAction("ReleaseUDPN", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

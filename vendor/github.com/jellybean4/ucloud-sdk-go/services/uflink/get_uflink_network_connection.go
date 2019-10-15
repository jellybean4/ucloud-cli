package uflink

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

type GetUFlinkNetworkConnectionRequest struct {
	request.CommonBase

	InstanceId *string `required:"true"`

	Ip *string `required:"true"`

	Port *int `required:"true"`
}

type GetUFlinkNetworkConnectionResponse struct {
	response.CommonBase

	Data NetworkConnectionResult
}

func (c *UFlinkClient) NewGetUFlinkNetworkConnectionRequest() *GetUFlinkNetworkConnectionRequest {
	req := &GetUFlinkNetworkConnectionRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(true)
	return req
}

func (c *UFlinkClient) GetUFlinkNetworkConnection(req* GetUFlinkNetworkConnectionRequest) (*GetUFlinkNetworkConnectionResponse, error) {
	var err error
	var res GetUFlinkNetworkConnectionResponse

	err = c.Client.InvokeAction("GetUFlinkNetworkConnection", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

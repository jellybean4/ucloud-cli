package uflink

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

type DescribeUFlinkApplicationResponse struct {
	response.CommonBase

	Data ApplicationDetailData
}

func (c *UFlinkClient) DescribeUFlinkApplication(req *ApplicationRequest) (*DescribeUFlinkApplicationResponse, error) {
	var err error
	var res DescribeUFlinkApplicationResponse

	err = c.Client.InvokeAction("DescribeUFlinkApplication", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

package uflink

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

type ListUFlinkSubmittedJobsRequest struct {
	request.CommonBase

	// 集群资源ID
	InstanceId *string `required:"true"`
}

type ListUFlinkSubmittedJobsResponse struct {
	response.CommonBase

	TotalCount int

	Data []SubmittedJob
}

func (c *UFlinkClient) NewListUFlinkSubmittedJobsRequest() *ListUFlinkSubmittedJobsRequest {
	req := &ListUFlinkSubmittedJobsRequest{}
	c.Client.SetupRequest(req)
	req.SetRetryable(true)
	return req
}

func (c *UFlinkClient) ListUFlinkSubmittedJobs(req *ListUFlinkSubmittedJobsRequest) (*ListUFlinkSubmittedJobsResponse, error) {
	var err error
	var res ListUFlinkSubmittedJobsResponse

	err = c.Client.InvokeAction("ListUFlinkSubmittedJob", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

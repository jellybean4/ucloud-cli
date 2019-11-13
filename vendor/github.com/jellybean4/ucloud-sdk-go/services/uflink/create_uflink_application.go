package uflink

import "github.com/ucloud/ucloud-sdk-go/ucloud/request"

type CreateUFlinkApplicationRequest struct {
	request.CommonBase

	InstanceId *string `required:"true"`

	Name *string `required:"true"`

	FlinkVersion *string `required:"true"`

	JobType *string `required:"true"`

	AppFile *string `required:"false"`

	AppParams *string `required:"false"`

	SubmitParams *string `required:"false"`

	RuntimeParams *string `required:"false"`

	ApplicationId *string `required:"false"`

	SQL *string `required:"false"`

	UDXFJars *string `required:"false"`

	Resources *string `required:"false"`
}

func (c *UFlinkClient) NewCreateUFlinkApplicationRequest() *CreateUFlinkApplicationRequest {
	req := &CreateUFlinkApplicationRequest{}
	c.Client.SetupRequest(req)
	req.SetRetryable(false)
	return req
}

func (c *UFlinkClient) CreateUFlinkApplication(req *CreateUFlinkApplicationRequest) (*ApplicationResponse3, error) {
	var err error
	var res ApplicationResponse3

	err = c.Client.InvokeAction("CreateUFlinkApplication", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

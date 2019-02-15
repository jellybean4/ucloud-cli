//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UAccount ModifyProject

package uaccount

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// ModifyProjectRequest is request schema for ModifyProject action
type ModifyProjectRequest struct {
	request.CommonBase

	// [公共参数] 项目ID。不填写为默认项目，子帐号必须填写。 请参考[GetProjectList接口](../summary/get_project_list.html)
	// ProjectId *string `required:"true"`

	// 新的项目名称
	ProjectName *string `required:"true"`
}

// ModifyProjectResponse is response schema for ModifyProject action
type ModifyProjectResponse struct {
	response.CommonBase
}

// NewModifyProjectRequest will create request of ModifyProject action.
func (c *UAccountClient) NewModifyProjectRequest() *ModifyProjectRequest {
	req := &ModifyProjectRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// ModifyProject - 修改项目
func (c *UAccountClient) ModifyProject(req *ModifyProjectRequest) (*ModifyProjectResponse, error) {
	var err error
	var res ModifyProjectResponse

	err = c.client.InvokeAction("ModifyProject", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

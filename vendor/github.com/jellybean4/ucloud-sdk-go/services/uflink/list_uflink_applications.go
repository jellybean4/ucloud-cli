package uflink

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// List all uflink applications reserved
type ListUFlinkApplicationsRequest struct {
	request.CommonBase

	// [公共参数] 地域。 参见 [地域和可用区列表](../summary/regionlist.html)
	// Region *string `required:"true"`

	// [公共参数] 可用区。参见 [可用区列表](../summary/regionlist.html)
	// Zone *string `required:"false"`

	// [公共参数] 项目ID。不填写为默认项目，子帐号必须填写。 请参考[GetProjectList接口](../summary/get_project_list.html)
	// ProjectId *string `required:"false"`

	// 集群资源ID
	InstanceId *string `required:"true"`

	// 分页起始条目数，默认为0
	Offset *int `required:"false"`

	// 每页显示数据条目上限，默认为0
	Limit *int `required:"false"`
}

type ListUFlinkApplicationsResponse struct {
	response.CommonBase

	// 应用总数
	TotalCount int

	// 应用详细数据
	Data []ApplicationData
}

func (c *UFlinkClient) NewListUFlinkApplicationsRequest() *ListUFlinkApplicationsRequest {
	req := &ListUFlinkApplicationsRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(true)
	return req
}

func (c *UFlinkClient) ListUFlinkApplications(req *ListUFlinkApplicationsRequest) (*ListUFlinkApplicationsResponse, error) {
	var err error
	var res ListUFlinkApplicationsResponse

	err = c.Client.InvokeAction("ListUFlinkApplications", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

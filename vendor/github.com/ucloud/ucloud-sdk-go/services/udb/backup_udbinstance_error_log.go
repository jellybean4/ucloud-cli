//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDB BackupUDBInstanceErrorLog

package udb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// BackupUDBInstanceErrorLogRequest is request schema for BackupUDBInstanceErrorLog action
type BackupUDBInstanceErrorLogRequest struct {
	request.CommonBase

	// [公共参数] 地域。 参见 [地域和可用区列表](../summary/regionlist.html)
	// Region *string `required:"true"`

	// [公共参数] 可用区。参见 [可用区列表](../summary/regionlist.html)
	// Zone *string `required:"false"`

	// [公共参数] 项目ID。不填写为默认项目，子帐号必须填写。 请参考[GetProjectList接口](../summary/get_project_list.html)
	// ProjectId *string `required:"false"`

	// DB实例Id,该值可以通过DescribeUDBInstance获取
	DBId *string `required:"true"`

	// 备份名称
	BackupName *string `required:"true"`
}

// BackupUDBInstanceErrorLogResponse is response schema for BackupUDBInstanceErrorLog action
type BackupUDBInstanceErrorLogResponse struct {
	response.CommonBase
}

// NewBackupUDBInstanceErrorLogRequest will create request of BackupUDBInstanceErrorLog action.
func (c *UDBClient) NewBackupUDBInstanceErrorLogRequest() *BackupUDBInstanceErrorLogRequest {
	req := &BackupUDBInstanceErrorLogRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// BackupUDBInstanceErrorLog - 备份UDB指定时间段的errorlog
func (c *UDBClient) BackupUDBInstanceErrorLog(req *BackupUDBInstanceErrorLogRequest) (*BackupUDBInstanceErrorLogResponse, error) {
	var err error
	var res BackupUDBInstanceErrorLogResponse

	err = c.client.InvokeAction("BackupUDBInstanceErrorLog", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

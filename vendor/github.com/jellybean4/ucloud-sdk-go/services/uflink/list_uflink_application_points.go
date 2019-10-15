package uflink

import "github.com/ucloud/ucloud-sdk-go/ucloud/response"

type ListUFlinkApplicationPointsResponse struct {
	response.CommonBase

	Data []HDFSFileStatus
}

func (c *UFlinkClient) ListUFlinkApplicationPoints(req *ApplicationRequest) (*ListUFlinkApplicationPointsResponse, error) {
	var err error
	var res ListUFlinkApplicationPointsResponse

	err = c.Client.InvokeAction("ListUFlinkApplicationPoints", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

package uflink

func (c *UFlinkClient) RunUFlinkApplicationSavepoint(req *ApplicationRequest) (*ApplicationResponse, error) {
	var err error
	var res ApplicationResponse

	err = c.Client.InvokeAction("RunUFlinkApplicationSavepoint", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

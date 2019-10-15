package uflink


func (c *UFlinkClient) NewApplicationRequest() *ApplicationRequest {
	req := &ApplicationRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(true)
	return req
}

func (c *UFlinkClient) KillUFlinkApplication(req *ApplicationRequest) (*ApplicationResponse, error) {
	var err error
	var res ApplicationResponse

	err = c.Client.InvokeAction("KillUFlinkApplication", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

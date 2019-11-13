package uflink

func (c *UFlinkClient) SendUFlinkResubmitJob(req *ApplicationRequest) (*ApplicationResponse3, error) {
	var err error
	var res ApplicationResponse3

	err = c.Client.InvokeAction("SendUFlinkResubmitJob", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

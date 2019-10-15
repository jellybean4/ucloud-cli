package uflink

func (c *UFlinkClient) SendUFlinkResubmitJob(req *ApplicationRequest) (*ApplicationResponse2, error) {
	var err error
	var res ApplicationResponse2

	err = c.Client.InvokeAction("SendUFlinkResubmitJob", req, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

package uflink

import "github.com/ucloud/ucloud-sdk-go/ucloud/response"

type ApplicationResponse struct {
	response.CommonBase

	Data bool
}

type ApplicationResponse2 struct {
	response.CommonBase

	Data string
}

type ApplicationResponse3 struct {
	response.CommonBase

	Data int64
}
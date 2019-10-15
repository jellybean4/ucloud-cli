package uflink

/*
UFlinkApplicationData - ListUFlinkApplications
*/

type ApplicationData struct {
	// 应用已经运行的时长
	ElapsedTime int

	// 应用的Yarn ApplicationId
	Id string

	// 应用的名称，全局唯一
	Name string

	// 应用的当前状态
	State string

	// 应用的启动时间
	StartedTime int

	// 应用的结束时间
	FinishedTime int

	// 应用分配的内存大小
	AllocatedMB float32

	// 应用分配的CPU数量
	AllocatedVCores int

	// Flink版本信息
	FlinkVersion string

	// 应用的原始名称
	OriginalName string

	// 应用提交的各种参数
	Request string
}

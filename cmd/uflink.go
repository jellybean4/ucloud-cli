package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jellybean4/ucloud-sdk-go/services/uflink"
	"github.com/spf13/cobra"
	"github.com/ucloud/ucloud-cli/base"
	"github.com/ucloud/ucloud-cli/ux"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

type UFlinkApplicationRow struct {
	ApplicationName string
	OriginalName    string
	Id              string
	State           string
	ElapsedTime     int
	StartedTime     int
	FinishedTime    int
	AllocatedMB     float32
	AllocatedVCores int
	FlinkVersion    string
	Request         string
}

type UFlinkPointFileRow struct {
	Path             string
	ModificationTime int
	Length           int
	BlockSize        int
}

type FlinkRunParams struct {
	Class                 string `json:"--class,omitempty"`
	Parallelism           int    `json:"--parallelism,omitempty"`
	FromSavepoint         string `json:"--fromSavepoint,omitempty"`
	YarnJobManagerMemory  string `json:"--yarnjobManagerMemory,omitempty"`
	YarnContainer         int    `json:"--yarncontainer,omitempty"`
	YarnSlots             int    `json:"--yarnslots,omitempty"`
	YarnTaskManagerMemory string `json:"--yarntaskManagerMemory,omitempty"`
}

func NewCmdUFlink() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uflink",
		Short: "list,create,kill,describe,resubmit uflink application",
		Long:  `list,create,kill,describe,resubmit uflink application`,
		Args:  cobra.NoArgs,
	}
	out := base.Cxt.GetWriter()
	cmd.AddCommand(NewCmdUFlinkApplicationList(out))
	cmd.AddCommand(NewCmdUFlinkApplicationCreate(out))
	cmd.AddCommand(NewCmdUFlinkSubmittedJobList(out))
	cmd.AddCommand(NewCmdUFlinkSubmittedJobLogGet(out))

	cmd.AddCommand(NewCmdUFlinkApplication(out, killUFlinkApplicationCmdBuilder))
	cmd.AddCommand(NewCmdUFlinkApplication(out, describeUFlinkApplicationCmdBuilder))
	cmd.AddCommand(NewCmdUFlinkApplication(out, resubmitUFlinkApplicationCmdBuilder))
	cmd.AddCommand(NewCmdUFlinkApplication(out, listUFlinkApplicationPointsCmdBuilder))
	cmd.AddCommand(NewCmdUFlinkApplication(out, runUFlinkApplicationSavepointCmdBuilder))
	return cmd
}

func NewCmdUFlinkApplicationCreate(out io.Writer) *cobra.Command {
	var submitParams FlinkRunParams
	var flinkConf string

	req := base.BizClient.NewCreateUFlinkApplicationRequest()
	cmd := &cobra.Command{
		Use:     "create",
		Short:   `Create uflink application`,
		Long:    `Create uflink application`,
		Example: `run [OPTIONS] <jar-file-path-in-ufile> <arguments>`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(len(cmd.Flags().Args()))
			submitParamsContent, err := json.Marshal(submitParams)
			appParamsContent := ""

			if err != nil {
				fmt.Fprintf(out, "Failed to read encode submit params for %s\n", base.ParseError(err))
				return
			}
			if len(cmd.Flags().Args()) < 1 {
				fmt.Fprint(out, "UFlink application ufile path not specified\n")
				return
			}
			appParams := cmd.Flags().Args()[1:]
			if len(appParams) > 0 {
				appParamsContent =strings.Join(appParams, " ")
			}

			encodeAppParams := Encode(appParamsContent)
			encodeSubmitParams := Encode(string(submitParamsContent))
			req.AppParams = &encodeAppParams
			req.SubmitParams = &encodeSubmitParams

			if flinkConf != "" {
				content, err := ioutil.ReadFile(flinkConf)
				if err != nil {
					base.LogError(fmt.Sprintf("Failed to read runtime configs from file[%s] for %s", flinkConf, base.ParseError(err)))
					return
				}
				configs := string(content[:])
				req.RuntimeParams = &configs
			} else {
				configs := ""
				req.RuntimeParams = &configs
			}

			_, err = base.BizClient.CreateUFlinkApplication(req)
			if err != nil {
				base.LogError(fmt.Sprintf("Failed to create new uflink application: %s", base.ParseError(err)))
				return
			}
			fmt.Fprintf(out, "Create application[%s] success\n", *req.Name)
		},
	}
	cmd.Flags().SortFlags = false
	req.ProjectId = cmd.Flags().String("project-id", base.ConfigIns.ProjectID, "Optional. Assign project-id")
	req.Region = cmd.Flags().String("region", base.ConfigIns.Region, "Optional. Assign region.")
	req.Zone = cmd.Flags().String("zone", base.ConfigIns.Zone, "Optional. Assign availability zone")
	req.InstanceId = cmd.Flags().String("instance-id", "", "Required. Resource ID of uflink instance")
	cmd.Flags().SetFlagValuesFunc("project-id", getProjectList)
	cmd.Flags().SetFlagValuesFunc("region", getRegionList)
	cmd.Flags().SetFlagValuesFunc("zone", func() []string {
		return getZoneList(req.GetRegion())
	})

	req.Name = cmd.Flags().String("name", "", "Required. Name of the new uflink application")
	req.FlinkVersion = cmd.Flags().String("flink-version", "1.6.4", "Required. Version of the new uflink application")
	//req.AppFile = cmd.Flags().String("app-file", "", "Optional. UFile path of the uflink application jar file")
	req.SQL = cmd.Flags().String("sql", "", "Optional. Flink sql used to create uflink application")
	req.Resources = cmd.Flags().String("resources", "", "Optional. Comma separated list of ufile resource files")

	// flink run params
	cmd.Flags().StringVar(&submitParams.Class, "classname", "", `Optional. Class with the program entry point ("main" method or "getPlan()" method. Only needed if the JAR file does not specify the class in its manifest.`)
	cmd.Flags().IntVar(&submitParams.Parallelism, "parallelism", 1, `Optional. The parallelism with which to run the program. Optional flag to override the default value specified in the configuration.`)
	cmd.Flags().StringVar(&submitParams.FromSavepoint, "fromSavepoint", "", `Optional. Path to a savepoint to restore the job from (for example hdfs:///flink/savepoint-1537)`)
	cmd.Flags().StringVar(&submitParams.YarnJobManagerMemory, "yarnjobManagerMemory", "", `Optional. Memory for JobManager Container with optional unit(default: MB)`)
	cmd.Flags().StringVar(&submitParams.YarnTaskManagerMemory, "yarntaskManagerMemory", "", `Optional. Memory per TaskManager Container with optional unit (default: MB)`)
	cmd.Flags().IntVar(&submitParams.YarnContainer, "yarncontainer", 1, `Optional. Number of YARN container to allocate (=Number of Task Managers)`)
	cmd.Flags().IntVar(&submitParams.YarnSlots, "yarnslots", 1, `Optional. Number of slots per TaskManager`)

	//cmd.Flags().StringVar(&appParams, "app-params", "", "Optional. Main method args of the new uflink application")
	cmd.Flags().StringVar(&flinkConf, "flink-conf", "", "Optional. File path of flink-conf.yaml override default flink conf")

	cmd.Flags().SetFlagValues("flink-version", "1.6.4_2.11", "1.7.2_2.11", "1.7.2_2.12", "1.8.0_2.11", "1.8.0_2.12")
	cmd.Flags().SetInterspersed(false)

	cmd.MarkFlagRequired("instance-id")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("flink-version")
	return cmd
}

func NewCmdUFlinkApplicationList(out io.Writer) *cobra.Command {
	var output string
	offset := 0
	limit := 100

	req := base.BizClient.NewListUFlinkApplicationsRequest()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all reserved uflink applications",
		Long:  `List all reserved uflink applications`,
		Run: func(cmd *cobra.Command, args []string) {
			applications, err := listAllApplications(req)
			if err != nil {
				base.HandleError(err)
				return
			}
			showApplications(applications, out, output)
		},
	}
	cmd.Flags().SortFlags = false
	req.ProjectId = cmd.Flags().String("project-id", base.ConfigIns.ProjectID, "Optional. Assign project-id")
	req.Region = cmd.Flags().String("region", base.ConfigIns.Region, "Optional. Assign region.")
	req.Zone = cmd.Flags().String("zone", base.ConfigIns.Zone, "Optional. Assign availability zone")
	req.InstanceId = cmd.Flags().String("instance-id", "", "Required. Resource ID of uflink instance")

	cmd.Flags().StringVarP(&output, "output", "o", "", "Optional. Accept values: wide. Display more information about uflink applications")

	req.Offset = &offset
	req.Limit = &limit

	cmd.Flags().SetFlagValues("output", "wide")
	cmd.Flags().SetFlagValuesFunc("project-id", getProjectList)
	cmd.Flags().SetFlagValuesFunc("region", getRegionList)
	cmd.Flags().SetFlagValuesFunc("zone", func() []string {
		return getZoneList(req.GetRegion())
	})

	cmd.MarkFlagRequired("instance-id")
	return cmd
}

func NewCmdUFlinkApplication(out io.Writer, builder func(req *uflink.ApplicationRequest, out io.Writer) *cobra.Command) *cobra.Command {
	req := base.BizClient.NewApplicationRequest()
	req.SetRetryable(true)
	cmd := builder(req, out)
	cmd.Flags().SortFlags = false
	req.ProjectId = cmd.Flags().String("project-id", base.ConfigIns.ProjectID, "Optional. Assign project-id")
	req.Region = cmd.Flags().String("region", base.ConfigIns.Region, "Optional. Assign region.")
	req.Zone = cmd.Flags().String("zone", base.ConfigIns.Zone, "Optional. Assign availability zone")
	req.InstanceId = cmd.Flags().String("instance-id", "", "Required. Resource ID of uflink instance")
	req.ApplicationId = cmd.Flags().String("application-id", "", "Required. Application Id of the specific job")

	cmd.Flags().SetFlagValuesFunc("project-id", getProjectList)
	cmd.Flags().SetFlagValuesFunc("region", getRegionList)
	cmd.Flags().SetFlagValuesFunc("zone", func() []string {
		return getZoneList(req.GetRegion())
	})
	cmd.MarkFlagRequired("instance-id")
	cmd.MarkFlagRequired("application-id")
	return cmd
}

func NewCmdUFlinkSubmittedJobList(out io.Writer) *cobra.Command {
	req := base.BizClient.NewListUFlinkSubmittedJobsRequest()
	cmd := &cobra.Command{
		Use:   "submitted",
		Short: "list all submitted history",
		Long:  `list all submitted history`,
		Run: func(cmd *cobra.Command, args []string) {
			info, err := base.BizClient.ListUFlinkSubmittedJobs(req)
			if err != nil {
				base.HandleError(err)
				return
			}
			showSubmittedJobs(info.Data, out)
		},
	}
	cmd.Flags().SortFlags = false
	req.ProjectId = cmd.Flags().String("project-id", base.ConfigIns.ProjectID, "Optional. Assign project-id")
	req.Region = cmd.Flags().String("region", base.ConfigIns.Region, "Optional. Assign region.")
	req.Zone = cmd.Flags().String("zone", base.ConfigIns.Zone, "Optional. Assign availability zone")
	req.InstanceId = cmd.Flags().String("instance-id", "", "Required. Resource ID of uflink instance")

	cmd.Flags().SetFlagValuesFunc("project-id", getProjectList)
	cmd.Flags().SetFlagValuesFunc("region", getRegionList)
	cmd.Flags().SetFlagValuesFunc("zone", func() []string {
		return getZoneList(req.GetRegion())
	})

	cmd.MarkFlagRequired("instance-id")
	return cmd
}

func NewCmdUFlinkSubmittedJobLogGet(out io.Writer) *cobra.Command {
	req := base.BizClient.NewGetUFlinkSubmittedJobLogRequest()
	cmd := &cobra.Command{
		Use:   "submitted-log",
		Short: "show the specific submitted job log",
		Long:  `show the specific submitted job log`,
		Run: func(cmd *cobra.Command, args []string) {
			info, err := base.BizClient.GetUFlinkSubmittedJobLog(req)
			if err != nil {
				base.HandleError(err)
				return
			}
			fmt.Fprint(out, info.Data, "\n")
		},
	}
	cmd.Flags().SortFlags = false
	req.ProjectId = cmd.Flags().String("project-id", base.ConfigIns.ProjectID, "Optional. Assign project-id")
	req.Region = cmd.Flags().String("region", base.ConfigIns.Region, "Optional. Assign region.")
	req.Zone = cmd.Flags().String("zone", base.ConfigIns.Zone, "Optional. Assign availability zone")
	req.InstanceId = cmd.Flags().String("instance-id", "", "Required. Resource ID of uflink instance")
	req.Id = cmd.Flags().Int("job-id", 0, "Required. Submitted Id of the job")

	cmd.Flags().SetFlagValuesFunc("project-id", getProjectList)
	cmd.Flags().SetFlagValuesFunc("region", getRegionList)
	cmd.Flags().SetFlagValuesFunc("zone", func() []string {
		return getZoneList(req.GetRegion())
	})

	cmd.MarkFlagRequired("instance-id")
	cmd.MarkFlagRequired("job-id")
	return cmd
}

func resubmitUFlinkApplicationCmdBuilder(req *uflink.ApplicationRequest, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "resubmit",
		Short: "Resubmit uflink application",
		Long:  "Resubmit uflink application",
		Run: func(cmd *cobra.Command, args []string) {
			var logs []string
			req.SetRetryable(false)
			res, logs := resubmitUFlinkApplication(req)
			logs = append(logs, fmt.Sprintf("Resubmit application[%s] %t", *req.ApplicationId, res))
			base.LogInfo(logs...)
			fmt.Fprintf(out, "Resubmit application[%s] %t\n", *req.ApplicationId, res)
		},
	}
}

func killUFlinkApplicationCmdBuilder(req *uflink.ApplicationRequest, out io.Writer) *cobra.Command {
	var yes *bool
	cmd := &cobra.Command{
		Use:   "kill",
		Short: "Kill uflink application",
		Long:  "Kill uflink application",
		Run: func(cmd *cobra.Command, args []string) {
			if !*yes {
				sure, err := ux.Prompt("Are you sure you want to kill the application?")
				if err != nil {
					base.Cxt.Println(err)
					return
				}
				if !sure {
					return
				}
			}
			res, logs := killApplication(req)
			logs = append(logs, fmt.Sprintf("kill application %t", res))
			base.LogInfo(logs...)
			fmt.Fprintf(out, "kill application %t\n", res)
		},
	}
	yes = cmd.Flags().BoolP("yes", "y", false, "Optional. Do not prompt for confirmation.")
	return cmd
}

func describeUFlinkApplicationCmdBuilder(req *uflink.ApplicationRequest, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "describe",
		Short: "Describe uflink application",
		Long:  "Describe uflink application",
		Run: func(cmd *cobra.Command, args []string) {
			application, err := describeUFlinkApplication(req)
			if err != nil {
				logs := fmt.Sprintf("Failed to describe application[%s] for %s", *req.ApplicationId, base.ParseError(err))
				base.LogError(logs)
				return
			}
			showDetailApplication(*application)
		},
	}
}

func listUFlinkApplicationPointsCmdBuilder(req *uflink.ApplicationRequest, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list-points",
		Short: "list-points of uflink application",
		Long:  "list-points of uflink application",
		Run: func(cmd *cobra.Command, args []string) {
			points, err := base.BizClient.ListUFlinkApplicationPoints(req)
			if err != nil {
				logs := fmt.Sprintf("Failed to list application[%s] points for %s", *req.ApplicationId, base.ParseError(err))
				base.LogError(logs)
				return
			}
			showUFlinkApplicationPoints(points.Data, out)
		},
	}
}

func runUFlinkApplicationSavepointCmdBuilder(req *uflink.ApplicationRequest, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "savepoint",
		Short: "trigger savepoint of the uflink application",
		Long:  "trigger savepoint of the uflink application",
		Run: func(cmd *cobra.Command, args [] string) {
			req.SetRetryable(false)
			resp, err := base.BizClient.RunUFlinkApplicationSavepoint(req)
			if err != nil {
				logs := fmt.Sprintf("Failed to trigger application[%s] savepoint for %s", *req.ApplicationId, base.ParseError(err))
				base.LogError(logs)
				return
			}
			fmt.Fprintf(out, "Trigger application[%s] savepoint: %t\n", *req.ApplicationId, resp.Data)
		},
	}
}

func showUFlinkApplicationPoints(points []uflink.HDFSFileStatus, out io.Writer) {
	if global.JSON {
		base.PrintJSON(points, out)
	} else {
		cols := []string{"Path", "ModificationTime", "Length", "BlockSize"}
		base.PrintTable(points, cols)
	}
}

func resubmitUFlinkApplication(req *uflink.ApplicationRequest) (bool, []string) {
	var logs []string
	logs = append(logs, fmt.Sprintf("api:SendUFlinkResubmitJob, request: %v", base.ToQueryMap(req)))
	_, err := base.BizClient.SendUFlinkResubmitJob(req)
	if err != nil {
		logs = append(logs, fmt.Sprintf("Resubmit uflink application[%s] failed: %s", *req.ApplicationId, base.ParseError(err)))
		return false, logs
	}
	logs = append(logs, fmt.Sprintf("application[%s] resubmitted", *req.ApplicationId))
	return true, logs
}

func killApplication(req *uflink.ApplicationRequest) (bool, []string) {
	var logs []string

	application, err := describeUFlinkApplication(req)
	if err != nil {
		logs = append(logs, fmt.Sprintf("describe application[%s] failed: %s", *req.ApplicationId, base.ParseError(err)))
		return false, logs
	}

	if application == nil {
		logs = append(logs, fmt.Sprintf("application[%s] does not exist", *req.ApplicationId))
		return false, logs
	}

	if application.State == "Running" {
		_req := base.BizClient.NewApplicationRequest()
		_req.ProjectId = req.ProjectId
		_req.Region = req.Region
		_req.Zone = req.Zone
		_req.InstanceId = req.InstanceId
		_req.ApplicationId = req.ApplicationId
	}
	logs = append(logs, fmt.Sprintf("api:KillUFlinkApplication, request:%v", base.ToQueryMap(req)))
	_, err = base.BizClient.KillUFlinkApplication(req)
	if err != nil {
		logs = append(logs, fmt.Sprintf("kill uflink application[%s] failed: %s", *req.ApplicationId, base.ParseError(err)))
		return false, logs
	}
	logs = append(logs, fmt.Sprintf("application[%s] killed", *req.ApplicationId))
	return true, logs
}

func describeUFlinkApplication(req *uflink.ApplicationRequest) (*uflink.ApplicationDetailData, error) {
	resp, err := base.BizClient.DescribeUFlinkApplication(req)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func showDetailApplication(application uflink.ApplicationDetailData) {
	var info []base.DescribeTableRow
	info = append(info, base.DescribeTableRow{"ApplicationId", application.Id})
	info = append(info, base.DescribeTableRow{"State", application.State})
	info = append(info, base.DescribeTableRow{"Name", application.Name})
	info = append(info, base.DescribeTableRow{"AllocatedMB", fmt.Sprintf("%f", application.AllocatedMB)})
	info = append(info, base.DescribeTableRow{"AllocatedVCores", fmt.Sprintf("%d", application.AllocatedVCores)})
	info = append(info, base.DescribeTableRow{"StartedTime", fmt.Sprintf("%d", application.StartedTime)})
	info = append(info, base.DescribeTableRow{"FinishedTime", fmt.Sprintf("%d", application.FinishedTime)})
	info = append(info, base.DescribeTableRow{"ElapsedTime", fmt.Sprintf("%d", application.ElapsedTime)})
	info = append(info, base.DescribeTableRow{"TrackingUrl", application.OriginalTrackingUrl})
	if global.JSON {
		base.PrintDescribe(info, true)
	} else {
		base.PrintDescribe(info, false)
	}
}

func showApplications(applications []uflink.ApplicationData, out io.Writer, output string) {
	list := make([]UFlinkApplicationRow, 0)
	for _, app := range applications {
		row := UFlinkApplicationRow{}
		row.ApplicationName = app.Name
		row.OriginalName = app.OriginalName
		row.Id = app.Id
		row.State = app.State
		row.StartedTime = app.StartedTime
		row.FinishedTime = app.FinishedTime
		row.AllocatedMB = app.AllocatedMB
		row.FlinkVersion = app.FlinkVersion
		row.Request = app.Request

		list = append(list, row)
	}

	if global.JSON {
		base.PrintJSON(list, out)
	} else {
		var cols []string
		if output != "wide" {
			cols = []string{"OriginalName", "ApplicationName", "Id", "FlinkVersion", "State"}
		} else {
			cols = []string{"OriginalName", "ApplicationName", "Id", "FlinkVersion", "State", "StartedTime", "FinishedTime", "AllocatedMB", "Request"}
		}
		base.PrintTable(list, cols)
	}
}

func showSubmittedJobs(jobs []uflink.SubmittedJob, out io.Writer) {
	if global.JSON {
		base.PrintJSON(jobs, out)
	} else {
		cols := []string{"Id", "Name", "SubmitStatus", "CreateTime", "JobType", "ApplicationId"}
		base.PrintTable(jobs, cols)
	}
}

func listAllApplications(req *uflink.ListUFlinkApplicationsRequest) ([]uflink.ApplicationData, error) {
	applications, _, err := listApplications(req)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func listApplications(req *uflink.ListUFlinkApplicationsRequest) ([]uflink.ApplicationData, int, error) {
	resp, err := base.BizClient.ListUFlinkApplications(req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Data, resp.TotalCount, nil
}

func Encode(v string) string {
	if v == "" {
		return ""
	}
	return url.QueryEscape(v)
}

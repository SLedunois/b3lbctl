package clusterinfo

import (
	"fmt"

	"github.com/SLedunois/b3lbctl/pkg/admin"
	"github.com/SLedunois/b3lbctl/pkg/render"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// ClusterInfoCmd struct represents the cluster info command object
type ClusterInfoCmd struct {
	Command *cobra.Command
}

// NewCmd return the instances add subcommand
func NewCmd() *cobra.Command {
	cmd := &ClusterInfoCmd{
		Command: &cobra.Command{
			Use:   "cluster-info",
			Short: "Get overall cluster information",
			Long:  `Get overall cluster information. It display all instances with %CPU, %MEM, Active meetings, Active paricipants and API status`,
		},
	}

	cmd.Command.RunE = cmd.process

	return cmd.Command
}

func colorizedAPIStatus(apiStatus string) string {
	if apiStatus == "Up" {
		return text.FgHiGreen.Sprint(apiStatus)
	}

	return text.FgHiRed.Sprint(apiStatus)
}

func colorizedMetrics(metric float64) string {
	value := fmt.Sprintf("%.2f %%", metric)
	if metric < 33.33 {
		return text.FgHiGreen.Sprint(value)
	} else if metric < 66.66 {
		return text.FgYellow.Sprint(value)
	} else {
		return text.FgHiRed.Sprint(value)
	}
}

func header() table.Row {
	return table.Row{
		text.Bold.Sprint("API"),
		text.Bold.Sprint("Host"),
		text.Bold.Sprint("CPU"),
		text.Bold.Sprint("Mem"),
		text.Bold.Sprint("Active meetings"),
		text.Bold.Sprint("Active participants"),
	}
}

func renderClusterHeaderTable(command *cobra.Command, b3lbStatus string, activeMeetings int64, activeParticipants int64, activeTenants int64) {
	t := table.NewWriter()
	t.SetStyle(render.TableStyle())
	t.AppendRow(table.Row{text.Bold.Sprint("B3LB API"), colorizedAPIStatus(b3lbStatus)})
	t.AppendRow(table.Row{text.Bold.Sprint("Active tenants"), activeTenants})
	t.AppendRow(table.Row{text.Bold.Sprint("Active meetings"), activeMeetings})
	t.AppendRow(table.Row{text.Bold.Sprint("Active participants"), activeParticipants})
	command.Println(t.Render())
	command.Println("") // print an empty line
}

func (cmd *ClusterInfoCmd) process(command *cobra.Command, args []string) error {
	status, err := admin.API.ClusterStatus()
	if err != nil {
		return fmt.Errorf("an error occurred while getting cluster status: %s", err)
	}

	b3lbStatus, err := admin.API.B3lbAPIStatus()
	if err != nil {
		return fmt.Errorf("an error occurred while getting b3lb status: %s", err)
	}

	tenants, err := admin.API.GetTenants()
	if err != nil {
		return fmt.Errorf("an error occured while getting b3lb tenants: %s", err.Error())
	}

	t := table.NewWriter()
	t.SetStyle(render.TableStyle())
	t.AppendHeader(header())

	activeMeetingSum := int64(0)
	activeParticipantsSum := int64(0)

	for _, instance := range status {
		t.AppendRow(table.Row{
			colorizedAPIStatus(instance.APIStatus),
			instance.Host,
			colorizedMetrics(instance.CPU),
			colorizedMetrics(instance.Mem),
			instance.ActiveMeeting,
			instance.ActiveParticipants,
		})

		activeMeetingSum += instance.ActiveMeeting
		activeParticipantsSum += instance.ActiveParticipants
	}

	renderClusterHeaderTable(command, b3lbStatus, activeMeetingSum, activeParticipantsSum, int64(len(tenants.Tenants)))
	command.Println(t.Render())
	return nil
}

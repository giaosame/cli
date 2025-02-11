package describe

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	corepb "github.com/projecteru2/core/rpc/gen"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Workloads describes a list of Workload
// output format can be json or yaml or table
func Workloads(workloads ...*corepb.Workload) {
	switch {
	case isJSON():
		describeAsJSON(workloads)
	case isYAML():
		describeAsYAML(workloads)
	default:
		describeWorkloads(workloads)
	}
}

// WorkloadsStatistics describes the statistics of the Workloads
func WorkloadsStatistics(workloads ...*corepb.Workload) {
	stat := struct {
		CPUs    float64
		Memory  int64
		Storage int64
	}{}
	for _, w := range workloads {
		stat.CPUs += w.Resource.CpuQuotaRequest
		stat.Memory += w.Resource.MemoryRequest
		stat.Storage += w.Resource.StorageRequest
	}

	describeStatistics := func() {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"CPUs", "Memory", "Storage"})

		rows := [][]string{
			{fmt.Sprintf("%f", stat.CPUs)},
			{fmt.Sprintf("%d", stat.Memory)},
			{fmt.Sprintf("%d", stat.Storage)},
		}
		t.AppendRows(toTableRows(rows))
		t.AppendSeparator()

		t.SetStyle(table.StyleLight)
		t.Render()
	}

	switch {
	case isJSON():
		describeAsJSON(stat)
	case isYAML():
		describeAsYAML(stat)
	default:
		describeStatistics()
	}
}

func describeWorkloads(workloads []*corepb.Workload) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name/ID/Pod/Node", "Status", "Volume", "Networks"})

	for _, c := range workloads {
		// networks
		ns := []string{}
		if c.Status != nil {
			for name, ip := range c.Status.Networks {
				if published, ok := c.Publish[name]; ok {
					addresses := strings.Split(published, ",")

					firstLine := fmt.Sprintf("%s: %s", name, addresses[0])
					ns = append(ns, firstLine)

					if len(addresses) > 1 {
						format := fmt.Sprintf("%%%ds", len(firstLine))
						for _, address := range addresses[1:] {
							ns = append(ns, fmt.Sprintf(format, address))
						}
					}
				} else {
					ns = append(ns, fmt.Sprintf("%s: %s", name, ip))
				}
			}
		}

		rows := [][]string{
			{c.Name, c.Id, c.Podname, c.Nodename},
			{
				fmt.Sprintf("CPUQuotaRequest: %f", c.Resource.CpuQuotaRequest),
				fmt.Sprintf("CPUQuotaLimit: %f", c.Resource.CpuQuotaLimit),
				fmt.Sprintf("CPUMap: %v", c.Resource.Cpu),
				fmt.Sprintf("MemoryRequest: %v", c.Resource.MemoryRequest),
				fmt.Sprintf("MemoryLimit: %v", c.Resource.MemoryLimit),
				fmt.Sprintf("StorageRequest: %v", c.Resource.StorageRequest),
				fmt.Sprintf("StorageLimit: %v", c.Resource.StorageLimit),
				fmt.Sprintf("Privileged: %v", c.Privileged),
			},
			{
				fmt.Sprintf("VolumesRequest: %+v", c.Resource.VolumesRequest),
				fmt.Sprintf("VolumesLimit: %+v", c.Resource.VolumesLimit),
				fmt.Sprintf("VolumePlanRequest: %+v", c.Resource.VolumePlanRequest),
				fmt.Sprintf("VolumePlanLimit: %+v", c.Resource.VolumePlanLimit),
			},
			ns,
		}
		t.AppendRows(toTableRows(rows))
		t.AppendSeparator()
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}

// WorkloadStatuses describes a list of WorkloadStatus
// output format can be json or yaml or table
func WorkloadStatuses(workloadStatuses ...*corepb.WorkloadStatus) {
	switch {
	case isJSON():
		describeAsJSON(workloadStatuses)
	case isYAML():
		describeAsYAML(workloadStatuses)
	default:
		describeWorkloadStatuses(workloadStatuses)
	}
}

func describeWorkloadStatuses(workloadStatuses []*corepb.WorkloadStatus) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Status", "Networks", "Extensions"})

	for _, s := range workloadStatuses {
		// networks
		ns := []string{}
		for name, ip := range s.Networks {
			ns = append(ns, fmt.Sprintf("%s: %s", name, ip))
		}

		// extensions
		extensions := map[string]string{}
		if len(s.Extension) != 0 {
			if err := json.Unmarshal(s.Extension, &extensions); err != nil {
				continue
			}
		}
		es := []string{}
		for k, v := range extensions {
			es = append(es, fmt.Sprintf("%s: %s", k, v))
		}

		rows := [][]string{
			{s.Id},
			{fmt.Sprintf("Running: %v", s.Running), fmt.Sprintf("Healthy: %v", s.Healthy)},
			ns,
			es,
		}
		t.AppendRows(toTableRows(rows))
		t.AppendSeparator()
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}

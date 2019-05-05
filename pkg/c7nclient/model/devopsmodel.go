package model

import (
	"fmt"
	"github.com/gosuri/uitable"
)

func formatText(results []EnvInfo, colWidth uint) string {

	table := uitable.New()
	table.MaxColWidth = colWidth
	table.AddRow("Name","Code","Status","Group","Cluster")
	for _, r := range results {
		table.AddRow(r.Name, r.Code, r.Status,"default",  r.Cluster)
	}

	return fmt.Sprintf("%s" ,table.String())
}
package label

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ListCommand = &cobra.Command{
		Use:   "list",
		Short: "list existing labels",
		Run:   ListLabels,
	}
)

func init() {
	LabelCmd.AddCommand(ListCommand)
}

func ListLabels(cmd *cobra.Command, args []string) {
	fmt.Println("labels:")
	for _, label := range allLabels {
		fmt.Printf("%s | %s\n", label.Name, label.re)
	}
}
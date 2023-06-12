package source

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ListCommand = &cobra.Command{
		Use:   "list",
		Short: "list existing sources",
		Run:   ListSources,
	}
)

func init() {
	SourceCmd.AddCommand(ListCommand)
}

func ListSources(cmd *cobra.Command, args []string) {
	fmt.Println("sources:")
	for labelName := range allSources {
		fmt.Println(labelName)
	}
}

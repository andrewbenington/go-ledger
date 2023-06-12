package label

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	CheckCommand = &cobra.Command{
		Use:               "check",
		Short:             "check label keywords",
		Run:               CheckLabel,
		ValidArgsFunction: autoComplete,
		Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func init() {
	LabelCmd.AddCommand(CheckCommand)
}

func autoComplete(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	labelNames := []string{}
	for _, l := range allLabels {
		labelNames = append(labelNames, l.Name)
	}
	return labelNames, cobra.ShellCompDirectiveNoFileComp
}

func CheckLabel(cmd *cobra.Command, args []string) {
	if label, ok := allLabels[args[0]]; ok {
		for _, keyword := range label.Keywords {
			fmt.Println(keyword)
		}
	}
}

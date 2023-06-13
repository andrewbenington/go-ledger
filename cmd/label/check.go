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
	var label Label
	for _, label = range allLabels {
		if label.Name == args[0] {
			break
		}
	}
	for _, keyword := range label.Keywords {
		fmt.Println(keyword)
	}
}

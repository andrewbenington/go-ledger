package label

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	AddKeywordCommand = &cobra.Command{
		Use:   "keyword",
		Short: "add keyword to label",
		Run:   AddLabelKeyword,
	}
)

func init() {
	LabelCmd.AddCommand(AddKeywordCommand)
}

func AddLabelKeyword(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("must specify label and at least one keyword")
		return
	}
	for i, label := range allLabels {
		if strings.EqualFold(label.Name, args[0]) {
			allLabels[i].Keywords = append(allLabels[i].Keywords, args[1:]...)
			err := saveLabels()
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	allLabels = append(allLabels, Label{Name: args[0], Keywords: args[1:]})
	err := saveLabels()
	if err != nil {
		fmt.Println(err)
	}
}

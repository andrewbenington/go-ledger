package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Command struct {
	Name         string
	Short        string
	Long         string
	SubCommands  []*Command
	ExpectedArgs []Argument
	Run          func(args []string) ([]Output, error)
	ShowLogs     bool
}

func (c *Command) ToCobra() *cobra.Command {
	use := strings.Join(strings.Split(c.Name, " "), "-")
	cobraCommand := &cobra.Command{
		Use:   use,
		Short: c.Short,
		Long:  c.Long,
	}
	for _, sub := range c.SubCommands {
		cobraCommand.AddCommand(sub.ToCobra())
	}
	if c.Run != nil {
		cobraCommand.Run = func(cmd *cobra.Command, args []string) {
			outputs, err := c.Run(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, o := range outputs {
				fmt.Println(o.String)
			}
		}
	}
	return cobraCommand
}

type Output struct {
	String    string
	Options   []OutputOption
	IsMessage bool
}

type OutputOption struct {
	Name   string
	Select *Command
	Args   []string
}

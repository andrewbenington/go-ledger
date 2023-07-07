package sources

import "github.com/spf13/cobra"

var (
	SourceCmd = &cobra.Command{
		Use:   "source",
		Short: "Manage sources",
		Long:  `go-ledger can pull transaction data from multiple sources`,
	}
)

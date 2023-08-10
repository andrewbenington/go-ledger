# go-ledger

go-ledger is a shell budgeting application that supports parsing transaction data
from CSVs, keyword-based categorizing the transactions, and compiling them into a
single Microsoft Excel spreadsheet.

The application can be used as a CLI, or as a terminal UI application if `go-ledger`
is executed without arguments. The CLI is uses the [Cobra](https://github.com/spf13/cobra) library and the UI uses
the [tview](https://github.com/rivo/tview) library.

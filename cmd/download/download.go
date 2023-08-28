package download

import (
	"os/exec"

	"github.com/andrewbenington/go-ledger/command"
)

var (
	DownloadCmd = &command.Command{
		Name:  "download",
		Short: "Visit sites to download transactions",
		SubCommands: []*command.Command{
			DownloadChaseCmd, DownloadVenmoCmd,
		},
	}
	DownloadChaseCmd = &command.Command{
		Name:  "chase",
		Short: "Visit chase.com",
		Run:   DownloadChase,
	}
	DownloadVenmoCmd = &command.Command{
		Name:  "download",
		Short: "Visit venmo.com/transactions",
		Run:   DownloadVenmo,
	}
	successOutput = command.Output{
		String:    "Successfully opened webpage",
		IsMessage: true,
	}
)

func DownloadChase(args []string) ([]command.Output, error) {
	openCmd := exec.Command("open", "https://chase.com")
	err := openCmd.Run()
	return []command.Output{}, err
}

func DownloadVenmo(args []string) ([]command.Output, error) {
	openCmd := exec.Command("open", "https://account.venmo.com/statement")
	err := openCmd.Run()
	return []command.Output{}, err
}

package ledger

var (
	allSources []Source
)

type SourceType string

const (
	VenmoSourceType SourceType = "VENMO"
	ChaseSourceType SourceType = "CHASE"
)

type Source interface {
	Name() string
	Type() SourceType
	GetLedgerEntries(year int) ([]Entry, error)
}

type FileSource interface {
}

func AllSources() []Source {
	return allSources
}

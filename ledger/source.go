package ledger

var (
	allSources []Source
)

type Source interface {
	Name() string
	GetLedgerEntries(year int) ([]Entry, error)
}

type FileSource interface {
}

func AllSources() []Source {
	return allSources
}

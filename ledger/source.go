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
	sources := make([]Source, 0, len(allSources))

	for _, value := range allSources {
		sources = append(sources, value)
	}
	return sources
}

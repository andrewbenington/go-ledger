package ledger

var (
	allSources []Source
)

type Source interface {
	Name() string
	GetLedgerEntries() ([]Entry, error)
}

type FileSource interface {
}

func All() []Source {
	sources := make([]Source, 0, len(allSources))

	for _, value := range allSources {
		sources = append(sources, value)
	}
	return sources
}

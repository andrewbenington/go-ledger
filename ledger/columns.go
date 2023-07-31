package ledger

const (
	IDIndex int = iota
	DateIndex
	SourceNameIndex
	SourceTypeIndex
	PersonIndex
	MemoIndex
	ValueIndex
	TypeIndex
	BalanceIndex
	LabelIndex
	NotesIndex
	FieldCount
	SwapTableStart
)

var (
	Columns = new([FieldCount]string)
)

func init() {
	Columns[IDIndex] = "ID"
	Columns[DateIndex] = "Date"
	Columns[SourceNameIndex] = "Source Name"
	Columns[SourceTypeIndex] = "Source Type"
	Columns[PersonIndex] = "Person"
	Columns[MemoIndex] = "Memo"
	Columns[ValueIndex] = "Value"
	Columns[TypeIndex] = "Type"
	Columns[BalanceIndex] = "Balance"
	Columns[LabelIndex] = "Label"
	Columns[NotesIndex] = "Notes"
}

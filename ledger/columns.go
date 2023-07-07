package ledger

const (
	ID_COLUMN               = 1
	DATE_COLUMN             = 2
	SOURCE_COLUMN           = 3
	PERSON_COLUMN           = 4
	MEMO_COLUMN             = 5
	VALUE_COLUMN            = 6
	TYPE_COLUMN             = 7
	BALANCE_COLUMN          = 8
	LABEL_COLUMN            = 9
	NOTES_COLUMN            = 10
	MAX_DATA_COLUMN         = 10
	SWAP_TABLE_START_COLUMN = 12
)

var (
	Columns = []string{
		"ID",
		"Date",
		"Source",
		"Person",
		"Memo",
		"Value",
		"Type",
		"Balance",
		"Label",
		"Notes",
	}
)

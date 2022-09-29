package enddatetimekind

// EndDateTimeKind describes how an end datetime shall be understood
//
//go:generate stringer --type EndDateTimeKind
//go:generate jsonenums --type EndDateTimeKind
type EndDateTimeKind int

const (
	// INCLUSIVE means, that the end date shall be understood as inclusive end date; e.g. "2022-10-31" for end of October
	INCLUSIVE EndDateTimeKind = iota + 1
	// EXCLUSIVE means, that the end date shall be understood as exclusive end date; e.g. "2022-11-01" for end of October
	EXCLUSIVE
)

// Code generated by "stringer --type EndDateTimeKind"; DO NOT EDIT.

package enddatetimekind

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INCLUSIVE-1]
	_ = x[EXCLUSIVE-2]
}

const _EndDateTimeKind_name = "INCLUSIVEEXCLUSIVE"

var _EndDateTimeKind_index = [...]uint8{0, 9, 18}

func (i EndDateTimeKind) String() string {
	i -= 1
	if i < 0 || i >= EndDateTimeKind(len(_EndDateTimeKind_index)-1) {
		return "EndDateTimeKind(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _EndDateTimeKind_name[_EndDateTimeKind_index[i]:_EndDateTimeKind_index[i+1]]
}

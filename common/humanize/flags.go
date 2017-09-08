package humanize

import (
	"strconv"
)

// Implements fmt.Stringer
func (f Flags) String() string {
	s := ""
	for i, name := range flagNames {
		t := Flags(1 << uint(i))
		if f&t != 0 {
			if len(s) > 0 {
				s += "|"
			}
			s += name
			f ^= t
		}
	}
	if f != 0 {
		if len(s) > 0 {
			s += "|"
		}
		s += "0x" + strconv.FormatUint(uint64(f), 16)
	}
	if len(s) == 0 {
		s = "0"
	}
	return s
}

package humanize

import (
	"fmt"
	"unicode"
)

// Int64 is a plain int64 but with a humanized implemtation of
// fmt.Formatter and fmt.Stringer.
//
// When used with printf style formating the following 'verbs' are supported
// and act as if the corresponding Flags are set:
//	%v	DefaultFlags
//	%T	a Go-syntax representation of the type of the value
//
//	%d	Decimal | Divisor1000 | NoSpace
//	%b	Decimal | Bytes | NoSpace
//	%s	Decimal | SIPrefixes | NoSpace
//
//	%#v	DefaultFlags & ^NoSpace
//	%#d	Decimal | Divisor1000
//	%#b	Decimal | Bytes
//	%#s	Decimal | SIPrefixes
//
// If a width is specified, it is used to specify the minimum string length and the output
// will be padded as required.
// If a precision is specified it is used to limit the string length.
// If a width is specified without a precision, the precision (the string limit) is
// taken to be the width.
// If neither a width or precision is specified they default to
// DefaultMinLen (0, no padding) and DefaultMaxLen (4).
//
// Other flags:
//
//	'+'	always print the sign ('+' or '-')
//	' '	(space) leave a space for elided sign (' ' or '-')
//	'-'	pad with spaces on the right rather than the left (left-justify the field)
//	'0'	pad with leading zeros rather than spaces
//	'#'	alternate format: put a space between the number and the prefix
type Int64 int64

// String implements fmt.Stringer
// and is shorthand for Format(int64(h), DefaultMinLen, DefaultMaxLen, DefaultFlags)
func (h Int64) String() string {
	return Format(int64(h), DefaultMinLen, DefaultMaxLen, DefaultFlags)
}

// Format implements fmt.Formatter
func (h Int64) Format(f fmt.State, r rune) {
	w, hasW := f.Width()
	p, hasP := f.Precision()
	//fmt.Printf("formatter: %p, %#c, %d, %d\n", f, r, w, p)
	switch {
	case hasW && !hasP:
		p = w
	case !hasW && !hasP:
		w, p = DefaultMinLen, DefaultMaxLen
	case !hasW:
		w = DefaultMinLen
	}

	var flags Flags
	// TODO, codes for these flags:
	//   flags = Decimal | Divisor1000 | Bytes | NoSpace
	//   flags = Decimal | SIPrefixes | Bytes | NoSpace
	switch unicode.ToLower(r) {
	case 'v':
		flags = DefaultFlags
	case 'd':
		flags = Decimal | Divisor1000 | NoSpace
	case 'b':
		flags = Decimal | Bytes | NoSpace
	case 's':
		flags = Decimal | SIPrefixes | NoSpace
	default:
		// Error: %!verb(type=value)
		// This will recurse into the 'v' case above.
		fmt.Fprintf(f, "%%!%c(%T=%v)", r, h, h)
		return
	}

	// XXX backward compatibility for depricated uppercase format verbs
	if f.Flag('#') || unicode.IsUpper(r) {
		flags &= ^NoSpace
	}
	if f.Flag('+') {
		flags |= AlwaysSign
	}
	if f.Flag(' ') {
		flags |= SpaceSign
	}
	if f.Flag('0') {
		flags |= ZeroPad
	}
	if f.Flag('-') {
		flags |= LeftJustify
	}

	str, _ := format(int64(h), w, p, autoScale, flags)
	f.Write([]byte(str))
}

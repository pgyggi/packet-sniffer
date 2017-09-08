package humanize

import (
	"errors"
)

// Set implements the flag.Value interface using Parse(s, DefaultFlags).
func (h *Int64) Set(s string) error {
	v, err := Parse(s, DefaultFlags)
	if err == nil {
		*h = v
	}
	return err
}

// Parsing based on ParseDuration in $GOROOT/src/pkg/time/format.go
// which is Copyright 2010 The Go Authors.

var errOverflow = errors.New("humanize: overflow") // Never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x >= (1<<63-10)/10 {
			// overflow
			return 0, "", errOverflow
		}
		x = x*10 + int64(c) - '0'
	}
	return x, s[i:], nil
}

// Parse parses a formated number string.
// A formated number string is a possibly signed number with optional
// unit prefix.
//
// The comibination of the flags argument plus the unit in the number string
// determine if the units are base2 or base10 (1024 or 1000 multipliers).
// If the unit is an SI prefix such as Ki, Mi, etc (lower case also accepted)
// then base2 multipliers are used without reguard to the flags argument.
// If the flags argument contains either the Divisor1000 or SIPrefixes flag
// then base10 multipliers are used, otherwise base2 multipliers are used.
func Parse(s string, flags Flags) (Int64, error) {
	// [-+]?([0-9]*(\.[0-9]*)? ?([kmgtpezyKMGTPEZY]i?)?
	var (
		orig      = `"` + s + `"`
		neg       = false
		x, frac   int64
		fracscale = 1
		err       error
	)

	// Consume [-+]?
	if s != "" {
		switch s[0] {
		case '-':
			neg = true
			fallthrough
		case '+':
			s = s[1:]
		}
	}

	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, errors.New("humanize: invalid size " + orig)
	}

	// The next character must be [0-9.]
	if !(s[0] == '.' || ('0' <= s[0] && s[0] <= '9')) {
		return 0, errors.New("humanize: invalid size " + orig)
	}
	// Consume [0-9]*
	pl := len(s)
	x, s, err = leadingInt(s)
	if err != nil {
		if err == errOverflow {
			return 0, errors.New("humanize: size too large " + orig)
		}
		return 0, errors.New("humanize: invalid size " + orig)
	}
	pre := pl != len(s) // whether we consumed anything before a period

	// Consume (\.[0-9]*)?
	post := false
	if s != "" && s[0] == '.' {
		s = s[1:]
		pl = len(s)
		frac, s, err = leadingInt(s)
		if err != nil {
			return 0, errors.New("humanize: invalid size " + orig)
		}
		for n := pl - len(s); n > 0; n-- {
			fracscale *= 10
		}
		post = pl != len(s)
	}
	if !pre && !post {
		// no digits (e.g. ".k" or "-.k")
		return 0, errors.New("humanize: invalid size " + orig)
	}

	// Consume optional space before unit.
	for s != "" && s[0] == ' ' { // just ' ', no '\t' etc allowed
		s = s[1:]
	}

	// All that remains is the unit.
	if s != "" {
		if len(s) > 2 || (len(s) > 1 && s[1] != 'i') {
			return 0, errors.New(`humanize: unknown unit "` + s + `" in size ` + orig)
		}
		mul := int64(1000)
		var max int64 = (1<<63 - 1000) / 1000
		base2 := flags&(Divisor1000|SIPrefixes) == 0
		if base2 || len(s) == 2 { // SI prefix: Ki, Mi, etc
			mul = 1024
			max = (1<<63 - 1024) / 1024
		}
		scale := 0
		switch s[0] {
		case 'b', 'B':
			scale = 0
		case 'k', 'K':
			scale = 1
		case 'm', 'M':
			scale = 2
		case 'g', 'G':
			scale = 3
		case 't', 'T':
			scale = 4
		case 'p', 'P':
			scale = 5
		case 'e', 'E':
			scale = 6
		case 'z', 'Z':
			scale = 7
		case 'y', 'Y':
			scale = 8
		default:
			return 0, errors.New(`humanize: unknown unit "` + s + `" in size ` + orig)
		}
		for ; scale > 0; scale-- {
			if frac >= max && fracscale > 1 {
				frac /= int64(fracscale)
				fracscale = 1
			}
			if x >= max || frac >= max {
				return 0, errors.New("humanize: size too large " + orig)
			}
			x *= mul
			frac *= mul
		}
	}

	// Combine parts, negate if required
	frac /= int64(fracscale)
	if x > 1<<63-1-frac {
		return 0, errors.New("humanize: size too large " + orig)
	}
	x += frac
	if neg {
		x = -x
	}

	return Int64(x), nil
}

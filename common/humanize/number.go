// The format function based on humanize_number.c

package humanize

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	autoScale    = -1 // Format using lowest multiplier possible
	getScale     = -1 << 20
	prefixes     = "B KiMiGiTiPiEiZiYi"
	decimalPoint = "." // TODO: localize
)

// Based on BSD's src/lib/libutil/humanize_number.c
//
// Note that the term prefix in the following refers to the unit; the original code
// supported adding a suffix string (e.g. "B", "B/s", etc) such that the unit had
// prefix and suffix components. We don't support suffixes since it's trivial in
// Go code to just use + and a string value in a wrapper if desired.
func format(value int64, minlen, maxlen int, scale int, flags Flags) (str string, cnt uint) {
	var (
		sign    = []rune{'+'}
		divisor int  // 1024 or 1000
		baselen = 2  // digit + prefix
		sep     = "" // "" or " "
		prefix  string
	)

	if scale >= int(MaxScale) {
		panic("humanize.Format: bad scale")
	}

	if flags&Divisor1000 != 0 {
		divisor = 1000
		flags &= ^SIPrefixes
	} else {
		divisor = 1 << 10
		if flags&SIPrefixes != 0 {
			baselen++ // for 'i' part of prefix
		}
	}

	switch {
	case value < 0:
		sign[0] = '-'
		value *= -1
	case flags&AlwaysSign != 0:
		sign[0] = '+'
	case flags&SpaceSign != 0:
		sign[0] = ' '
	default:
		sign = sign[:0]
	}
	baselen += len(sign)

	if flags&NoSpace == 0 {
		sep = " "
		baselen++
	}

	// Not enough room for even a single digit + sep + prefix
	if maxlen < baselen {
		panic(fmt.Errorf("humanize.Format: maxlen %d < %d", maxlen, baselen))
	}

	remainder := 0
	if scale < 0 {
		if maxlen-baselen >= 17 {
			// 1e18 is the largest power of 10 that will fit in int64
			// so we can't compute max below;
			// just leave cnt as 0 and we won't divide the value at all.
		} else {
			// See if additional columns can be used.
			var max int64 = 1
			for i := maxlen - baselen; i >= 0; i-- {
				max *= 10
			}

			// Divide the value until it fits the given space.
			// If there will be an overflow in the later rounding step,
			// divide one extra time.
			for ; (value >= max || (value == max-1 && remainder >= divisor/2)) && cnt < uint(MaxScale); cnt++ {
				remainder = int(value % int64(divisor))
				value /= int64(divisor)
			}
		}
		if scale == getScale {
			return "", cnt
		}
	} else {
		for ; cnt < uint(scale) && cnt < uint(MaxScale); cnt++ {
			remainder = int(value % int64(divisor))
			value /= int64(divisor)
		}
	}

	if cnt == 0 && flags&Bytes == 0 {
		// No 'B' prefix
		prefix = prefixes[:0]
	} else {
		if flags&SIPrefixes != 0 && divisor == 1024 {
			if cnt == 0 {
				// Just 'B'
				prefix = prefixes[:1]
			} else {
				// Two character SI prefixes
				prefix = prefixes[2*cnt : 2*cnt+2]
			}
		} else {
			if cnt == 1 && divisor == 1024 {
				prefix = "k" // Instead of "K"
			} else {
				// Single character prefixes
				prefix = prefixes[2*cnt : 2*cnt+1]
			}
		}
	}

	// if value will round to <= 9.9 ...
	p05 := divisor / 20
	p95 := divisor - p05 // p95*divisor >= 0.95, the point we'd round to 1.0
	if (value < 9 || (value == 9 && remainder < p95)) && cnt > 0 && flags&Decimal != 0 && maxlen >= baselen+2 {
		s2 := (remainder*10 + divisor/2) / divisor
		s1 := int(value) + s2/10
		s2 %= 10

		str = strconv.Itoa(s1) + decimalPoint +
			strconv.Itoa(s2) + sep + prefix
	} else {
		str = strconv.FormatInt(value+int64((remainder+divisor/2)/divisor), 10) +
			sep + prefix
	}

	// Add sign and padding
	if padsize := minlen - len(str) - len(sign); padsize > 0 || len(sign) > 0 {
		switch {
		case padsize <= 0:
			str = string(sign) + str
		case flags&ZeroPad != 0:
			str = string(sign) + strings.Repeat("0", padsize) + str
		case flags&LeftJustify != 0:
			str = string(sign) + str + strings.Repeat(" ", padsize)
		default:
			str = strings.Repeat(" ", padsize) + string(sign) + str
		}
	}

	return str, cnt
}

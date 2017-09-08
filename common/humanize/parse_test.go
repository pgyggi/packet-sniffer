package humanize

import (
	"fmt"
	"math"
	"testing"
)

func TestParse(t *testing.T) {
	const invalid = "humanize: invalid size %q"
	const unit = "humanize: unknown unit %q in size %q"
	const toolarge = "humanize: size too large %q"
	data := []struct {
		in    string
		flags Flags
		want  int64
		err   error
	}{
		{"", 0, 0, fmt.Errorf(invalid, "")},
		{"0", 0, 0, nil},
		{"-0", 0, 0, nil},
		{"-", 0, 0, fmt.Errorf(invalid, "-")},
		{".", 0, 0, fmt.Errorf(invalid, ".")},
		{".k", 0, 0, fmt.Errorf(invalid, ".k")},
		{"-.k", 0, 0, fmt.Errorf(invalid, "-.k")},
		{"1k", 0, 1024, nil},
		{"1K", 0, 1024, nil},
		{"1 K", 0, 1024, nil},
		{"1  K", 0, 1024, nil},
		{"1K ", 0, 1024, fmt.Errorf(unit, "K ", "1K ")},
		{" 1K", 0, 1024, fmt.Errorf(invalid, " 1K")},
		{"1k", Divisor1000, 1000, nil},
		{"1Ki", Divisor1000, 1024, nil},
		{"1.Ki", 0, 1024, nil},
		{"1.0Ki", 0, 1024, nil},
		{"1.0000Ki", 0, 1024, nil},
		{"010.000Ki", 0, 10 * 1024, nil},
		{".5Ki", 0, 512, nil},
		{"0.5Ki", 0, 512, nil},
		{"0.05Ki", 0, 51, nil},
		{"0.05Mi", 0, 52428, nil},
		{"1.5Ki", 0, 1024 + 512, nil},
		{"1.05Ki", 0, 1024 + 51, nil},
		{"5B", Divisor1000, 5e0, nil},
		{"5b", Divisor1000, 5e0, nil},
		{"5K", Divisor1000, 5e3, nil},
		{"5k", Divisor1000, 5e3, nil},
		{"5M", Divisor1000, 5e6, nil},
		{"5m", Divisor1000, 5e6, nil},
		{"5G", Divisor1000, 5e9, nil},
		{"5g", Divisor1000, 5e9, nil},
		{"5T", Divisor1000, 5e12, nil},
		{"5t", Divisor1000, 5e12, nil},
		{"5P", Divisor1000, 5e15, nil},
		{"5p", Divisor1000, 5e15, nil},
		{"5E", Divisor1000, 5e18, nil},
		{"5e", Divisor1000, 5e18, nil},
		{"9.22337e", Divisor1000, 9.22337e18, nil},
		{".5E", 0, 1 << 59, nil},
		{".9P", 0, 1013309916158361, nil},
		{"7E", 0, 7 * 1 << 60, nil},
		{"7.5E", 0, 15 * 1 << 59, nil},
		{"7.75E", 0, 31 * 1 << 58, nil},
		{"7.875E", 0, 63 * 1 << 57, nil},
		{"7.9375E", 0, 127 * 1 << 56, nil},
		{"7.96875E", 0, 255 * 1 << 55, nil},
		{"7.984375E", 0, 511 * 1 << 54, nil},
		{"-7.984375E", 0, -511 * 1 << 54, nil},
		{"10e", Divisor1000, 0, fmt.Errorf(toolarge, "10e")},
		{"5Z", Divisor1000, 0, fmt.Errorf(toolarge, "5Z")},
		{"5z", Divisor1000, 0, fmt.Errorf(toolarge, "5z")},
		{"5Y", Divisor1000, 0, fmt.Errorf(toolarge, "5Y")},
		{"5y", Divisor1000, 0, fmt.Errorf(toolarge, "5y")},
		{"4611686018427387903", 0, math.MaxInt64 / 2, nil},
		{"-4611686018427387903", 0, -math.MaxInt64 / 2, nil},
		{"9223372036854775789", 0, 9223372036854775789, nil},
		{"-9223372036854775789", 0, -9223372036854775789, nil},
		// XXX these two should be made to work
		//{"9223372036854775790", 0, 9223372036854775790, nil},
		//{"9223372036854775807", 0, math.MaxInt64, nil},
		{"9223372036854775790", 0, 9223372036854775790, fmt.Errorf(toolarge, "9223372036854775790")},
		{"9223372036854775807", 0, math.MaxInt64, fmt.Errorf(toolarge, "9223372036854775807")},
	}
	for i, d := range data {
		val, err := Parse(d.in, d.flags)
		if err != nil {
			if d.err == nil {
				t.Errorf("%2d: Parse(%q,%v) failed: %v, wanted %d",
					i, d.in, d.flags, err, d.want)
			} else if err.Error() != d.err.Error() {
				t.Errorf("%2d: Parse(%q,%v) failed: %s, wanted %s",
					i, d.in, d.flags, err.Error(), d.err.Error())
			} else if false {
				t.Logf("%2d: Parse(%q,%v) = %v, %v (as expected)",
					i, d.in, d.flags, int64(val), err)
			}
			continue
		} else if d.err != nil {
			t.Errorf("%2d: Parse(%q,%v) = %d, expected failure",
				i, d.in, d.flags, int64(val))
			continue
		}
		if int64(val) != d.want {
			t.Errorf("%2d: Parse(%q,%v) = %d want %d",
				i, d.in, d.flags, int64(val), d.want)
		} else if false {
			t.Logf("%2d: Parse(%q,%v) = %v, %v (as expected)",
				i, d.in, d.flags, int64(val), err)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("4.5T", 0)
	}
}

// Copyright (C) Cyrill@Schumacher.fm @SchumacherFM Twitter/GitHub
// Wanderlust - a cache warmer for your web app with priorities
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package helpers

import (
	"bytes"
	"unicode/utf8"
)

type FfjsonIf interface {
	MarshalJSON() ([]byte, error)
	MarshalJSONBuf(buf *bytes.Buffer) error
}

func Ffjson_WriteJsonString(buf *bytes.Buffer, s string) {
	const hex = "0123456789abcdef"

	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			default:

				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}

		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
}

func Ffjson_FormatBits(dst *bytes.Buffer, u uint64, base int, neg bool) {
	const (
		digits   = "0123456789abcdefghijklmnopqrstuvwxyz"
		digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
		digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
	)

	var shifts = [len(digits) + 1]uint{
		1 << 1: 1,
		1 << 2: 2,
		1 << 3: 3,
		1 << 4: 4,
		1 << 5: 5,
	}

	if base < 2 || base > len(digits) {
		panic("strconv: illegal AppendInt/FormatInt base")
	}

	var a [64 + 1]byte
	i := len(a)

	if neg {
		u = -u
	}

	if base == 10 {

		for u >= 100 {
			i -= 2
			q := u / 100
			j := uintptr(u - q*100)
			a[i+1] = digits01[j]
			a[i+0] = digits10[j]
			u = q
		}
		if u >= 10 {
			i--
			q := u / 10
			a[i] = digits[uintptr(u-q*10)]
			u = q
		}

	} else if s := shifts[base]; s > 0 {

		b := uint64(base)
		m := uintptr(b) - 1
		for u >= b {
			i--
			a[i] = digits[uintptr(u)&m]
			u >>= s
		}

	} else {

		b := uint64(base)
		for u >= b {
			i--
			a[i] = digits[uintptr(u%b)]
			u /= b
		}
	}

	i--
	a[i] = digits[uintptr(u)]

	if neg {
		i--
		a[i] = '-'
	}

	dst.Write(a[i:])

	return
}

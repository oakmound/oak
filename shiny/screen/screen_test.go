// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package screen

import (
	"testing"
)

func TestSanitizeUTF8(t *testing.T) {
	const n = 8

	testCases := []struct {
		s, want string
	}{
		{"", ""},
		{"a", "a"},
		{"a\x00", "a"},
		{"a\x80", "a"},
		{"\x00a", ""},
		{"\x80a", ""},
		{"abc", "abc"},
		{"foo b\x00r qux", "foo b"},
		{"foo b\x80r qux", "foo b"},
		{"foo b\xffr qux", "foo b"},

		// "\xc3\xa0" is U+00E0 LATIN SMALL LETTER A WITH GRAVE.
		{"\xc3\xa0pqrs", "\u00e0pqrs"},
		{"a\xc3\xa0pqrs", "a\u00e0pqrs"},
		{"ab\xc3\xa0pqrs", "ab\u00e0pqrs"},
		{"abc\xc3\xa0pqrs", "abc\u00e0pqr"},
		{"abcd\xc3\xa0pqrs", "abcd\u00e0pq"},
		{"abcde\xc3\xa0pqrs", "abcde\u00e0p"},
		{"abcdef\xc3\xa0pqrs", "abcdef\u00e0"},
		{"abcdefg\xc3\xa0pqrs", "abcdefg"},
		{"abcdefgh\xc3\xa0pqrs", "abcdefgh"},
		{"abcdefghi\xc3\xa0pqrs", "abcdefgh"},
		{"abcdefghij\xc3\xa0pqrs", "abcdefgh"},

		// "世" is "\xe4\xb8\x96".
		// "界" is "\xe7\x95\x8c".
		{"H 世界", "H 世界"},
		{"Hi 世界", "Hi 世"},
		{"Hello 世界", "Hello "},
	}

	for _, tc := range testCases {
		if got := sanitizeUTF8(tc.s, n); got != tc.want {
			t.Errorf("sanitizeUTF8(%q): got %q, want %q", tc.s, got, tc.want)
		}
	}
}

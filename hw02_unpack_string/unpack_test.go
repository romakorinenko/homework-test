package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "🙃5", expected: "🙃🙃🙃🙃🙃"},
		{input: "a", expected: "a"},
		{input: "🙃", expected: "🙃"},
		{input: "aaф0b", expected: "aab"},
		{input: "aaф9b", expected: "aaфффффффффb"},
		{input: "g?3c", expected: "g???c"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "界", expected: "界"},
		{input: "Hello, 世界", expected: "Hello, 世界"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\3\4\5`, expected: `345`},
		{input: `\`, expected: `\`},
		{input: `\\`, expected: `\`},
		{input: `b\`, expected: `b\`},
		{input: `\0`, expected: `0`},
		{input: `\00`, expected: ``},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `\\55`, `a\a`}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

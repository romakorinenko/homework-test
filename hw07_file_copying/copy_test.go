package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		name   string
		offset int64
		limit  int64
		file   string
	}{
		{"out_offset0_limit0", 0, 0, "testdata/out_offset0_limit0.txt"},
		{"out_offset0_limit10", 0, 10, "testdata/out_offset0_limit10.txt"},
		{"out_offset0_limit1000", 0, 1000, "testdata/out_offset0_limit1000.txt"},
		{"out_offset0_limit10000", 0, 10000, "testdata/out_offset0_limit10000.txt"},
		{"out_offset100_limit1000", 100, 1000, "testdata/out_offset100_limit1000.txt"},
		{"out_offset6000_limit1000", 6000, 1000, "testdata/out_offset6000_limit1000.txt"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out := "out.txt"
			err := Copy("testdata/input.txt", out, testCase.offset, testCase.limit)
			require.NoError(t, err)

			actual, err := os.ReadFile(out)
			require.NoError(t, err)
			expected, err := os.ReadFile(testCase.file)
			require.NoError(t, err)
			require.Equal(t, expected, actual)

			err = os.Remove(out)
			require.NoError(t, err)
		})
	}
}

func TestCopy_ErrOffsetExceedsFileSize(t *testing.T) {
	err := Copy("testdata/out_offset0_limit1000.txt", "out.txt", 99999999999999, 10)
	require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
}

func TestCopy_ErrNotExist(t *testing.T) {
	err := Copy("notExist.txt", "out.txt", 0, 0)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestCopy_ErrUnsupportedFile(t *testing.T) {
	dir := t.TempDir()
	err := Copy(dir, "out.txt", 0, 0)
	require.ErrorIs(t, err, ErrUnsupportedFile)
}

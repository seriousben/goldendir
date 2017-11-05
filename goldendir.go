package goldendir

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var getPathnames = func(t require.TestingT, dirpath string) []string {
	expectedPathnames := []string{}
	err := godirwalk.Walk(dirpath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() || de.IsSymlink() {
				return nil
			}
			relpath, err := filepath.Rel(dirpath, osPathname)
			require.NoError(t, err)
			expectedPathnames = append(expectedPathnames, relpath)
			return nil
		},
	})
	require.NoError(t, err)
	return expectedPathnames
}

type compFlag uint8

const (
	compMissing compFlag = iota
	compSame
	compAdd
)

func compare(actual, expected []string) map[string]compFlag {
	comparison := map[string]compFlag{}
	for _, exp := range expected {
		comparison[exp] = compMissing
	}
	for _, act := range actual {
		if _, ok := comparison[act]; ok {
			comparison[act] = compSame
		} else {
			comparison[act] = compAdd
		}
	}
	return comparison
}

// Assert compares the actual dir to the expected content in the golden dir.
// Returns whether the assertion was successful (true) or not (false)
func Assert(t require.TestingT, actualDir, expectedDir string, msgAndArgs ...interface{}) bool {
	actualPathnames := getPathnames(t, actualDir)
	expectedPathnames := getPathnames(t, expectedDir)

	comp := compare(actualPathnames, expectedPathnames)

	details := []string{}
	for filename, flag := range comp {
		switch flag {
		case compMissing:
			details = append(details, fmt.Sprintf("- %s", filename))
		case compSame:
			actual, err := ioutil.ReadFile(filepath.Join(actualDir, filename))
			require.NoError(t, err, msgAndArgs...)
			expected, err := ioutil.ReadFile(filepath.Join(expectedDir, filename))
			require.NoError(t, err, msgAndArgs...)

			if assert.ObjectsAreEqual(expected, actual) {
				break
			}

			diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				A:        difflib.SplitLines(string(expected)),
				B:        difflib.SplitLines(string(actual)),
				FromFile: "Expected",
				ToFile:   "Actual",
				Context:  3,
			})
			require.NoError(t, err, msgAndArgs...)
			details = append(details, fmt.Sprintf("~ %s [ \n%s ]", filename, diff))
		case compAdd:
			details = append(details, fmt.Sprintf("+ %s", filename))
		default:
			require.Fail(t, "unknown compare flag %+v", flag)
		}
	}

	if len(details) != 0 {
		require.Fail(t, fmt.Sprintf("Not Equal: \n%s", strings.Join(details, "\n")), msgAndArgs...)
		return false
	}

	return true
}

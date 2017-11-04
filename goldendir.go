package goldendir

import (
	"fmt"
	"reflect"

	"github.com/karrick/godirwalk"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var getPathnames = func(dirpath string) ([]string, error) {
	expectedPathnames := []string{}
	err := godirwalk.Walk(dirpath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			expectedPathnames = append(expectedPathnames, osPathname)
			return nil
		},
	})
	return expectedPathnames, err
}

func addEol(pathnames []string) []string {
	for idx, name := range pathnames {
		pathnames[idx] = name + "\n"
	}
	return pathnames
}

// Assert compares the actual dir to the expected content in the golden dir.
// Returns whether the assertion was successful (true) or not (false)
func Assert(t require.TestingT, actualDir, expectedDir string, msgAndArgs ...interface{}) bool {
	actualPathnames, err := getPathnames(actualDir)
	if err != nil {
		require.NoError(t, err, msgAndArgs...)
		return false
	}
	expectedPathnames, err := getPathnames(expectedDir)
	if err != nil {
		require.NoError(t, err, msgAndArgs...)
		return false
	}

	sameFiles := reflect.DeepEqual(actualPathnames, expectedPathnames)
	if !sameFiles {
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        addEol(expectedPathnames),
			B:        addEol(actualPathnames),
			FromFile: "Expected",
			ToFile:   "Actual",
			Context:  0,
		})
		require.NoError(t, err, msgAndArgs...)
		return assert.Fail(t, fmt.Sprintf("Not Equal: \n%s", diff), msgAndArgs...)
	}

	return true
}

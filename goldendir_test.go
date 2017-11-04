package goldendir

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type noopT struct {
	errors []string
	failed bool
}

func (t *noopT) Errorf(format string, args ...interface{}) {
	t.errors = append(t.errors, fmt.Sprintf(format, args))
}
func (t *noopT) FailNow() {
	t.failed = true
}

func TestActualPathnameFails(t *testing.T) {
	tmpPathnames := getPathnames
	getPathnames = func(dirpath string) ([]string, error) { return []string{}, errors.New("getPathnames error") }
	defer func() { getPathnames = tmpPathnames }()
	nt := noopT{}
	success := Assert(&nt, ".", ".", "")
	assert.False(t, success, "getpathname error fails assert")
	assert.True(t, nt.failed, "assert called failNow")
}

func TestExpectedPathnameFails(t *testing.T) {
	numCalled := 0
	tmpPathnames := getPathnames
	getPathnames = func(dirpath string) ([]string, error) {
		numCalled++
		if numCalled == 2 {
			return []string{}, errors.New("getPathnames error")
		}
		return tmpPathnames(dirpath)
	}
	defer func() { getPathnames = tmpPathnames }()
	nt := noopT{}
	success := Assert(&nt, ".", ".", "")
	assert.False(t, success, "getpathname error fails assert")
	assert.True(t, nt.failed, "assert called failNow")
}

func TestAssertSuccess(t *testing.T) {
	assert.True(t, Assert(&noopT{}, ".", ".", "same dir"))
}

func containsMatch(elems []string, match string) bool {
	for _, elem := range elems {
		if strings.Contains(elem, match) {
			return true
		}
	}
	return false
}

func TestAssertFail(t *testing.T) {
	nt := noopT{}
	failed := Assert(&nt, "vendor/github.com/karrick", "vendor/github.com/pkg", "same dir")
	assert.False(t, failed, "directories are different")
	assert.False(t, nt.failed, "should not have failed now")
	assert.True(t, containsMatch(nt.errors, "-vendor/github.com/pkg"), "missing files from pkg")
	assert.True(t, containsMatch(nt.errors, "+vendor/github.com/karrick"), "extra files from karrick")
}

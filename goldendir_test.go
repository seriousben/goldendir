package goldendir

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopT struct {
	errors []string
	failed bool
}

func (t *noopT) Errorf(format string, args ...interface{}) {
	t.errors = append(t.errors, fmt.Sprintf(format, args...))
}
func (t *noopT) FailNow() {
	t.failed = true
}

func TestActualPathnameFails(t *testing.T) {
	tmpPathnames := getPathnames
	getPathnames = func(t require.TestingT, dirpath string) []string {
		t.FailNow()
		return []string{}
	}
	defer func() { getPathnames = tmpPathnames }()
	nt := noopT{}
	Assert(&nt, ".", ".", "")
	assert.True(t, nt.failed, "assert called failNow")
}

func TestExpectedPathnameFails(t *testing.T) {
	numCalled := 0
	tmpPathnames := getPathnames
	getPathnames = func(t require.TestingT, dirpath string) []string {
		numCalled++
		if numCalled == 2 {
			t.FailNow()
			return []string{}
		}
		return tmpPathnames(t, dirpath)
	}
	defer func() { getPathnames = tmpPathnames }()
	nt := noopT{}
	Assert(&nt, ".", ".", "")
	assert.True(t, nt.failed, "assert called failNow")
}

func TestCompare(t *testing.T) {
	compareTests := []struct {
		testCase      string
		inputActual   []string
		inputExpected []string
		expected      map[string]compFlag
	}{
		{"empty", []string{}, []string{}, map[string]compFlag{}},
		{"same", []string{"file"}, []string{"file"}, map[string]compFlag{
			"file": compSame,
		}},
		{"missing", []string{}, []string{"file"}, map[string]compFlag{
			"file": compMissing,
		}},
		{"add", []string{"file"}, []string{}, map[string]compFlag{
			"file": compAdd,
		}},
	}

	for _, test := range compareTests {
		t.Run(test.testCase, func(t *testing.T) {
			comp := compare(test.inputActual, test.inputExpected)
			assert.Equal(t, test.expected, comp)
		})
	}

}

func containsMatch(elems []string, match string) bool {
	for _, elem := range elems {
		if strings.Contains(elem, match) {
			return true
		}
	}
	return false
}

func TestAssertSuccess(t *testing.T) {
	assertTests := []struct {
		testCase    string
		actualOps   []fs.PathOp
		expectedOps []fs.PathOp
	}{
		{"empty", []fs.PathOp{}, []fs.PathOp{}},
		{"simple", []fs.PathOp{
			fs.WithFile("file1", "content\n"),
		}, []fs.PathOp{
			fs.WithFile("file1", "content\n"),
		}},
		{"multiple files", []fs.PathOp{
			fs.WithFile("file1", "foo\n"),
			fs.WithFile("file2", "bar\n"),
		}, []fs.PathOp{
			fs.WithFile("file1", "foo\n"),
			fs.WithFile("file2", "bar\n"),
		}},
		{"nested directories", []fs.PathOp{
			fs.WithFile("file1", "foo\n"),
			fs.WithFile("file2", "bar\n"),
			fs.WithDir("dir1", fs.WithFile("file3", "content\n")),
		}, []fs.PathOp{
			fs.WithFile("file1", "foo\n"),
			fs.WithFile("file2", "bar\n"),
			fs.WithDir("dir1", fs.WithFile("file3", "content\n")),
		}},
	}

	for _, test := range assertTests {
		t.Run(test.testCase, func(t *testing.T) {
			actualDir := fs.NewDir(t, "", test.actualOps...)
			defer actualDir.Remove()
			expectedDir := fs.NewDir(t, "", test.expectedOps...)
			defer expectedDir.Remove()
			Assert(t, actualDir.Path(), expectedDir.Path(), "same content")
		})
	}

}

func TestAssertFailMissingFiles(t *testing.T) {
	nt := new(noopT)

	expectedDir := fs.NewDir(t, "expected", fs.WithFile("file1", "content\n"))
	defer expectedDir.Remove()

	actualDir := fs.NewDir(t, "actual")
	defer actualDir.Remove()

	Assert(nt, actualDir.Path(), expectedDir.Path())
	assert.True(t, nt.failed, "failed for missing files")
	assert.True(t, containsMatch(nt.errors, "- file1"), "missing file1")
}

func TestAssertFailExtraFiles(t *testing.T) {
	nt := new(noopT)

	expectedDir := fs.NewDir(t, "expected", fs.WithFile("file1", "content\n"))
	defer expectedDir.Remove()

	actualDir := fs.NewDir(t, "actual",
		fs.WithFile("file1", "content\n"),
		fs.WithFile("file2", "foo\n"))

	defer actualDir.Remove()

	Assert(nt, actualDir.Path(), expectedDir.Path())
	assert.True(t, nt.failed, "failed for new files")
	assert.True(t, containsMatch(nt.errors, "+ file2"), "new file2")
}

func TestAssertFailDifferentContent(t *testing.T) {
	nt := new(noopT)

	expectedDir := fs.NewDir(t, "expected", fs.WithFile("file1", "content\nfoo\n"))
	defer expectedDir.Remove()

	actualDir := fs.NewDir(t, "actual", fs.WithFile("file1", "content\nbar\n"))
	defer actualDir.Remove()

	Assert(nt, actualDir.Path(), expectedDir.Path())
	assert.True(t, nt.failed, "fail for different content")
	assert.True(t, containsMatch(nt.errors, "~ file1"), "different file1")
	assert.True(t, containsMatch(nt.errors, "-foo"), "file1 has missing content")
	assert.True(t, containsMatch(nt.errors, "+bar"), "file1 has added content")
}

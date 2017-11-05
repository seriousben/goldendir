package goldendir_test

import (
	"testing"

	"github.com/seriousben/goldendir"
)

var t = &testing.T{}

func ExampleAssert() {
	goldendir.Assert(t, "/path/to/foo-dir", "foo-dir.golden")
}

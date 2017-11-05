# goldendir

A package compatible with `go test` to support golden file testing patterns but on directories.

[![GoDoc](https://godoc.org/github.com/seriousben/goldendir?status.svg)](https://godoc.org/github.com/seriousben/goldendir)
[![CircleCI](https://circleci.com/gh/seriousben/goldendir/tree/master.svg?style=shield)](https://circleci.com/gh/seriousben/goldendir/tree/master)
[![Go Reportcard](https://goreportcard.com/badge/github.com/seriousben/goldendir)](https://goreportcard.com/report/github.com/seriousben/goldendir)
[![codecov](https://codecov.io/gh/seriousben/goldendir/branch/master/graph/badge.svg)](https://codecov.io/gh/seriousben/goldendir)

## Detailed Output

`goldendir` makes it easy to spot differences.

For example:
```console
+ extrafile
- missingfile
~ changedfile [
  --- Expected
  +++ Actual
  @@ -1,3 +1,3 @@
   content
  -bar
  +foo

]
+ dir/extrafile1
- dir/missingfile1
```

## Example

```go
import (
	"testing"

	"github.com/seriousben/goldendir"
)

func TestOutput(t *testing.T) {
	goldendir.Assert(t, "/path/to/foo-dir", "foo-dir.golden")
}
```

## Install

`go get -u github.com/seriousben/goldendir`

## Related

* [gotestyourself/golden](https://godoc.org/github.com/gotestyourself/gotestyourself/golden) - compare large multi-line strings
* [testify/assert](https://godoc.org/github.com/stretchr/testify/assert) and
  [testify/require](https://godoc.org/github.com/stretchr/testify/require) -
  assertion libraries with common assertions

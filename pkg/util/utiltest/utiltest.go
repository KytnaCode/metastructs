// Package utiltest contains a bunch of helpers for testing.
package utiltest

import (
	"go/types"
	"testing"

	"github.com/kytnacode/metastructs/pkg/util"
)

// GetType returns a type by its name on current package, if none found test fail inmediantly.
func GetType(t *testing.T, typName string) *types.Named {
	t.Helper()

	named, _, err := util.LoadType(t.Context(), ".", typName)
	if err != nil {
		t.Fatal(err)
	}

	return named
}

// PrintDiff prints difference between `expected` and `got`.
func PrintDiff(t *testing.T, expected, got string) {
	t.Helper()

	t.Logf("expected:\n%v\n\n----------\n\ngot:\n%v\n", expected, got)
}

package main

import (
	"go/types"
	"testing"
)

func getType(t *testing.T, typName string) *types.Named {
	t.Helper()

	named, _, err := loadType(t.Context(), ".", typName)
	if err != nil {
		t.Fatal(err)
	}

	return named
}

func printDiff(t *testing.T, expected, got string) {
	t.Helper()

	t.Logf("expected:\n%v\n\n----------\n\ngot:\n%v\n", expected, got)
}

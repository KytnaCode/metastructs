package main

import (
	"go/types"
	"io"
	"strings"
	"testing"
	"text/template"
)

const (
	testStructName        = "testStruct"
	nonStructTestTypeName = "nonStructTestType"
)

var expectedTmpl *template.Template

type expectedData struct {
	Package    string
	Comment    string
	Receiver   string
	RecvType   string
	MethodName string
	NameKey    string
	CountKey   string
	SliceKey   string
}

type testStruct struct {
	Name  string
	Count int
	Slice []float64
}

type nonStructTestType string

func init() {
	var err error

	expectedTmpl, err = template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() map[string]any {
	return map[string]any{
		"{{.CountKey}}": {{.Receiver}}.Count,
		"{{.NameKey}}":  {{.Receiver}}.Name,
		"{{.SliceKey}}": {{.Receiver}}.Slice,
	}
}
`)
	if err != nil {
		panic(err)
	}
}

func printDiff(t *testing.T, expected, got string) {
	t.Helper()

	t.Logf("expected:\n%v\n\n----------\n\ngot:\n%v\n", expected, got)
}

func getType(t *testing.T, typName string) *types.Named {
	t.Helper()

	pkgs, err := loadPackages(t.Context(), ".")
	if err != nil {
		t.Fatal(err)
	}

	_ = testStruct{}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(typName)
		if obj == nil {
			continue
		}

		typ, ok := obj.Type().(*types.Named)
		if !ok {
			panic("object type must be named")
		}

		return typ
	}

	t.Fatalf("`%v` not found", typName)

	return nil
}

func getTestStructType(t *testing.T) *types.Named {
	t.Helper()

	return getType(t, testStructName)
}

func TestToMap_Defaults(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "main"

	cfg := ToMapConfig{
		PkgName: pkgName,
		Typ:     typ,
	}

	var res strings.Builder

	if err := ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	var expected strings.Builder

	err := expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgName,
		Comment:    PackageComment,
		Receiver:   MethodReceiver,
		RecvType:   typ.Obj().Name(),
		MethodName: DefaultMethodName,
		NameKey:    "Name",
		CountKey:   "Count",
		SliceKey:   "Slice",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		printDiff(t, expected.String(), res.String())
		t.Fail()
	}
}

func TestToMap_CustomMethodName(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "abc"

	const methodName = "CustomMethodName"

	cfg := ToMapConfig{
		PkgName:    pkgName,
		Typ:        typ,
		MethodName: methodName,
	}

	var res strings.Builder

	if err := ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	var expected strings.Builder

	err := expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgName,
		Comment:    PackageComment,
		Receiver:   MethodReceiver,
		RecvType:   typ.Obj().Name(),
		MethodName: methodName,
		NameKey:    "Name",
		CountKey:   "Count",
		SliceKey:   "Slice",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		printDiff(t, expected.String(), res.String())
		t.Fail()
	}
}

func TestToMap_PointerReceiver(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "models"

	cfg := ToMapConfig{
		PkgName: pkgName,
		Typ:     typ,
		Pointer: true,
	}

	var res strings.Builder

	if err := ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	var expected strings.Builder

	err := expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgName,
		Comment:    PackageComment,
		Receiver:   MethodReceiver,
		RecvType:   "*" + typ.Obj().Name(),
		MethodName: DefaultMethodName,
		NameKey:    "Name",
		CountKey:   "Count",
		SliceKey:   "Slice",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		printDiff(t, expected.String(), res.String())

		t.Fail()
	}
}

func TestToMap_NilType(t *testing.T) {
	t.Parallel()

	cfg := ToMapConfig{
		PkgName: "main",
	}

	err := ToMap(io.Discard, cfg)
	if err == nil {
		t.Fatal("expected a non-nil error")
	}
}

func TestToMap_MissingPackageName(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	cfg := ToMapConfig{
		Typ: typ,
	}

	err := ToMap(io.Discard, cfg)
	if err == nil {
		t.Fatal("expected to map return an error when it cannot get package name")
	}
}

func TestToMap_NoStructType(t *testing.T) {
	t.Parallel()

	typ := getType(t, nonStructTestTypeName)

	cfg := ToMapConfig{
		PkgName: "main",
		Typ:     typ,
	}

	err := ToMap(io.Discard, cfg)
	if err == nil {
		t.Fatal("expected ToMap return error when called with a non struct type")
	}
}

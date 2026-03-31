package main

import (
	"go/types"
	"strings"
	"testing"
	"text/template"
)

const testStructName = "testStruct"

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

func getTestStructType(t *testing.T) *types.Named {
	t.Helper()

	pkgs, err := loadPackages(t.Context(), ".")
	if err != nil {
		t.Fatal(err)
	}

	_ = testStruct{}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(testStructName)
		if obj == nil {
			continue
		}

		typ, ok := obj.Type().(*types.Named)
		if !ok {
			panic("object type must be named")
		}

		return typ
	}

	t.Fatalf("`%v` not found", testStructName)

	return nil
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

	if err := ToMap(t.Context(), &res, cfg); err != nil {
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

	if err := ToMap(t.Context(), &res, cfg); err != nil {
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

	if err := ToMap(t.Context(), &res, cfg); err != nil {
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

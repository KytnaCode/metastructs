package main

import (
	"go/types"
	"io"
	"strings"
	"testing"
	"text/template"
)

const (
	testStructName          = "testStruct"
	taggedTestStructName    = "taggedTestStruct"
	omitemptyTestStructName = "omitemptyTestStruct"
	nonStructTestTypeName   = "nonStructTestType"
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

var (
	_ = testStruct{}
	_ = taggedTestStruct{}
	_ = omitemptyTestStruct{}
	_ = nonStructTestType("")
)

type testStruct struct {
	Name  string
	Count int
	Slice []float64
}

type taggedTestStruct struct {
	Name  string    `to-map:"name"`
	Count int       `to-map:"count"`
	Slice []float64 `to-map:"-"`
}

type omitemptyTestStruct struct {
	Name  string `to-map:",omitempty"`
	Count int
	Slice []float64 `to-map:"slice,omitempty"`
}

type nonStructTestType string

func init() {
	var err error

	expectedTmpl, err = template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() map[string]any {
	structMap := map[string]any{
		"{{.CountKey}}": {{.Receiver}}.Count,
		"{{.NameKey}}":  {{.Receiver}}.Name,
		"{{.SliceKey}}": {{.Receiver}}.Slice,
	}
	return structMap
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

func TestToMap_Tags(t *testing.T) {
	t.Parallel()

	typ := getType(t, taggedTestStructName)

	const pkgMain = "main"

	cfg := ToMapConfig{
		PkgName: pkgMain,
		Typ:     typ,
	}

	var res strings.Builder

	if err := ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	expectedTmpl, err := template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() map[string]any {
	structMap := map[string]any{
		"{{.CountKey}}": {{.Receiver}}.Count,
		"{{.NameKey}}":  {{.Receiver}}.Name,
	}
	return structMap
}
`)
	if err != nil {
		panic(err)
	}

	var expected strings.Builder

	err = expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgMain,
		Comment:    PackageComment,
		Receiver:   MethodReceiver,
		RecvType:   typ.Obj().Name(),
		MethodName: DefaultMethodName,
		NameKey:    "name",
		CountKey:   "count",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		printDiff(t, expected.String(), res.String())

		t.Fail()
	}
}

func TestToMap_OmitEmptyFields(t *testing.T) {
	t.Parallel()

	typ := getType(t, omitemptyTestStructName)

	const pkgName = "main"

	cfg := ToMapConfig{
		PkgName: pkgName,
		Typ:     typ,
	}

	var res strings.Builder

	if err := ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	expectedTmpl, err := template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() map[string]any {
	structMap := map[string]any{"{{.CountKey}}": {{.Receiver}}.Count}
	var _NameEmpty string
	if {{.Receiver}}.Name != _NameEmpty {
		structMap["{{.NameKey}}"] = {{.Receiver}}.Name
	}
	var _SliceEmpty []float64
	if {{.Receiver}}.Slice != _SliceEmpty {
		structMap["{{.SliceKey}}"] = {{.Receiver}}.Slice
	}
	return structMap
}
`)
	if err != nil {
		panic(err)
	}

	var expected strings.Builder

	err = expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgName,
		Comment:    PackageComment,
		Receiver:   MethodReceiver,
		RecvType:   typ.Obj().Name(),
		MethodName: DefaultMethodName,
		CountKey:   "Count",
		NameKey:    "Name",
		SliceKey:   "slice",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		printDiff(t, expected.String(), res.String())

		t.Fail()
	}
}

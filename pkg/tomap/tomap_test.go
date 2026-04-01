package tomap_test

import (
	"go/types"
	"io"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/kytnacode/metastructs"
	"github.com/kytnacode/metastructs/pkg/tomap"
	"github.com/kytnacode/metastructs/pkg/util/utiltest"
)

var (
	testStructName          = reflect.TypeFor[testStruct]().Name()
	taggedTestStructName    = reflect.TypeFor[taggedTestStruct]().Name()
	omitemptyTestStructName = reflect.TypeFor[omitemptyTestStruct]().Name()
	nonStructTestTypeName   = reflect.TypeFor[nonStructTestType]().Name()
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

func getTestStructType(t *testing.T) *types.Named {
	t.Helper()

	return utiltest.GetType(t, testStructName)
}

func TestToMap_Defaults(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "main"

	cfg := tomap.Config{
		PkgName: pkgName,
		Typ:     typ,
	}

	var res strings.Builder

	if err := tomap.ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	var expected strings.Builder

	err := expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgName,
		Comment:    metastructs.PackageComment,
		Receiver:   metastructs.MethodReceiver,
		RecvType:   typ.Obj().Name(),
		MethodName: tomap.DefaultMethodName,
		NameKey:    "Name",
		CountKey:   "Count",
		SliceKey:   "Slice",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res.String())
		t.Fail()
	}
}

func TestToMap_CustomMethodName(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "abc"

	const methodName = "CustomMethodName"

	cfg := tomap.Config{
		PkgName:    pkgName,
		Typ:        typ,
		MethodName: methodName,
	}

	var res strings.Builder

	if err := tomap.ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	var expected strings.Builder

	err := expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgName,
		Comment:    metastructs.PackageComment,
		Receiver:   metastructs.MethodReceiver,
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
		utiltest.PrintDiff(t, expected.String(), res.String())
		t.Fail()
	}
}

func TestToMap_PointerReceiver(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "models"

	cfg := tomap.Config{
		PkgName: pkgName,
		Typ:     typ,
		Pointer: true,
	}

	var res strings.Builder

	if err := tomap.ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	var expected strings.Builder

	err := expectedTmpl.Execute(&expected, expectedData{
		Package:    pkgName,
		Comment:    metastructs.PackageComment,
		Receiver:   metastructs.MethodReceiver,
		RecvType:   "*" + typ.Obj().Name(),
		MethodName: tomap.DefaultMethodName,
		NameKey:    "Name",
		CountKey:   "Count",
		SliceKey:   "Slice",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res.String())

		t.Fail()
	}
}

func TestToMap_NilType(t *testing.T) {
	t.Parallel()

	cfg := tomap.Config{
		PkgName: "main",
	}

	err := tomap.ToMap(io.Discard, cfg)
	if err == nil {
		t.Fatal("expected a non-nil error")
	}
}

func TestToMap_MissingPackageName(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	cfg := tomap.Config{
		Typ: typ,
	}

	err := tomap.ToMap(io.Discard, cfg)
	if err == nil {
		t.Fatal("expected to map return an error when it cannot get package name")
	}
}

func TestToMap_NoStructType(t *testing.T) {
	t.Parallel()

	typ := utiltest.GetType(t, nonStructTestTypeName)

	cfg := tomap.Config{
		PkgName: "main",
		Typ:     typ,
	}

	err := tomap.ToMap(io.Discard, cfg)
	if err == nil {
		t.Fatal("expected ToMap return error when called with a non struct type")
	}
}

func TestToMap_Tags(t *testing.T) {
	t.Parallel()

	typ := utiltest.GetType(t, taggedTestStructName)

	const pkgMain = "main"

	cfg := tomap.Config{
		PkgName: pkgMain,
		Typ:     typ,
	}

	var res strings.Builder

	if err := tomap.ToMap(&res, cfg); err != nil {
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
		Comment:    metastructs.PackageComment,
		Receiver:   metastructs.MethodReceiver,
		RecvType:   typ.Obj().Name(),
		MethodName: tomap.DefaultMethodName,
		NameKey:    "name",
		CountKey:   "count",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res.String())

		t.Fail()
	}
}

func TestToMap_OmitEmptyFields(t *testing.T) {
	t.Parallel()

	typ := utiltest.GetType(t, omitemptyTestStructName)

	const pkgName = "main"

	cfg := tomap.Config{
		PkgName: pkgName,
		Typ:     typ,
	}

	var res strings.Builder

	if err := tomap.ToMap(&res, cfg); err != nil {
		t.Fatal(err)
	}

	expectedTmpl, err := template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() map[string]any {
	structMap := map[string]any{"{{.CountKey}}": {{.Receiver}}.Count}
	if {{.Receiver}}.Name != "" {
		structMap["{{.NameKey}}"] = {{.Receiver}}.Name
	}
	if len({{.Receiver}}.Slice) > 0 {
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
		Comment:    metastructs.PackageComment,
		Receiver:   metastructs.MethodReceiver,
		RecvType:   typ.Obj().Name(),
		MethodName: tomap.DefaultMethodName,
		CountKey:   "Count",
		NameKey:    "Name",
		SliceKey:   "slice",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.String() != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res.String())

		t.Fail()
	}
}

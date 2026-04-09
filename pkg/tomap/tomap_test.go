package tomap_test

import (
	"go/types"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/kytnacode/metastructs"
	"github.com/kytnacode/metastructs/pkg/tomap"
	"github.com/kytnacode/metastructs/pkg/util"
	"github.com/kytnacode/metastructs/pkg/util/utiltest"
	"github.com/kytnacode/metastructs/testdata"
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
	structMap := make(map[string]any, 3)
	structMap["{{.CountKey}}"] = {{.Receiver}}.Count
	structMap["{{.NameKey}}"] = {{.Receiver}}.Name
	structMap["{{.SliceKey}}"] = {{.Receiver}}.Slice
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
		Typ: typ,
	}

	f := util.NewFile(pkgName)

	if err := tomap.ToMap(f, cfg); err != nil {
		t.Fatal(err)
	}

	res := f.GoString()

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

	if res != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res)
		t.Fail()
	}
}

func TestToMap_CustomMethodName(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "abc"

	const methodName = "CustomMethodName"

	cfg := tomap.Config{
		Typ:        typ,
		MethodName: methodName,
	}

	f := util.NewFile(pkgName)

	if err := tomap.ToMap(f, cfg); err != nil {
		t.Fatal(err)
	}

	res := f.GoString()

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

	if res != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res)
		t.Fail()
	}
}

func TestToMap_PointerReceiver(t *testing.T) {
	t.Parallel()

	typ := getTestStructType(t)

	const pkgName = "models"

	cfg := tomap.Config{
		Typ:     typ,
		Pointer: true,
	}

	f := util.NewFile(pkgName)

	if err := tomap.ToMap(f, cfg); err != nil {
		t.Fatal(err)
	}

	res := f.GoString()

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

	if res != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res)

		t.Fail()
	}
}

func TestToMap_NilType(t *testing.T) {
	t.Parallel()

	f := util.NewFile("main")

	err := tomap.ToMap(f, tomap.Config{})
	if err == nil {
		t.Fatal("expected a non-nil error")
	}
}

func TestToMap_NoStructType(t *testing.T) {
	t.Parallel()

	typ := utiltest.GetType(t, nonStructTestTypeName)

	cfg := tomap.Config{
		Typ: typ,
	}

	f := util.NewFile("main")

	err := tomap.ToMap(f, cfg)
	if err == nil {
		t.Fatal("expected ToMap return error when called with a non struct type")
	}
}

func TestToMap_Tags(t *testing.T) {
	t.Parallel()

	typ := utiltest.GetType(t, taggedTestStructName)

	const pkgMain = "main"

	cfg := tomap.Config{
		Typ: typ,
	}

	f := util.NewFile(pkgMain)

	if err := tomap.ToMap(f, cfg); err != nil {
		t.Fatal(err)
	}

	res := f.GoString()

	expectedTmpl, err := template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() map[string]any {
	structMap := make(map[string]any, 2)
	structMap["{{.CountKey}}"] = {{.Receiver}}.Count
	structMap["{{.NameKey}}"] = {{.Receiver}}.Name
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

	if res != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res)

		t.Fail()
	}
}

func TestToMap_OmitEmptyFields(t *testing.T) {
	t.Parallel()

	typ := utiltest.GetType(t, omitemptyTestStructName)

	const pkgName = "main"

	cfg := tomap.Config{
		Typ: typ,
	}

	f := util.NewFile(pkgName)

	if err := tomap.ToMap(f, cfg); err != nil {
		t.Fatal(err)
	}

	res := f.GoString()

	expectedTmpl, err := template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() map[string]any {
	structMap := make(map[string]any, 3)
	structMap["{{.CountKey}}"] = {{.Receiver}}.Count
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

	if res != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res)

		t.Fail()
	}
}

func BenchmarkToMap(b *testing.B) {
	data := testdata.ToMapBench{
		Str1:  "hello world",
		Str2:  "My name is Arya!",
		Str3:  "What's yours?",
		Bool:  true,
		Float: 3.1416,
		Int:   8191, // Fifth Mersenne Prime.
		Uint:  57,   // The best prime number.
		// 09-F9 code.
		Slice: []string{"09", "F9", "11", "02", "9D", "74", "E3", "5B", "D8", "41", "56", "C5", "63", "56", "88", "C0"},
		Map: map[float64]string{
			2.717:  "Bob",
			3.1416: "Alice",
			1.414:  "Eva",
		},
		AnonymousStruct: struct {
			Field1 string
			Field2 float64
		}{
			Field1: "Hi, my hame is Arya!",
			Field2: 1 + 1/(1+(1/(1+1))),
		},
	}

	b.ResetTimer()

	b.Run("ToMap_GeneratedMethod", func(b *testing.B) {
		for b.Loop() {
			_ = data.ToMap()
		}
	})

	b.Run("ToMap_Reflection", func(b *testing.B) {
		for b.Loop() {
			v := reflect.ValueOf(data)

			out := make(map[string]any, v.NumField())

			for i := range v.NumField() {
				out[v.Type().Field(i).Name] = v.Field(i).Interface()
			}

			_ = out
		}
	})
}

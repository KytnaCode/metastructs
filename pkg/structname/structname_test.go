package structname_test

import (
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/kytnacode/metastructs"
	"github.com/kytnacode/metastructs/pkg/structname"
	"github.com/kytnacode/metastructs/pkg/util/utiltest"
)

var (
	testStructName = reflect.TypeFor[testStruct]().Name()
	testIntName    = reflect.TypeFor[testInt]().Name()
)

type testStruct struct{}

type testInt int

var (
	_ = testStruct{}
	_ = testInt(0)
)

var expectedTmpl *template.Template

func init() {
	var err error

	expectedTmpl, err = template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

func ({{.Receiver}} {{.RecvType}}) {{.MethodName}}() string {
	return "{{.Name}}"
}
`)
	if err != nil {
		panic(err)
	}
}

func TestName(t *testing.T) {
	t.Parallel()

	type testData struct {
		pkg        string
		pointer    bool
		recvType   string
		methodName string
		name       string
	}

	testCases := map[string]testData{
		"pkg-main_no-pointer_def-method_struct": {
			pkg:      "main",
			pointer:  false,
			recvType: testStructName,
			name:     testStructName,
		},
		"pkg-abc_no-pointer_def-method_struct": {
			pkg:      "abc",
			pointer:  false,
			recvType: testStructName,
			name:     testStructName,
		},
		"pkg-main_pointer_def-method_struct": {
			pkg:      "main",
			pointer:  true,
			recvType: testStructName,
			name:     testStructName,
		},
		"pkg-main_no-pointer_custom-method_int": {
			pkg:        "main",
			pointer:    false,
			recvType:   testIntName,
			name:       testIntName,
			methodName: "IntName",
		},
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			named := utiltest.GetType(t, data.name)

			cfg := structname.Config{
				Typ:        named,
				PkgName:    data.pkg,
				MethodName: data.methodName,
				Pointer:    data.pointer,
			}

			var res strings.Builder

			if err := structname.StructName(&res, cfg); err != nil {
				t.Fatal(err)
			}

			var expected strings.Builder

			var recvType string

			if data.pointer {
				recvType += "*"
			}

			methodName := structname.DefaultMethodName
			if data.methodName != "" {
				methodName = data.methodName
			}

			recvType += data.recvType

			err := expectedTmpl.Execute(&expected, map[string]any{
				"Name":       data.name,
				"Comment":    metastructs.PackageComment,
				"Receiver":   metastructs.MethodReceiver,
				"MethodName": methodName,
				"RecvType":   recvType,
				"Package":    data.pkg,
			})
			if err != nil {
				t.Fatal(err)
			}

			if res.String() != expected.String() {
				utiltest.PrintDiff(t, expected.String(), res.String())

				t.FailNow()
			}
		})
	}
}

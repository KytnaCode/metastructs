package main

import (
	"strings"
	"testing"
	"text/template"
)

const (
	testNameStructName = "testNameStruct"
	testNameIntName    = "testNameInt"
)

type testNameStruct struct{}

type testNameInt int

var (
	_ = testNameStruct{}
	_ = testNameInt(0)
)

var expectedNameTmpl *template.Template

func init() {
	var err error

	expectedNameTmpl, err = template.New("expected").Parse(`// {{.Comment}}
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
			recvType: testNameStructName,
			name:     testNameStructName,
		},
		"pkg-abc_no-pointer_def-method_struct": {
			pkg:      "abc",
			pointer:  false,
			recvType: testNameStructName,
			name:     testNameStructName,
		},
		"pkg-main_pointer_def-method_struct": {
			pkg:      "main",
			pointer:  true,
			recvType: testNameStructName,
			name:     testNameStructName,
		},
		"pkg-main_no-pointer_custom-method_int": {
			pkg:        "main",
			pointer:    false,
			recvType:   testNameIntName,
			name:       testNameIntName,
			methodName: "IntName",
		},
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			named := getType(t, data.name)

			cfg := StructNameConfig{
				Typ:        named,
				PkgName:    data.pkg,
				MethodName: data.methodName,
				Pointer:    data.pointer,
			}

			var res strings.Builder

			if err := StructName(&res, cfg); err != nil {
				t.Fatal(err)
			}

			var expected strings.Builder

			var recvType string

			if data.pointer {
				recvType += "*"
			}

			methodName := DefaultNameMethodName
			if data.methodName != "" {
				methodName = data.methodName
			}

			recvType += data.recvType

			err := expectedNameTmpl.Execute(&expected, map[string]any{
				"Name":       data.name,
				"Comment":    PackageComment,
				"Receiver":   MethodReceiver,
				"MethodName": methodName,
				"RecvType":   recvType,
				"Package":    data.pkg,
			})
			if err != nil {
				t.Fatal(err)
			}

			if res.String() != expected.String() {
				printDiff(t, expected.String(), res.String())

				t.FailNow()
			}
		})
	}
}

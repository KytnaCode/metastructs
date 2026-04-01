package main

import (
	"strings"
	"testing"
	"text/template"
	"time"
	"unicode/utf8"
)

const testPartialUserName = "testPartialUser"

var expectedPartialTmpl *template.Template

var _ = testPartialUser{}

type testType string

type testPartialUser struct {
	Name       string
	Pass       string
	Registered time.Time
	ExtraData  testType
}

func init() {
	var err error

	expectedPartialTmpl, err = template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

import "time"

type {{.StructName}} struct {
	Name       *string
	Pass       *string
	Registered *time.Time
	ExtraData  *testType
}
`)
	if err != nil {
		panic(err)
	}
}

func TestPartial(t *testing.T) {
	t.Parallel()

	type testData struct {
		pkgName    string
		structName string
		prefix     string
		suffix     string
	}

	testCases := map[string]testData{
		"pkg-main_def-struct-name_no-prefix_def-suffix": {
			pkgName: "main",
		},
		"pkg-main_def-struct-name_no-prefix_custom-suffix": {
			pkgName: "main",
			suffix:  "Merge",
		},
		"pkg-other_def-struct-name_with-prefix_no-suffix": {
			pkgName: "other",
			suffix:  "",
			prefix:  "Patch",
		},
		"pkg-other_def-struct-name_with-prefix_custom_suffix": {
			pkgName: "other",
			suffix:  "Dto",
			prefix:  "Partial",
		},
		"pkg-other_custom-struct-name": {
			pkgName:    "other",
			structName: "CustomTypeName",
		},
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			typ := getType(t, testPartialUserName)

			cfg := PartialConfig{
				Typ:        typ,
				PkgName:    data.pkgName,
				StructName: data.structName,
				Suffix:     &data.prefix,
				Prefix:     data.prefix,
			}

			var res strings.Builder

			if err := Partial(&res, cfg); err != nil {
				t.Fatal(err)
			}

			var expected strings.Builder

			structName := data.structName
			if structName == "" {
				name := cfg.Typ.Obj().Name()
				r, size := utf8.DecodeRuneInString(name)

				structName = cfg.Prefix + strings.ToUpper(string(r)) + name[size:] + *cfg.Suffix
			}

			err := expectedPartialTmpl.Execute(&expected, map[string]any{
				"Comment":    PackageComment,
				"Package":    data.pkgName,
				"StructName": structName,
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

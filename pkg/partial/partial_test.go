package partial_test

import (
	"go/types"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/kytnacode/metastructs"
	"github.com/kytnacode/metastructs/pkg/partial"
	"github.com/kytnacode/metastructs/pkg/util"
	"github.com/kytnacode/metastructs/pkg/util/utiltest"
)

var testStructName = reflect.TypeFor[testUser]().Name()

var (
	expectedTmpl       *template.Template
	expectedTmplTagged *template.Template
)

var _ = testUser{}

type testType string

type testUser struct {
	Name        string       `db:"name,omitempty" json:"name"`
	Pass        string       `db:"pass"           json:"pass"`
	Registered  time.Time    `db:"-"              json:"registered"`
	ExtraData   testType     `db:",omitempty"     json:",omitempty"`
	PointerType *types.Named `db:"ptr"            json:"ptr"`
}

func init() {
	var err error

	expectedTmpl, err = template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

import (
	"go/types"
	"time"
)

type {{.StructName}} struct {
	Name        *string
	Pass        *string
	Registered  *time.Time
	ExtraData   *testType
	PointerType *types.Named
}
`)
	if err != nil {
		panic(err)
	}

	expectedTmplTagged, err = template.New("expected").Parse(`// {{.Comment}}
package {{.Package}}

import (
	"go/types"
	"time"
)

type {{.StructName}} struct {
	Name        *string      {{.NameTag}}
	Pass        *string      {{.PassTag}}
	Registered  *time.Time   {{.RegisteredTag}}
	ExtraData   *testType    {{.ExtraTag}}
	PointerType *types.Named {{.PtrTag}}
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
			typ := utiltest.GetType(t, testStructName)

			cfg := partial.Config{
				Typ:        typ,
				StructName: data.structName,
				Suffix:     data.suffix,
				Prefix:     &data.prefix,
			}

			f := util.NewFile(data.pkgName)

			if err := partial.Partial(f, cfg); err != nil {
				t.Fatal(err)
			}

			res := f.GoString()

			var expected strings.Builder

			structName := data.structName
			if structName == "" {
				name := cfg.Typ.Obj().Name()
				r, size := utf8.DecodeRuneInString(name)

				structName = *cfg.Prefix + strings.ToUpper(string(r)) + name[size:] + cfg.Suffix
			}

			err := expectedTmpl.Execute(&expected, map[string]any{
				"Comment":    metastructs.PackageComment,
				"Package":    data.pkgName,
				"StructName": structName,
			})
			if err != nil {
				t.Fatal(err)
			}

			if res != expected.String() {
				utiltest.PrintDiff(t, expected.String(), res)

				t.FailNow()
			}
		})
	}
}

func TestPartial_Tagget(t *testing.T) {
	t.Parallel()

	typ := utiltest.GetType(t, testStructName)

	const pkgName = "main"

	const structName = "MyStructName"

	cfg := partial.Config{
		Typ:          typ,
		StructName:   structName,
		PreserveTags: true,
	}

	f := util.NewFile(pkgName)

	if err := partial.Partial(f, cfg); err != nil {
		t.Fatal(err)
	}

	res := f.GoString()

	var expected strings.Builder

	//nolint:gosec // No hardcoded credentials, PassTag is just dummy data.
	err := expectedTmplTagged.Execute(&expected, map[string]any{
		"Comment":       metastructs.PackageComment,
		"Package":       pkgName,
		"StructName":    structName,
		"NameTag":       "`db:\"name,omitempty\" json:\"name\"`",
		"PassTag":       "`db:\"pass\"           json:\"pass\"`",
		"RegisteredTag": "`db:\"-\"              json:\"registered\"`",
		"ExtraTag":      "`db:\",omitempty\"     json:\",omitempty\"`",
		"PtrTag":        "`db:\"ptr\"            json:\"ptr\"`",
	})
	if err != nil {
		t.Fatal(err)
	}

	if res != expected.String() {
		utiltest.PrintDiff(t, expected.String(), res)

		t.FailNow()
	}
}

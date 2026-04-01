// Package util contains a bunch of utility functions.
package util

import (
	"context"
	"errors"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/kytnacode/metastructs"
	"golang.org/x/tools/go/packages"
)

// FieldData contains metadata for struct fields.
type FieldData struct {
	Name string
	Typ  types.Type
	Tag  reflect.StructTag
}

// NewFile creates a new [jen.File] with [metastructs.PackageComment].
func NewFile(pkgName string) *jen.File {
	f := jen.NewFile(pkgName)

	f.PackageComment(metastructs.PackageComment)

	return f
}

// LoadPackages load packages given a list of paths.
func LoadPackages(ctx context.Context, paths ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:    packages.NeedTypes | packages.NeedImports | packages.NeedTypesInfo,
		Context: ctx,
		Tests:   true,
	}

	return packages.Load(cfg, paths...)
}

// GetStructFields returns fields metadata for a given struct.
func GetStructFields(structType *types.Struct) []FieldData {
	fields := make([]FieldData, 0, structType.NumFields())

	for i := range structType.NumFields() {
		field := structType.Field(i)
		tag := reflect.StructTag(structType.Tag(i))

		data := FieldData{
			Name: field.Name(),
			Typ:  field.Type(),
			Tag:  tag,
		}

		fields = append(fields, data)
	}

	return fields
}

// LoadType loads a single type, and return it together with a bool that indicates if type is part of a test package.
func LoadType(ctx context.Context, pkg string, typ string) (*types.Named, bool, error) {
	pkgs, err := LoadPackages(ctx, pkg)
	if err != nil {
		return nil, false, err
	}

	if len(pkgs) == 0 {
		return nil, false, errors.New("package not found")
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(typ)
		if obj == nil {
			continue
		}

		if named, ok := obj.Type().(*types.Named); ok {
			return named, pkg.ForTest != "", nil
		}

		return nil, false, fmt.Errorf("`%v` is not a struct", obj.Name())
	}

	return nil, false, errors.New("type not found")
}

// Filename creates a filename with a given suffix (without underscores). If `test` is true an additional _test will
// be added to filename. The filename will be of the form <lowercase_struct_name>_<suffix>.go or, if test is true
// <lowercase_struct_name>_<suffix>_test.go.
func Filename(strct, suffix string, test bool) string {
	var testSuffix string
	if test {
		testSuffix = "_test"
	}

	path := fmt.Sprintf("%v_%v%v.go", strings.ToLower(strct), suffix, testSuffix)

	return filepath.Join("./", filepath.Clean(path))
}

// GetPkgName returns value of GOPACKAGE environmente variable if `pkg` is equal to `.`, otherwise it's returned as is.
func GetPkgName(pkg string) string {
	p := pkg
	if p == "." {
		return os.Getenv("GOPACKAGE")
	}

	return p
}

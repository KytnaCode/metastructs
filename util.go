package main

import (
	"context"
	"errors"
	"fmt"
	"go/types"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

type fieldData struct {
	name string
	typ  types.Type
	tag  reflect.StructTag
}

func newFile(pkgName string) *jen.File {
	f := jen.NewFile(pkgName)

	f.PackageComment(PackageComment)

	return f
}

func loadPackages(ctx context.Context, paths ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:    packages.NeedTypes | packages.NeedImports | packages.NeedTypesInfo,
		Context: ctx,
		Tests:   true,
	}

	return packages.Load(cfg, paths...)
}

func getStructFields(structType *types.Struct) []fieldData {
	fields := make([]fieldData, 0, structType.NumFields())

	for i := range structType.NumFields() {
		field := structType.Field(i)
		tag := reflect.StructTag(structType.Tag(i))

		data := fieldData{
			name: field.Name(),
			typ:  field.Type(),
			tag:  tag,
		}

		fields = append(fields, data)
	}

	return fields
}

func loadType(ctx context.Context, pkg string, typ string) (*types.Named, bool, error) {
	pkgs, err := loadPackages(ctx, pkg)
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

func filename(strct, suffix string, test bool) string {
	var testSuffix string
	if test {
		testSuffix = "_test"
	}

	path := fmt.Sprintf("%v_%v%v.go", strings.ToLower(strct), suffix, testSuffix)

	return filepath.Join("./", filepath.Clean(path))
}

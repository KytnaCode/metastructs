package main

import (
	"context"
	"go/types"
	"reflect"

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

	f.PackageComment(PACKAGE_COMMENT)

	return f
}

func loadPackages(ctx context.Context, paths ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:    packages.NeedTypes | packages.NeedImports,
		Context: ctx,
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

package main

import (
	"context"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

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

package main

import (
	"github.com/dave/jennifer/jen"
)

func NewFile(pkgName string) *jen.File {
	f := jen.NewFile(pkgName)

	f.PackageComment(PACKAGE_COMMENT)

	return f
}

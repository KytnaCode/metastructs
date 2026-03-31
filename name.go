package main

import (
	"errors"
	"go/types"
	"io"

	"github.com/dave/jennifer/jen"
)

const DefaultNameMethodName = "Name"

type StructNameConfig struct {
	Typ        *types.Named
	PkgName    string
	MethodName string
	Pointer    bool
}

func StructName(w io.Writer, cfg StructNameConfig) error {
	if cfg.PkgName == "" {
		return errors.New("cfg.PkgName is required ")
	}

	if cfg.MethodName == "" {
		cfg.MethodName = DefaultNameMethodName
	}

	f := newFile(cfg.PkgName)

	recv := jen.Id(MethodReceiver)

	if cfg.Pointer {
		recv.Op("*")
	}

	recv.Id(cfg.Typ.Obj().Name())

	f.Func().Params(recv).Id(cfg.MethodName).Params().String().Block(
		jen.Return(jen.Lit(cfg.Typ.Obj().Name())),
	)

	return f.Render(w)
}

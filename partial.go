package main

import (
	"errors"
	"fmt"
	"go/types"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/dave/jennifer/jen"
)

const DefaultPartialSuffix = "Partial"

type PartialConfig struct {
	Typ        *types.Named
	PkgName    string
	StructName string
	Suffix     *string
	Prefix     string
}

func Partial(w io.Writer, cfg PartialConfig) error {
	if cfg.PkgName == "" {
		return errors.New("cfg.PkgName is required")
	}

	if cfg.Typ == nil {
		return errors.New("cfg.Typ is required")
	}

	if cfg.Suffix == nil {
		cfg.Suffix = new(string)
		*cfg.Suffix = DefaultPartialSuffix
	}

	if cfg.StructName == "" {
		name := cfg.Typ.Obj().Name()
		r, size := utf8.DecodeRuneInString(name)

		cfg.StructName = cfg.Prefix + strings.ToUpper(string(r)) + name[size:] + *cfg.Suffix
	}

	if cfg.StructName == cfg.Typ.Obj().Name() {
		return errors.New("output struct name is the same as the target, try set cfg>Suffix, cfg.Prefix or cfg.StructName")
	}

	structType, ok := cfg.Typ.Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("`%v` is not a struct", cfg.Typ.Obj().Name())
	}

	fields := make([]jen.Code, 0, structType.NumFields())

	for field := range structType.Fields() {
		structField := jen.Id(field.Name()).Op("*")

		if named, ok := field.Type().(*types.Named); ok {
			if named.Obj().Pkg().Path() == cfg.Typ.Obj().Pkg().Path() {
				structField.Id(named.Obj().Name())
			} else {
				structField.Qual(named.Obj().Pkg().Path(), named.Obj().Name())
			}
		} else {
			structField.Id(field.Type().String())
		}

		fields = append(fields, structField)
	}

	f := newFile(cfg.PkgName)

	f.Type().Id(cfg.StructName).Struct(fields...)

	return f.Render(w)
}

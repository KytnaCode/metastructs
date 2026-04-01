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

const DefaultPartialPrefix = "Partial"

type PartialConfig struct {
	Typ          *types.Named
	PkgName      string
	StructName   string
	Suffix       string
	Prefix       *string
	PreserveTags bool
}

func Partial(w io.Writer, cfg PartialConfig) error {
	if cfg.PkgName == "" {
		return errors.New("cfg.PkgName is required")
	}

	if cfg.Typ == nil {
		return errors.New("cfg.Typ is required")
	}

	if cfg.Prefix == nil {
		cfg.Prefix = new(string)
		*cfg.Prefix = DefaultPartialPrefix
	}

	if cfg.StructName == "" {
		name := cfg.Typ.Obj().Name()
		r, size := utf8.DecodeRuneInString(name)

		cfg.StructName = *cfg.Prefix + strings.ToUpper(string(r)) + name[size:] + cfg.Suffix
	}

	if cfg.StructName == cfg.Typ.Obj().Name() {
		return errors.New("output struct name is the same as the target, try set cfg>Suffix, cfg.Prefix or cfg.StructName")
	}

	structType, ok := cfg.Typ.Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("`%v` is not a struct", cfg.Typ.Obj().Name())
	}

	fields := make([]jen.Code, 0, structType.NumFields())

	for i := range structType.NumFields() {
		field := structType.Field(i)
		structField := jen.Id(field.Name()).Op("*")

		typ := field.Type()

		for {
			if named, ok := typ.(*types.Named); ok {
				if named.Obj().Pkg().Path() == cfg.Typ.Obj().Pkg().Path() {
					structField.Id(named.Obj().Name())
				} else {
					structField.Qual(named.Obj().Pkg().Path(), named.Obj().Name())
				}

				break
			} else if basic, ok := typ.(*types.Basic); ok {
				structField.Id(basic.String())

				break
			} else if pointer, ok := typ.(*types.Pointer); ok {
				typ = pointer.Elem()

				continue
			}

			return fmt.Errorf("unsupported type: `%v`", typ.String())
		}

		if cfg.PreserveTags {
			structField.Op("`" + structType.Tag(i) + "`")
		}

		fields = append(fields, structField)
	}

	f := newFile(cfg.PkgName)

	f.Type().Id(cfg.StructName).Struct(fields...)

	return f.Render(w)
}

package main

import (
	"errors"
	"fmt"
	"go/types"
	"io"

	"github.com/dave/jennifer/jen"
)

const DefaultMethodName = "ToMap"

// ToMapConfig is the configuration struct for [ToMap].
type ToMapConfig struct {
	// Output's method name, defaults to [DefaultMethodName].
	MethodName string

	// Generated file's package.
	PkgName string

	// Target type, must be a struct.
	Typ *types.Named

	// Should receiver be a pointer to Typ.
	Pointer bool
}

// ToMap generates a file defining a method to convert a struct given a config into a map and writes the content
// into `w`.
func ToMap(w io.Writer, cfg ToMapConfig) error {
	if cfg.MethodName == "" {
		cfg.MethodName = DefaultMethodName
	}

	if cfg.Typ == nil {
		return errors.New("type must not be nil")
	}

	structType, ok := cfg.Typ.Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("type %v is not a struct", cfg.Typ.Obj().Name())
	}

	fields := getStructFields(structType)

	values := make(map[jen.Code]jen.Code, len(fields)) // Map values.

	for _, field := range fields {
		values[jen.Lit(field.name)] = jen.Id(MethodReceiver).Dot(field.name)
	}

	f := newFile(cfg.PkgName)

	recvType := jen.Id(cfg.Typ.Obj().Name())

	if cfg.Pointer {
		recvType = jen.Op("*").Add(recvType)
	}

	f.Func().Params(jen.Id(MethodReceiver).Add(recvType)).Id(cfg.MethodName).Params().Map(jen.String()).Any().Block(
		jen.Return(
			jen.Map(jen.String()).Any().Values(
				jen.Dict(values),
			),
		),
	)

	return f.Render(w)
}

package main

import (
	"errors"
	"fmt"
	"go/types"
	"io"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
)

const (
	DefaultMethodName = "ToMap"
	DefaultTagName    = "to-map"
)

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

	// TagName is the tag from which ToMap will retrieve field metadata, works like [encoding/json] tags. Default is
	// [DefaultTagName].
	//
	// ```
	// type MyStruct struct {
	//   Name        string                         // map key will be "Name".
	//   Description string   `to-map:"desc"`       // map key will be "desc".
	//   Authors     []string `to-map:",omitempty"` // map key will be "Authors", omitted if empty (per encoding/json omitempty rules).
	//   Completed   bool     `to-map:"-"`          // will not be included in map.
	// }
	// ```
	TagName *string
}

// ToMap generates a file defining a method to convert a struct given a config into a map and writes the content
// into `w`.
//
// cfg.Typ and cfg.PkgName are required, cfg.MethodName defaults to [DefaultMethodName].
func ToMap(w io.Writer, cfg ToMapConfig) error {
	if cfg.MethodName == "" {
		cfg.MethodName = DefaultMethodName
	}

	if cfg.Typ == nil {
		return errors.New("type must not be nil")
	}

	if cfg.PkgName == "" {
		return errors.New("cfg.PackageName is required")
	}

	if cfg.TagName == nil {
		cfg.TagName = new(string)
		*cfg.TagName = DefaultTagName
	}

	structType, ok := cfg.Typ.Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("type %v is not a struct", cfg.Typ.Obj().Name())
	}

	fields := getStructFields(structType)

	values := make(map[jen.Code]jen.Code, len(fields)) // Map values.

	omitemptyFields := make(map[string]fieldData, len(fields))

	for _, field := range fields {
		name := field.name

		tag, ok := field.tag.Lookup(*cfg.TagName)
		if ok {
			if tag == "-" {
				continue
			}

			parts := strings.Split(tag, ",")

			if parts[0] != "" {
				name = parts[0]
			}

			if len(parts) > 1 && parts[1] == "omitempty" {
				omitemptyFields[name] = field

				continue
			}
		}

		values[jen.Lit(name)] = jen.Id(MethodReceiver).Dot(field.name)
	}

	const mapID = "structMap"

	stmts := make([]jen.Code, 0, len(omitemptyFields)*2+2)

	stmts = append(stmts,
		jen.Id(mapID).Op(":=").Map(jen.String()).Any().Values(jen.Dict(values)),
	)

	keys := make([]string, 0, len(omitemptyFields))
	for k := range omitemptyFields {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		field := omitemptyFields[key]

		emptyID := "_" + field.name + "Empty"

		stmt1 := jen.Var().Id(emptyID).Id(field.typ.String())
		stmt2 := jen.If(jen.Id(MethodReceiver).Dot(field.name).Op("!=").Id(emptyID)).Block(
			jen.Id(mapID).Index(jen.Lit(key)).Op("=").Id(MethodReceiver).Dot(field.name),
		)

		stmts = append(stmts, stmt1, stmt2)
	}

	stmts = append(stmts, jen.Return(jen.Id(mapID)))

	f := newFile(cfg.PkgName)

	recvType := jen.Id(cfg.Typ.Obj().Name())

	if cfg.Pointer {
		recvType = jen.Op("*").Add(recvType)
	}

	f.Func().Params(jen.Id(MethodReceiver).Add(recvType)).Id(cfg.MethodName).Params().Map(jen.String()).Any().Block(
		stmts...,
	)

	return f.Render(w)
}

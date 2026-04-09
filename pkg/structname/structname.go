// Package structname contains generator for the name method.
package structname

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/kytnacode/metastructs"
)

// DefaultMethodName is used if no one is specified in [Config].
const DefaultMethodName = "StructName"

// Config is the config for [StructName].
type Config struct {
	Typ        *types.Named
	MethodName string
	Pointer    bool
}

// StructName is the generator for the StructName method.
func StructName(f *jen.File, cfg Config) error {
	if cfg.MethodName == "" {
		cfg.MethodName = DefaultMethodName
	}

	recv := jen.Id(metastructs.MethodReceiver)

	if cfg.Pointer {
		recv.Op("*")
	}

	recv.Id(cfg.Typ.Obj().Name())

	f.Func().Params(recv).Id(cfg.MethodName).Params().String().Block(
		jen.Return(jen.Lit(cfg.Typ.Obj().Name())),
	)

	return nil
}

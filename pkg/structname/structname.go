// Package structname contains generator for the name method.
package structname

import (
	"errors"
	"go/types"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/kytnacode/metastructs"
	"github.com/kytnacode/metastructs/pkg/util"
)

// DefaultMethodName is used if no one is specified in [Config].
const DefaultMethodName = "StructName"

// Config is the config for [StructName].
type Config struct {
	Typ        *types.Named
	PkgName    string
	MethodName string
	Pointer    bool
}

// StructName is the generator for the StructName method.
func StructName(w io.Writer, cfg Config) error {
	if cfg.PkgName == "" {
		return errors.New("cfg.PkgName is required ")
	}

	if cfg.MethodName == "" {
		cfg.MethodName = DefaultMethodName
	}

	f := util.NewFile(cfg.PkgName)

	recv := jen.Id(metastructs.MethodReceiver)

	if cfg.Pointer {
		recv.Op("*")
	}

	recv.Id(cfg.Typ.Obj().Name())

	f.Func().Params(recv).Id(cfg.MethodName).Params().String().Block(
		jen.Return(jen.Lit(cfg.Typ.Obj().Name())),
	)

	return f.Render(w)
}

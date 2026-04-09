package main

import (
	"path/filepath"
	"strings"

	"github.com/kytnacode/metastructs/pkg/structname"
	"github.com/kytnacode/metastructs/pkg/util"
	"github.com/spf13/cobra"
)

var structNameMethodName string

var nameCmd = &cobra.Command{
	Use:   "structname",
	Short: "generate a method that returns struct's name",
	RunE: func(cmd *cobra.Command, _ []string) error {
		sourceType, pointer := strings.CutPrefix(target, "*")

		typ, test, err := util.LoadType(cmd.Context(), pkgName, sourceType)
		if err != nil {
			return err
		}

		cfg := structname.Config{
			Typ:        typ,
			MethodName: structNameMethodName,
			Pointer:    pointer,
		}

		path := filepath.Clean(util.Filename(sourceType, "name", test))

		f := util.NewFile(pkgName)

		if err := structname.StructName(f, cfg); err != nil {
			return err
		}

		return f.Save(path)
	},
}

func init() {
	nameCmd.Flags().StringVarP(&structNameMethodName, "method", "m", structname.DefaultMethodName, "out method name")
}

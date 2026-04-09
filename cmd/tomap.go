package main

import (
	"path/filepath"
	"strings"

	"github.com/kytnacode/metastructs/pkg/tomap"
	"github.com/kytnacode/metastructs/pkg/util"
	"github.com/spf13/cobra"
)

var (
	toMapMethodName string
	toMapTag        string
)

var toMapCmd = &cobra.Command{
	Use:   "to-map",
	Short: "create a method to convert a struct into a map",
	RunE: func(cmd *cobra.Command, _ []string) error {
		sourceType, pointer := strings.CutPrefix(target, "*")

		typ, test, err := util.LoadType(cmd.Context(), pkgName, sourceType)
		if err != nil {
			return err
		}

		cfg := tomap.Config{
			MethodName: toMapMethodName,
			TagName:    &toMapTag,
			Typ:        typ,
			Pointer:    pointer,
		}

		path := filepath.Clean(util.Filename(sourceType, "map", test))

		f := util.NewFile(util.GetPkgName(pkgName))

		if err := tomap.ToMap(f, cfg); err != nil {
			return err
		}

		return f.Save(path)
	},
}

func init() {
	toMapCmd.Flags().StringVarP(&toMapMethodName, "method", "m", tomap.DefaultMethodName, "out method name")
	toMapCmd.Flags().StringVarP(&toMapTag, "tag", "s", tomap.DefaultTag, "tag from which read metadata")
}

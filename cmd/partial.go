package main

import (
	"path/filepath"

	"github.com/kytnacode/metastructs/pkg/partial"
	"github.com/kytnacode/metastructs/pkg/util"
	"github.com/spf13/cobra"
)

var (
	partialPrefix, partialSuffix string
	partialStructName            string
	partialPreserveTags          bool
	partialForcePointer          bool
)

var partialCmd = &cobra.Command{
	Use:   "partial",
	Short: "generate a partial struct",
	RunE: func(cmd *cobra.Command, _ []string) error {
		typ, test, err := util.LoadType(cmd.Context(), pkgName, target)
		if err != nil {
			return err
		}

		cfg := partial.Config{
			Typ:          typ,
			StructName:   partialStructName,
			Suffix:       partialSuffix,
			Prefix:       &partialPrefix,
			PreserveTags: partialPreserveTags,
			ForcePointer: partialForcePointer,
		}

		path := filepath.Clean(util.Filename(typ.Obj().Name(), "partial", test))

		f := util.NewFile(util.GetPkgName(pkgName))

		if err := partial.Partial(f, cfg); err != nil {
			return err
		}

		return f.Save(path)
	},
}

func init() {
	partialCmd.Flags().StringVarP(&partialPrefix, "prefix", "e", partial.DefaultPrefix,
		"struct name prefix, no effect when used with --structname")
	partialCmd.Flags().StringVarP(&partialSuffix, "suffix", "s", "",
		"struct name suffix, no effect when used with --structname")
	partialCmd.Flags().StringVarP(&partialStructName, "structname", "n", "",
		"struct name, if set suffix and prefix options will be ignored")
	partialCmd.Flags().BoolVarP(&partialPreserveTags, "preserve-tags", "r", false, "should partial struct preserve tags")
	partialCmd.Flags().BoolVarP(&partialForcePointer, "force-pointer", "f", false,
		"force using a pointer for already nilable types (pointers and interfaces)")
}

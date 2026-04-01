package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	target  string
	pkgName string
)

var (
	toMapMethodName string
	toMapTag        string
)

var nameMethodName string

var (
	partialPrefix, partialSuffix string
	partialStructName            string
)

var rootCmd = &cobra.Command{
	Use:   "metastructs",
	Short: "metastructs group boilerplate generators for structs",
}

var toMapCmd = &cobra.Command{
	Use:   "to-map",
	Short: "create a method to convert a struct into a map",
	RunE: func(cmd *cobra.Command, _ []string) error {
		sourceType, pointer := strings.CutPrefix(target, "*")

		typ, test, err := loadType(cmd.Context(), pkgName, sourceType)
		if err != nil {
			return err
		}

		cfg := ToMapConfig{
			PkgName:    pkgName,
			MethodName: toMapMethodName,
			TagName:    &toMapTag,
			Typ:        typ,
			Pointer:    pointer,
		}

		file := filepath.Clean(filename(sourceType, "map", test))

		f, err := os.Create(file)
		if err != nil {
			return err
		}

		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()

		if err := ToMap(f, cfg); err != nil {
			return err
		}

		return f.Sync()
	},
}

var nameCmd = &cobra.Command{
	Use:   "name",
	Short: "generate a method that returns struct's name",
	RunE: func(cmd *cobra.Command, _ []string) error {
		sourceType, pointer := strings.CutPrefix(target, "*")

		typ, test, err := loadType(cmd.Context(), pkgName, sourceType)
		if err != nil {
			return err
		}

		cfg := StructNameConfig{
			Typ:        typ,
			PkgName:    pkgName,
			MethodName: nameMethodName,
			Pointer:    pointer,
		}

		file := filepath.Clean(filename(sourceType, "name", test))

		f, err := os.Create(file)
		if err != nil {
			return err
		}

		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()

		if err := StructName(f, cfg); err != nil {
			return err
		}

		return f.Sync()
	},
}

var partialCmd = &cobra.Command{
	Use:   "partial",
	Short: "generate a partial struct",
	RunE: func(cmd *cobra.Command, _ []string) error {
		typ, test, err := loadType(cmd.Context(), pkgName, target)
		if err != nil {
			return err
		}

		cfg := PartialConfig{
			Typ:        typ,
			PkgName:    pkgName,
			StructName: partialStructName,
			Suffix:     &partialSuffix,
			Prefix:     partialPrefix,
		}

		file := filepath.Clean(filename(typ.Obj().Name(), "partial", test))

		f, err := os.Create(file)
		if err != nil {
			return err
		}

		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()

		if err := Partial(f, cfg); err != nil {
			return err
		}

		return f.Sync()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&pkgName,
		"pkg",
		"p",
		os.Getenv("GOPACKAGE"),
		"package on which generate code, defaults to GOPACKAGE")
	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target type")

	toMapCmd.Flags().StringVarP(&toMapMethodName, "method", "m", "ToMap", "out method name")
	toMapCmd.Flags().StringVarP(&toMapTag, "tag", "s", "to-map", "tag from which read metadata")

	nameCmd.Flags().StringVarP(&nameMethodName, "method", "m", DefaultNameMethodName, "out method name")

	partialCmd.Flags().StringVarP(&partialPrefix, "prefix", "e", "Partial",
		"struct name prefix, no effect when used with --structname")
	partialCmd.Flags().StringVarP(&partialSuffix, "suffix", "s", "",
		"struct name suffix, no effect when used with --structname")
	partialCmd.Flags().StringVarP(&partialStructName, "structname", "n", "",
		"struct name, if set suffix and prefix options will be ignored")

	rootCmd.AddCommand(toMapCmd)
	rootCmd.AddCommand(nameCmd)
	rootCmd.AddCommand(partialCmd)
}

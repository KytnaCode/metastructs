package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	target     string
	pkgName    string
	methodName string
	tag        string
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
			MethodName: methodName,
			TagName:    &tag,
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
			MethodName: methodName,
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

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&pkgName,
		"pkg",
		"p",
		os.Getenv("GOPACKAGE"),
		"package on which generate code, defaults to GOPACKAGE")

	toMapCmd.Flags().StringVarP(&target, "target", "t", "", "target type")
	toMapCmd.Flags().StringVarP(&methodName, "method", "m", "ToMap", "out method name")
	toMapCmd.Flags().StringVarP(&tag, "tag", "s", "to-map", "tag from which read metadata")

	nameCmd.Flags().StringVarP(&target, "target", "t", "", "target type")
	nameCmd.Flags().StringVarP(&methodName, "method", "m", DefaultNameMethodName, "out method name")

	rootCmd.AddCommand(toMapCmd)
	rootCmd.AddCommand(nameCmd)
}

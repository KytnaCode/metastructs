package main

import (
	"github.com/spf13/cobra"
)

var (
	target  string
	pkgName string
)

var rootCmd = &cobra.Command{
	Use:   "metastructs",
	Short: "metastructs group boilerplate generators for structs",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&pkgName,
		"pkg",
		"p",
		".",
		"package on which generate code, defaults to current package")
	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target type")

	rootCmd.AddCommand(toMapCmd)
	rootCmd.AddCommand(nameCmd)
	rootCmd.AddCommand(partialCmd)
}

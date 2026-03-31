package main

import (
	"os"

	"github.com/spf13/cobra"
)

var pkgName string

var rootCmd = &cobra.Command{
	Use:   "metastructs",
	Short: "metastructs group boilerplate generators for structs",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&pkgName,
		"pkg",
		"p",
		os.Getenv("GOPACKAGE"),
		"package on which generate code, defaults to GOPACKAGE")
}

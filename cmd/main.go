// Package main contains the entrypoint of metastructs and its command definitions.
package main

import (
	"log"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}

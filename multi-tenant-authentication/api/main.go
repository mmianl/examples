package main

import (
	"os"

	"github.com/mmianl/api/v2/cmd"
	"github.com/spf13/cobra"
)

func main() {
	args := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin, args)
	cobra.CheckErr(rootCmd.Execute())
}

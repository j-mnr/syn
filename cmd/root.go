package cmd

import (
	"fmt"
	"os"

	"github.com/j-mnr/syn/cmd/search"
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(search.Cmd)
}

var root = &cobra.Command{
	Use:   "syn",
	Short: "Syn gives you synonyms of a word from the command line",
	Long:  "",
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/timsexperiments/gcat/internal/cli"
	"github.com/timsexperiments/gcat/internal/clipboard"
	"github.com/timsexperiments/gcat/pkg/gcat"
)

var copyOutput bool

func main() {
	rootCmd := &cobra.Command{
		Use:   "gcat <source>",
		Short: "gcat concatenates files from a repository or local folder",
		Args:  cobra.ExactArgs(1),
		Run:   runGcat,
	}

	rootCmd.Flags().BoolVarP(&copyOutput, "copy", "c", false, "Copy output to clipboard instead of printing")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runGcat(cmd *cobra.Command, args []string) {
	source := args[0]

	repo, err := gcat.OpenRepository(source)
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	files, err := repo.GetFiles()
	if err != nil {
		log.Fatalf("Error retrieving files: %v", err)
	}

	selectedFiles, err := cli.SimpleSelector(files)
	if err != nil {
		log.Fatalf("Error during file selection: %v", err)
	}

	output, err := repo.ConcatFiles(selectedFiles)
	if err != nil {
		log.Fatalf("Error concatenating files: %v", err)
	}

	if copyOutput {
		clipboard.WriteText(output)
		fmt.Println("\nOutput copied to clipboard")
	} else {
		fmt.Println("\n=== Concatenated Output ===")
		fmt.Println(output)
	}
}

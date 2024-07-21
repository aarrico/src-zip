package cmd

import (
	"path/filepath"

	"github.com/aarrico/src-zip/internal/compressdir"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "src-zip [path to src directory]",
	Short: "src directory compression",
	Long: `compresses a src directory into a zip file considering the .gitignore file`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := filepath.Clean(args[0])
		compressType := "zip" // os.Args[2]
		target := source + "." + compressType
		compressdir.CompressDir(source, target)
	},
  }
  
 func Init () {
	rootCmd.Execute();
 }
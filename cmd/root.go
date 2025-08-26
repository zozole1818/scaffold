package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var outputDirName = "out"

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Creates scaffolding for a new project",
	Long: `Should be able to create new go project with initial code structure that enables 
you to quickly build a new project.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ordersys.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Int64P("port", "p", 8080, "port to listen on")
	rootCmd.PersistentFlags().StringP("projectName", "n", "test", "project name")
	rootCmd.PersistentFlags().StringP("output", "o", "out", "location of where to create project")
	rootCmd.PersistentFlags().StringP("moduleName", "m", "com/example", "name of Go module")
	rootCmd.PersistentFlags().String("goVersion", "1.24.5", "version of a Go")
	rootCmd.PersistentFlags().Bool("debugLogger", false, "enable debug logging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

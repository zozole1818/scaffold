package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zozole1818/scaffold/internal"
	"github.com/zozole1818/scaffold/internal/tmpl"
	"log/slog"
	"path/filepath"
)

var plainCmd = &cobra.Command{
	Use:     "plainGo",
	Aliases: []string{"plain"},
	Short:   "Generate a basic Go module",
	Long:    `Generate a basic Go module`,
	Run: func(cmd *cobra.Command, args []string) {
		parsedFlags, err := internal.GetFlags(cmd)
		if err != nil {
			slog.Error("error getting flags", "error", err)
			return
		}

		projectPath := internal.GetProjectPath(parsedFlags.Output, outputDirName, parsedFlags.ProjectName)

		templateInfoMap := map[string]internal.TemplateInfo{
			"module": {
				TmplString: tmpl.GetPlainModTemplate(),
				FilePath:   filepath.Join(projectPath, "go.mod"),
			},
			"main": {
				TmplString: tmpl.GetPlainMainTemplate(),
				FilePath:   filepath.Join(projectPath, "main.go"),
			},
		}

		err = internal.NewGenerator(
			"plainGo",
			parsedFlags,
			templateInfoMap,
			filepath.Join(projectPath)).
			Generate()
		if err != nil {
			slog.Error("error generating basic HTTP server", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(plainCmd)

}

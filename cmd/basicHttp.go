package cmd

import (
	"com/zuzanna/scaffold/internal"
	"com/zuzanna/scaffold/internal/tmpl"
	"github.com/spf13/cobra"
	"log/slog"
	"path/filepath"
)

var basicHttpCmd = &cobra.Command{
	Use:     "basicHttp",
	Aliases: []string{"http"},
	Short:   "Generate a basic HTTP server",
	Long:    `Generate a basic HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		parsedFlags, err := internal.GetFlags(cmd)
		if err != nil {
			slog.Error("error getting flags", "error", err)
			return
		}

		projectPath := internal.GetProjectPath(parsedFlags.Output, outputDirName, parsedFlags.ProjectName)

		templateInfoMap := map[string]internal.TemplateInfo{
			"module": {
				TmplString: tmpl.GetBasicHttpModTemplate(),
				FilePath:   filepath.Join(projectPath, "go.mod"),
			},
			"main": {
				TmplString: tmpl.GetBasicHttpMainTemplate(),
				FilePath:   filepath.Join(projectPath, "cmd/main.go"),
			},
			"dto": {
				TmplString: tmpl.GetBasicHttpDTOTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/dto.go"),
			},
			"handler": {
				TmplString: tmpl.GetBasicHttpHandlerTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/handler.go"),
			},
		}

		err = internal.NewGenerator(
			"basicHttp",
			parsedFlags,
			templateInfoMap,
			filepath.Join(projectPath, "cmd"),
			filepath.Join(projectPath, "internal")).
			Generate()
		if err != nil {
			slog.Error("error generating basic HTTP server", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(basicHttpCmd)

}

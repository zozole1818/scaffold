package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zozole1818/scaffold/internal"
	"github.com/zozole1818/scaffold/internal/tmpl"
	"log/slog"
	"path/filepath"
)

var websocketCmd = &cobra.Command{
	Use:     "websocket",
	Aliases: []string{"ws"},
	Short:   "Generate a WebSocket server",
	Long:    `Generate a WebSocket server`,
	Run: func(cmd *cobra.Command, args []string) {
		parsedFlags, err := internal.GetFlags(cmd)
		if err != nil {
			slog.Error("error getting flags", "error", err)
			return
		}

		projectPath := internal.GetProjectPath(parsedFlags.Output, outputDirName, parsedFlags.ProjectName)

		templateInfoMap := map[string]internal.TemplateInfo{
			"module": {
				TmplString: tmpl.GetWsModTemplate(),
				FilePath:   filepath.Join(projectPath, "go.mod"),
			},
			"main": {
				TmplString: tmpl.GetWsMainTemplate(),
				FilePath:   filepath.Join(projectPath, "cmd/main.go"),
			},
			"dto": {
				TmplString: tmpl.GetWsDTOTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/dto.go"),
			},
			"handler": {
				TmplString: tmpl.GetWsHandlerTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/handler.go"),
			},
		}

		err = internal.NewGenerator(
			"websocket",
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
	rootCmd.AddCommand(websocketCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zozole1818/scaffold/internal"
	"github.com/zozole1818/scaffold/internal/tmpl"
	"log/slog"
	"path/filepath"
)

var expandedEchoHttpCmd = &cobra.Command{
	Use:     "expandedEchoHttp",
	Aliases: []string{"fullEcho"},
	Short:   "Generate a Echo HTTP server with more code",
	Long: `Generate a Echo HTTP server with more code. The same as in echoHttp plus:
- main: loading envs from file .env 
- dto: ErrorResponse
- basic service layer
- basic db layer - db is map with mutex,  loading data from file (at the startup) and saving data to file (when app is closing).`,
	Run: func(cmd *cobra.Command, args []string) {
		parsedFlags, err := internal.GetFlags(cmd)
		if err != nil {
			slog.Error("error getting flags", "error", err)
			return
		}

		projectPath := internal.GetProjectPath(parsedFlags.Output, outputDirName, parsedFlags.ProjectName)

		templateInfoMap := map[string]internal.TemplateInfo{
			"module": {
				TmplString: tmpl.GetExpandedEchoHttpModTemplate(),
				FilePath:   filepath.Join(projectPath, "go.mod"),
			},
			"main": {
				TmplString: tmpl.GetExpandedEchoHttpMainTemplate(),
				FilePath:   filepath.Join(projectPath, "cmd/main.go"),
			},
			"dto": {
				TmplString: tmpl.GetExpandedEchoHttpDTOTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/dto.go"),
			},
			"handler": {
				TmplString: tmpl.GetExpandedEchoHttpHandlerTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/handler.go"),
			},
			"service": {
				TmplString: tmpl.GetExpandedEchoHttpServiceTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/service.go"),
			},
			"misc": {
				TmplString: tmpl.GetExpandedEchoHttpMiscTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/misc.go"),
			},
			"db": {
				TmplString: tmpl.GetExpandedEchoHttpDBTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/db.go"),
			},
			"db_map": {
				TmplString: tmpl.GetExpandedEchoHttpDBMapTemplate(),
				FilePath:   filepath.Join(projectPath, "internal/db_map.go"),
			},
		}

		err = internal.NewGenerator(
			"expandedEchoHttp",
			parsedFlags,
			templateInfoMap,
			filepath.Join(projectPath, "cmd"),
			filepath.Join(projectPath, "internal")).
			Generate()
		if err != nil {
			slog.Error("error generating project", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(expandedEchoHttpCmd)

}

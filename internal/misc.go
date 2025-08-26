package internal

import (
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
)

type ParsedFlags struct {
	ProjectName string
	ModuleName  string
	GoVersion   string
	Port        int64
	DebugLogger bool
	Output      string
	Verbose     bool
}

func GetFlags(cmd *cobra.Command) (ParsedFlags, error) {
	d := ParsedFlags{}
	verbose, _ := cmd.Flags().GetBool("verbose")
	d.Verbose = verbose

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return ParsedFlags{}, fmt.Errorf("error getting output flag: %w", err)
	}
	d.Output = output

	projectName, err := cmd.Flags().GetString("projectName")
	if err != nil {
		return ParsedFlags{}, fmt.Errorf("error getting project name flag: %w", err)
	}
	d.ProjectName = projectName

	moduleName, err := cmd.Flags().GetString("moduleName")
	if err != nil {
		return ParsedFlags{}, fmt.Errorf("error getting module name flag: %w", err)
	}
	d.ModuleName = moduleName

	goVersion, err := cmd.Flags().GetString("goVersion")
	if err != nil {
		return ParsedFlags{}, fmt.Errorf("error getting go version flag: %w", err)
	}
	d.GoVersion = goVersion

	port, err := cmd.Flags().GetInt64("port")
	if err != nil {
		return ParsedFlags{}, fmt.Errorf("error getting port flag: %w", err)
	}
	d.Port = port

	debugLogger, err := cmd.Flags().GetBool("debugLogger")
	if err != nil {
		return ParsedFlags{}, fmt.Errorf("error getting debug logger flag: %w", err)
	}
	d.DebugLogger = debugLogger

	return d, nil
}

func GetProjectPath(parsedPath string, defaultPath string, projectName string) string {
	projectPath := ""
	if parsedPath != "" {
		if parsedPath == "." {
			projectPath = projectName
		} else {
			projectPath = filepath.Join(parsedPath, projectName)
		}
	} else {
		projectPath = filepath.Join(defaultPath, projectName)
	}
	return projectPath
}

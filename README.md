# scaffold
Go project that helps to get started with new Go project. Creates small simple Go project that you can use to easily start coding.

## Install
```bash
$ go install github.com/zozole1818/scaffold@latest
```

### Usage
To see all available commands run `scaffold -h`. Example of an output:
```bash
$ scaffold -h
Should be able to create new go project with initial code structure that enables
you to quickly build a new project.

Usage:
   [command]

Available Commands:
  basicHttp   Generate a basic HTTP server
  completion  Generate the autocompletion script for the specified shell
  echoHttp    Generate a Echo HTTP server
  help        Help about any command
  websocket   Generate a WebSocket server

Flags:
      --debugLogger          enable debug logging
      --goVersion string     version of a Go (default "1.24.5")
  -h, --help                 help for this command
  -m, --moduleName string    name of Go module (default "com/example")
  -o, --output string        location of where to create project (default "out")
  -p, --port int             port to listen on (default 8080)
  -n, --projectName string   project name (default "test")
  -v, --verbose              verbose output

Use " [command] --help" for more information about a command.
```
package tmpl

func GetPlainModTemplate() string {
	return `module {{.ModuleName}}

go {{.GoVersion}}
`
}

func GetPlainMainTemplate() string {
	return `package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
{{ if .DebugLogger }}
	slog.SetLogLoggerLevel(slog.LevelDebug)
{{ end }}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Println("Hello on this lovely day! :)")

	// Wait for context cancellation
	<-ctx.Done()
}`
}

package tmpl

func GetEchoHttpModTemplate() string {
	return `module {{.ModuleName}}

go {{.GoVersion}}

require github.com/labstack/echo/v4 v4.13.4

require (
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)
`
}
func GetEchoHttpMainTemplate() string {
	return `package main

import (
	"{{ .ModuleName }}/internal"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
{{ if .DebugLogger }}
	slog.SetLogLoggerLevel(slog.LevelDebug)
{{ end }}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := {{ .Port }}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.Static("public/static")) // uncomment to serve static files
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"http://localhost:5173"},
	//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	//})) // uncomment to define CORS

	h := internal.NewHandler(ctx)

	e.GET("/quotes", h.GetQuote())

	go func() {
		if err := e.Start(":" + strconv.Itoa(port)); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server Start error: " + err.Error())
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	<-shutdownCtx.Done()
	slog.Info("Server shutdown complete.")
}
`
}

func GetEchoHttpDTOTemplate() string {
	return `package internal

type QuoteResponse struct {
	Quote string ` + "`" + `json:"quote"` + "`" + `
}`
}

func GetEchoHttpHandlerTemplate() string {
	return `package internal

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		ctx: ctx,
	}
}

func (h *Handler) GetQuote() func(echo.Context) error {
	return func(c echo.Context) error {
		quote := QuoteResponse{
			Quote: "Life is like a box of chocolates, you never know what you're going to get.",
		}

		return c.JSON(http.StatusOK, quote)
	}
}`
}

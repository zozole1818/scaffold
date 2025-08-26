package tmpl

func GetBasicHttpModTemplate() string {
	return `module {{.ModuleName}}

go {{.GoVersion}}
`
}

func GetBasicHttpMainTemplate() string {
	return `package main

import (
	"{{ .ModuleName }}/internal"
	"context"
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

	server := &http.Server{
		Addr: ":" + strconv.Itoa(port),
	}

	h := internal.NewHandler(ctx)

	go func() {
		http.HandleFunc("/quotes", internal.MethodWrapper(http.MethodGet, h.GetQuote()))

		slog.Info("Starting HTTP server :" + strconv.Itoa(port) + "...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server ListenAndServe: " + err.Error())
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	<-shutdownCtx.Done()
	slog.Info("Server shutdown complete.")
}`
}

func GetBasicHttpDTOTemplate() string {
	return `package internal

type QuoteResponse struct {
	Quote string ` + "`" + `json:"quote"` + "`" + `
}

type ErrResponse struct {
	Error     string ` + "`" + `json:"error"` + "`" + `
	Timestamp string ` + "`" + `json:"timestamp"` + "`" + `
}`
}

func GetBasicHttpHandlerTemplate() string {
	return `package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		ctx: ctx,
	}
}

func (h *Handler) GetQuote() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		quote := QuoteResponse{
			Quote: "Life is like a box of chocolates, you never know what you're going to get.",
		}

		writeResponse(w, http.StatusOK, quote)
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(ErrResponse{Error: errorMessage, Timestamp: time.Now().Format(time.RFC3339)})
	if err != nil {
		slog.Error("error marshalling json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeResponse[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("error marshalling json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// MethodWrapper wraps a handler function to validate HTTP method
func MethodWrapper(method string, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handler(w, r)
	}
}`
}

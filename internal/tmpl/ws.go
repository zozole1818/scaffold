package tmpl

func GetWsModTemplate() string {
	return `module {{.ModuleName}}

go {{.GoVersion}}

require github.com/gorilla/websocket v1.5.3
`
}

func GetWsMainTemplate() string {
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
		http.HandleFunc("/quotes", h.GetQuote(ctx))

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
}
`
}

func GetWsDTOTemplate() string {
	return `package internal

type QuoteResponse struct {
	Quote string ` + "`" + `json:"quote"` + "`" + `
}`
}

func GetWsHandlerTemplate() string {
	return `package internal

import (
	"context"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"strconv"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) GetQuote(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wsCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		tick := time.Tick(time.Second * 2)
		counter := 0
		quote := QuoteResponse{
			Quote: "Life is like a box of chocolates, you never know what you're going to get.",
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("upgrader error:" + err.Error())
			return
		}
		defer conn.Close()
		slog.Info("Client connected")

		go func() {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						slog.Error("client closing ws", "error", err)
					}
					cancel()
					return
				}
			}
		}()

		timeout := time.After(time.Second * 10)
		for {
			select {
			case msg := <-tick:
				counter++
				err = conn.WriteJSON(map[string]string{
					"counter":   strconv.Itoa(counter),
					"quote":     quote.Quote,
					"timestamp": msg.Format(time.RFC3339),
				})
				if err != nil {
					slog.Error("writer error:" + err.Error())
					return
				}
			case <-wsCtx.Done():
				slog.Warn("Client closed ws. Cleanup.")
				return
			case <-timeout:
				slog.Warn("Timeout. Cleanup.")
				return
			case <-ctx.Done():
				slog.Warn("Context cancelled. Closing down current ws connection.")
				return
			}
		}

	}
}
`
}

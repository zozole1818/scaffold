package tmpl

func GetExpandedEchoHttpModTemplate() string {
	return `module {{.ModuleName}}

go {{.GoVersion}}

require github.com/labstack/echo/v4 v4.13.4

require (
	github.com/google/uuid v1.6.0 // indirect
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

func GetExpandedEchoHttpMainTemplate() string {
	return `package main

import (
	"{{ .ModuleName }}/internal"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v3"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
{{ if .DebugLogger }}
	slog.SetLogLoggerLevel(slog.LevelDebug)
{{ end }}

	err := loadEnvs(".env")
	if err != nil {
		slog.Error("Error when loading env file", "error", err)
		return
	}

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

	filePath := "db.json"
	db := internal.NewMapDB(filePath)
	svc := internal.NewService(ctx, db)
	h := internal.NewHandler(ctx, svc)

	e.GET("/quotes", h.GetQuotes())

	go func() {
		if err := e.Start(":" + strconv.Itoa(port)); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server Start error: " + err.Error())
		}
	}()

	// Wait for context cancellation
	select {
	case <-ctx.Done():
		err := db.Flush(filePath)
		if err != nil {
			slog.Error("error when flushing db: " + err.Error())
		}
		break
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	<-shutdownCtx.Done()
	slog.Info("Server shutdown complete.")
}

func loadEnvs(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error when reading %s file: %v", path, err)
	}
	result := make(map[string]string)
	err = yaml.Unmarshal(b, &result)
	if err != nil {
		return fmt.Errorf("error when Unmarshal %s file: %v", path, err)
	}
	var errors []string
	for k, v := range result {
		err = os.Setenv(k, v)
		if err != nil {
			errors = append(errors, fmt.Errorf("error when setting %s env var: %v", k, err).Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("errors when setting env vars: %s", strings.Join(errors, ", "))
	}
	return nil
}
`
}

func GetExpandedEchoHttpDTOTemplate() string {
	return `package internal

import "time"

type QuoteResponse struct {
	Quote string ` + "`" + `json:"quote"` + "`" + `
}

type ErrorResponse struct {
	Error     string ` + "`" + `json:"error"` + "`" + `
	Timestamp string ` + "`" + `json:"timestamp"` + "`" + `
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:     err.Error(),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
`
}

func GetExpandedEchoHttpHandlerTemplate() string {
	return `package internal

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	ctx context.Context
	svc Service
}

func NewHandler(ctx context.Context, svc Service) *Handler {
	return &Handler{
		ctx: ctx,
		svc: svc,
	}
}

func (h *Handler) GetQuotes() func(echo.Context) error {
	return func(c echo.Context) error {
		q, err := h.svc.GetQuotes()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("error when getting quote: %w", err)))
		}

		return c.JSON(http.StatusOK, q)
	}
}`
}

func GetExpandedEchoHttpServiceTemplate() string {
	return `package internal

import "context"

type Service interface {
	GetQuotes() ([]Quote, error)
}

type Quote struct {
	id    string
	Value string
}

func (q *Quote) SetID(id string) {
	q.id = id
}

func (q *Quote) GetID() string {
	return q.id
}

type MyService struct {
	ctx context.Context
	db  DB
}

func NewService(ctx context.Context, db DB) Service {
	return &MyService{
		ctx: ctx,
		db:  db,
	}
}

func (s *MyService) GetQuotes() ([]Quote, error) {
	items := s.db.GetAll()
	return mapToQuotes(items), nil
}

func mapToQuotes(items []Item) []Quote {
	q := make([]Quote, 0, len(items))
	for _, it := range items {
		q = append(q, *it.(*Quote))
	}
	return q
}`
}

func GetExpandedEchoHttpDBTemplate() string {
	return `package internal

type DB interface {
	Add(s Item) string
	Get(id string) Item
	Delete(id string) string
	GetAll() []Item
	Replace(s Item)
	Clear()
	Flush(path string) error
}

type Item interface {
	SetID(id string)
	GetID() string
}`
}

func GetExpandedEchoHttpDBMapTemplate() string {
	return `package internal

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"sync"
)

type MapDB struct {
	items map[string]Item
	mu    *sync.Mutex
}

func NewMapDB(filePath string) DB {
	b, err := os.ReadFile(filePath)
	if err != nil {
		slog.Debug("Failed to read file", "path", filePath)
		return &MapDB{
			items: make(map[string]Item),
			mu:    &sync.Mutex{},
		}
	}

	var quotes map[string]*Quote
	if err := json.Unmarshal(b, &quotes); err != nil {
		slog.Debug("Failed to unmarshal items", "error", err)
		return &MapDB{
			items: make(map[string]Item),
			mu:    &sync.Mutex{},
		}
	}
	items := make(map[string]Item, len(quotes))
	for k, v := range quotes {
		items[k] = v
	}
	return &MapDB{
		items: items,
		mu:    &sync.Mutex{},
	}
}

func (db *MapDB) Add(s Item) string {
	db.mu.Lock()
	defer db.mu.Unlock()
	id := uuid.New().String()
	s.SetID(id)
	db.items[id] = s
	return id
}

func (db *MapDB) Get(id string) Item {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.items[id]
}

func (db *MapDB) Delete(id string) string {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.items, id)
	return id
}

func (db *MapDB) GetAll() []Item {
	db.mu.Lock()
	defer db.mu.Unlock()
	items := make([]Item, 0, len(db.items))
	for _, s := range db.items {
		items = append(items, s)
	}
	return items
}

func (db *MapDB) Replace(s Item) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items[s.GetID()] = s
}

func (db *MapDB) Clear() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items = make(map[string]Item)
}

func (db *MapDB) Flush(path string) error {
	slog.Debug("Flushing MapDB to file", "path", path)
	db.mu.Lock()
	defer db.mu.Unlock()

	b, err := json.Marshal(db.items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %s", err)
	}

	// Write to the file
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %s", err)
	}
	return nil
}
`
}

func GetExpandedEchoHttpMiscTemplate() string {
	return `package internal

import "fmt"

func ToConcreteSlice[T Item](items []Item) ([]T, error) {
	result := make([]T, 0, len(items))
	for _, it := range items {
		t, ok := it.(T)
		if !ok {
			return nil, fmt.Errorf("item is not of type %T: %#v", *new(T), it)
		}
		result = append(result, t)
	}
	return result, nil
}`
}

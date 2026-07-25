package http

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"db-studio-go/internal/db"
	"db-studio-go/internal/http/handlers"
	"db-studio-go/internal/http/routes"
	"db-studio-go/internal/http/services"
)

type Server struct {
	router     *chi.Mux
	driver     db.Database
	webFS      embed.FS
	port       int
	actualPort int
}

func NewServer(driver db.Database, webFS embed.FS, port int) *Server {
	s := &Server{
		router:     chi.NewRouter(),
		driver:     driver,
		webFS:      webFS,
		port:       port,
		actualPort: port,
	}

	s.setupMiddleware()
	s.setupRoutes()
	return s
}

func (s *Server) GetPort() int {
	return s.actualPort
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	// CORS Middleware for development
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}

func (s *Server) setupRoutes() {
	// Initialize services
	connService := services.NewConnectionService(s.driver)
	tableService := services.NewTableService(s.driver)
	queryService := services.NewQueryService(s.driver)

	// Initialize handlers
	connHandler := handlers.NewConnectionHandler(connService)
	tableHandler := handlers.NewTableHandler(tableService)
	queryHandler := handlers.NewQueryHandler(queryService)

	// Register API routes
	routes.RegisterAPIRoutes(s.router, connHandler, tableHandler, queryHandler)

	// Embedded Static Assets Handler for SvelteKit SPA
	subFS, err := fs.Sub(s.webFS, "web/build")
	if err == nil {
		fileServer := http.FileServer(http.FS(subFS))
		s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				http.Error(w, "API Endpoint Not Found", http.StatusNotFound)
				return
			}
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path != "" {
				if _, err := subFS.Open(path); err == nil {
					fileServer.ServeHTTP(w, r)
					return
				}
			}
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
		})
	}
}

// ListenAndServe tries binding ports starting from s.port up to s.port+10 (Automatic Port Fallback)
func (s *Server) ListenAndServe(onServerReady func(port int)) error {
	var listener net.Listener
	var err error
	foundPort := s.port

	for p := s.port; p <= s.port+10; p++ {
		addr := fmt.Sprintf(":%d", p)
		listener, err = net.Listen("tcp", addr)
		if err == nil {
			foundPort = p
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to bind HTTP server to any port between %d and %d: %w", s.port, s.port+10, err)
	}

	s.actualPort = foundPort
	fmt.Printf("🚀 DBStudio Web Server running at http://localhost:%d\n", foundPort)

	if onServerReady != nil {
		onServerReady(foundPort)
	}

	return http.Serve(listener, s.router)
}

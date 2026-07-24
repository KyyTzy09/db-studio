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
	h := NewHandler(s.driver)

	s.router.Route("/api", func(r chi.Router) {
		r.Get("/connection/status", h.HandleGetConnectionStatus)
		r.Get("/tables", h.HandleGetTables)
		r.Get("/tables/{name}/schema", h.HandleGetSchema)
		r.Get("/tables/{name}/data", h.HandleGetData)
		r.Post("/tables/{name}", h.HandleInsertRow)
		r.Patch("/tables/{name}", h.HandleUpdateRow)
		r.Delete("/tables/{name}", h.HandleDeleteRow)
		r.Post("/query", h.HandleExecuteQuery)
	})

	// Embedded Static Assets Handler for SvelteKit SPA
	subFS, err := fs.Sub(s.webFS, "web/build")
	if err == nil {
		fileServer := http.FileServer(http.FS(subFS))
		s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				http.Error(w, "API Endpoint Not Found", http.StatusNotFound)
				return
			}
			// Serve static file if exists, else fallback to index.html for SPA routing
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path != "" {
				if _, err := subFS.Open(path); err == nil {
					fileServer.ServeHTTP(w, r)
					return
				}
			}
			// Fallback to index.html
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

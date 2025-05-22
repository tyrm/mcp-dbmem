package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tyrm/mcp-dbmem/internal/logic"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const serverTimeout = 5 * time.Second

type Server struct {
	engine *gin.Engine
	logic  logic.Logic
	srv    *http.Server
}

type ServerConfig struct {
	Logic           logic.Logic
	Store           sessions.Store
	ApplicationName string
	HttpBind        string
}

// NewServer creates a new http web server.
func NewServer(_ context.Context, cnf ServerConfig) (*Server, error) {
	engine := gin.Default()

	engine.Use(otelgin.Middleware(cnf.ApplicationName))
	engine.Use(sessions.Sessions(cnf.ApplicationName, cnf.Store))

	return &Server{
		engine: engine,
		logic:  cnf.Logic,
		srv: &http.Server{
			Addr:         cnf.HttpBind,
			Handler:      engine,
			WriteTimeout: serverTimeout,
			ReadTimeout:  serverTimeout,
		},
	}, nil
}

func (s *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.engine.Group(relativePath, handlers...)
}

func (s *Server) Static(relativePath string, root string) {
	s.engine.Static(relativePath, root)
}

func (s *Server) StaticFS(relativePath string, fs http.FileSystem) {
	s.engine.StaticFS(relativePath, fs)
}

// Start starts the web server.
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Stop shuts down the web server.
func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Use(handlers ...gin.HandlerFunc) {
	s.engine.Use(handlers...)
}

func (s *Server) Delete(relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.DELETE(relativePath, handlers...)
}

func (s *Server) Get(relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.GET(relativePath, handlers...)
}

func (s *Server) Head(relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.HEAD(relativePath, handlers...)
}

func (s *Server) Options(relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.OPTIONS(relativePath, handlers...)
}

func (s *Server) Path(relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.PATCH(relativePath, handlers...)
}

func (s *Server) Post(relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.POST(relativePath, handlers...)
}

func (s *Server) Put(relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.PUT(relativePath, handlers...)
}

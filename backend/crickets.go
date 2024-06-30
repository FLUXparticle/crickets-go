package main

import (
	"context"
	"crickets-go/handler"
	"crickets-go/service"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"go.uber.org/fx"

	"github.com/gin-gonic/gin"
)

// NewLogger initialisiert den Standard-Logger von Go
func NewLogger() *log.Logger {
	return log.New(os.Stdout, "[Crickets] ", log.LstdFlags)
}

// Middleware für statische Dateien basierend auf dem Pfad
func staticFileServer() gin.HandlerFunc {
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir("./app")))
	staticHandler := http.StripPrefix("/", http.FileServer(http.Dir("./static")))
	return func(c *gin.Context) {
		requestURI := c.Request.RequestURI

		var handler http.Handler
		if strings.HasPrefix(requestURI, "/app/") {
			handler = appHandler
		} else if !strings.HasPrefix(requestURI, "/api/") {
			handler = staticHandler
		}

		if handler != nil {
			handler.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func NewGinHandler(userHandler *handler.UserHandler) http.Handler {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(userHandler.Auth)

	r.Use(staticFileServer())

	// API-Routes
	api := r.Group("/api")
	api.POST("/login", userHandler.Login)
	api.GET("/username", userHandler.Username)

	return r
}

func NewHTTPServer(lc fx.Lifecycle, logger *log.Logger, handler http.Handler) *http.Server {
	// Einstellungen für die Server-Adresse über Umgebungsvariable
	addr := ":8080"
	if env, found := os.LookupEnv("ADDR"); found {
		addr = env
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Printf("Server on %s running...", addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			NewGinHandler,
			NewLogger,
			handler.NewUserHandler,
			service.NewUserService,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

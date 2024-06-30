package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"go.uber.org/fx"

	"github.com/gin-gonic/gin"
)

const sessionCookieName = "session_token"

// In-Memory-Storage für Sessions (in einer echten Anwendung sollte dies persistent sein)
var sessions = map[string]string{}

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

func loginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// Überprüfe Benutzername und Passwort (dies ist ein einfaches Beispiel, in einer echten Anwendung sollte dies sicherer sein)
	if loginData.Username == "user" && loginData.Password == "pass" {
		// Generiere ein Session-Token
		sessionToken := "some-session-token" // Generiere hier ein echtes Token, z.B. UUID
		sessions[sessionToken] = loginData.Username

		// Setze den Session-Cookie
		c.SetCookie(sessionCookieName, sessionToken, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func authMiddleware(c *gin.Context) {
	sessionToken, err := c.Cookie(sessionCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if _, exists := sessions[sessionToken]; !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		c.Abort()
		return
	}

	c.Next()
}

func NewGinHandler() http.Handler {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(staticFileServer())

	r.Use(authMiddleware)

	// API-Routes
	api := r.Group("/api")
	api.POST("/login", loginHandler)

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
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

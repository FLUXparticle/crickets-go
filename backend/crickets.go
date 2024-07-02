//go:generate protoc --go_out=. --go-grpc_out=. --proto_path=proto proto/timeline.proto
package main

import (
	"context"
	"crickets-go/gen/timeline"
	"crickets-go/handler"
	"crickets-go/repository"
	"crickets-go/service"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

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

		var useHandler http.Handler
		if strings.HasPrefix(requestURI, "/app/") {
			useHandler = appHandler
		} else if !strings.HasPrefix(requestURI, "/api/") {
			useHandler = staticHandler
		}

		if useHandler != nil {
			useHandler.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func NewGinHandler(userHandler *handler.UserHandler, profileHandler *handler.ProfileHandler, timelineHandler *handler.TimelineHandler, chatHandler *handler.ChatHandler, internalHandler *handler.InternalHandler, debugHandler *handler.DebugHandler) http.Handler {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(userHandler.Auth)

	r.Use(staticFileServer())

	// API-Routes
	{
		api := r.Group("/api")

		api.POST("/login", userHandler.Login)
		api.GET("/username", userHandler.Username)

		api.GET("/profile", profileHandler.Profile)
		api.POST("/subscribe", profileHandler.Subscribe)

		api.GET("/search", timelineHandler.Search)
		api.POST("/post", timelineHandler.Post)
		api.GET("/timeline", timelineHandler.Timeline)

		api.GET("/chatWS", chatHandler.ChatWebSocket)

		// Internal-Routes
		{
			internal := api.Group("/internal")

			internal.POST("/subscribe", internalHandler.Subscribe)
		}

		// Debug-Routes (müssen natürlich für den Release deaktiviert werden)
		{
			debug := api.Group("/debug")

			debug.GET("/subscriptions", debugHandler.Subscriptions)
		}
	}

	return r
}

// NewHTTPServer initialisiert und startet den HTTP-Server (Gin)
func NewHTTPServer(lc fx.Lifecycle, logger *log.Logger, handler http.Handler) *http.Server {
	// Einstellungen für die Server-Adresse über Umgebungsvariable
	localhost := os.Getenv("LOCALHOST")
	addr := localhost + ":8080"
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
			logger.Printf("HTTP server on %s running...", addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

type server struct {
	timeline.UnimplementedTimelineServiceServer
	timelineService *service.TimelineService
}

func (s *server) Search(ctx context.Context, req *timeline.SearchRequest) (*timeline.SearchResponse, error) {
	posts := s.timelineService.Search("", req.Query)
	responsePosts := make([]*timeline.Post, len(posts))
	for i, post := range posts {
		responsePosts[i] = convertPost(post)
	}
	return &timeline.SearchResponse{Posts: responsePosts}, nil
}

func (s *server) TimelineUpdates(req *timeline.TimelineUpdateRequest, stream timeline.TimelineService_TimelineUpdatesServer) error {
	// Fetch updates from the TimelineService
	updatesChan := s.timelineService.LocalTimelineUpdates(req.CreatorIds)

	for {
		select {
		case update, ok := <-updatesChan:
			if !ok {
				return nil
			}
			post := convertPost(update)
			if err := stream.Send(&timeline.TimelineUpdateResponse{Post: post}); err != nil {
				log.Printf("Error sending update: %v", err)
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func convertPost(update *repository.Post) *timeline.Post {
	post := &timeline.Post{
		Username:  update.Creator.Username,
		Content:   update.Content,
		CreatedAt: timestamppb.New(update.CreatedAt),
	}
	return post
}

// NewGRPCServer initialisiert und startet den gRPC-Server
func NewGRPCServer(lc fx.Lifecycle, logger *log.Logger, timelineService *service.TimelineService) *grpc.Server {
	localhost := os.Getenv("LOCALHOST")
	addr := localhost + ":50051"
	grpcServer := grpc.NewServer()
	timeline.RegisterTimelineServiceServer(grpcServer, &server{timelineService: timelineService})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			logger.Printf("gRPC server on %s running...", addr)
			go grpcServer.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
	return grpcServer
}

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			NewGRPCServer,
			NewGinHandler,
			NewLogger,

			handler.NewUserHandler,
			handler.NewProfileHandler,
			handler.NewTimelineHandler,
			handler.NewChatHandler,
			handler.NewInternalHandler,
			handler.NewDebugHandler,

			service.NewUserService,
			service.NewProfileService,
			service.NewTimelineService,
			service.NewChatService,
			service.NewPubSub,
			service.NewMessageQueueProvider,

			repository.NewUserRepository,
			repository.NewSubscriptionRepository,
			repository.NewPostRepository,
		),
		fx.Invoke(func(*http.Server, *grpc.Server) {}),
	).Run()
}

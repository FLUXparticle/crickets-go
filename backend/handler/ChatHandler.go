package handler

import (
	"crickets-go/data"
	"crickets-go/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ChatHandler struct {
	logger      *log.Logger
	userHandler *UserHandler
	chatService *service.ChatService
}

func NewChatHandler(logger *log.Logger, userHandler *UserHandler, chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		logger:      logger,
		userHandler: userHandler,
		chatService: chatService,
	}
}

func (h *ChatHandler) ChatWebSocket(c *gin.Context) {
	user := h.userHandler.getUser(c)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return
	}
	defer conn.Close()

	// Subscribe to chat updates
	chatUpdates := h.chatService.ChatUpdates()
	// TODO defer h.chatService.Unsubscribe(chatUpdates)

	go func() {
		for post := range chatUpdates {
			if err := conn.WriteJSON(displayPost(post)); err != nil {
				log.Println("Failed to send message:", err)
				return
			}
		}
	}()

	for {
		var post data.Post
		if err := conn.ReadJSON(&post); err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("Unexpected WebSocket close:", err)
			}
			break
		}

		post.Creator = user
		post.CreatedAt = time.Now()

		// Post the chat message
		if err := h.chatService.SendChatMessage(&post); err != nil {
			h.logger.Println("Failed to send message:", err)
		}
	}
}

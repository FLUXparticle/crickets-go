package handler

import (
	"crickets-go/repository"
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
	userHandler *UserHandler
	chatService *service.ChatService
}

func NewChatHandler(userHandler *UserHandler, chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
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
		var data struct {
			Content string `json:"content"`
		}
		if err := conn.ReadJSON(&data); err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("Unexpected WebSocket close:", err)
			}
			break
		}

		post := &repository.Post{
			Creator:   user,
			Content:   data.Content,
			CreatedAt: time.Now(),
		}

		// Post the chat message
		h.chatService.SendChatMessage(post)
	}
}

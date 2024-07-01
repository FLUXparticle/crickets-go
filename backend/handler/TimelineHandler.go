package handler

import (
	"crickets-go/repository"
	"crickets-go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type TimelineHandler struct {
	userHandler     *UserHandler
	timelineService *service.TimelineService
}

func NewTimelineHandler(userHandler *UserHandler, timelineService *service.TimelineService) *TimelineHandler {
	return &TimelineHandler{
		userHandler:     userHandler,
		timelineService: timelineService,
	}
}

func (h *TimelineHandler) Search(c *gin.Context) {
	query := c.Param("q")

	posts := h.timelineService.Search("localhost", query)

	result := make([]map[string]any, len(posts))
	for i, post := range posts {
		result[i] = displayPost(post)
	}

	c.JSON(http.StatusOK, result)
}

func (h *TimelineHandler) Post(c *gin.Context) {
	var data struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creator := h.userHandler.getUser(c)
	h.timelineService.Post(creator, data.Content)
}

func (h *TimelineHandler) Timeline(c *gin.Context) {
	subscriber := h.userHandler.getUser(c)

	updates := h.timelineService.TimelineUpdates(subscriber)

	c.Stream(func(w io.Writer) bool {
		select {
		case post := <-updates:
			fmt.Println("Received:", post.Content)
			c.SSEvent("", displayPost(post))
			return true
		case <-c.Writer.CloseNotify():
			return false
		}
	})
}

func displayPost(post *repository.Post) map[string]any {
	creator := post.Creator
	username := creator.Username
	if len(creator.Server) > 0 {
		username += "@" + creator.Server
	}
	return map[string]any{
		"username":  username,
		"content":   post.Content,
		"createdAt": post.CreatedAt,
	}
}

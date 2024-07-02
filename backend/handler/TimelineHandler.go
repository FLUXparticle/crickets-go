package handler

import (
	"crickets-go/data"
	"crickets-go/service"
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
	server := c.Param("s")
	query := c.Param("q")

	posts := h.timelineService.Search(server, query)

	result := make([]map[string]any, len(posts))
	for i, post := range posts {
		result[i] = displayPost(post)
	}

	c.JSON(http.StatusOK, result)
}

func (h *TimelineHandler) Post(c *gin.Context) {
	var body struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creator := h.userHandler.getUser(c)
	h.timelineService.Post(creator, body.Content)
}

func (h *TimelineHandler) Timeline(c *gin.Context) {
	subscriber := h.userHandler.getUser(c)

	updates := h.timelineService.TimelineUpdates(subscriber.ID)

	c.Stream(func(w io.Writer) bool {
		select {
		case post := <-updates:
			//fmt.Println("Received:", post.Content)
			c.SSEvent("", displayPost(post))
			return true
		case <-c.Writer.CloseNotify():
			return false
		}
	})
}

func displayPost(post *data.Post) map[string]any {
	creator := post.Creator
	creatorName := creator.Username
	if len(creator.Server) > 0 {
		creatorName += "@" + creator.Server
	}
	return map[string]any{
		"creatorName": creatorName,
		"content":     post.Content,
		"createdAt":   post.CreatedAt,
	}
}

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
	server := c.Query("s")
	query := c.Query("q")

	posts := h.timelineService.Search(server, query)

	searchResults := make([]map[string]any, len(posts))
	for i, post := range posts {
		searchResults[i] = displayPost(post)
	}

	if len(searchResults) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "nothing found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"searchResults": searchResults})
}

func (h *TimelineHandler) Post(c *gin.Context) {
	var body struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	creator := h.userHandler.getUser(c)
	err := h.timelineService.Post(creator, body.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
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

func (h *TimelineHandler) LikePost(c *gin.Context) {
	var data struct {
		PostID      int64  `json:"postId" binding:"required"`
		CreatorName string `json:"creatorName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.timelineService.LikePost(data.PostID, data.CreatorName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func displayPost(post *data.Post) map[string]any {
	creator := post.Creator
	creatorName := creator.Username
	if len(creator.Server) > 0 {
		creatorName += "@" + creator.Server
	}
	return map[string]any{
		"id":          post.ID,
		"creatorName": creatorName,
		"content":     post.Content,
		"createdAt":   post.CreatedAt,
		"likes":       post.Likes,
	}
}

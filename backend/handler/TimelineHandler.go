package handler

import (
	"crickets-go/service"
	"github.com/gin-gonic/gin"
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

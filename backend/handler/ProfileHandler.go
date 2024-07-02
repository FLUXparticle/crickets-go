package handler

import (
	"crickets-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfileHandler struct {
	userHandler    *UserHandler
	profileService *service.ProfileService
}

func NewProfileHandler(userHandler *UserHandler, profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		userHandler:    userHandler,
		profileService: profileService,
	}
}

func (h *ProfileHandler) Profile(c *gin.Context) {
	user := h.userHandler.getUser(c)
	subscriberCount := h.profileService.SubscriberCount(user.ID)
	c.JSON(http.StatusOK, gin.H{"subscriberCount": subscriberCount})
}

func (h *ProfileHandler) Subscribe(c *gin.Context) {
	var data struct {
		Server      string `json:"server"`
		CreatorName string `json:"creatorName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	subscriber := h.userHandler.getUser(c)
	success, err := h.profileService.Subscribe(subscriber, data.Server, data.CreatorName)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"success": success})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}
}

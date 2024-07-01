package handler

import (
	"crickets-go/repository"
	"crickets-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerHandler struct {
	profileService *service.ProfileService
}

func NewServerHandler(profileService *service.ProfileService) *ServerHandler {
	return &ServerHandler{profileService: profileService}
}

func (h *ServerHandler) Subscribe(c *gin.Context) {
	var data struct {
		SubscriberServer string `json:"subscriberServer" binding:"required"`
		SubscriberName   string `json:"subscriberName" binding:"required"`
		CreatorName      string `json:"creatorName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscriber := &repository.User{
		Server:   data.SubscriberServer,
		Username: data.SubscriberName,
	}
	success, err := h.profileService.Subscribe(subscriber, "", data.CreatorName)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"successes": []string{success}})
	} else {
		c.JSON(http.StatusOK, gin.H{"errors": []string{err.Error()}})
	}
}

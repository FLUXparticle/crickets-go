package handler

import (
	"crickets-go/common"
	"crickets-go/repository"
	"crickets-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InternalHandler struct {
	profileService *service.ProfileService
}

func NewInternalHandler(profileService *service.ProfileService) *InternalHandler {
	return &InternalHandler{profileService: profileService}
}

func (h *InternalHandler) Subscribe(c *gin.Context) {
	// TODO gemeinsame Datenstruktur
	var data common.SubscribeRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creator, err := h.profileService.LocalSubscribe(data.Subscriber, data.CreatorName)
	if err == nil {
		c.JSON(http.StatusOK, &common.SubscribeResponse{
			User: &repository.User{
				ID:       creator.ID,
				Username: creator.Username,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}
}

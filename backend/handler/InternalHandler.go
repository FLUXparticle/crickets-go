package handler

import (
	"crickets-go/common"
	"crickets-go/data"
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
	var request common.SubscribeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creator, err := h.profileService.LocalSubscribe(request.Subscriber, request.CreatorName)
	var response *common.SubscribeResponse
	if err == nil {
		response = &common.SubscribeResponse{
			User: &data.User{
				ID:       creator.ID,
				Username: creator.Username,
			},
		}
	} else {
		response = &common.SubscribeResponse{
			Error: err.Error(),
		}
	}
	c.JSON(http.StatusOK, response)
}

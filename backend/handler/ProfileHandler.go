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

}

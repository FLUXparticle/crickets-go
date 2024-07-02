package handler

import (
	"crickets-go/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DebugHandler struct {
	subscriptionRepository *repository.SubscriptionRepository
}

func NewDebugHandler(subscriptionRepository *repository.SubscriptionRepository) *DebugHandler {
	return &DebugHandler{subscriptionRepository: subscriptionRepository}
}

func (h *DebugHandler) Subscriptions(c *gin.Context) {
	subscriptions := h.subscriptionRepository.FindAll()
	c.JSON(http.StatusOK, subscriptions)
}

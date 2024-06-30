package handler

import (
	"crickets-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const sessionCookieName = "session_token"

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) Auth(c *gin.Context) {
	if c.Request.RequestURI == "/api/login" {
		c.Next()
		return
	}

	sessionToken, err := c.Cookie(sessionCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if !h.userService.CheckSession(sessionToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		c.Abort()
		return
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
		return
	}

	// Überprüfe Benutzername und Passwort
	sessionToken, err := h.userService.Login(loginData.Username, loginData.Password)
	if err == nil {
		// Setze den Session-Cookie
		c.SetCookie(sessionCookieName, sessionToken, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "login successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
}

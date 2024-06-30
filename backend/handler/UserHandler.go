package handler

import (
	"crickets-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const sessionCookieName = "sessionToken"

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) Auth(c *gin.Context) {
	requestURI := c.Request.RequestURI

	sessionToken, err := c.Cookie(sessionCookieName)

	if strings.Contains(requestURI, "/app/") {
		if err != nil || !h.userService.CheckSession(sessionToken) {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
	}

	if strings.Contains(requestURI, "/api/") {
		if requestURI == "/api/login" {
			c.Next()
			return
		}

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

	c.Set(sessionCookieName, sessionToken)
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

func (h *UserHandler) Username(c *gin.Context) {
	sessionToken := c.GetString(sessionCookieName)
	username := h.userService.Username(sessionToken)
	c.JSON(http.StatusOK, gin.H{"username": username})
}

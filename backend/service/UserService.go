package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
)

type UserService struct {
	logger *log.Logger

	// User-Tabelle (in einer echten Anwendung sollten diese aus einer Datenbank oder z.B. LDAP kommen)
	users map[string]string
	// In-Memory-Storage für Sessions (in einer echten Anwendung sollte dies persistent sein)
	sessions map[string]string
}

func NewUserService(logger *log.Logger) *UserService {
	return &UserService{
		logger: logger,
		users: map[string]string{
			"admin":    "Secret123",
			"helpdesk": "Secret123",
			"employee": "Secret123",
			"manager":  "Secret123",
		},
		sessions: make(map[string]string),
	}
}

func (s *UserService) CheckSession(sessionToken string) bool {
	_, valid := s.sessions[sessionToken]
	return valid
}

func (s *UserService) Login(username, password string) (string, error) {
	if foundPassword, found := s.users[username]; found {
		if foundPassword == password {
			sessionToken, err := s.generateSessionToken()
			if err != nil {
				return "", err
			}
			s.sessions[sessionToken] = username
			return sessionToken, nil
		}
	}
	return "", errors.New("invalid credentials")
}

func (s *UserService) Username(sessionToken string) string {
	return s.sessions[sessionToken]
}

func (s *UserService) generateSessionToken() (string, error) {
	token := make([]byte, 4)
	_, err := rand.Read(token)
	if err != nil {
		s.logger.Printf("Error generating session token: %v", err)
		return "", err
	}
	return hex.EncodeToString(token), nil
}

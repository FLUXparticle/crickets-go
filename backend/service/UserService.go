package service

import (
	"crickets-go/repository"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
)

type UserService struct {
	logger         *log.Logger
	userRepository *repository.UserRepository

	// In-Memory-Storage für Sessions (in einer echten Anwendung sollte dies persistent sein)
	sessions map[string]*repository.User
}

func NewUserService(logger *log.Logger, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		logger:         logger,
		userRepository: userRepository,
		sessions:       make(map[string]*repository.User),
	}
}

func (s *UserService) CheckSession(sessionToken string) bool {
	_, valid := s.sessions[sessionToken]
	return valid
}

func (s *UserService) Login(username, password string) (string, error) {
	if user := s.userRepository.FindByUsername(username); user != nil {
		if user.Password == password {
			sessionToken, err := s.generateSessionToken()
			if err != nil {
				return "", err
			}
			s.sessions[sessionToken] = user
			return sessionToken, nil
		}
	}

	return "", errors.New("invalid credentials")
}

func (s *UserService) User(sessionToken string) *repository.User {
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

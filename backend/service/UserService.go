package service

import (
	"crickets-go/data"
	"crickets-go/repository"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
)

var validAPIKeys = map[string]bool{
	"fi4thee4kieyahhei3Chahth3iek6eib": true,
	"eeGix6Ooceew4booVeele6VeeTa1ahWu": true,
	"Ue4Aeghei4hagei1Tai4axoothooJam3": true,
}

type UserService struct {
	logger         *log.Logger
	userRepository *repository.UserRepository

	// In-Memory-Storage für Sessions (in einer echten Anwendung sollte dies persistent sein)
	sessions map[string]*data.User
}

func NewUserService(logger *log.Logger, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		logger:         logger,
		userRepository: userRepository,
		sessions:       make(map[string]*data.User),
	}
}

func (s *UserService) CheckApiKey(apiKey string) bool {
	valid, exists := validAPIKeys[apiKey]
	return exists && valid
}

func (s *UserService) CheckSession(sessionToken string) bool {
	_, exists := s.sessions[sessionToken]
	return exists
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

func (s *UserService) User(sessionToken string) *data.User {
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

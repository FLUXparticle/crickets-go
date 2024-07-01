package service

import (
	"crickets-go/common"
	"crickets-go/repository"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

type ProfileService struct {
	userRepository         *repository.UserRepository
	subscriptionRepository *repository.SubscriptionRepository
}

func NewProfileService(userRepository *repository.UserRepository, subscriptionRepository *repository.SubscriptionRepository) *ProfileService {
	return &ProfileService{
		userRepository:         userRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *ProfileService) SubscriberCount(creatorID int) int {
	return len(s.subscriptionRepository.FindByCreatorID(creatorID))
}

func (s *ProfileService) Subscribe(subscriber *repository.User, creatorServer string, creatorName string) (string, error) {
	if creatorServer == "" {
		creator := s.userRepository.FindByUsername(creatorName)
		if creator == nil {
			return "", errors.New("user not found")
		}
		s.subscriptionRepository.Save(&repository.Subscription{
			Creator:    creator,
			Subscriber: subscriber,
		})
	} else {
		s.sendSubscribe(creatorServer, creatorName, subscriber.Username)
	}
	return "success", nil
}

func (s *ProfileService) sendSubscribe(creatorServer string, creatorName string, subscriberName string) {
	client := resty.New()

	hostname, err := os.Hostname()
	apiKey := os.Getenv("API_KEY")

	// Daten für die Anfrage
	data := map[string]string{
		"subscriberServer": hostname,
		"subscriberName":   subscriberName,
		"creatorName":      creatorName,
	}

	// Anfrage an den Endpunkt senden
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-API-KEY", apiKey).
		SetBody(data).
		SetResult(&common.SubscribeResponse{}).
		Post(fmt.Sprintf("http://%s:8080/api/internal/subscribe", creatorServer))

	// Fehlerbehandlung
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}

	// Auswertung der Antwort
	subscribeResponse := resp.Result().(*common.SubscribeResponse)
	if len(subscribeResponse.Errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range subscribeResponse.Errors {
			fmt.Printf(" - %s\n", err)
		}
	} else {
		fmt.Println("Successes:")
		for _, success := range subscribeResponse.Successes {
			fmt.Printf(" - %s\n", success)
		}
	}
}

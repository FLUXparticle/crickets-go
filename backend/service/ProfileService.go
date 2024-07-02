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

func (s *ProfileService) SubscriberCount(creatorID int32) int {
	return len(s.subscriptionRepository.FindByCreatorID(creatorID))
}

func (s *ProfileService) Subscribe(subscriber *repository.User, creatorServer string, creatorName string) (string, error) {
	if creatorServer == "" {
		creator, err := s.LocalSubscribe(subscriber, creatorName)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Successfully subscribed to user '%s'", creator.Username), nil
	} else {
		creator, err := s.remoteSubscribe(creatorServer, creatorName, subscriber)
		if err != nil {
			return "", err
		}
		s.subscriptionRepository.Save(&repository.Subscription{
			Creator:    creator,
			Subscriber: subscriber,
		})
		return fmt.Sprintf("Successfully subscribed to user '%s' on server '%s'", creator.Username, creatorServer), nil
	}
}

func (s *ProfileService) LocalSubscribe(subscriber *repository.User, creatorName string) (*repository.User, error) {
	creator := s.userRepository.FindByUsername(creatorName)
	if creator == nil {
		return nil, errors.New(fmt.Sprintf("user '%s' not found", creatorName))
	}
	subscription := &repository.Subscription{
		Creator:    creator,
		Subscriber: subscriber,
	}
	s.subscriptionRepository.Save(subscription)
	return creator, nil
}

func (s *ProfileService) remoteSubscribe(creatorServer string, creatorName string, subscriber *repository.User) (*repository.User, error) {
	client := resty.New()

	// TODO Hostname über Fx (Config) default: localhost
	hostname, _ := os.Hostname()
	apiKey := os.Getenv("API_KEY")

	// Daten für die Anfrage
	data := &common.SubscribeRequest{
		Subscriber: &repository.User{
			ID:       subscriber.ID,
			Server:   hostname,
			Username: subscriber.Username,
		},
		CreatorName: creatorName,
	}

	// Anfrage an den Endpunkt senden
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-API-KEY", apiKey).
		SetBody(data).
		SetResult(&common.SubscribeResponse{}).
		Post(fmt.Sprintf("http://%s:8080/api/internal/subscribe", creatorServer))

	// Fehlerbehandlung
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, errors.New("error making request")
	}

	// Auswertung der Antwort
	subscribeResponse := response.Result().(*common.SubscribeResponse)
	if len(subscribeResponse.Error) > 0 {
		return nil, errors.New(subscribeResponse.Error)
	} else {
		creator := subscribeResponse.User
		return &repository.User{
			ID:       creator.ID,
			Server:   creatorServer,
			Username: creator.Username,
		}, nil
	}
}

package service

import (
	"crickets-go/common"
	"crickets-go/config"
	"crickets-go/data"
	"crickets-go/repository"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type ProfileService struct {
	config                 *config.Config
	userRepository         *repository.UserRepository
	subscriptionRepository *repository.SubscriptionRepository
}

func NewProfileService(config *config.Config, userRepository *repository.UserRepository, subscriptionRepository *repository.SubscriptionRepository) *ProfileService {
	return &ProfileService{
		config:                 config,
		userRepository:         userRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *ProfileService) SubscriberCount(creatorID int32) int {
	return len(s.subscriptionRepository.FindByCreatorID(creatorID))
}

func (s *ProfileService) Subscribe(subscriber *data.User, creatorServer string, creatorName string) (string, error) {
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
		s.subscriptionRepository.Save(&data.Subscription{
			Creator:    creator,
			Subscriber: subscriber,
		})
		return fmt.Sprintf("Successfully subscribed to user '%s' on server '%s'", creator.Username, creatorServer), nil
	}
}

func (s *ProfileService) LocalSubscribe(subscriber *data.User, creatorName string) (*data.User, error) {
	creator := s.userRepository.FindByUsername(creatorName)
	if creator == nil {
		return nil, errors.New(fmt.Sprintf("user '%s' not found", creatorName))
	}
	subscription := &data.Subscription{
		Creator:    creator,
		Subscriber: subscriber,
	}
	s.subscriptionRepository.Save(subscription)
	return creator, nil
}

func (s *ProfileService) remoteSubscribe(creatorServer string, creatorName string, subscriber *data.User) (*data.User, error) {
	client := resty.New()

	// Daten für die Anfrage
	request := &common.SubscribeRequest{
		Subscriber: &data.User{
			ID:       subscriber.ID,
			Server:   s.config.Hostname,
			Username: subscriber.Username,
		},
		CreatorName: creatorName,
	}

	// Anfrage an den Endpunkt senden
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-API-KEY", s.config.ApiKey).
		SetBody(request).
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
		return &data.User{
			ID:       creator.ID,
			Server:   creatorServer,
			Username: creator.Username,
		}, nil
	}
}

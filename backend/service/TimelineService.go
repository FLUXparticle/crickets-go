package service

import (
	"context"
	"crickets-go/gen/timeline"
	"crickets-go/repository"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type TimelineService struct {
	postRepository         *repository.PostRepository
	subscriptionRepository *repository.SubscriptionRepository
	pubSub                 *PubSub
}

func NewTimelineService(postRepository *repository.PostRepository, subscriptionRepository *repository.SubscriptionRepository, pubSub *PubSub) *TimelineService {
	return &TimelineService{
		postRepository:         postRepository,
		subscriptionRepository: subscriptionRepository,
		pubSub:                 pubSub,
	}
}

func (s *TimelineService) TimelineUpdates(subscriber *repository.User) chan *repository.Post {
	subscriptions := s.subscriptionRepository.FindBySubscriberID(subscriber.ID)

	creators := make([]*repository.User, len(subscriptions))
	for i, sub := range subscriptions {
		creators[i] = sub.Creator
	}

	aggregator := make(chan *repository.Post)

	for _, creator := range creators {
		ch := s.pubSub.Subscribe(userID(creator))
		go func() {
			for post := range ch {
				aggregator <- post
			}
		}()
	}

	return aggregator
}

func (s *TimelineService) Post(creator *repository.User, content string) {
	post := &repository.Post{
		Creator:   creator,
		Content:   content,
		CreatedAt: time.Now(),
	}

	s.postRepository.Save(post)
	s.pubSub.Publish(userID(creator), post)
}

func (s *TimelineService) Search(server string, query string) []*repository.Post {
	if len(server) == 0 {
		return s.postRepository.FindByContentContains(query)
	} else {
		// Erstelle eine gRPC-Verbindung zum entfernten Server
		conn, err := grpc.NewClient("dns:"+server+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := timeline.NewTimelineServiceClient(conn)

		// Erstelle einen Kontext mit Timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Führe die gRPC-Suche durch
		resp, err := client.Search(ctx, &timeline.SearchRequest{Query: query})
		if err != nil {
			log.Fatalf("could not search: %v", err)
		}

		// Konvertiere die Antwort in das interne Format
		posts := make([]*repository.Post, len(resp.Posts))
		for idx, post := range resp.Posts {
			// TODO Fehlerbehandlung
			parsedTime, _ := time.Parse(time.RFC3339, post.CreatedAt)
			posts[idx] = &repository.Post{
				Creator: &repository.User{
					Username: post.Username,
					Server:   server,
				},
				Content:   post.Content,
				CreatedAt: parsedTime,
			}
		}
		return posts
	}
}

func userID(user *repository.User) string {
	return fmt.Sprintf("creator%d", user.ID)
}

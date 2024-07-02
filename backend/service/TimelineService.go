package service

import (
	"context"
	"crickets-go/data"
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
	clientMap              map[string]timeline.TimelineServiceClient
}

func NewTimelineService(postRepository *repository.PostRepository, subscriptionRepository *repository.SubscriptionRepository, pubSub *PubSub) *TimelineService {
	return &TimelineService{
		postRepository:         postRepository,
		subscriptionRepository: subscriptionRepository,
		pubSub:                 pubSub,
		clientMap:              make(map[string]timeline.TimelineServiceClient),
	}
}

func (s *TimelineService) TimelineUpdates(subscriberID int32) chan *data.Post {
	subscriptions := s.subscriptionRepository.FindBySubscriberServerAndSubscriberID("", subscriberID)

	creatorIDsMap := make(map[string][]int32)

	// Eigene Posts automatisch sehen
	creatorIDsMap[""] = append(creatorIDsMap[""], subscriberID)

	for _, sub := range subscriptions {
		creator := sub.Creator
		creatorServer := creator.Server
		creatorIDsMap[creatorServer] = append(creatorIDsMap[creatorServer], creator.ID)
	}

	channels := make([]chan *data.Post, 0)

	for server, creatorIDs := range creatorIDsMap {
		var ch chan *data.Post
		var err error
		if len(server) == 0 {
			ch = s.LocalTimelineUpdates(creatorIDs)
		} else {
			ch, err = s.remoteTimelineUpdates(server, creatorIDs)
			if err != nil {
				log.Printf("Error fetching updates from server %s: %v", server, err)
				continue
			}
		}
		channels = append(channels, ch)
	}

	return aggregate(channels)
}

func (s *TimelineService) LocalTimelineUpdates(creatorIDs []int32) chan *data.Post {
	channels := make([]chan *data.Post, 0)
	for _, creatorID := range creatorIDs {
		ch := s.pubSub.Subscribe(userID(creatorID))
		channels = append(channels, ch)
	}
	return aggregate(channels)
}

func (s *TimelineService) remoteTimelineUpdates(server string, creatorIDs []int32) (chan *data.Post, error) {
	client, err := s.getClient(server)
	if err != nil {
		return nil, err
	}

	request := &timeline.TimelineUpdateRequest{
		CreatorIds: creatorIDs,
	}

	stream, err := client.TimelineUpdates(context.Background(), request)
	if err != nil {
		return nil, err
	}

	ch := make(chan *data.Post)
	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				close(ch)
				return
			}

			ch <- convertPost(server, response.Post)
		}
	}()

	return ch, nil
}

func (s *TimelineService) Post(creator *data.User, content string) error {
	post := &data.Post{
		Creator:   creator,
		Content:   content,
		CreatedAt: time.Now(),
	}

	s.postRepository.Save(post)
	return s.pubSub.Publish(userID(creator.ID), post)
}

func (s *TimelineService) Search(server string, query string) []*data.Post {
	if len(server) == 0 {
		return s.postRepository.FindByContentContains(query)
	} else {
		// Erstelle eine gRPC-Verbindung zum entfernten Server
		client, err := s.getClient(server)
		if err != nil {
			log.Printf("Error connecting to server %s: %v", server, err)
			return nil
		}

		// Erstelle einen Kontext mit Timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Führe die gRPC-Suche durch
		response, err := client.Search(ctx, &timeline.SearchRequest{Query: query})
		if err != nil {
			log.Fatalf("Error search posts from server %s: %v", server, err)
			return nil
		}

		// Konvertiere die Antwort in das interne Format
		posts := make([]*data.Post, len(response.Posts))
		for idx, post := range response.Posts {
			posts[idx] = convertPost(server, post)
		}
		return posts
	}
}

func (s *TimelineService) getClient(server string) (timeline.TimelineServiceClient, error) {
	timelineClient, exists := s.clientMap[server]

	if !exists {
		client, err := grpc.NewClient("dns:"+server+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}

		timelineClient = timeline.NewTimelineServiceClient(client)
		s.clientMap[server] = timelineClient
	}

	return timelineClient, nil
}

func convertPost(server string, post *timeline.Post) *data.Post {
	r := &data.Post{
		Creator: &data.User{
			Username: post.Username,
			Server:   server,
		},
		Content:   post.Content,
		CreatedAt: post.CreatedAt.AsTime(),
	}
	return r
}

func aggregate(channels []chan *data.Post) chan *data.Post {
	aggregator := make(chan *data.Post)
	for _, ch := range channels {
		go func() {
			for post := range ch {
				aggregator <- post
			}
		}()
	}
	return aggregator
}

func userID(userID int32) string {
	return fmt.Sprintf("creator%d", userID)
}

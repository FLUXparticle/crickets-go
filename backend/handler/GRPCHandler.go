package handler

import (
	"context"
	"crickets-go/data"
	"crickets-go/gen/timeline"
	"crickets-go/service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type GRPCHandler struct {
	timeline.UnimplementedTimelineServiceServer
	timelineService *service.TimelineService
}

func NewGRPCHandler(timelineService *service.TimelineService) *GRPCHandler {
	return &GRPCHandler{timelineService: timelineService}
}

func (s *GRPCHandler) Search(ctx context.Context, request *timeline.SearchRequest) (*timeline.SearchResponse, error) {
	posts := s.timelineService.Search("", request.Query)
	responsePosts := make([]*timeline.Post, len(posts))
	for i, post := range posts {
		responsePosts[i] = convertPost(post)
	}
	return &timeline.SearchResponse{Posts: responsePosts}, nil
}

func (s *GRPCHandler) TimelineUpdates(request *timeline.TimelineUpdateRequest, stream timeline.TimelineService_TimelineUpdatesServer) error {
	// Fetch updates from the TimelineService
	updatesChan := s.timelineService.LocalTimelineUpdates(request.CreatorIds)

	for {
		select {
		case update, ok := <-updatesChan:
			if !ok {
				return nil
			}
			post := convertPost(update)
			if err := stream.Send(&timeline.TimelineUpdateResponse{Post: post}); err != nil {
				log.Printf("Error sending update: %v", err)
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func (s *GRPCHandler) LikePost(ctx context.Context, request *timeline.LikePostRequest) (*timeline.Empty, error) {
	err := s.timelineService.LikePost(request.PostId, "")
	if err != nil {
		return nil, err
	}

	return &timeline.Empty{}, nil
}

func convertPost(update *data.Post) *timeline.Post {
	post := &timeline.Post{
		PostId:    update.ID,
		Username:  update.Creator.Username,
		Content:   update.Content,
		CreatedAt: timestamppb.New(update.CreatedAt),
	}
	return post
}

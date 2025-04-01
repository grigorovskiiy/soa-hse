package service

import (
	"auth/posts_service/internal/infrastructure/models"
	pb "auth/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type PostsRepository interface {
	CreatePost(*models.DbPost) error
	DeletePost(int32) error
	UpdatePost(post *models.DbPost) error
	GetPost(int32) (*models.DbPost, error)
	GetPostList(int32, int32, int32) ([]*models.DbPost, error)
}
type PService struct {
	repository PostsRepository
}

func NewUService(repository PostsRepository) *PService {
	return &PService{repository: repository}
}

func (s *PService) CreatePost(pb *pb.PostDataRequest, userId int32) error {
	post := models.DbPost{
		Name:         pb.PostName,
		Tags:         pb.Tags,
		Description:  pb.PostDescription,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserId:       int(userId),
		SecurityFlag: pb.SecurityFlag,
	}

	return s.repository.CreatePost(&post)
}

func (s *PService) DeletePost(pb *pb.PostID) error {
	return s.repository.DeletePost(pb.PostId)
}

func (s *PService) UpdatePost(pb *pb.UpdatePostRequest) error {
	post := models.DbPost{
		UpdatedAt:    time.Now(),
		Id:           int(pb.PostId.PostId),
		Name:         pb.PostData.PostName,
		Description:  pb.PostData.PostDescription,
		SecurityFlag: pb.PostData.SecurityFlag,
		Tags:         pb.PostData.Tags,
	}

	return s.repository.UpdatePost(&post)
}

func (s *PService) GetPost(p *pb.PostID) (*pb.PostDataResponse, error) {
	postInfo, err := s.repository.GetPost(p.PostId)
	if err != nil {
		return nil, err
	}

	pb := pb.PostDataResponse{
		PostId:          &pb.PostID{PostId: int32(postInfo.Id)},
		Tags:            postInfo.Tags,
		PostName:        postInfo.Name,
		PostDescription: postInfo.Description,
		SecurityFlag:    postInfo.SecurityFlag,
		CreatedAt:       timestamppb.New(postInfo.CreatedAt),
		UpdatedAt:       timestamppb.New(postInfo.UpdatedAt),
	}

	return &pb, nil
}

func (s *PService) GetPostList(p *pb.ListPostsRequest, userId int32) (*pb.ListPostsResponse, error) {
	posts, err := s.repository.GetPostList(p.Page, p.PageSize, userId)
	if err != nil {
		return nil, err
	}
	pbPosts := pb.ListPostsResponse{
		Posts: make([]*pb.PostDataResponse, len(posts)),
	}
	for ind, _ := range posts {
		pbPosts.Posts[ind] = &pb.PostDataResponse{
			PostId:          &pb.PostID{PostId: int32(posts[ind].Id)},
			Tags:            posts[ind].Tags,
			PostName:        posts[ind].Name,
			PostDescription: posts[ind].Description,
			SecurityFlag:    posts[ind].SecurityFlag,
			CreatedAt:       timestamppb.New(posts[ind].CreatedAt),
			UpdatedAt:       timestamppb.New(posts[ind].UpdatedAt),
		}
	}

	return &pbPosts, nil
}

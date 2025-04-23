package service

import (
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/models"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type PostsRepository interface {
	CreatePost(*models.DbPost) error
	DeletePost(int32, int32) error
	UpdatePost(*models.DbPost) error
	GetPost(int32, int32) (*models.DbPost, error)
	GetPostList(int32, int32, int32) ([]*models.DbPost, error)
}
type PService struct {
	repository PostsRepository
}

func NewPService(repository PostsRepository) *PService {
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

	if err := s.repository.CreatePost(&post); err != nil {
		infrastructure.Logger.Error("create post error", "error", err.Error())
		return err
	}

	return nil
}

func (s *PService) DeletePost(pb *pb.PostID, userId int32) error {
	if err := s.repository.DeletePost(pb.PostId, userId); err != nil {
		infrastructure.Logger.Error("delete post error", "error", err.Error())
		return err
	}

	return nil
}

func (s *PService) UpdatePost(pb *pb.UpdatePostRequest, userId int32) error {
	post := models.DbPost{
		UpdatedAt:    time.Now(),
		Id:           int(pb.PostId),
		Name:         pb.PostData.PostName,
		Description:  pb.PostData.PostDescription,
		SecurityFlag: pb.PostData.SecurityFlag,
		Tags:         pb.PostData.Tags,
		UserId:       int(userId),
	}

	if err := s.repository.UpdatePost(&post); err != nil {
		infrastructure.Logger.Error("update post error", "error", err.Error())
		return err
	}

	return nil
}

func (s *PService) GetPost(p *pb.PostID, userId int32) (*pb.PostDataResponse, error) {
	postInfo, err := s.repository.GetPost(p.PostId, userId)
	if err != nil {
		infrastructure.Logger.Error("get post info error", "error", err.Error())
		return nil, err
	}

	pb := pb.PostDataResponse{
		PostId:          int32(postInfo.Id),
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
		infrastructure.Logger.Error("get post list error", "error", err.Error())
		return nil, err
	}
	pbPosts := pb.ListPostsResponse{
		Posts: make([]*pb.PostDataResponse, len(posts)),
	}
	for ind, _ := range posts {
		pbPosts.Posts[ind] = &pb.PostDataResponse{
			PostId:          int32(posts[ind].Id),
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

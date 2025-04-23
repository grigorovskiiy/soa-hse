package models

import (
	pb "github.com/grigorovskiiy/soa-hse/protos"
)

func (m *GetPostListRequest) ToProto() *pb.ListPostsRequest {
	return &pb.ListPostsRequest{
		Page:     m.Page,
		PageSize: m.PageSize,
	}
}

func (m *CreatePostRequest) ToProto() *pb.PostDataRequest {
	return &pb.PostDataRequest{
		PostName:        m.PostName,
		PostDescription: m.PostDescription,
		Tags:            m.Tags,
		SecurityFlag:    m.SecurityFlag,
	}
}

func (m *GetDeletePostRequest) ToProto() *pb.PostID {
	return &pb.PostID{
		PostId: int32(m.PostID),
	}
}

func (m *UpdatePostRequest) ToProto() *pb.UpdatePostRequest {
	return &pb.UpdatePostRequest{
		PostId: int32(m.PostID),
		PostData: &pb.PostDataRequest{
			PostName:        m.PostName,
			PostDescription: m.PostDescription,
			SecurityFlag:    m.SecurityFlag,
			Tags:            m.Tags,
		},
	}
}

func FromProtoPostResponse(pb *pb.PostDataResponse) *GetPostResponse {
	return &GetPostResponse{
		PostId:          int(pb.GetPostId()),
		PostName:        pb.GetPostName(),
		PostDescription: pb.GetPostDescription(),
		SecurityFlag:    pb.GetSecurityFlag(),
		CreatedAt:       pb.GetCreatedAt().AsTime().Local(),
		UpdatedAt:       pb.GetUpdatedAt().AsTime().Local(),
		Tags:            pb.GetTags(),
		UserId:          int(pb.GetUserId()),
	}
}

func FromProtoListPostResponse(pb *pb.ListPostsResponse) *GetPostListResponse {
	posts := make([]*GetPostResponse, len(pb.Posts))
	for i, post := range pb.Posts {
		posts[i] = FromProtoPostResponse(post)
	}

	return &GetPostListResponse{
		Posts: posts,
	}
}

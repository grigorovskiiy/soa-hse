package models

import (
	pb "github.com/grigorovskiiy/soa-hse/protos"
)

func (m *CreatePostRequest) ToProto() *pb.PostDataRequest {
	return &pb.PostDataRequest{
		PostName:        m.PostName,
		PostDescription: m.PostDescription,
		Tags:            m.Tags,
		SecurityFlag:    m.SecurityFlag,
	}
}

func (m *PostID) ToProto() *pb.PostID {
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
		PostID:          int(pb.GetPostId()),
		PostName:        pb.GetPostName(),
		PostDescription: pb.GetPostDescription(),
		SecurityFlag:    pb.GetSecurityFlag(),
		CreatedAt:       pb.GetCreatedAt().AsTime().Local(),
		UpdatedAt:       pb.GetUpdatedAt().AsTime().Local(),
		Tags:            pb.GetTags(),
		UserID:          int(pb.GetUserId()),
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

func (m *PostCommentRequest) ToProto() *pb.PostCommentRequest {
	return &pb.PostCommentRequest{
		PostId:             int32(m.PostID),
		CommentDescription: m.Description,
	}
}

func FromProtoPostCommentResponse(pb *pb.CommentDataResponse) *GetCommentResponse {
	return &GetCommentResponse{
		CommentID:   int(pb.CommentId),
		PostID:      int(pb.PostId),
		UserID:      int(pb.UserId),
		Description: pb.CommentDescription,
	}
}

func FromProtoListCommentResponse(pb *pb.ListCommentsResponse) *GetCommentListResponse {
	comments := make([]*GetCommentResponse, len(pb.Comments))
	for i, comment := range pb.Comments {
		comments[i] = FromProtoPostCommentResponse(comment)
	}

	return &GetCommentListResponse{
		Comments: comments,
	}
}

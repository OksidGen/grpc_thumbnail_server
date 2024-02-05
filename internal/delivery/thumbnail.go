package delivery

import (
	"context"
	"github.com/OksidGen/grpc_thumbnail/server/internal/usecase"
	"github.com/OksidGen/grpc_thumbnail/server/proto"
)

type ThumbnailHandler struct {
	proto.UnimplementedThumbnailServiceServer
	usecase usecase.ThumbnailUsecase
}

func NewThumbnailHandler(usecase usecase.ThumbnailUsecase) *ThumbnailHandler {
	return &ThumbnailHandler{
		usecase: usecase,
	}
}

func (h *ThumbnailHandler) GetThumbnail(ctx context.Context, req *proto.ThumbnailRequest) (*proto.ThumbnailResponse, error) {
	url := req.VideoUrl
	thumbnailData, err := h.usecase.GetThumbnail(url)
	if err != nil {
		return nil, err
	}
	return &proto.ThumbnailResponse{ThumbnailData: thumbnailData}, nil
}

package recognition

import (
	"context"
	"go-scrfd-api/pkg/grpc/recognition/pb"
)

type handlerGrpc struct {
	pb.RecognitionServer
	uc UseCase
}

func NewGRPCHandler(recognitionUC UseCase) pb.RecognitionServer {
	return &handlerGrpc{uc: recognitionUC}
}

func (h *handlerGrpc) Recognize(ctx context.Context, req *pb.RecognizeRequest) (*pb.RecognizeResponse, error) {
	boxes, err := h.uc.Detect(req.GetImage())
	if err != nil {
		return nil, err
	}

	var res pb.RecognizeResponse
	for _, box := range boxes {
		res.Boxes = append(res.Boxes, &pb.PredictionBox{
			Score:     box.Score,
			BoxLeft:   box.BoxLeft,
			BoxTop:    box.BoxTop,
			BoxRight:  box.BoxRight,
			BoxBottom: box.BoxBottom,
			KeyPoints: box.KeyPoints,
		})
	}
	return &res, nil
}
